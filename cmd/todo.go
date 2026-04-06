package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const todoListFile = ".clionly/TODO/TODO_LIST.md"
const todoHistoryFile = ".clionly/TODO/TODO_HISTORY.md"

var todoHeader = "| 주소 (Address) | 발생 시간 (Time) | 작업 요약 (Summary) | 상태 (Status) | 선행 조건 (Depends_On) |"
var todoSep = "| :--- | :--- | :--- | :--- | :--- |"

var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Manage TODO items in TODO_LIST",
	Long: `Manage the project TODO list stored in .loadstar/.clionly/TODO/TODO_LIST.md.
Do NOT edit that file directly — always use these commands.

Subcommands:
  add      Add a new TODO item (status: PENDING)
  list     Show current PENDING/ACTIVE items (BLOCKED auto-detected)
  update   Change status to PENDING / ACTIVE / BLOCKED
  done     Mark as completed and move to TODO_HISTORY
  delete   Remove without completion record
  history  Show TODO_HISTORY (all completed/deleted events)

Quick start:
  loadstar todo add W://root/cli/cmd_log "log/findlog 구현"
  loadstar todo list
  loadstar todo update W://root/cli/cmd_log ACTIVE
  loadstar todo done W://root/cli/cmd_log`,
}

var todoAddCmd = &cobra.Command{
	Use:   "add [ADDRESS] [SUMMARY]",
	Short: "Add a new TODO item",
	Long: `Add a new TODO item to the project TODO list (status: PENDING).

  ADDRESS  W:// address of the WayPoint for this task
  SUMMARY  One-line description of the task

Examples:
  loadstar todo add W://root/cli/cmd_log "log/findlog 구현"
  loadstar todo add W://root/cli/cmd_sync "sync 구현" --depends W://root/cli/cmd_log`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		depends, _ := cmd.Flags().GetString("depends")
		if depends == "" {
			depends = "-"
		}
		address, summary := args[0], args[1]
		now := time.Now().Format("2006-01-02 15:04")

		todoPath := filepath.Join(fs.AvcsPath(""), todoListFile)
		ensureTodoFile(todoPath)

		content, err := os.ReadFile(todoPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not read TODO list: %v\n", err)
			os.Exit(1)
		}

		newRow := fmt.Sprintf("| %s | %s | %s | PENDING | %s |",
			address, now, summary, depends)

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
		fmt.Printf("todo added: %s\n", address)
	},
}

var todoDoneCmd = &cobra.Command{
	Use:   "done [ADDRESS]",
	Short: "Mark a TODO item as completed and move to TODO_HISTORY",
	Long: `Mark a TODO item as COMPLETED.
The item is removed from TODO_LIST and recorded in TODO_HISTORY.

Examples:
  loadstar todo done W://root/cli/cmd_log`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		address := args[0]
		todoPath := filepath.Join(fs.AvcsPath(""), todoListFile)

		content, err := os.ReadFile(todoPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not read TODO list: %v\n", err)
			os.Exit(1)
		}

		lines := strings.Split(string(content), "\n")
		found := false
		var doneRow string
		var kept []string
		for _, line := range lines {
			if isDataRow(line) && extractCol(line, 0) == address {
				found = true
				doneRow = line
				continue
			}
			kept = append(kept, line)
		}

		if !found {
			fmt.Fprintf(os.Stderr, "error: TODO item not found for address: %s\n", address)
			os.Exit(1)
		}

		// Clear Depends_On references to the completed address
		for i, line := range kept {
			if isDataRow(line) && extractCol(line, 4) == address {
				parts := strings.Split(line, "|")
				if len(parts) >= 6 {
					parts[5] = " - "
					kept[i] = strings.Join(parts, "|")
				}
			}
		}

		// Remove from TODO_LIST
		if err := os.WriteFile(todoPath, []byte(strings.Join(kept, "\n")), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error: could not write TODO list: %v\n", err)
			os.Exit(1)
		}

		// Append to TODO_HISTORY
		histPath := filepath.Join(fs.AvcsPath(""), todoHistoryFile)
		appendTodoHistory(histPath, doneRow, "DONE")

		fmt.Printf("todo done: %s\n", address)
	},
}

// allowedUpdateStatuses lists the status values that todo update accepts.
var allowedUpdateStatuses = map[string]bool{
	"PENDING": true, "ACTIVE": true,
}

