package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

const todoListFile = ".clionly/TODO/TODO_LIST.md"
const wpSnapshotFile = ".clionly/TODO/WP_SNAPSHOT.json"

// wpSnapshot holds cached file info for change detection.
type wpSnapshot struct {
	ModTime string `json:"modTime"`
	Size    int64  `json:"size"`
	Status  string `json:"status"`
}

// todoItem represents a single TODO list entry.
type todoItem struct {
	Address string
	Status  string // PENDING, ACTIVE, [BLOCKED]
	Summary string
}

var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Manage TODO items via WayPoint sync",
	Long: `Manage the project TODO list by syncing with WayPoint STATUS.

Subcommands:
  sync     Scan WayPoint files and update TODO_LIST
  list     Show current PENDING/ACTIVE/BLOCKED items
  history  Show completed TECH_SPEC items from WayPoints

Examples:
  loadstar todo sync
  loadstar todo sync W://root/cli/cmd_log
  loadstar todo list
  loadstar todo history
  loadstar todo history M://root/cli`,
}

// ===== sync =====

var todoSyncCmd = &cobra.Command{
	Use:   "sync [ADDRESS]",
	Short: "Sync TODO_LIST with WayPoint STATUS",
	Long: `Scan WayPoint files and update TODO_LIST based on their STATUS.

  Without arguments: scan all WayPoints via MAP traversal.
  With ADDRESS: sync a single WayPoint.

Examples:
  loadstar todo sync
  loadstar todo sync W://root/cli/cmd_log`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		loadstarBase := fs.AvcsPath("")

		if len(args) == 1 {
			// Single WP sync
			syncSingleWP(loadstarBase, args[0])
			return
		}

		// Full sync: collect all WP addresses via MAP traversal
		wpAddrs := collectAllWaypoints(loadstarBase)
		if len(wpAddrs) == 0 {
			fmt.Println("no waypoints found")
			return
		}

		snapshot := loadSnapshot(loadstarBase)
		items := loadTodoList(loadstarBase)
		itemMap := make(map[string]*todoItem)
		for i := range items {
			itemMap[items[i].Address] = &items[i]
		}

		// Track which addresses are still valid
		validAddrs := make(map[string]bool)
		added, updated, removed := 0, 0, 0

		for _, addr := range wpAddrs {
			validAddrs[addr] = true
			wpFile := addressToFilePath(loadstarBase, addr)

			info, err := os.Stat(wpFile)
			if err != nil {
				continue
			}

			modTime := info.ModTime().Format(time.RFC3339)
			size := info.Size()

			// Check if changed since last snapshot
			cached, hasCached := snapshot[addr]
			if hasCached && cached.ModTime == modTime && cached.Size == size {
				continue // unchanged
			}

			// Read STATUS and SUMMARY from WP file
			status, summary := readWPStatusAndSummary(wpFile)

			// Update snapshot
			snapshot[addr] = wpSnapshot{ModTime: modTime, Size: size, Status: status}

			todoStatus := wpStatusToTodoStatus(status)

			if todoStatus == "" {
				// S_STB → remove from TODO
				if _, exists := itemMap[addr]; exists {
					delete(itemMap, addr)
					removed++
				}
			} else {
				if existing, exists := itemMap[addr]; exists {
					if existing.Status != todoStatus || existing.Summary != summary {
						existing.Status = todoStatus
						existing.Summary = summary
						updated++
					}
				} else {
					itemMap[addr] = &todoItem{Address: addr, Status: todoStatus, Summary: summary}
					added++
				}
			}
		}

		// Remove TODO entries for deleted WPs
		for addr := range itemMap {
			if !validAddrs[addr] {
				delete(itemMap, addr)
				removed++
			}
		}

		// Apply BLOCKED detection via REFERENCE
		applyBlocked(loadstarBase, itemMap)

		// Rebuild items list
		var result []todoItem
		for _, item := range itemMap {
			result = append(result, *item)
		}

		saveTodoList(loadstarBase, result)
		saveSnapshot(loadstarBase, snapshot)

		fmt.Printf("sync complete: %d added, %d updated, %d removed (%d total)\n",
			added, updated, removed, len(result))
	},
}

