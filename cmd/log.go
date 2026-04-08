package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var validLogKinds = map[string]bool{
	"NOTE": true, "DECISION": true, "ISSUE": true,
	"RESOLVED": true, "PROGRESS": true, "MODIFIED": true,
}

// logEntry represents a parsed log record from a .clionly/LOG/ file.
type logEntry struct {
	Timestamp string
	Kind      string
	Content   string
	Address   string // W:// or M:// address of the target element
}

var logCmd = &cobra.Command{
	Use:   "log [ADDRESS_OR_ID] [KIND] [CONTENT]",
	Short: "Append a log entry for an element",
	Long: `Append a structured log entry for the specified element.

ADDRESS can be a full address (W://root/cli/cmd_log) or a short ID (cmd_log).
If the ID matches multiple elements, candidates are listed for disambiguation.
KIND defaults to NOTE if omitted (2-arg form).

Use --list to browse available elements and their IDs.

KIND values:
  NOTE      General observation or memo
  DECISION  Design or implementation decision and its rationale
  ISSUE     Discovered problem, constraint, or bug
  RESOLVED  Resolution of a previously logged ISSUE
  PROGRESS  Implementation progress checkpoint
  MODIFIED  Direct md edit by AI or developer (record what changed)

Examples:
  loadstar log cmd_log "CODE_MAP 갱신 필요"                       # ID + NOTE (default)
  loadstar log cmd_log ISSUE "multiline 파싱 실패"                # ID + KIND
  loadstar log W://root/cli/cmd_log NOTE "full address 방식"      # full address
  loadstar log --list                                             # show all elements
  loadstar log --list calc                                        # filter by keyword`,
	Args: cobra.RangeArgs(0, 3),
	Run: func(cmd *cobra.Command, args []string) {
		listFlag, _ := cmd.Flags().GetBool("list")

		// --list mode: show element index
		if listFlag {
			filter := ""
			page := 1
			for _, a := range args {
				if n, err := strconv.Atoi(a); err == nil && n > 0 {
					page = n
				} else {
					filter = a
				}
			}
			showElementIndex(filter, page)
			return
		}

		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "error: at least 2 arguments required: [ADDRESS_OR_ID] [CONTENT]")
			fmt.Fprintln(os.Stderr, "  or use --list to browse elements")
			os.Exit(1)
		}

		// Parse flexible args: 2-arg (ID CONTENT) or 3-arg (ID KIND CONTENT)
		var addrInput, kind, content string
		if len(args) == 2 {
			addrInput = args[0]
			kind = "NOTE"
			content = args[1]
		} else {
			addrInput = args[0]
			maybeKind := strings.ToUpper(args[1])
			if validLogKinds[maybeKind] {
				kind = maybeKind
				content = args[2]
			} else {
				fmt.Fprintf(os.Stderr, "error: invalid KIND %q — allowed: NOTE, DECISION, ISSUE, RESOLVED, PROGRESS, MODIFIED\n", args[1])
				os.Exit(1)
			}
		}

		// Resolve address: full address or ID lookup
		addrStr := resolveAddress(addrInput)
		if addrStr == "" {
			os.Exit(1)
		}

		_, err := svc.ParseAddress(addrStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		loadstarBase := fs.AvcsPath("")
		writeLogChangeLog(loadstarBase, addrStr, kind, content)
		fmt.Printf("logged: [%s] %s\n", kind, addrStr)
	},
}

var findlogCmd = &cobra.Command{
	Use:   "findlog [OFFSET] [LIMIT]",
	Short: "Query log entries (newest first)",
	Long: `Scan log files and output entries sorted newest-first.
  OFFSET: number of entries to skip (0 = most recent)
  LIMIT:  maximum number of entries to output

Examples:
  loadstar findlog 0 10                                   # latest 10 entries globally
  loadstar findlog 0 20 --address W://root/cli/cmd_create # latest 20 for one element
  loadstar findlog 0 5 --kind ISSUE                       # latest 5 ISSUE entries
  loadstar findlog 10 10                                  # entries 11-20 (paging)`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		offset, err1 := strconv.Atoi(args[0])
		limit, err2 := strconv.Atoi(args[1])
		if err1 != nil || err2 != nil || offset < 0 || limit <= 0 {
			fmt.Fprintln(os.Stderr, "error: OFFSET must be >= 0 and LIMIT must be >= 1")
			os.Exit(1)
		}

		addrFilter, _ := cmd.Flags().GetString("address")
		kindFlag, _ := cmd.Flags().GetString("kind")
		kindFilter := strings.ToUpper(kindFlag)

		loadstarBase := fs.AvcsPath("")
		logDir := filepath.Join(loadstarBase, ".clionly", "LOG")

		entries, err := collectLogEntries(logDir, addrFilter, kindFilter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Timestamp > entries[j].Timestamp
		})

		if offset >= len(entries) {
			fmt.Println("no log entries found")
			return
		}
		end := offset + limit
		if end > len(entries) {
			end = len(entries)
		}
		entries = entries[offset:end]

		if len(entries) == 0 {
			fmt.Println("no log entries found")
			return
		}

		for _, e := range entries {
			fmt.Printf("[%s] [%s] %s\n  → %s\n", e.Timestamp, e.Kind, e.Content, e.Address)
		}
	},
}

