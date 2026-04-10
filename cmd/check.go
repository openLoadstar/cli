package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check WayPoint sync status against git history",
	Long: `Compare each WayPoint's SYNCED_AT date with the latest git commit time.

Reports WayPoints that may need updating because git has newer changes.

Examples:
  loadstar check`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		loadstarBase := fs.AvcsPath("")
		projectRoot := filepath.Dir(loadstarBase)

		// 1. git 최신 커밋 시간 가져오기
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
			SyncedAt string
			Status   string
			State    string // OK, OUTDATED, NO_SYNC
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

			// S_STB는 완료된 WP — 동기화 체크 불필요
			if status == "S_STB" {
				continue
			}

			// SYNCED_AT 추출
			syncedAt := ""
			if m := syncRe.FindStringSubmatch(text); len(m) > 1 {
				syncedAt = strings.TrimSpace(m[1])
			}

			if syncedAt == "" {
				results = append(results, wpStatus{addr, "-", status, "NO_SYNC"})
				noSyncCount++
				continue
			}

			// SYNCED_AT과 git 최신 커밋 비교
			syncDate, err := time.Parse("2006-01-02", syncedAt)
			if err != nil {
				results = append(results, wpStatus{addr, syncedAt, status, "NO_SYNC"})
				noSyncCount++
				continue
			}

			if !lastCommit.IsZero() && lastCommit.After(syncDate.Add(24*time.Hour)) {
				results = append(results, wpStatus{addr, syncedAt, status, "OUTDATED"})
				outdatedCount++
			}
		}

		// 출력
		if len(results) == 0 {
			if lastCommit.IsZero() {
				fmt.Println("all active waypoints synced (no git history found)")
			} else {
				fmt.Printf("all active waypoints synced (last commit: %s)\n", lastCommit.Format("2006-01-02 15:04"))
			}
			return
		}

		fmt.Printf("last git commit: %s\n\n", lastCommit.Format("2006-01-02 15:04"))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ADDRESS\tSYNCED_AT\tSTATUS\tSTATE")
		fmt.Fprintln(w, "-------\t---------\t------\t-----")
		for _, r := range results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.Address, r.SyncedAt, r.Status, r.State)
		}
		w.Flush()

		fmt.Printf("\n%d outdated, %d no sync date (%d total)\n", outdatedCount, noSyncCount, len(results))
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
