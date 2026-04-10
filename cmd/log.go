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

type logEntry struct {
	Timestamp string
	Kind      string
	Content   string
	Address   string
}

const maxLogLines = 1000

// ===== log (parent command) =====

var logCmd = &cobra.Command{
	Use:   "log [TIME_RANGE] [FILTER]",
	Short: "Search log entries, or use 'log add' to write",
	Long: `Search log entries. All arguments are optional.

TIME_RANGE: Nd (days) or Nh (hours). If omitted, shows all entries.
FILTER: keyword to match against address, KIND, or content.

Examples:
  loadstar log                          # all logs
  loadstar log 7d                       # last 7 days
  loadstar log 3h                       # last 3 hours
  loadstar log cmd_show                 # filter by keyword, all time
  loadstar log 2d cmd_show              # keyword + time range
  loadstar log ISSUE                    # KIND filter

Subcommand:
  add    Append a log entry for an element

Use --list to browse available elements.`,
	Args: cobra.MaximumNArgs(2),
}

// ===== log add =====

var logAddCmd = &cobra.Command{
	Use:   "add [ADDRESS_OR_ID] [KIND] [CONTENT]",
	Short: "Append a log entry for an element",
	Long: `Append a structured log entry for the specified element.

ADDRESS can be a full address (W://root/cli/cmd_log) or a short ID (cmd_log).
KIND defaults to NOTE if omitted (2-arg form).

Examples:
  loadstar log add cmd_log "CODE_MAP 갱신 필요"
  loadstar log add cmd_log ISSUE "multiline 파싱 실패"
  loadstar log add W://root/cli/cmd_log NOTE "full address"`,
	Args: cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
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

		addrStr := resolveAddress(addrInput)
		if addrStr == "" {
			os.Exit(1)
		}
		if _, err := svc.ParseAddress(addrStr); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		loadstarBase := fs.AvcsPath("")
		writeLogEntry(loadstarBase, addrStr, kind, content)
		fmt.Printf("logged: [%s] %s\n", kind, addrStr)
	},
}

// ===== log find =====

var logFindCmd = &cobra.Command{
	Use:   "find [FILTER] [TIME_RANGE]",
	Short: "Search log entries with filter and time range",
	Long: `Search log entries. All arguments are optional.

FILTER: keyword to match against address, KIND, or content.
TIME_RANGE: Nd (days) or Nh (hours). Default: 1d.

Examples:
  loadstar log find                     # all logs, last 24h
  loadstar log find 7d                  # all logs, last 7 days
  loadstar log find 3h                  # all logs, last 3 hours
  loadstar log find cmd_show            # filter by keyword, last 24h
  loadstar log find cmd_show 3d         # keyword + time range
  loadstar log find ISSUE 2d            # KIND filter + time range

Max output: 1000 lines.`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filter := ""
		duration := 24 * time.Hour

		for _, a := range args {
			if d, ok := parseTimeRange(a); ok {
				duration = d
			} else {
				filter = a
			}
		}

		loadstarBase := fs.AvcsPath("")
		logDir := filepath.Join(loadstarBase, ".clionly", "LOG")
		cutoff := time.Now().Add(-duration)

		entries, err := collectLogEntries(logDir, filter, cutoff)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Timestamp > entries[j].Timestamp
		})

		if len(entries) > maxLogLines {
			entries = entries[:maxLogLines]
		}

		if len(entries) == 0 {
			fmt.Println("no log entries found")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for _, e := range entries {
			fmt.Fprintf(w, "[%s]\t[%s]\t%s\t→ %s\n", e.Timestamp, e.Kind, truncate(e.Content, 60), e.Address)
		}
		w.Flush()
		fmt.Printf("\n%d entries\n", len(entries))
	},
}

func init() {
	logCmd.AddCommand(logAddCmd)
	logCmd.Flags().Bool("list", false, "Show all element IDs and addresses")
	logCmd.Run = func(cmd *cobra.Command, args []string) {
		listFlag, _ := cmd.Flags().GetBool("list")
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

		// log [TIME_RANGE] [FILTER] — 직접 조회
		filter := ""
		var cutoff time.Time // zero value = no cutoff (all entries)
		hasDuration := false

		for _, a := range args {
			if d, ok := parseTimeRange(a); ok {
				cutoff = time.Now().Add(-d)
				hasDuration = true
			} else {
				filter = a
			}
		}
		_ = hasDuration

		loadstarBase := fs.AvcsPath("")
		logDir := filepath.Join(loadstarBase, ".clionly", "LOG")

		entries, err := collectLogEntries(logDir, filter, cutoff)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Timestamp > entries[j].Timestamp
		})

		if len(entries) > maxLogLines {
			entries = entries[:maxLogLines]
		}

		if len(entries) == 0 {
			fmt.Println("no log entries found")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for _, e := range entries {
			fmt.Fprintf(w, "[%s]\t[%s]\t%s\t→ %s\n", e.Timestamp, e.Kind, truncate(e.Content, 60), e.Address)
		}
		w.Flush()
		fmt.Printf("\n%d entries\n", len(entries))
	}
}