func init() {
	logCmd.Flags().Bool("list", false, "Show all element IDs and addresses")
	findlogCmd.Flags().String("address", "", "Filter by element address (e.g. W://root/cli/cmd_create)")
	findlogCmd.Flags().String("kind", "", "Filter by KIND (NOTE, DECISION, ISSUE, RESOLVED, PROGRESS, MODIFIED)")
}

// elementInfo holds address + summary for --list display.
type elementInfo struct {
	ID      string // last segment of path (e.g. "cmd_log")
	Address string // full address (e.g. "W://root/cli/cmd_log")
	Summary string
}

const listPageSize = 100

// showElementIndex prints a table of all MAP/WAYPOINT elements with IDs.
func showElementIndex(filter string, page int) {
	loadstarBase := fs.AvcsPath("")
	var elements []elementInfo

	for _, typeDir := range []string{"MAP", "WAYPOINT"} {
		dir := filepath.Join(loadstarBase, typeDir)
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		prefix := "M"
		if typeDir == "WAYPOINT" {
			prefix = "W"
		}
		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
				continue
			}
			dotName := strings.TrimSuffix(f.Name(), ".md")
			path := strings.ReplaceAll(dotName, ".", "/")
			addr := prefix + "://" + path

			parts := strings.Split(path, "/")
			id := parts[len(parts)-1]

			summary := extractSummaryFromFile(filepath.Join(dir, f.Name()))

			elements = append(elements, elementInfo{ID: id, Address: addr, Summary: summary})
		}
	}

	if filter != "" {
		lower := strings.ToLower(filter)
		var filtered []elementInfo
		for _, e := range elements {
			if strings.Contains(strings.ToLower(e.ID), lower) ||
				strings.Contains(strings.ToLower(e.Address), lower) ||
				strings.Contains(strings.ToLower(e.Summary), lower) {
				filtered = append(filtered, e)
			}
		}
		elements = filtered
	}

	if len(elements) == 0 {
		fmt.Println("no elements found")
		return
	}

	sort.Slice(elements, func(i, j int) bool {
		return elements[i].Address < elements[j].Address
	})

	total := len(elements)
	totalPages := (total + listPageSize - 1) / listPageSize
	if page > totalPages {
		page = totalPages
	}
	start := (page - 1) * listPageSize
	end := start + listPageSize
	if end > total {
		end = total
	}
	pageItems := elements[start:end]

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tADDRESS\tSUMMARY")
	fmt.Fprintln(w, "--\t-------\t-------")
	for _, e := range pageItems {
		fmt.Fprintf(w, "%s\t%s\t%s\n", e.ID, e.Address, e.Summary)
	}
	w.Flush()

	if total > listPageSize {
		fmt.Printf("\n[page %d/%d] showing %d-%d of %d elements\n", page, totalPages, start+1, end, total)
	}
}

// extractSummaryFromFile reads a LOADSTAR element file and returns the SUMMARY value.
func extractSummaryFromFile(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- SUMMARY:") {
			val := strings.TrimSpace(strings.TrimPrefix(trimmed, "- SUMMARY:"))
			if val != "" && len(val) > 50 {
				val = val[:50] + "..."
			}
			return val
		}
	}
	return ""
}

