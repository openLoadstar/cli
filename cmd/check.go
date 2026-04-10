package cmd

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

const syncGracePeriod = 30 * time.Minute
const maxDisplay = 10

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Show recently modified WayPoints with git sync status",
	Long: `Show the 10 most recently modified WayPoint files and compare with the latest git commit time.

Helps identify which WayPoints were recently touched and whether they are in sync with git.

Examples:
  loadstar check`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		loadstarBase := fs.AvcsPath("")
		projectRoot := filepath.Dir(loadstarBase)

		lastCommit := getLastGitCommitTime(projectRoot)

		wpDir := filepath.Join(loadstarBase, "WAYPOINT")
		files, err := os.ReadDir(wpDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: cannot read WAYPOINT dir: %v\n", err)
			os.Exit(1)
		}

		type wpEntry struct {
			Address string
			ModTime time.Time
			Status  string
			State   string
			Gap     string
		}

		statusRe := regexp.MustCompile(`(?m)^## \[STATUS\]\s+(\S+)`)

		var all []wpEntry

		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
				continue
			}

			filePath := filepath.Join(wpDir, f.Name())
			info, err := os.Stat(filePath)
			if err != nil {
				continue
			}
			modTime := info.ModTime()

			content, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			dotName := strings.TrimSuffix(f.Name(), ".md")
			path := strings.ReplaceAll(dotName, ".", "/")
			addr := "W://" + path

			status := "?"
			if m := statusRe.FindStringSubmatch(string(content)); len(m) > 1 {
				status = m[1]
			}

			state := "SYNCED"
			gap := "-"
			if !lastCommit.IsZero() && lastCommit.After(modTime.Add(syncGracePeriod)) {
				state = "OUTDATED"
				gap = formatDuration(lastCommit.Sub(modTime))
			}

			all = append(all, wpEntry{addr, modTime, status, state, gap})
		}

		// modTime 최신순 정렬
		sort.Slice(all, func(i, j int) bool {
			return all[i].ModTime.After(all[j].ModTime)
		})

		// 최근 10개만
		display := all
		if len(display) > maxDisplay {
			display = display[:maxDisplay]
		}

		// 통계
		outdatedCount := 0
		for _, e := range all {
			if e.State == "OUTDATED" {
				outdatedCount++
			}
		}

		if lastCommit.IsZero() {
			fmt.Println("no git history found")
		} else {
			fmt.Printf("last git commit: %s\n", lastCommit.Format("2006-01-02 15:04"))
		}
		fmt.Printf("waypoints: %d total, %d outdated\n\n", len(all), outdatedCount)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ADDRESS\tMODIFIED\tSTATUS\tSTATE\tGAP")
		fmt.Fprintln(w, "-------\t--------\t------\t-----\t---")
		for _, r := range display {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				r.Address, r.ModTime.Format("2006-01-02 15:04"), r.Status, r.State, r.Gap)
		}
		w.Flush()

		if len(all) > maxDisplay {
			fmt.Printf("\n(showing %d of %d, sorted by last modified)\n", maxDisplay, len(all))
		}
	},
}

func getLastGitCommitTime(projectRoot string) time.Time {
	cmd := exec.Command("git", "log", "-1", "--format=%aI")
	cmd.Dir = projectRoot
	out, err := cmd.Output()
	if err != nil {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, strings.TrimSpace(string(out)))
	if err != nil {
		return time.Time{}
	}
	return t
}

func formatDuration(d time.Duration) string {
	hours := d.Hours()
	if hours < 1 {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if hours < 24 {
		return fmt.Sprintf("%.0fh", math.Floor(hours))
	}
	days := int(hours / 24)
	remainH := int(hours) % 24
	if remainH == 0 {
		return fmt.Sprintf("%dd", days)
	}
	return fmt.Sprintf("%dd%dh", days, remainH)
}