// ===== time range parsing =====

var timeRangeRe = regexp.MustCompile(`^(\d+)(d|h)$`)

func parseTimeRange(s string) (time.Duration, bool) {
	m := timeRangeRe.FindStringSubmatch(s)
	if m == nil {
		return 0, false
	}
	n, _ := strconv.Atoi(m[1])
	if m[2] == "d" {
		return time.Duration(n) * 24 * time.Hour, true
	}
	return time.Duration(n) * time.Hour, true
}

// ===== file I/O: daily log files =====

func writeLogEntry(loadstarBase, targetAddr, kind, content string) {
	logDir := filepath.Join(loadstarBase, ".clionly", "LOG")
	_ = os.MkdirAll(logDir, 0755)

	ts := time.Now()
	fileName := ts.Format("2006-01-02") + ".log"
	logPath := filepath.Join(logDir, fileName)

	line := fmt.Sprintf("%s|%s|%s|%s\n",
		ts.Format("2006-01-02T15:04:05"), kind, targetAddr, content)

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.WriteString(line)
}

func collectLogEntries(logDir, filter string, cutoff time.Time) ([]logEntry, error) {
	files, err := os.ReadDir(logDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	filterLower := strings.ToLower(filter)
	isKindFilter := validLogKinds[strings.ToUpper(filter)]

	var results []logEntry

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		name := f.Name()
		if strings.HasSuffix(name, ".log") || strings.HasSuffix(name, ".log.md") {
			fullPath := filepath.Join(logDir, name)
			data, err := os.ReadFile(fullPath)
			if err != nil {
				continue
			}

			if strings.HasSuffix(name, ".log.md") {
				entry := parseLegacyLogFile(string(data))
				if entry != nil && matchesFilter(entry, filterLower, isKindFilter) && !beforeCutoff(entry.Timestamp, cutoff) {
					results = append(results, *entry)
				}
			} else {
				for _, line := range strings.Split(string(data), "\n") {
					line = strings.TrimSpace(line)
					if line == "" {
						continue
					}
					entry := parseDailyLogLine(line)
					if entry != nil && matchesFilter(entry, filterLower, isKindFilter) && !beforeCutoff(entry.Timestamp, cutoff) {
						results = append(results, *entry)
					}
				}
			}
		}
	}

	return results, nil
}

func parseDailyLogLine(line string) *logEntry {
	parts := strings.SplitN(line, "|", 4)
	if len(parts) < 4 {
		return nil
	}
	return &logEntry{
		Timestamp: strings.TrimSpace(parts[0]),
		Kind:      strings.TrimSpace(parts[1]),
		Address:   strings.TrimSpace(parts[2]),
		Content:   strings.TrimSpace(parts[3]),
	}
}

func parseLegacyLogFile(content string) *logEntry {
	var entry logEntry
	targetRe := regexp.MustCompile(`^-\s+TARGET:\s+(.+)$`)
	kindRe := regexp.MustCompile(`^-\s+KIND:\s+(.+)$`)
	atRe := regexp.MustCompile(`^-\s+AT:\s+(.+)$`)
	contentRe := regexp.MustCompile(`^-\s+CONTENT:\s+(.+)$`)

	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if m := targetRe.FindStringSubmatch(trimmed); m != nil {
			entry.Address = strings.TrimSpace(m[1])
		}
		if m := kindRe.FindStringSubmatch(trimmed); m != nil {
			entry.Kind = strings.TrimSpace(m[1])
		}
		if m := atRe.FindStringSubmatch(trimmed); m != nil {
			entry.Timestamp = strings.TrimSpace(m[1])
		}
		if m := contentRe.FindStringSubmatch(trimmed); m != nil {
			entry.Content = strings.TrimSpace(m[1])
		}
	}

	if strings.HasPrefix(entry.Address, "B://") {
		entry.Address = "W://" + strings.TrimPrefix(entry.Address, "B://")
	}
	if entry.Timestamp == "" || entry.Kind == "" {
		return nil
	}
	return &entry
}

func matchesFilter(entry *logEntry, filterLower string, isKindFilter bool) bool {
	if filterLower == "" {
		return true
	}
	if isKindFilter {
		return strings.ToUpper(entry.Kind) == strings.ToUpper(filterLower)
	}
	return strings.Contains(strings.ToLower(entry.Address), filterLower) ||
		strings.Contains(strings.ToLower(entry.Content), filterLower)
}

func beforeCutoff(timestamp string, cutoff time.Time) bool {
	if cutoff.IsZero() {
		return false // no cutoff → keep all entries
	}
	t, err := time.Parse("2006-01-02T15:04:05", timestamp)
	if err != nil {
		return false
	}
	return t.Before(cutoff)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

// ===== element index (--list) =====

type elementInfo struct {
	ID      string
	Address string
	Summary string
}

const listPageSize = 100

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