// resolveAddress resolves an address input that may be a full address or a short ID.
func resolveAddress(input string) string {
	if strings.Contains(input, "://") {
		return input
	}

	loadstarBase := fs.AvcsPath("")
	var matches []string

	for _, typeDir := range []string{"MAP", "WAYPOINT"} {
		dir := filepath.Join(loadstarBase, typeDir)
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		prefix := "M"
		if typeDir == "WAYPOINT" {
			prefix = "W"
		}
		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
				continue
			}
			dotName := strings.TrimSuffix(f.Name(), ".md")
			path := strings.ReplaceAll(dotName, ".", "/")

			parts := strings.Split(path, "/")
			id := parts[len(parts)-1]
			if id == input || strings.HasSuffix(path, input) {
				matches = append(matches, prefix+"://"+path)
			}
		}
	}

	if len(matches) == 0 {
		fmt.Fprintf(os.Stderr, "error: no element found matching %q\n", input)
		fmt.Fprintln(os.Stderr, "  use 'loadstar log --list' to see available elements")
		return ""
	}
	if len(matches) == 1 {
		return matches[0]
	}

	fmt.Fprintf(os.Stderr, "error: %q matches multiple elements:\n", input)
	for _, m := range matches {
		fmt.Fprintf(os.Stderr, "  %s\n", m)
	}
	fmt.Fprintln(os.Stderr, "  specify a longer path to disambiguate")
	return ""
}

// writeLogChangeLog writes a LOG record to .clionly/LOG/.
func writeLogChangeLog(loadstarBase, targetAddr, kind, content string) {
	clDir := filepath.Join(loadstarBase, ".clionly", "LOG")
	_ = os.MkdirAll(clDir, 0755)

	ts := time.Now()
	fileName := fmt.Sprintf("CL.%s.log.md", ts.Format("20060102.150405"))
	clPath := filepath.Join(clDir, fileName)

	record := fmt.Sprintf("## LOG\n- TARGET: %s\n- KIND: %s\n- AT: %s\n- CONTENT: %s\n",
		targetAddr, kind, ts.Format("2006-01-02T15:04:05"), content)
	_ = os.WriteFile(clPath, []byte(record), 0644)
}

// logRecordRe parses structured fields from .clionly/LOG/ files.
var logTargetRe = regexp.MustCompile(`^-\s+TARGET:\s+(.+)$`)
var logKindRe = regexp.MustCompile(`^-\s+KIND:\s+(.+)$`)
var logAtRe = regexp.MustCompile(`^-\s+AT:\s+(.+)$`)
var logContentRe = regexp.MustCompile(`^-\s+CONTENT:\s+(.+)$`)

// collectLogEntries scans .clionly/LOG/ directory for log records and applies optional filters.
func collectLogEntries(logDir, addrFilter, kindFilter string) ([]logEntry, error) {
	files, err := os.ReadDir(logDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var results []logEntry
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".log.md") {
			continue
		}

		fullPath := filepath.Join(logDir, f.Name())
		data, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		entry := parseLogFile(string(data))
		if entry == nil {
			continue
		}

		// Normalize B:// addresses to W:// for backward compatibility
		if strings.HasPrefix(entry.Address, "B://") {
			entry.Address = "W://" + strings.TrimPrefix(entry.Address, "B://")
		}

		// Apply address filter
		if addrFilter != "" {
			filterPath := addrFilter
			if parts := strings.SplitN(addrFilter, "://", 2); len(parts) == 2 {
				filterPath = parts[1]
			}
			entryPath := entry.Address
			if parts := strings.SplitN(entry.Address, "://", 2); len(parts) == 2 {
				entryPath = parts[1]
			}
			if filterPath != entryPath {
				continue
			}
		}

		// Apply kind filter
		if kindFilter != "" && entry.Kind != kindFilter {
			continue
		}

		results = append(results, *entry)
	}
	return results, nil
}

// parseLogFile parses a single .clionly/LOG/ file into a logEntry.
func parseLogFile(content string) *logEntry {
	var entry logEntry
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if m := logTargetRe.FindStringSubmatch(trimmed); m != nil {
			entry.Address = strings.TrimSpace(m[1])
		}
		if m := logKindRe.FindStringSubmatch(trimmed); m != nil {
			entry.Kind = strings.TrimSpace(m[1])
		}
		if m := logAtRe.FindStringSubmatch(trimmed); m != nil {
			entry.Timestamp = strings.TrimSpace(m[1])
		}
		if m := logContentRe.FindStringSubmatch(trimmed); m != nil {
			entry.Content = strings.TrimSpace(m[1])
		}
	}
	if entry.Timestamp == "" || entry.Kind == "" {
		return nil
	}
	return &entry
}