func syncSingleWP(loadstarBase, addr string) {
	wpFile := addressToFilePath(loadstarBase, addr)
	if _, err := os.Stat(wpFile); err != nil {
		fmt.Fprintf(os.Stderr, "error: WayPoint not found: %s\n", addr)
		os.Exit(1)
	}

	snapshot := loadSnapshot(loadstarBase)
	items := loadTodoList(loadstarBase)
	itemMap := make(map[string]*todoItem)
	for i := range items {
		itemMap[items[i].Address] = &items[i]
	}

	info, _ := os.Stat(wpFile)
	status, summary := readWPStatusAndSummary(wpFile)
	snapshot[addr] = wpSnapshot{
		ModTime: info.ModTime().Format(time.RFC3339),
		Size:    info.Size(),
		Status:  status,
	}

	todoStatus := wpStatusToTodoStatus(status)
	if todoStatus == "" {
		delete(itemMap, addr)
	} else {
		itemMap[addr] = &todoItem{Address: addr, Status: todoStatus, Summary: summary}
	}

	applyBlocked(loadstarBase, itemMap)

	var result []todoItem
	for _, item := range itemMap {
		result = append(result, *item)
	}

	saveTodoList(loadstarBase, result)
	saveSnapshot(loadstarBase, snapshot)

	fmt.Printf("synced: %s [%s]\n", addr, status)
}

// ===== list =====

var todoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all PENDING/ACTIVE TODO items",
	Run: func(cmd *cobra.Command, args []string) {
		loadstarBase := fs.AvcsPath("")
		items := loadTodoList(loadstarBase)

		if len(items) == 0 {
			fmt.Println("no TODO items")
			return
		}

		// Sort: ACTIVE first, then PENDING, then BLOCKED
		statusOrder := map[string]int{"ACTIVE": 0, "PENDING": 1, "[BLOCKED]": 2}
		sort.Slice(items, func(i, j int) bool {
			oi, oj := statusOrder[items[i].Status], statusOrder[items[j].Status]
			if oi != oj {
				return oi < oj
			}
			return items[i].Address < items[j].Address
		})

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ADDRESS\tSTATUS\tSUMMARY")
		fmt.Fprintln(w, "-------\t------\t-------")
		for _, item := range items {
			summary := item.Summary
			if len(summary) > 60 {
				summary = summary[:60] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", item.Address, item.Status, summary)
		}
		w.Flush()
		fmt.Printf("\n%d item(s)\n", len(items))
	},
}

// ===== history =====

var todoHistoryCmd = &cobra.Command{
	Use:   "history [MAP_ADDRESS]",
	Short: "Show completed TECH_SPEC items from WayPoints",
	Long: `Collect completed TECH_SPEC items ([x] YYYY-MM-DD ...) from WayPoints
and display them sorted newest-first.

Without arguments: scan all WayPoints.
With MAP_ADDRESS: scan only WayPoints under that Map.

Examples:
  loadstar todo history
  loadstar todo history M://root/cli`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		loadstarBase := fs.AvcsPath("")

		var wpAddrs []string
		if len(args) == 1 {
			wpAddrs = collectWaypointsUnderMap(loadstarBase, args[0])
		} else {
			wpAddrs = collectAllWaypoints(loadstarBase)
		}

		type histEntry struct {
			Address string
			Date    string
			Item    string
		}

		doneRe := regexp.MustCompile(`^\s*-\s*\[x\]\s*(\d{4}-\d{2}-\d{2})\s+(.+)$`)
		var entries []histEntry

		for _, addr := range wpAddrs {
			wpFile := addressToFilePath(loadstarBase, addr)
			data, err := os.ReadFile(wpFile)
			if err != nil {
				continue
			}

			inTodo := false
			for _, line := range strings.Split(string(data), "\n") {
				trimmed := strings.TrimSpace(line)
				if trimmed == "### TODO" || trimmed == "- TECH_SPEC:" {
					inTodo = true
					continue
				}
				if inTodo && strings.HasPrefix(trimmed, "###") {
					inTodo = false
					continue
				}
				if inTodo {
					m := doneRe.FindStringSubmatch(line)
					if m != nil {
						entries = append(entries, histEntry{Address: addr, Date: m[1], Item: m[2]})
					}
				}
			}
		}

		if len(entries) == 0 {
			fmt.Println("no completed items found")
			return
		}

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Date > entries[j].Date
		})

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ADDRESS\tDATE\tITEM")
		fmt.Fprintln(w, "-------\t----\t----")
		for _, e := range entries {
			item := e.Item
			if len(item) > 60 {
				item = item[:60] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", e.Address, e.Date, item)
		}
		w.Flush()
		fmt.Printf("\n%d completed item(s)\n", len(entries))
	},
}

