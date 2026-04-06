package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bono/loadstar/internal/git"
	"github.com/spf13/cobra"
)

var checkpointCmd = &cobra.Command{
	Use:   "checkpoint",
	Short: "Commit .loadstar/ metadata atomically",
	Long: `Atomically commit .loadstar/ metadata in a single git commit.
Changed element addresses are automatically listed in the commit message.
If a remote is configured (loadstar git set), the commit is pushed automatically.

Use --auto to mark the commit as an automatic checkpoint (adds [AUTO-CHECKPOINT] prefix).

Examples:
  loadstar checkpoint -m "implement cmd_log"
  loadstar checkpoint -m "fix: appendToContains multiline parsing"
  loadstar checkpoint --auto -m "periodic auto save"`,
	Run: func(cmd *cobra.Command, args []string) {
		message, _ := cmd.Flags().GetString("message")
		auto, _ := cmd.Flags().GetBool("auto")
		loadstarBase := fs.AvcsPath("")

		// Collect changed .loadstar/ files for commit message enrichment
		changedFiles, _ := gitClient.ChangedLoadstarFiles()
		commitMsg := buildCheckpointMessage(message, auto, changedFiles)

		// Commit via git — Atomic: abort if commit fails
		hash, err := gitClient.Commit(commitMsg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: git commit failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("committed: %s\n", hash)

		// Delete checkpoint_needed.flag if exists
		flagFile := filepath.Join(loadstarBase, ".clionly", "MONITOR", "checkpoint_needed.flag")
		if fs.Exists(flagFile) {
			_ = os.Remove(flagFile)
			fmt.Println("cleared: checkpoint_needed.flag")
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

// buildCheckpointMessage constructs a commit message with the user message,
// optional [AUTO-CHECKPOINT] prefix, and a list of changed .loadstar/ elements.
func buildCheckpointMessage(userMsg string, auto bool, changedFiles []string) string {
	var sb strings.Builder

	if auto {
		sb.WriteString("[AUTO-CHECKPOINT] ")
	}
	sb.WriteString(userMsg)

	// Filter to element directories only
	elementDirs := map[string]bool{
		"WAYPOINT": true, "BLACKBOX": true, "MAP": true,
	}
	var elements []string
	for _, f := range changedFiles {
		parts := strings.SplitN(f, "/", 3)
		if len(parts) >= 3 && elementDirs[parts[1]] {
			elements = append(elements, parts[1]+"/"+parts[2])
		}
	}

	if len(elements) > 0 {
		sort.Strings(elements)
		sb.WriteString("\n\n변경 요소:\n")
		for _, e := range elements {
			sb.WriteString("- ")
			sb.WriteString(e)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func init() {
	checkpointCmd.Flags().StringP("message", "m", "", "Checkpoint message")
	checkpointCmd.MarkFlagRequired("message")
	checkpointCmd.Flags().Bool("auto", false, "Mark as automatic checkpoint (adds [AUTO-CHECKPOINT] prefix)")
}
