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
	Long: `Compare each WayPoint's SYNCED_AT date with the latest git commit time.

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
			Address   string
			SyncedAt  string
			SyncDate  time.Time // 정렬용
			Status    string
			State     string // OUTDATED, NO_SYNC
			Gap       string // 시간 차이 표시
		}

		var results []wpStatus
		outdatedCount := 0
		noSyncCount := 0

		syncRe := regexp.MustCompile(`(?m)^-\s*SYNCED_AT:\s*(.+)$`)
		statusRe := regexp.MustCompile(`(?m)^## \[STATUS\]\s+(\S+)`)

		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
				continue
			}

			content, err := os.ReadFile(filepath.Join(wpDir, f.Name()))
			if err != nil {
				continue
			}
			text := string(content)

			dotName := strings.TrimSuffix(f.Name(), ".md")
			path := strings.ReplaceAll(dotName, ".", "/")
			addr := "W://" + path

			// STATUS 추출
			status := "?"
			if m := statusRe.FindStringSubmatch(text); len(m) > 1 {
				status = m[1]
			}

			// SYNCED_AT 추출
			syncedAt := ""
			if m := syncRe.FindStringSubmatch(text); len(m) > 1 {
				syncedAt = strings.TrimSpace(m[1])
			}

			if syncedAt == "" {
				results = append(results, wpStatus{addr, "-", time.Time{}, status, "NO_SYNC", "-"})
				noSyncCount++
				continue
			}

			// SYNCED_AT과 git 최신 커밋 비교
			syncDate, err := time.Parse("2006-01-02", syncedAt)
			if err != nil {
				results = append(results, wpStatus{addr, syncedAt, time.Time{}, status, "NO_SYNC", "-"})
				noSyncCount++
				continue
			}

			// SYNCED_AT은 날짜만 있으므로 해당일 끝(23:59:59)으로 간주
			syncEnd := syncDate.Add(24*time.Hour - time.Second)

			if !lastCommit.IsZero() && lastCommit.After(syncEnd.Add(syncGracePeriod)) {
				gap := formatDuration(lastCommit.Sub(syncEnd))
				results = append(results, wpStatus{addr, syncedAt, syncDate, status, "OUTDATED", gap})
				outdatedCount++
			}
		}

		// 출력
		if len(results) == 0 {
			if lastCommit.IsZero() {
				fmt.Println("all waypoints synced (no git history found)")
			} else {
				fmt.Printf("all waypoints synced (last commit: %s)\n", lastCommit.Format("2006-01-02 15:04"))
			}
			return
		}

		// SYNCED_AT 최신순 정렬 (NO_SYNC는 뒤로)
		sort.Slice(results, func(i, j int) bool {
			if results[i].SyncDate.IsZero() && !results[j].SyncDate.IsZero() {
				return false
			}
			if !results[i].SyncDate.IsZero() && results[j].SyncDate.IsZero() {
				return true
			}
			return results[i].SyncDate.After(results[j].SyncDate)
		})

		// 최대 10개만 표시
		const maxDisplay = 10
		totalCount := len(results)
		if len(results) > maxDisplay {
			results = results[:maxDisplay]
		}

		fmt.Printf("last git commit: %s\n\n", lastCommit.Format("2006-01-02 15:04"))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ADDRESS\tSYNCED_AT\tSTATUS\tSTATE\tGAP")
		fmt.Fprintln(w, "-------\t---------\t------\t-----\t---")
		for _, r := range results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", r.Address, r.SyncedAt, r.Status, r.State, r.Gap)
		}
		w.Flush()

		if totalCount > maxDisplay {
			fmt.Printf("\n... and %d more (showing top %d by SYNCED_AT)\n", totalCount-maxDisplay, maxDisplay)
		}
		fmt.Printf("\n%d outdated, %d no sync date (%d total)\n", outdatedCount, noSyncCount, totalCount)
		if outdatedCount > 0 {
			fmt.Println("\n⚠ OUTDATED WayPoints may need SYNCED_AT update or TECH_SPEC review.")
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