func init() {
	todoCmd.AddCommand(todoSyncCmd)
	todoCmd.AddCommand(todoListCmd)
	todoCmd.AddCommand(todoHistoryCmd)
}

// ===== helpers =====

func wpStatusToTodoStatus(wpStatus string) string {
	switch wpStatus {
	case "S_IDL":
		return "PENDING"
	case "S_PRG":
		return "ACTIVE"
	case "S_ERR", "S_REV":
		return "ACTIVE"
	default: // S_STB or unknown
		return ""
	}
}

func addressToFilePath(loadstarBase, addr string) string {
	var typeDir, pathPart string
	if strings.HasPrefix(addr, "M://") {
		typeDir = "MAP"
		pathPart = addr[4:]
	} else if strings.HasPrefix(addr, "W://") {
		typeDir = "WAYPOINT"
		pathPart = addr[4:]
	} else {
		return ""
	}
	dotName := strings.ReplaceAll(pathPart, "/", ".")
	return filepath.Join(loadstarBase, typeDir, dotName+".md")
}

func readWPStatusAndSummary(filePath string) (status, summary string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "?", ""
	}
	statusRe := regexp.MustCompile(`##\s*\[STATUS\]\s*(\S+)`)
	summaryRe := regexp.MustCompile(`(?m)^-\s*SUMMARY:\s*(.*)$`)

	if m := statusRe.FindStringSubmatch(string(data)); len(m) >= 2 {
		status = m[1]
	} else {
		status = "?"
	}
	if m := summaryRe.FindStringSubmatch(string(data)); len(m) >= 2 {
		summary = strings.TrimSpace(m[1])
	}
	return
}

// collectAllWaypoints traverses MAP files to collect all WP addresses.
func collectAllWaypoints(loadstarBase string) []string {
	var result []string
	mapDir := filepath.Join(loadstarBase, "MAP")
	files, err := os.ReadDir(mapDir)
	if err != nil {
		return nil
	}
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(mapDir, f.Name()))
		if err != nil {
			continue
		}
		result = append(result, extractWaypointAddresses(string(data))...)
	}
	return result
}

// collectWaypointsUnderMap collects WP addresses from a specific Map and its sub-Maps.
func collectWaypointsUnderMap(loadstarBase, mapAddr string) []string {
	mapFile := addressToFilePath(loadstarBase, mapAddr)
	data, err := os.ReadFile(mapFile)
	if err != nil {
		return nil
	}

	var result []string
	inWP := false
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "### WAYPOINTS" {
			inWP = true
			continue
		}
		if inWP && strings.HasPrefix(trimmed, "###") {
			break
		}
		if inWP && strings.HasPrefix(trimmed, "- ") {
			addr := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
			if strings.HasPrefix(addr, "W://") {
				result = append(result, addr)
			} else if strings.HasPrefix(addr, "M://") {
				// Recurse into sub-Map
				result = append(result, collectWaypointsUnderMap(loadstarBase, addr)...)
			}
		}
	}
	return result
}

