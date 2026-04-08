package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---- isDataRow ----

func TestIsDataRow_Header(t *testing.T) {
	header := "| 주소 (Address) | 상태 (Status) | 작업 요약 (Summary) |"
	if isDataRow(header) {
		t.Error("header row should not be a data row")
	}
}

func TestIsDataRow_Separator(t *testing.T) {
	sep := "| :--- | :--- | :--- |"
	if isDataRow(sep) {
		t.Error("separator row should not be a data row")
	}
}

func TestIsDataRow_DataRow(t *testing.T) {
	row := "| W://root/cli/cmd_show | ACTIVE | show 명령 구현 |"
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
	row := "| W://root/cmd | ACTIVE | summary text |"
	cases := []struct {
		col  int
		want string
	}{
		{0, "W://root/cmd"},
		{1, "ACTIVE"},
		{2, "summary text"},
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

// ---- wpStatusToTodoStatus ----

func TestWpStatusToTodoStatus(t *testing.T) {
	cases := []struct {
		wpStatus string
		want     string
	}{
		{"S_IDL", "PENDING"},
		{"S_PRG", "ACTIVE"},
		{"S_ERR", "ACTIVE"},
		{"S_REV", "ACTIVE"},
		{"S_STB", ""},
	}
	for _, c := range cases {
		got := wpStatusToTodoStatus(c.wpStatus)
		if got != c.want {
			t.Errorf("wpStatusToTodoStatus(%q) = %q, want %q", c.wpStatus, got, c.want)
		}
	}
}

// ---- saveTodoList / loadTodoList ----

func TestSaveAndLoadTodoList(t *testing.T) {
	loadstarBase := setupCmdTest(t)

	items := []todoItem{
		{Address: "W://root/a", Status: "ACTIVE", Summary: "task a"},
		{Address: "W://root/b", Status: "PENDING", Summary: "task b"},
	}

	saveTodoList(loadstarBase, items)

	loaded := loadTodoList(loadstarBase)
	if len(loaded) != 2 {
		t.Fatalf("expected 2 items, got %d", len(loaded))
	}
	if loaded[0].Address != "W://root/a" || loaded[0].Status != "ACTIVE" {
		t.Errorf("unexpected first item: %+v", loaded[0])
	}
	if loaded[1].Address != "W://root/b" || loaded[1].Status != "PENDING" {
		t.Errorf("unexpected second item: %+v", loaded[1])
	}
}

// ---- snapshot ----

func TestSaveAndLoadSnapshot(t *testing.T) {
	loadstarBase := setupCmdTest(t)

	snap := map[string]wpSnapshot{
		"W://root/a": {ModTime: "2026-04-08T12:00:00Z", Size: 100, Status: "S_PRG"},
	}

	saveSnapshot(loadstarBase, snap)
	loaded := loadSnapshot(loadstarBase)

	if loaded["W://root/a"].Status != "S_PRG" {
		t.Errorf("snapshot status mismatch: %+v", loaded["W://root/a"])
	}
}

// ---- readWPStatusAndSummary ----

func TestReadWPStatusAndSummary(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	writeElement(t, loadstarBase, "W://root/test",
		"<WAYPOINT>\n## [ADDRESS] W://root/test\n## [STATUS] S_PRG\n\n### IDENTITY\n- SUMMARY: test summary\n</WAYPOINT>\n")

	status, summary := readWPStatusAndSummary(filepath.Join(loadstarBase, "WAYPOINT", "root.test.md"))
	if status != "S_PRG" {
		t.Errorf("expected S_PRG, got %q", status)
	}
	if summary != "test summary" {
		t.Errorf("expected 'test summary', got %q", summary)
	}
}

// ---- sync: full sync ----

func TestTodoSync_FullSync(t *testing.T) {
	loadstarBase := setupCmdTest(t)

	// Create MAP with 2 WPs
	writeElement(t, loadstarBase, "M://root",
		"<MAP>\n## [ADDRESS] M://root\n## [STATUS] S_PRG\n\n### WAYPOINTS\n- W://root/a\n- W://root/b\n\n### COMMENT\n</MAP>\n")
	writeElement(t, loadstarBase, "W://root/a",
		"<WAYPOINT>\n## [ADDRESS] W://root/a\n## [STATUS] S_PRG\n\n### IDENTITY\n- SUMMARY: task a\n\n### CONNECTIONS\n- PARENT: M://root\n- CHILDREN: []\n- REFERENCE: []\n</WAYPOINT>\n")
	writeElement(t, loadstarBase, "W://root/b",
		"<WAYPOINT>\n## [ADDRESS] W://root/b\n## [STATUS] S_STB\n\n### IDENTITY\n- SUMMARY: task b done\n\n### CONNECTIONS\n- PARENT: M://root\n- CHILDREN: []\n- REFERENCE: []\n</WAYPOINT>\n")

	// Run sync
	wpAddrs := collectAllWaypoints(loadstarBase)
	if len(wpAddrs) != 2 {
		t.Fatalf("expected 2 WP addresses, got %d", len(wpAddrs))
	}

	// Simulate sync logic
	snapshot := loadSnapshot(loadstarBase)
	items := loadTodoList(loadstarBase)
	itemMap := make(map[string]*todoItem)
	for i := range items {
		itemMap[items[i].Address] = &items[i]
	}

	for _, addr := range wpAddrs {
		wpFile := addressToFilePath(loadstarBase, addr)
		status, summary := readWPStatusAndSummary(wpFile)
		info, _ := os.Stat(wpFile)
		snapshot[addr] = wpSnapshot{ModTime: info.ModTime().String(), Size: info.Size(), Status: status}

		todoStatus := wpStatusToTodoStatus(status)
		if todoStatus == "" {
			delete(itemMap, addr)
		} else {
			itemMap[addr] = &todoItem{Address: addr, Status: todoStatus, Summary: summary}
		}
	}

	var result []todoItem
	for _, item := range itemMap {
		result = append(result, *item)
	}

	saveTodoList(loadstarBase, result)

	// Verify: only S_PRG should be in list (S_STB excluded)
	loaded := loadTodoList(loadstarBase)
	if len(loaded) != 1 {
		t.Fatalf("expected 1 item (S_STB excluded), got %d", len(loaded))
	}
	if loaded[0].Address != "W://root/a" {
		t.Errorf("expected W://root/a, got %q", loaded[0].Address)
	}
}

// ---- history: TECH_SPEC [x] parsing ----

func TestTodoHistory_ParsesTechSpec(t *testing.T) {
	loadstarBase := setupCmdTest(t)

	writeElement(t, loadstarBase, "M://root",
		"<MAP>\n## [ADDRESS] M://root\n## [STATUS] S_STB\n\n### WAYPOINTS\n- W://root/done\n\n### COMMENT\n</MAP>\n")
	writeElement(t, loadstarBase, "W://root/done",
		"<WAYPOINT>\n## [ADDRESS] W://root/done\n## [STATUS] S_STB\n\n### IDENTITY\n- SUMMARY: completed\n\n### TODO\n- TECH_SPEC:\n  - [x] 2026-04-08 implement feature A\n  - [x] 2026-04-07 design schema\n  - [ ] not yet done\n</WAYPOINT>\n")

	wpAddrs := collectAllWaypoints(loadstarBase)
	if len(wpAddrs) != 1 {
		t.Fatalf("expected 1 WP, got %d", len(wpAddrs))
	}

	// Parse completed items manually (same logic as history cmd)
	wpFile := addressToFilePath(loadstarBase, "W://root/done")
	data, _ := os.ReadFile(wpFile)
	content := string(data)

	count := strings.Count(content, "[x]")
	if count != 2 {
		t.Errorf("expected 2 [x] items, got %d", count)
	}
}

// ---- extractReferences ----

func TestExtractReferences(t *testing.T) {
	content := "### CONNECTIONS\n- PARENT: M://root\n- CHILDREN: []\n- REFERENCE: [W://root/a, W://root/b]\n"
	refs := extractReferences(content)
	if len(refs) != 2 {
		t.Errorf("expected 2 refs, got %d: %v", len(refs), refs)
	}
}

func TestExtractReferences_Empty(t *testing.T) {
	content := "- REFERENCE: []\n"
	refs := extractReferences(content)
	if len(refs) != 0 {
		t.Errorf("expected 0 refs, got %d", len(refs))
	}
}
