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

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check WayPoint sync status against git history",
	Long: `Compare each WayPoint file's last modified time with the latest git commit time.

Reports WayPoints that may need updating because git has newer changes.
Grace period: 30 minutes — changes within this window are considered synced.

Examples:
  loadstar check`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		loadstarBase := fs.AvcsPath("")
		projectRoot := filepath.Dir(loadstarBase)

		// 1. git 최신 커밋 시간
		lastCommit := getLastGitCommitTime(projectRoot)

		// 2. WP 스캔
		wpDir := filepath.Join(loadstarBase, "WAYPOINT")
		files, err := os.ReadDir(wpDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: cannot read WAYPOINT dir: %v\n", err)
			os.Exit(1)
		}

		type wpStatus struct {
			Address  string
			ModTime  time.Time
			Status   string
			State    string // OUTDATED, SYNCED
			Gap      string
		}

		statusRe := regexp.MustCompile(`(?m)^## \[STATUS\]\s+(\S+)`)

		var outdated []wpStatus
		var synced []wpStatus
		totalWP := 0

		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
				continue
			}
			totalWP++

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

			if !lastCommit.IsZero() && lastCommit.After(modTime.Add(syncGracePeriod)) {
				gap := formatDuration(lastCommit.Sub(modTime))
				outdated = append(outdated, wpStatus{addr, modTime, status, "OUTDATED", gap})
			} else {
				synced = append(synced, wpStatus{addr, modTime, status, "SYNCED", "-"})
			}
		}

		// 결과 합산: OUTDATED 먼저 (modTime 최신순), 그 다음 SYNCED
		sort.Slice(outdated, func(i, j int) bool {
			return outdated[i].ModTime.After(outdated[j].ModTime)
		})
		sort.Slice(synced, func(i, j int) bool {
			return synced[i].ModTime.After(synced[j].ModTime)
		})

		results := append(outdated, synced...)

		// 최대 10개 표시
		const maxDisplay = 10
		totalResults := len(results)
		displayResults := results
		if len(displayResults) > maxDisplay {
			displayResults = displayResults[:maxDisplay]
		}

		if lastCommit.IsZero() {
			fmt.Println("no git history found")
		} else {
			fmt.Printf("last git commit: %s\n", lastCommit.Format("2006-01-02 15:04"))
		}
		fmt.Printf("waypoints: %d total, %d outdated, %d synced\n\n", totalWP, len(outdated), len(synced))

		if len(displayResults) == 0 {
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ADDRESS\tMODIFIED\tSTATUS\tSTATE\tGAP")
		fmt.Fprintln(w, "-------\t--------\t------\t-----\t---")
		for _, r := range displayResults {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				r.Address, r.ModTime.Format("2006-01-02 15:04"), r.Status, r.State, r.Gap)
		}
		w.Flush()

		if totalResults > maxDisplay {
			fmt.Printf("\n... and %d more (showing top %d)\n", totalResults-maxDisplay, maxDisplay)
		}
		if len(outdated) > 0 {
			fmt.Println("\n⚠ OUTDATED WayPoints may need TECH_SPEC review.")
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
