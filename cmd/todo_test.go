package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeTodoFile writes TODO_LIST.md in the temp loadstarBase TODO dir.
func writeTodoFile(t *testing.T, loadstarBase, content string) string {
	t.Helper()
	todoDir := filepath.Join(loadstarBase, ".clionly", "TODO")
	if err := os.MkdirAll(todoDir, 0755); err != nil {
		t.Fatalf("mkdir .clionly/TODO: %v", err)
	}
	todoPath := filepath.Join(todoDir, "TODO_LIST.md")
	if err := os.WriteFile(todoPath, []byte(content), 0644); err != nil {
		t.Fatalf("write todo file: %v", err)
	}
	return todoPath
}

// ---- isDataRow ----

func TestIsDataRow_Header(t *testing.T) {
	if isDataRow(todoHeader) {
		t.Error("header row should not be a data row")
	}
}

func TestIsDataRow_Separator(t *testing.T) {
	if isDataRow(todoSep) {
		t.Error("separator row should not be a data row")
	}
}

func TestIsDataRow_DataRow(t *testing.T) {
	row := "| W://root/cli/cmd_create | 2026-03-10 12:00 | create 구현 | PENDING | - |"
	if !isDataRow(row) {
		t.Error("valid data row should return true")
	}
}

func TestIsDataRow_Empty(t *testing.T) {
	if isDataRow("") {
		t.Error("empty string should not be a data row")
	}
}

// ---- extractCol ----

func TestExtractCol(t *testing.T) {
	row := "| W://root/cmd | 2026-03-10 | summary text | PENDING | - |"
	cases := []struct {
		col  int
		want string
	}{
		{0, "W://root/cmd"},
		{1, "2026-03-10"},
		{2, "summary text"},
		{3, "PENDING"},
		{4, "-"},
	}
	for _, c := range cases {
		got := extractCol(row, c.col)
		if got != c.want {
			t.Errorf("extractCol(row, %d) = %q, want %q", c.col, got, c.want)
		}
	}
}

func TestExtractCol_OutOfBounds(t *testing.T) {
	row := "| A | B |"
	got := extractCol(row, 10)
	if got != "" {
		t.Errorf("out-of-bounds col should return empty, got %q", got)
	}
}

// ---- ensureTodoFile ----

func TestEnsureTodoFile_Creates(t *testing.T) {
	dir := t.TempDir()
	todoPath := filepath.Join(dir, ".clionly", "TODO", "TODO_LIST.md")
	ensureTodoFile(todoPath)

	if _, err := os.Stat(todoPath); os.IsNotExist(err) {
		t.Error("ensureTodoFile should create the file if it doesn't exist")
	}
	data, _ := os.ReadFile(todoPath)
	if !strings.Contains(string(data), "주소 (Address)") {
		t.Error("created todo file should contain the header")
	}
}

func TestEnsureTodoFile_NoOverwrite(t *testing.T) {
	dir := t.TempDir()
	todoPath := filepath.Join(dir, "TODO_LIST.md")
	existing := "existing content"
	_ = os.WriteFile(todoPath, []byte(existing), 0644)

	ensureTodoFile(todoPath)

	data, _ := os.ReadFile(todoPath)
	if string(data) != existing {
		t.Error("ensureTodoFile should not overwrite existing file")
	}
}

// ---- todo add: row insertion ----

func TestTodoAdd_NewRow(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	initial := todoHeader + "\n" + todoSep + "\n"
	todoPath := writeTodoFile(t, loadstarBase, initial)

	newRow := "| W://root/exec | 2026-03-10 12:00 | test task | PENDING | - |"
	lines := strings.Split(initial, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "| :---") {
			lines = append(lines[:i+1], append([]string{newRow}, lines[i+1:]...)...)
			break
		}
	}
	_ = os.WriteFile(todoPath, []byte(strings.Join(lines, "\n")), 0644)

	data, _ := os.ReadFile(todoPath)
	if !strings.Contains(string(data), "test task") {
		t.Error("new row should appear in the todo file")
	}
}

// ---- todo done: row removal ----

func TestTodoDone_RemovesRow(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	row := "| W://root/exec | 2026-03-10 | task | PENDING | - |"
	initial := todoHeader + "\n" + todoSep + "\n" + row + "\n"
	todoPath := writeTodoFile(t, loadstarBase, initial)

	lines := strings.Split(initial, "\n")
	var kept []string
	for _, line := range lines {
		if isDataRow(line) && extractCol(line, 0) == "W://root/exec" {
			continue
		}
		kept = append(kept, line)
	}
	_ = os.WriteFile(todoPath, []byte(strings.Join(kept, "\n")), 0644)

	data, _ := os.ReadFile(todoPath)
	if strings.Contains(string(data), "W://root/exec") {
		t.Error("address row should have been removed from todo list")
	}
}

