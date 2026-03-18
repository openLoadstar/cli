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
	if isDataRow("| 실행 요소 (Executor) | 요청 요소 (Requester) | 발생 시간 (Time) | 작업 요약 (Summary) | 상태 (Status) | 선행 조건 (Depends_On) |") {
		t.Error("header row should not be a data row")
	}
}

func TestIsDataRow_Separator(t *testing.T) {
	if isDataRow("| :--- | :--- | :--- | :--- | :--- | :--- |") {
		t.Error("separator row should not be a data row")
	}
}

func TestIsDataRow_DataRow(t *testing.T) {
	row := "| W://root/cli/cmd_create | M://root/cli | 2026-03-10 12:00 | create 구현 | PENDING | - |"
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
	row := "| W://root/cmd | M://root | 2026-03-10 | summary text | PENDING | - |"
	cases := []struct {
		col  int
		want string
	}{
		{0, "W://root/cmd"},
		{1, "M://root"},
		{2, "2026-03-10"},
		{3, "summary text"},
		{4, "PENDING"},
		{5, "-"},
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
	if !strings.Contains(string(data), "실행 요소") {
		t.Error("created todo file should contain the header")
	}
}

func TestEnsureTodoFile_NoOverwrite(t *testing.T) {
	dir := t.TempDir()
	todoPath := filepath.Join(dir, "GLOBAL_TODO_LIST.md")
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

	newRow := "| W://root/exec | M://root | 2026-03-10 12:00 | test task | PENDING | - |"
	lines := strings.Split(initial, "\n")
	insertIdx := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "| :---") {
			insertIdx = i + 1
			break
		}
	}
	if insertIdx >= 0 && insertIdx <= len(lines) {
		lines = append(lines[:insertIdx], append([]string{newRow}, lines[insertIdx:]...)...)
	}
	_ = os.WriteFile(todoPath, []byte(strings.Join(lines, "\n")), 0644)

	data, _ := os.ReadFile(todoPath)
	if !strings.Contains(string(data), "test task") {
		t.Error("new row should appear in the todo file")
	}
	if !strings.Contains(string(data), "W://root/exec") {
		t.Error("executor should appear in the todo file")
	}
}

func TestTodoAdd_WithDepends(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	initial := todoHeader + "\n" + todoSep + "\n"
	todoPath := writeTodoFile(t, loadstarBase, initial)

	depends := "W://root/prereq"
	newRow := "| W://root/exec | M://root | 2026-03-10 | task | PENDING | " + depends + " |"
	content, _ := os.ReadFile(todoPath)
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "| :---") {
			lines = append(lines[:i+1], append([]string{newRow}, lines[i+1:]...)...)
			break
		}
	}
	_ = os.WriteFile(todoPath, []byte(strings.Join(lines, "\n")), 0644)

	data, _ := os.ReadFile(todoPath)
	if !strings.Contains(string(data), depends) {
		t.Error("depends_on column should appear in the todo row")
	}
}

// ---- todo done: row removal ----

func TestTodoDone_RemovesRow(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	row := "| W://root/exec | M://root | 2026-03-10 | task | PENDING | - |"
	initial := todoHeader + "\n" + todoSep + "\n" + row + "\n"
	todoPath := writeTodoFile(t, loadstarBase, initial)

	// Simulate todoDoneCmd logic
	lines := strings.Split(string(initial), "\n")
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
		t.Error("executor row should have been removed from todo list")
	}
}