var todoUpdateCmd = &cobra.Command{
	Use:   "update [ADDRESS] [STATUS]",
	Short: "Update the status of a TODO item (PENDING, ACTIVE, BLOCKED)",
	Long: `Change the status of an existing TODO item.
Allowed values: PENDING, ACTIVE
BLOCKED is auto-calculated from Depends_On — not manually settable.
Use 'loadstar todo done' to mark as COMPLETED.

Examples:
  loadstar todo update W://root/cli/cmd_log ACTIVE
  loadstar todo update W://root/cli/cmd_log PENDING`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		address := args[0]
		newStatus := strings.ToUpper(args[1])

		if !allowedUpdateStatuses[newStatus] {
			fmt.Fprintf(os.Stderr, "error: invalid status %q — allowed: PENDING, ACTIVE\n", newStatus)
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
			if !isDataRow(line) || extractCol(line, 0) != address {
				continue
			}
			found = true
			originalRow = line
			oldStatus = extractCol(line, 3)
			parts := strings.Split(line, "|")
			if len(parts) >= 6 {
				parts[4] = " " + newStatus + " "
				lines[i] = strings.Join(parts, "|")
			}
			fmt.Printf("updated: %s  %s → %s\n", address, oldStatus, newStatus)
			break
		}

		if !found {
			fmt.Fprintf(os.Stderr, "error: TODO item not found for address: %s\n", address)
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
	Use:   "delete [ADDRESS]",
	Short: "Delete a TODO item and record it in TODO_HISTORY as DELETED",
	Long: `Cancel and remove a TODO item.
Use this for tasks that are cancelled or no longer needed.

Examples:
  loadstar todo delete W://root/cli/cmd_log`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		address := args[0]
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
			if isDataRow(line) && extractCol(line, 0) == address {
				found = true
				deletedRow = line
				continue
			}
			kept = append(kept, line)
		}

		if !found {
			fmt.Fprintf(os.Stderr, "error: TODO item not found for address: %s\n", address)
			os.Exit(1)
		}

		if err := os.WriteFile(todoPath, []byte(strings.Join(kept, "\n")), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error: could not write TODO list: %v\n", err)
			os.Exit(1)
		}

		// Append to TODO_HISTORY
		histPath := filepath.Join(fs.AvcsPath(""), todoHistoryFile)
		appendTodoHistory(histPath, deletedRow, "DELETED")

		fmt.Printf("todo deleted: %s\n", address)
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
		fmt.Println(todoHeader)
		fmt.Println(todoSep)
		for _, line := range lines {
			if !isDataRow(line) {
				continue
			}
			depends := extractCol(line, 4)
			status := extractCol(line, 3)
			if status != "PENDING" && status != "ACTIVE" {
				continue
			}
			displayLine := line
			if depends != "-" && depends != "" {
				displayLine = strings.Replace(line, "| PENDING |", "| [BLOCKED] |", 1)
			}
			fmt.Println(displayLine)
		}
	},
}

var todoHistoryCmd = &cobra.Command{
	Use:   "history [ADDRESS]",
	Short: "Show TODO_HISTORY. Optionally filter by address.",
	Long: `Show all completed and deleted TODO events from TODO_HISTORY.
Optionally filter by a specific address.

Examples:
  loadstar todo history
  loadstar todo history W://root/cli/cmd_log`,
	Args: cobra.MaximumNArgs(1),
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

		histHeader := "| 주소 (Address) | 발생 시간 (Time) | 작업 요약 (Summary) | 액션 (Action) | 처리 시각 (At) | 선행 조건 (Depends_On) |"
		histSep := "| :--- | :--- | :--- | :--- | :--- | :--- |"

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
				fmt.Printf("no history found for address: %s\n", filter)
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
	todoAddCmd.Flags().String("depends", "", "Prerequisite address")
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
func appendTodoHistory(histPath, row, action string) {
	_ = os.MkdirAll(filepath.Dir(histPath), 0755)

	now := time.Now().Format("2006-01-02 15:04")
	histHeader := "| 주소 (Address) | 발생 시간 (Time) | 작업 요약 (Summary) | 액션 (Action) | 처리 시각 (At) | 선행 조건 (Depends_On) |"
	histSep := "| :--- | :--- | :--- | :--- | :--- | :--- |"

	// Build history row from original row:
	// original cols: Address(0) Time(1) Summary(2) Status(3) Depends_On(4)
	// history cols:  Address    Time    Summary    Action     At           Depends_On
	parts := strings.Split(row, "|")
	var histRow string
	if len(parts) >= 6 {
		address := strings.TrimSpace(parts[1])
		addedAt := strings.TrimSpace(parts[2])
		summary := strings.TrimSpace(parts[3])
		dependsOn := strings.TrimSpace(parts[5])
		histRow = fmt.Sprintf("| %s | %s | %s | %s | %s | %s |",
			address, addedAt, summary, action, now, dependsOn)
	} else {
		histRow = row
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
	if strings.Contains(line, "주소 (Address)") || strings.Contains(line, "실행 요소") {
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