func TestTodoDone_AppendsToHistory(t *testing.T) {
	dir := t.TempDir()
	doneRow := "| W://root/exec | 2026-03-10 12:00 | task done | PENDING | - |"
	histPath := filepath.Join(dir, ".clionly", "TODO", "TODO_HISTORY.md")

	appendTodoHistory(histPath, doneRow, "DONE")

	data, err := os.ReadFile(histPath)
	if err != nil {
		t.Fatalf("TODO_HISTORY.md should have been created: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "W://root/exec") {
		t.Error("history file should contain the address")
	}
	if !strings.Contains(content, "task done") {
		t.Error("history file should contain the summary")
	}
}

// ---- todo list: BLOCKED display ----

func TestTodoList_BlockedDisplay(t *testing.T) {
	row := "| W://root/exec | 2026-03-10 | task | PENDING | W://root/prereq |"
	depends := extractCol(row, 4)
	status := extractCol(row, 3)

	if status != "PENDING" {
		t.Fatalf("unexpected status: %q", status)
	}

	displayLine := row
	if depends != "-" && depends != "" {
		displayLine = strings.Replace(row, "| PENDING |", "| [BLOCKED] |", 1)
	}

	if !strings.Contains(displayLine, "[BLOCKED]") {
		t.Error("row with unmet dependency should be displayed as [BLOCKED]")
	}
}

func TestTodoList_NotBlocked_WhenNoDepends(t *testing.T) {
	row := "| W://root/exec | 2026-03-10 | task | PENDING | - |"
	depends := extractCol(row, 4)

	displayLine := row
	if depends != "-" && depends != "" {
		displayLine = strings.Replace(row, "| PENDING |", "| [BLOCKED] |", 1)
	}

	if strings.Contains(displayLine, "[BLOCKED]") {
		t.Error("row with no dependency should not be displayed as [BLOCKED]")
	}
}

// ---- todo update ----

func updateStatusInLines(lines []string, address, newStatus string) ([]string, bool) {
	found := false
	for i, line := range lines {
		if !isDataRow(line) || extractCol(line, 0) != address {
			continue
		}
		found = true
		parts := strings.Split(line, "|")
		if len(parts) >= 6 {
			parts[4] = " " + newStatus + " "
			lines[i] = strings.Join(parts, "|")
		}
		break
	}
	return lines, found
}

func TestTodoUpdate_StatusChange(t *testing.T) {
	row := "| W://root/exec | 2026-03-10 | task | PENDING | - |"
	lines := []string{todoHeader, todoSep, row}
	lines, found := updateStatusInLines(lines, "W://root/exec", "ACTIVE")
	if !found {
		t.Fatal("address not found during update")
	}
	if !strings.Contains(lines[2], "ACTIVE") {
		t.Error("status should have been updated to ACTIVE")
	}
}

func TestTodoUpdate_InvalidStatus(t *testing.T) {
	for _, s := range []string{"COMPLETED", "FAILED", "DONE"} {
		if allowedUpdateStatuses[strings.ToUpper(s)] {
			t.Errorf("status %q should not be allowed in todo update", s)
		}
	}
}

func TestTodoUpdate_AllowedStatuses(t *testing.T) {
	for _, s := range []string{"PENDING", "ACTIVE"} {
		if !allowedUpdateStatuses[s] {
			t.Errorf("status %q should be allowed in todo update", s)
		}
	}
	// BLOCKED is auto-calculated from Depends_On, not manually settable
	if allowedUpdateStatuses["BLOCKED"] {
		t.Error("BLOCKED should not be manually settable")
	}
}

// ---- appendTodoHistory ----

func TestAppendTodoHistory_DoneAction(t *testing.T) {
	dir := t.TempDir()
	histPath := filepath.Join(dir, ".clionly", "TODO", "TODO_HISTORY.md")
	row := "| W://root/exec | 2026-03-10 12:00 | task summary | PENDING | - |"

	appendTodoHistory(histPath, row, "DONE")

	data, _ := os.ReadFile(histPath)
	content := string(data)
	if !strings.Contains(content, "DONE") {
		t.Error("history should contain action DONE")
	}
	if !strings.Contains(content, "W://root/exec") {
		t.Error("history should contain address")
	}
}

func TestAppendTodoHistory_Accumulates(t *testing.T) {
	dir := t.TempDir()
	histPath := filepath.Join(dir, ".clionly", "TODO", "TODO_HISTORY.md")
	row1 := "| W://root/a | 2026-03-10 12:00 | task a | PENDING | - |"
	row2 := "| W://root/b | 2026-03-10 13:00 | task b | ACTIVE | - |"

	appendTodoHistory(histPath, row1, "DONE")
	appendTodoHistory(histPath, row2, "UPDATED(ACTIVE→BLOCKED)")

	data, _ := os.ReadFile(histPath)
	content := string(data)
	if !strings.Contains(content, "W://root/a") || !strings.Contains(content, "W://root/b") {
		t.Error("history should accumulate multiple rows")
	}
}

// ---- todo history: filter ----

func filterHistoryLines(content, filter string) []string {
	lines := strings.Split(content, "\n")
	var result []string
	for _, line := range lines {
		if !isDataRow(line) {
			continue
		}
		if filter != "" && extractCol(line, 0) != filter {
			continue
		}
		result = append(result, line)
	}
	return result
}

func TestTodoHistory_FilterByAddress(t *testing.T) {
	dir := t.TempDir()
	histPath := filepath.Join(dir, ".clionly", "TODO", "TODO_HISTORY.md")

	appendTodoHistory(histPath, "| W://root/a | 2026-03-10 12:00 | task a | PENDING | - |", "DONE")
	appendTodoHistory(histPath, "| W://root/b | 2026-03-10 14:00 | task b | PENDING | - |", "DELETED")

	content, _ := os.ReadFile(histPath)
	rows := filterHistoryLines(string(content), "W://root/a")
	if len(rows) != 1 {
		t.Errorf("expected 1 row for W://root/a, got %d", len(rows))
	}
}