func TestTodoDone_AppendsToHistory(t *testing.T) {
	dir := t.TempDir()
	doneRow := "| W://root/exec | M://root | 2026-03-10 12:00 | task done | PENDING | - |"
	histPath := filepath.Join(dir, ".clionly", "TODO", "TODO_HISTORY.md")

	appendTodoHistory(histPath, doneRow, "DONE")

	data, err := os.ReadFile(histPath)
	if err != nil {
		t.Fatalf("TODO_HISTORY.md should have been created: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "W://root/exec") {
		t.Error("history file should contain the executor")
	}
	if !strings.Contains(content, "task done") {
		t.Error("history file should contain the summary")
	}
	if !strings.Contains(content, "액션 (Action)") {
		t.Error("history file should contain the Action header")
	}
}

func TestTodoDone_HistoryAccumulates(t *testing.T) {
	dir := t.TempDir()
	histPath := filepath.Join(dir, ".clionly", "TODO", "TODO_HISTORY.md")
	row1 := "| W://root/a | NONE | 2026-03-10 12:00 | task a | PENDING | - |"
	row2 := "| W://root/b | NONE | 2026-03-10 13:00 | task b | PENDING | - |"

	appendTodoHistory(histPath, row1, "DONE")
	appendTodoHistory(histPath, row2, "DONE")

	data, _ := os.ReadFile(histPath)
	content := string(data)
	if !strings.Contains(content, "W://root/a") || !strings.Contains(content, "W://root/b") {
		t.Error("history file should accumulate multiple completed rows")
	}
}

func TestTodoDone_NotFound(t *testing.T) {
	initial := todoHeader + "\n" + todoSep + "\n"
	lines := strings.Split(initial, "\n")
	found := false
	for _, line := range lines {
		if isDataRow(line) && extractCol(line, 0) == "W://root/nonexistent" {
			found = true
		}
	}
	if found {
		t.Error("nonexistent executor should not be found")
	}
}

// ---- todo done: EXECUTION_HISTORY ----

func TestTodoDone_ExecutionHistory(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	elemPath := writeElement(t, loadstarBase, "W://root/exec", buildTemplate("W", "W://root/exec", "M://root"))

	appendExecutionHistory(elemPath, "task completed")

	data, _ := os.ReadFile(elemPath)
	content := string(data)
	if !strings.Contains(content, "task completed") {
		t.Error("execution history should contain the task summary")
	}
	if !strings.Contains(content, "[COMPLETED]") {
		t.Error("execution history should contain [COMPLETED] marker")
	}
}

// ---- todo list: BLOCKED display ----

func TestTodoList_BlockedDisplay(t *testing.T) {
	row := "| W://root/exec | M://root | 2026-03-10 | task | PENDING | W://root/prereq |"
	// Depends on W://root/prereq which is NOT in completedSet
	completedSet := make(map[string]bool)
	depends := extractCol(row, 5)
	status := extractCol(row, 4)

	if status != "PENDING" {
		t.Fatalf("unexpected status: %q", status)
	}

	displayLine := row
	if depends != "-" && !completedSet[depends] {
		displayLine = strings.Replace(row, "| PENDING |", "| [BLOCKED] |", 1)
	}

	if !strings.Contains(displayLine, "[BLOCKED]") {
		t.Error("row with unmet dependency should be displayed as [BLOCKED]")
	}
}

func TestTodoList_NotBlocked_WhenNoDepends(t *testing.T) {
	row := "| W://root/exec | M://root | 2026-03-10 | task | PENDING | - |"
	completedSet := make(map[string]bool)
	depends := extractCol(row, 5)

	displayLine := row
	if depends != "-" && !completedSet[depends] {
		displayLine = strings.Replace(row, "| PENDING |", "| [BLOCKED] |", 1)
	}

	if strings.Contains(displayLine, "[BLOCKED]") {
		t.Error("row with no dependency should not be displayed as [BLOCKED]")
	}
}

func TestTodoList_Empty(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	initial := todoHeader + "\n" + todoSep + "\n"
	writeTodoFile(t, loadstarBase, initial)

	lines := strings.Split(initial, "\n")
	count := 0
	for _, line := range lines {
		if isDataRow(line) {
			count++
		}
	}
	if count != 0 {
		t.Errorf("expected 0 todo items, got %d", count)
	}
}

func TestTodoList_FiltersNonPending(t *testing.T) {
	rows := []string{
		"| W://root/a | M://root | 2026-03-10 | task a | PENDING | - |",
		"| W://root/b | M://root | 2026-03-10 | task b | DONE | - |",
		"| W://root/c | M://root | 2026-03-10 | task c | ACTIVE | - |",
	}

	var displayed []string
	for _, row := range rows {
		status := extractCol(row, 4)
		if status == "PENDING" || status == "ACTIVE" {
			displayed = append(displayed, row)
		}
	}

	if len(displayed) != 2 {
		t.Errorf("expected 2 displayed rows (PENDING+ACTIVE), got %d", len(displayed))
	}
}

// ---- todo update ----

// updateStatusInLines simulates the todoUpdateCmd row-edit logic for tests.
func updateStatusInLines(lines []string, executor, newStatus string) ([]string, bool) {
	found := false
	for i, line := range lines {
		if !isDataRow(line) || extractCol(line, 0) != executor {
			continue
		}
		found = true
		parts := strings.Split(line, "|")
		if len(parts) >= 7 {
			parts[5] = " " + newStatus + " "
			lines[i] = strings.Join(parts, "|")
		}
		break
	}
	return lines, found
}

func TestTodoUpdate_StatusChange(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	row := "| W://root/exec | M://root | 2026-03-10 | task | PENDING | - |"
	initial := todoHeader + "\n" + todoSep + "\n" + row + "\n"
	todoPath := writeTodoFile(t, loadstarBase, initial)

	lines := strings.Split(initial, "\n")
	lines, found := updateStatusInLines(lines, "W://root/exec", "ACTIVE")
	if !found {
		t.Fatal("executor not found during update")
	}
	_ = os.WriteFile(todoPath, []byte(strings.Join(lines, "\n")), 0644)

	data, _ := os.ReadFile(todoPath)
	if !strings.Contains(string(data), "ACTIVE") {
		t.Error("status should have been updated to ACTIVE")
	}
	if strings.Contains(string(data), "| PENDING |") {
		t.Error("old PENDING status should have been replaced")
	}
}

func TestTodoUpdate_PendingToBlocked(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	row := "| W://root/exec | M://root | 2026-03-10 | task | PENDING | - |"
	initial := todoHeader + "\n" + todoSep + "\n" + row + "\n"
	todoPath := writeTodoFile(t, loadstarBase, initial)

	lines := strings.Split(initial, "\n")
	lines, found := updateStatusInLines(lines, "W://root/exec", "BLOCKED")
	if !found {
		t.Fatal("executor not found during update")
	}
	_ = os.WriteFile(todoPath, []byte(strings.Join(lines, "\n")), 0644)

	data, _ := os.ReadFile(todoPath)
	if !strings.Contains(string(data), "BLOCKED") {
		t.Error("status should have been updated to BLOCKED")
	}
}

func TestTodoUpdate_InvalidStatus(t *testing.T) {
	invalidStatuses := []string{"COMPLETED", "FAILED", "DONE", "unknown"}
	for _, s := range invalidStatuses {
		if allowedUpdateStatuses[strings.ToUpper(s)] {
			t.Errorf("status %q should not be allowed in todo update", s)
		}
	}
}

func TestTodoUpdate_AllowedStatuses(t *testing.T) {
	allowed := []string{"PENDING", "ACTIVE", "BLOCKED"}
	for _, s := range allowed {
		if !allowedUpdateStatuses[s] {
			t.Errorf("status %q should be allowed in todo update", s)
		}
	}
}

func TestTodoUpdate_NotFound(t *testing.T) {
	initial := todoHeader + "\n" + todoSep + "\n"
	lines := strings.Split(initial, "\n")
	_, found := updateStatusInLines(lines, "W://root/nonexistent", "ACTIVE")
	if found {
		t.Error("update on nonexistent executor should return not found")
	}
}

func TestTodoUpdate_PreservesOtherCols(t *testing.T) {
	row := "| W://root/exec | M://root | 2026-03-10 12:00 | my summary | PENDING | W://root/dep |"
	lines := []string{todoHeader, todoSep, row}
	lines, found := updateStatusInLines(lines, "W://root/exec", "ACTIVE")
	if !found {
		t.Fatal("executor not found")
	}
	updated := lines[2]

	// All other columns must be preserved
	for _, want := range []string{"W://root/exec", "M://root", "2026-03-10 12:00", "my summary", "W://root/dep"} {
		if !strings.Contains(updated, want) {
			t.Errorf("updated row missing expected value %q: %s", want, updated)
		}
	}
}

// ---- appendTodoHistory: Action column ----

func TestAppendTodoHistory_DoneAction(t *testing.T) {
	dir := t.TempDir()
	histPath := filepath.Join(dir, ".clionly", "TODO", "TODO_HISTORY.md")
	row := "| W://root/exec | M://root | 2026-03-10 12:00 | task summary | PENDING | - |"

	appendTodoHistory(histPath, row, "DONE")

	data, err := os.ReadFile(histPath)
	if err != nil {
		t.Fatalf("TODO_HISTORY.md should have been created: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "DONE") {
		t.Error("history should contain action DONE")
	}
	if !strings.Contains(content, "W://root/exec") {
		t.Error("history should contain executor")
	}
	if !strings.Contains(content, "task summary") {
		t.Error("history should contain summary")
	}
	if !strings.Contains(content, "액션 (Action)") {
		t.Error("history header should contain Action column")
	}
}

func TestAppendTodoHistory_UpdatedAction(t *testing.T) {
	dir := t.TempDir()
	histPath := filepath.Join(dir, ".clionly", "TODO", "TODO_HISTORY.md")
	row := "| W://root/exec | M://root | 2026-03-10 12:00 | task summary | PENDING | - |"

	appendTodoHistory(histPath, row, "UPDATED(PENDING→ACTIVE)")

	data, _ := os.ReadFile(histPath)
	content := string(data)
	if !strings.Contains(content, "UPDATED(PENDING→ACTIVE)") {
		t.Error("history should contain UPDATED action with old and new status")
	}
}

func TestAppendTodoHistory_DeletedAction(t *testing.T) {
	dir := t.TempDir()
	histPath := filepath.Join(dir, ".clionly", "TODO", "TODO_HISTORY.md")
	row := "| W://root/exec | M://root | 2026-03-10 12:00 | task summary | PENDING | - |"

	appendTodoHistory(histPath, row, "DELETED")

	data, _ := os.ReadFile(histPath)
	if !strings.Contains(string(data), "DELETED") {
		t.Error("history should contain action DELETED")
	}
}

func TestAppendTodoHistory_Accumulates(t *testing.T) {
	dir := t.TempDir()
	histPath := filepath.Join(dir, ".clionly", "TODO", "TODO_HISTORY.md")
	row1 := "| W://root/a | NONE | 2026-03-10 12:00 | task a | PENDING | - |"
	row2 := "| W://root/b | NONE | 2026-03-10 13:00 | task b | ACTIVE | - |"

	appendTodoHistory(histPath, row1, "DONE")
	appendTodoHistory(histPath, row2, "UPDATED(ACTIVE→BLOCKED)")

	data, _ := os.ReadFile(histPath)
	content := string(data)
	if !strings.Contains(content, "W://root/a") || !strings.Contains(content, "W://root/b") {
		t.Error("history should accumulate multiple rows")
	}
	if !strings.Contains(content, "DONE") || !strings.Contains(content, "UPDATED(ACTIVE→BLOCKED)") {
		t.Error("history should contain both actions")
	}
}

// ---- todo delete ----

func TestTodoDelete_RemovesRow(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	row := "| W://root/exec | M://root | 2026-03-10 12:00 | task | PENDING | - |"
	initial := todoHeader + "\n" + todoSep + "\n" + row + "\n"
	todoPath := writeTodoFile(t, loadstarBase, initial)

	// Simulate todoDeleteCmd logic
	lines := strings.Split(initial, "\n")
	var kept []string
	var deletedRow string
	for _, line := range lines {
		if isDataRow(line) && extractCol(line, 0) == "W://root/exec" {
			deletedRow = line
			continue
		}
		kept = append(kept, line)
	}
	_ = os.WriteFile(todoPath, []byte(strings.Join(kept, "\n")), 0644)

	data, _ := os.ReadFile(todoPath)
	if strings.Contains(string(data), "W://root/exec") {
		t.Error("deleted row should be removed from TODO_LIST")
	}
	if deletedRow == "" {
		t.Error("deletedRow should have been captured")
	}
}

func TestTodoDelete_RecordsHistory(t *testing.T) {
	dir := t.TempDir()
	histPath := filepath.Join(dir, ".clionly", "TODO", "TODO_HISTORY.md")
	row := "| W://root/exec | M://root | 2026-03-10 12:00 | task to delete | PENDING | - |"

	appendTodoHistory(histPath, row, "DELETED")

	data, _ := os.ReadFile(histPath)
	content := string(data)
	if !strings.Contains(content, "DELETED") {
		t.Error("history should record DELETED action")
	}
	if !strings.Contains(content, "task to delete") {
		t.Error("history should contain the original summary")
	}
}

func TestTodoDelete_NotFound(t *testing.T) {
	initial := todoHeader + "\n" + todoSep + "\n"
	lines := strings.Split(initial, "\n")
	found := false
	for _, line := range lines {
		if isDataRow(line) && extractCol(line, 0) == "W://root/nonexistent" {
			found = true
		}
	}
	if found {
		t.Error("nonexistent executor should not be found for delete")
	}
}

// ---- todo history ----

// filterHistoryLines simulates todoHistoryCmd logic for tests.
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

func TestTodoHistory_AllRows(t *testing.T) {
	dir := t.TempDir()
	histPath := filepath.Join(dir, ".clionly", "TODO", "TODO_HISTORY.md")
	row1 := "| W://root/a | NONE | 2026-03-10 12:00 | task a | DONE | 2026-03-10 13:00 | - |"
	row2 := "| W://root/b | NONE | 2026-03-10 14:00 | task b | DELETED | 2026-03-10 15:00 | - |"

	appendTodoHistory(histPath, "| W://root/a | NONE | 2026-03-10 12:00 | task a | PENDING | - |", "DONE")
	appendTodoHistory(histPath, "| W://root/b | NONE | 2026-03-10 14:00 | task b | PENDING | - |", "DELETED")
	_ = row1
	_ = row2

	content, _ := os.ReadFile(histPath)
	rows := filterHistoryLines(string(content), "")
	if len(rows) != 2 {
		t.Errorf("expected 2 history rows, got %d", len(rows))
	}
}

func TestTodoHistory_FilterByExecutor(t *testing.T) {
	dir := t.TempDir()
	histPath := filepath.Join(dir, ".clionly", "TODO", "TODO_HISTORY.md")

	appendTodoHistory(histPath, "| W://root/a | NONE | 2026-03-10 12:00 | task a | PENDING | - |", "DONE")
	appendTodoHistory(histPath, "| W://root/b | NONE | 2026-03-10 14:00 | task b | PENDING | - |", "DELETED")
	appendTodoHistory(histPath, "| W://root/a | NONE | 2026-03-10 12:00 | task a | PENDING | - |", "UPDATED(PENDING→ACTIVE)")

	content, _ := os.ReadFile(histPath)
	rows := filterHistoryLines(string(content), "W://root/a")
	if len(rows) != 2 {
		t.Errorf("expected 2 rows for W://root/a, got %d", len(rows))
	}
	for _, row := range rows {
		if extractCol(row, 0) != "W://root/a" {
			t.Errorf("filtered row should only contain W://root/a, got: %s", row)
		}
	}
}

func TestTodoHistory_FilterNoMatch(t *testing.T) {
	dir := t.TempDir()
	histPath := filepath.Join(dir, ".clionly", "TODO", "TODO_HISTORY.md")

	appendTodoHistory(histPath, "| W://root/a | NONE | 2026-03-10 12:00 | task a | PENDING | - |", "DONE")

	content, _ := os.ReadFile(histPath)
	rows := filterHistoryLines(string(content), "W://root/nonexistent")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows for nonexistent executor, got %d", len(rows))
	}
}

func TestTodoHistory_Empty(t *testing.T) {
	content := ""
	rows := filterHistoryLines(content, "")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows for empty history, got %d", len(rows))
	}
}
