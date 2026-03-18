package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const todoListFile = ".clionly/TODO/TODO_LIST.md"
const todoHistoryFile = ".clionly/TODO/TODO_HISTORY.md"

var todoHeader = "| 실행 요소 (Executor) | 요청 요소 (Requester) | 발생 시간 (Time) | 작업 요약 (Summary) | 상태 (Status) | 선행 조건 (Depends_On) |"
var todoSep = "| :--- | :--- | :--- | :--- | :--- | :--- |"

var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Manage TODO items in TODO_LIST",
}

var todoAddCmd = &cobra.Command{
	Use:   "add [EXECUTOR] [REQUESTER] [SUMMARY]",
	Short: "Add a new TODO item",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		depends, _ := cmd.Flags().GetString("depends")
		if depends == "" {
			depends = "-"
		}
		executor, requester, summary := args[0], args[1], args[2]
		now := time.Now().Format("2006-01-02 15:04")

		todoPath := filepath.Join(fs.AvcsPath(""), todoListFile)
		ensureTodoFile(todoPath)

		content, err := os.ReadFile(todoPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not read TODO list: %v\n", err)
			os.Exit(1)
		}

		newRow := fmt.Sprintf("| %s | %s | %s | %s | PENDING | %s |",
			executor, requester, now, summary, depends)

		lines := strings.Split(string(content), "\n")
		insertIdx := -1
		for i, line := range lines {
			if strings.HasPrefix(line, "| :---") {
				insertIdx = i + 1
				break
			}
		}
		if insertIdx < 0 || insertIdx > len(lines) {
			lines = append(lines, newRow)
		} else {
			lines = append(lines[:insertIdx], append([]string{newRow}, lines[insertIdx:]...)...)
		}

		if err := os.WriteFile(todoPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error: could not write TODO list: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("todo added: %s\n", executor)
	},
}

var todoDoneCmd = &cobra.Command{
	Use:   "done [EXECUTOR]",
	Short: "Mark a TODO item as completed and move to TODO_HISTORY",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		executor := args[0]
		todoPath := filepath.Join(fs.AvcsPath(""), todoListFile)

		content, err := os.ReadFile(todoPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not read TODO list: %v\n", err)
			os.Exit(1)
		}

		lines := strings.Split(string(content), "\n")
		found := false
		var doneRow string
		var summary string
		var kept []string
		for _, line := range lines {
			if isDataRow(line) && extractCol(line, 0) == executor {
				found = true
				summary = extractCol(line, 3)
				doneRow = line
				continue
			}
			kept = append(kept, line)
		}

		if !found {
			fmt.Fprintf(os.Stderr, "error: TODO item not found for executor: %s\n", executor)
			os.Exit(1)
		}

		// Remove from TODO_LIST
		if err := os.WriteFile(todoPath, []byte(strings.Join(kept, "\n")), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error: could not write TODO list: %v\n", err)
			os.Exit(1)
		}

		// Append to TODO_HISTORY
		histPath := filepath.Join(fs.AvcsPath(""), todoHistoryFile)
		appendTodoHistory(histPath, doneRow, "DONE")

		// Append to executor element's EXECUTION_HISTORY
		addr, err := svc.ParseAddress(executor)
		if err == nil {
			loadstarBase := fs.AvcsPath("")
			elemFile := addr.ToFilePath(loadstarBase)
			if fs.Exists(elemFile) {
				appendExecutionHistory(elemFile, summary)
			}
		}

		fmt.Printf("todo done: %s\n", executor)
	},
}

// allowedUpdateStatuses lists the status values that todo update accepts.
// COMPLETED and FAILED are intentionally excluded — use `todo done` instead.
var allowedUpdateStatuses = map[string]bool{
	"PENDING": true, "ACTIVE": true, "BLOCKED": true,
}

