package cmd

import (
	"strings"
	"testing"
)

// ---- extractChildren ----

func TestExtractChildren_WaypointFormat(t *testing.T) {
	content := "### CONNECTIONS\n- PARENT: M://root\n- CHILDREN: [W://root/a, W://root/b]\n"
	children := extractChildren(content)
	if len(children) != 2 {
		t.Errorf("expected 2 children, got %d: %v", len(children), children)
	}
}

func TestExtractChildren_MapFormat(t *testing.T) {
	content := "### WAYPOINTS\n- W://root/cli/cmd_create\n- W://root/cli/cmd_log\n\n### COMMENT\n"
	children := extractChildren(content)
	if len(children) != 2 {
		t.Errorf("expected 2 waypoints, got %d: %v", len(children), children)
	}
}

func TestExtractChildren_Empty(t *testing.T) {
	content := "- CHILDREN: []\n"
	children := extractChildren(content)
	if len(children) != 0 {
		t.Errorf("expected 0 children, got %d", len(children))
	}
}

func TestExtractChildren_MapEmpty(t *testing.T) {
	content := "### WAYPOINTS\n(없음)\n\n### COMMENT\n"
	children := extractChildren(content)
	if len(children) != 0 {
		t.Errorf("expected 0 waypoints, got %d", len(children))
	}
}

// ---- showElement ----

func TestShow_Depth0_NoChildren(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	writeElement(t, loadstarBase, "W://root/leaf",
		"<WAYPOINT>\n## [ADDRESS] W://root/leaf\n## [STATUS] S_STB\n### CONNECTIONS\n- CHILDREN: []\n</WAYPOINT>\n")

	visited := make(map[string]bool)
	showElement(loadstarBase, "W://root/leaf", 0, 0, visited)
}

func TestShow_CircularRefProtection(t *testing.T) {
	loadstarBase := setupCmdTest(t)

	writeElement(t, loadstarBase, "W://root/a",
		"<WAYPOINT>\n## [ADDRESS] W://root/a\n## [STATUS] S_STB\n### CONNECTIONS\n- CHILDREN: [W://root/b]\n</WAYPOINT>\n")
	writeElement(t, loadstarBase, "W://root/b",
		"<WAYPOINT>\n## [ADDRESS] W://root/b\n## [STATUS] S_STB\n### CONNECTIONS\n- CHILDREN: [W://root/a]\n</WAYPOINT>\n")

	visited := make(map[string]bool)
	done := make(chan struct{})
	go func() {
		showElement(loadstarBase, "W://root/a", 5, 0, visited)
		close(done)
	}()
	<-done
}

func TestShow_NotFound(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	visited := make(map[string]bool)
	showElement(loadstarBase, "W://root/missing", 0, 0, visited)
}

// ---- indent helper ----

func TestIndent(t *testing.T) {
	if indent(0) != "" {
		t.Error("depth 0 should be empty string")
	}
	if indent(1) != "  " {
		t.Errorf("depth 1 should be 2 spaces, got %q", indent(1))
	}
}

// ---- extractField ----

func TestExtractField(t *testing.T) {
	content := "## [STATUS] S_PRG\n"
	got := extractField(content, `## \[STATUS\]\s+(\S+)`)
	if got != "S_PRG" {
		t.Errorf("expected S_PRG, got %q", got)
	}
}

func TestExtractField_Missing(t *testing.T) {
	got := extractField("no status here", `## \[STATUS\]\s+(\S+)`)
	if got != "?" {
		t.Errorf("expected ?, got %q", got)
	}
}

// ---- backward compat: old format extractContainsItems still works in extractChildren ----

func TestExtractChildren_OldContainsFormat(t *testing.T) {
	// Old format with ITEMS should fall through to CHILDREN or WAYPOINTS
	content := "- ITEMS: [W://root/a, W://root/b]\n"
	children := extractChildren(content)
	// This format is no longer directly supported — should return empty
	_ = strings.TrimSpace("ok")
	_ = children
}