func extractWaypointAddresses(content string) []string {
	var result []string
	inWP := false
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "### WAYPOINTS" {
			inWP = true
			continue
		}
		if inWP && strings.HasPrefix(trimmed, "###") {
			break
		}
		if inWP && strings.HasPrefix(trimmed, "- W://") {
			addr := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
			result = append(result, addr)
		}
	}
	return result
}

// applyBlocked checks REFERENCE fields and marks items as [BLOCKED] if ref targets are not S_STB.
func applyBlocked(loadstarBase string, itemMap map[string]*todoItem) {
	for addr, item := range itemMap {
		if item.Status == "" {
			continue
		}
		wpFile := addressToFilePath(loadstarBase, addr)
		data, err := os.ReadFile(wpFile)
		if err != nil {
			continue
		}

		refs := extractReferences(string(data))
		blocked := false
		for _, ref := range refs {
			refFile := addressToFilePath(loadstarBase, ref)
			refStatus, _ := readWPStatusAndSummary(refFile)
			if refStatus != "S_STB" {
				blocked = true
				break
			}
		}

		if blocked && item.Status != "[BLOCKED]" {
			item.Status = "[BLOCKED]"
		} else if !blocked && item.Status == "[BLOCKED]" {
			// Un-block: revert to status based on WP
			status, _ := readWPStatusAndSummary(wpFile)
			item.Status = wpStatusToTodoStatus(status)
		}
	}
}

func extractReferences(content string) []string {
	re := regexp.MustCompile(`(?m)^-\s*REFERENCE:\s*\[([^\]]*)\]`)
	m := re.FindStringSubmatch(content)
	if m == nil || strings.TrimSpace(m[1]) == "" {
		return nil
	}
	var result []string
	for _, ref := range strings.Split(m[1], ",") {
		ref = strings.TrimSpace(ref)
		if ref != "" && strings.Contains(ref, "://") {
			result = append(result, ref)
		}
	}
	return result
}

// ===== file I/O =====

func loadSnapshot(loadstarBase string) map[string]wpSnapshot {
	path := filepath.Join(loadstarBase, wpSnapshotFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return make(map[string]wpSnapshot)
	}
	var result map[string]wpSnapshot
	if err := json.Unmarshal(data, &result); err != nil {
		return make(map[string]wpSnapshot)
	}
	return result
}

func saveSnapshot(loadstarBase string, snapshot map[string]wpSnapshot) {
	path := filepath.Join(loadstarBase, wpSnapshotFile)
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	data, _ := json.MarshalIndent(snapshot, "", "  ")
	_ = os.WriteFile(path, data, 0644)
}

func loadTodoList(loadstarBase string) []todoItem {
	path := filepath.Join(loadstarBase, todoListFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var items []todoItem
	for _, line := range strings.Split(string(data), "\n") {
		if !isDataRow(line) {
			continue
		}
		addr := extractCol(line, 0)
		status := extractCol(line, 1)
		summary := extractCol(line, 2)
		if addr != "" {
			items = append(items, todoItem{Address: addr, Status: status, Summary: summary})
		}
	}
	return items
}

func saveTodoList(loadstarBase string, items []todoItem) {
	path := filepath.Join(loadstarBase, todoListFile)
	_ = os.MkdirAll(filepath.Dir(path), 0755)

	header := "| 주소 (Address) | 상태 (Status) | 작업 요약 (Summary) |"
	sep := "| :--- | :--- | :--- |"

	var lines []string
	lines = append(lines, header, sep)
	for _, item := range items {
		lines = append(lines, fmt.Sprintf("| %s | %s | %s |", item.Address, item.Status, item.Summary))
	}

	_ = os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

// isDataRow returns true for markdown table data rows (not header/separator).
func isDataRow(line string) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") {
		return false
	}
	if strings.Contains(line, ":---") {
		return false
	}
	if strings.Contains(line, "주소 (Address)") {
		return false
	}
	return true
}

// extractCol extracts the Nth column (0-indexed) from a markdown table row.
func extractCol(line string, n int) string {
	parts := strings.Split(line, "|")
	idx := n + 1
	if idx >= len(parts) {
		return ""
	}
	return strings.TrimSpace(parts[idx])
}