var todoUpdateCmd = &cobra.Command{
	Use:   "update [EXECUTOR] [STATUS]",
	Short: "Update the status of a TODO item (PENDING, ACTIVE, BLOCKED)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		executor := args[0]
		newStatus := strings.ToUpper(args[1])

		if !allowedUpdateStatuses[newStatus] {
			fmt.Fprintf(os.Stderr, "error: invalid status %q — allowed: PENDING, ACTIVE, BLOCKED\n", newStatus)
			os.Exit(1)
		}

		todoPath := filepath.Join(fs.AvcsPath(""), todoListFile)
		content, err := os.ReadFile(todoPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not read TODO list: %v\n", err)
			os.Exit(1)
		}

		lines := strings.Split(string(content), "\n")
		found := false
		var originalRow string
		var oldStatus string
		for i, line := range lines {
			if !isDataRow(line) || extractCol(line, 0) != executor {
				continue
			}
			found = true
			originalRow = line
			oldStatus = extractCol(line, 4)
			parts := strings.Split(line, "|")
			if len(parts) >= 7 {
				parts[5] = " " + newStatus + " "
				lines[i] = strings.Join(parts, "|")
			}
			fmt.Printf("updated: %s  %s → %s\n", executor, oldStatus, newStatus)
			break
		}

		if !found {
			fmt.Fprintf(os.Stderr, "error: TODO item not found for executor: %s\n", executor)
			os.Exit(1)
		}

		if err := os.WriteFile(todoPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error: could not write TODO list: %v\n", err)
			os.Exit(1)
		}

		// Append to TODO_HISTORY
		histPath := filepath.Join(fs.AvcsPath(""), todoHistoryFile)
		action := fmt.Sprintf("UPDATED(%s→%s)", oldStatus, newStatus)
		appendTodoHistory(histPath, originalRow, action)
	},
}

var todoDeleteCmd = &cobra.Command{
	Use:   "delete [EXECUTOR]",
	Short: "Delete a TODO item and record it in TODO_HISTORY as DELETED",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		executor := args[0]
		todoPath := filepath.Join(fs.AvcsPath(""), todoListFile)

		content, err := os.ReadFile(todoPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not read TODO list: %v\n", err)
			os.Exit(1)
		}

		lines := strings.Split(string(content), "\n")
		found := false
		var deletedRow string
		var kept []string
		for _, line := range lines {
			if isDataRow(line) && extractCol(line, 0) == executor {
				found = true
				deletedRow = line
				continue
			}
			kept = append(kept, line)
		}

		if !found {
			fmt.Fprintf(os.Stderr, "error: TODO item not found for executor: %s\n", executor)
			os.Exit(1)
		}

		if err := os.WriteFile(todoPath, []byte(strings.Join(kept, "\n")), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error: could not write TODO list: %v\n", err)
			os.Exit(1)
		}

		// Append to TODO_HISTORY
		histPath := filepath.Join(fs.AvcsPath(""), todoHistoryFile)
		appendTodoHistory(histPath, deletedRow, "DELETED")

		fmt.Printf("todo deleted: %s\n", executor)
	},
}

var todoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all PENDING/ACTIVE TODO items",
	Run: func(cmd *cobra.Command, args []string) {
		todoPath := filepath.Join(fs.AvcsPath(""), todoListFile)
		content, err := os.ReadFile(todoPath)
		if err != nil {
			fmt.Println("no TODO list found")
			return
		}

		lines := strings.Split(string(content), "\n")
		completedSet := make(map[string]bool)

		fmt.Println(todoHeader)
		fmt.Println(todoSep)
		for _, line := range lines {
			if !isDataRow(line) {
				continue
			}
			executor := extractCol(line, 0)
			depends := extractCol(line, 5)
			status := extractCol(line, 4)
			if status != "PENDING" && status != "ACTIVE" {
				continue
			}
			displayLine := line
			if depends != "-" && !completedSet[depends] {
				displayLine = strings.Replace(line, "| PENDING |", "| [BLOCKED] |", 1)
			}
			fmt.Println(displayLine)
			_ = executor
		}
	},
}

