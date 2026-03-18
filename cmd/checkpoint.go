package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/bono/loadstar/internal/git"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

var checkpointCmd = &cobra.Command{
	Use:   "checkpoint",
	Short: "Commit code and metadata atomically, and update SavePoint with git hash",
	Run: func(cmd *cobra.Command, args []string) {
		message, _ := cmd.Flags().GetString("message")
		loadstarBase := fs.AvcsPath("")

		// Commit via git — Atomic: abort if commit fails
		hash, err := gitClient.Commit(message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: git commit failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("committed: %s\n", hash)

		// Update ACTIVE SavePoint files with commit hash
		spDir := filepath.Join(loadstarBase, "SAVEPOINT")
		files, _ := fs.ListByPrefix(spDir, "")
		for _, f := range files {
			content, err := fs.Read(f)
			if err != nil {
				continue
			}
			if !strings.Contains(content, "S_ACT") {
				continue
			}
			updated := content + fmt.Sprintf("- git: %s\n", hash)
			_ = fs.Write(f, updated)
		}

		fmt.Printf("checkpoint complete: %s\n", message)

		// Auto-push if remote is configured and LOADSTAR_GIT_TOKEN is set.
		gc := git.NewClient(fs.Root)
		if err := gc.Push(); err != nil {
			msg := err.Error()
			if strings.Contains(msg, "no remote configured") ||
				strings.Contains(msg, "LOADSTAR_GIT_TOKEN") {
				// No remote / no token — skip silently.
			} else {
				fmt.Fprintf(os.Stderr, "warning: push failed: %v\n", err)
			}
		} else {
			fmt.Println("pushed to remote")
		}
	},
}

var historyCmd = &cobra.Command{
	Use:   "history [ADDRESS]",
	Short: "Show change history of an element",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addr, err := svc.ParseAddress(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		loadstarBase := fs.AvcsPath("")
		histDir := filepath.Join(loadstarBase, "HISTORY")
		dotName := strings.ReplaceAll(addr.Path, "/", ".")
		prefix := dotName + "_"

		entries, err := fs.ListByPrefix(histDir, prefix)
		if err != nil || len(entries) == 0 {
			fmt.Printf("no history found for %s\n", args[0])
			return
		}

		// Sort by filename descending (newest first)
		sort.Slice(entries, func(i, j int) bool {
			return entries[i] > entries[j]
		})

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "H_ID\tSIZE")
		fmt.Fprintln(w, "----\t----")
		for _, e := range entries {
			base := filepath.Base(e)
			hID := strings.TrimSuffix(base, ".md")
			info, _ := os.Stat(e)
			size := int64(0)
			if info != nil {
				size = info.Size()
			}
			fmt.Fprintf(w, "%s\t%d bytes\n", hID, size)
		}
		w.Flush()
	},
}

var diffCmd = &cobra.Command{
	Use:   "diff [ADDRESS] [H_ID]",
	Short: "Compare current element with a history snapshot",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		addr, err := svc.ParseAddress(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		loadstarBase := fs.AvcsPath("")
		currentFile := addr.ToFilePath(loadstarBase)
		histFile := filepath.Join(loadstarBase, "HISTORY", args[1]+".md")

		if !fs.Exists(currentFile) {
			fmt.Fprintf(os.Stderr, "error: element not found: %s\n", args[0])
			os.Exit(1)
		}
		if !fs.Exists(histFile) {
			fmt.Fprintf(os.Stderr, "error: history snapshot not found: %s\n", args[1])
			os.Exit(1)
		}

		current, _ := fs.Read(currentFile)
		hist, _ := fs.Read(histFile)

		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(hist, current, false)
		diffs = dmp.DiffCleanupSemantic(diffs)

		for _, d := range diffs {
			lines := strings.Split(d.Text, "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}
				switch d.Type {
				case diffmatchpatch.DiffInsert:
					fmt.Printf("\033[32m+ %s\033[0m\n", line)
				case diffmatchpatch.DiffDelete:
					fmt.Printf("\033[31m- %s\033[0m\n", line)
				case diffmatchpatch.DiffEqual:
					fmt.Printf("  %s\n", line)
				}
			}
		}
	},
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback [ADDRESS] [H_ID]",
	Short: "Rollback an element to a previous history snapshot",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		addr, err := svc.ParseAddress(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		loadstarBase := fs.AvcsPath("")
		currentFile := addr.ToFilePath(loadstarBase)
		histFile := filepath.Join(loadstarBase, "HISTORY", args[1]+".md")

		if !fs.Exists(currentFile) {
			fmt.Fprintf(os.Stderr, "error: element not found: %s\n", args[0])
			os.Exit(1)
		}
		if !fs.Exists(histFile) {
			fmt.Fprintf(os.Stderr, "error: history snapshot not found: %s\n", args[1])
			os.Exit(1)
		}

		if !force {
			fmt.Printf("rollback %s to %s? [y/N] ", args[0], args[1])
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" && answer != "yes" {
				fmt.Println("aborted")
				return
			}
		}

		// Pre-rollback backup
		ts := time.Now().Format("20060102T150405")
		dotName := strings.ReplaceAll(addr.Path, "/", ".")
		preBackup := filepath.Join(loadstarBase, "HISTORY", dotName+"_"+ts+"_pre_rollback.md")
		if err := fs.CopyFile(currentFile, preBackup); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not create pre-rollback backup: %v\n", err)
		}

		// Restore from snapshot
		if err := fs.CopyFile(histFile, currentFile); err != nil {
			fmt.Fprintf(os.Stderr, "error: rollback failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("rolled back: %s -> %s\n", args[0], args[1])
		fmt.Println("note: CONTAINS, LINEAGE, and LINKS are not automatically restored — verify manually")
	},
}

func init() {
	checkpointCmd.Flags().StringP("message", "m", "", "Checkpoint message")
	checkpointCmd.MarkFlagRequired("message")
	rollbackCmd.Flags().Bool("force", false, "Skip confirmation prompt")
}