var todoHistoryCmd = &cobra.Command{
	Use:   "history [EXECUTOR]",
	Short: "Show TODO_HISTORY. Optionally filter by executor address.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		histPath := filepath.Join(fs.AvcsPath(""), todoHistoryFile)
		content, err := os.ReadFile(histPath)
		if err != nil {
			fmt.Println("no TODO history found")
			return
		}

		filter := ""
		if len(args) == 1 {
			filter = args[0]
		}

		histHeader := "| 실행 요소 (Executor) | 요청 요소 (Requester) | 발생 시간 (Time) | 작업 요약 (Summary) | 액션 (Action) | 처리 시각 (At) | 선행 조건 (Depends_On) |"
		histSep := "| :--- | :--- | :--- | :--- | :--- | :--- | :--- |"

		lines := strings.Split(string(content), "\n")
		printed := false
		for _, line := range lines {
			if !isDataRow(line) {
				continue
			}
			if filter != "" && extractCol(line, 0) != filter {
				continue
			}
			if !printed {
				fmt.Println(histHeader)
				fmt.Println(histSep)
				printed = true
			}
			fmt.Println(line)
		}

		if !printed {
			if filter != "" {
				fmt.Printf("no history found for executor: %s\n", filter)
			} else {
				fmt.Println("no TODO history found")
			}
		}
	},
}

func init() {
	todoCmd.AddCommand(todoAddCmd)
	todoCmd.AddCommand(todoDoneCmd)
	todoCmd.AddCommand(todoDeleteCmd)
	todoCmd.AddCommand(todoListCmd)
	todoCmd.AddCommand(todoUpdateCmd)
	todoCmd.AddCommand(todoHistoryCmd)
	todoAddCmd.Flags().String("depends", "", "Prerequisite executor address")
}

// ensureTodoFile creates TODO_LIST.md if it doesn't exist.
func ensureTodoFile(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.MkdirAll(filepath.Dir(path), 0755)
		initial := todoHeader + "\n" + todoSep + "\n"
		_ = os.WriteFile(path, []byte(initial), 0644)
	}
}

// appendTodoHistory appends an action record to TODO_HISTORY.md.
// action: "DONE", "UPDATED(OLD→NEW)", "DELETED"
// row: the original TODO_LIST data row at the time of the action.
func appendTodoHistory(histPath, row, action string) {
	_ = os.MkdirAll(filepath.Dir(histPath), 0755)

	now := time.Now().Format("2006-01-02 15:04")
	histHeader := "| 실행 요소 (Executor) | 요청 요소 (Requester) | 발생 시간 (Time) | 작업 요약 (Summary) | 액션 (Action) | 처리 시각 (At) | 선행 조건 (Depends_On) |"
	histSep := "| :--- | :--- | :--- | :--- | :--- | :--- | :--- |"

	// Build history row from original row:
	// original cols: Executor(0) Requester(1) Time(2) Summary(3) Status(4) Depends_On(5)
	// history cols:  Executor    Requester    Time    Summary    Action     At           Depends_On
	parts := strings.Split(row, "|")
	var histRow string
	if len(parts) >= 7 {
		executor := strings.TrimSpace(parts[1])
		requester := strings.TrimSpace(parts[2])
		addedAt := strings.TrimSpace(parts[3])
		summary := strings.TrimSpace(parts[4])
		dependsOn := strings.TrimSpace(parts[6])
		histRow = fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s |",
			executor, requester, addedAt, summary, action, now, dependsOn)
	} else {
		histRow = row // fallback
	}

	var content []byte
	if _, err := os.Stat(histPath); os.IsNotExist(err) {
		content = []byte(histHeader + "\n" + histSep + "\n" + histRow + "\n")
	} else {
		existing, err := os.ReadFile(histPath)
		if err != nil {
			return
		}
		content = []byte(string(existing) + histRow + "\n")
	}

	_ = os.WriteFile(histPath, content, 0644)
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
	if strings.Contains(line, "실행 요소") {
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

// reevaluateDependents marks dependents of completedExecutor as no longer blocked (no-op for now, handled in list display).
func reevaluateDependents(lines []string, completedExecutor string) {
	_ = lines
	_ = completedExecutor
}

// appendExecutionHistory appends a completion record to the element's EXECUTION_HISTORY.
func appendExecutionHistory(filePath, summary string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}
	now := time.Now().Format("2006-01-02 15:04")
	re := regexp.MustCompile(`(- EXECUTION_HISTORY:\s*\[)(\s*)(\])`)
	entry := fmt.Sprintf("\n    * %s: [COMPLETED] %s\n  ", now, summary)
	updated := re.ReplaceAllString(string(content), "${1}"+entry+"${3}")
	_ = os.WriteFile(filePath, []byte(updated), 0644)
}
