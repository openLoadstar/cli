package cmd

import (
	"testing"
)

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

// ---- listWaypoints ----

func TestListWaypoints_All(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	writeElement(t, loadstarBase, "W://root/a",
		"<WAYPOINT>\n## [ADDRESS] W://root/a\n## [STATUS] S_STB\n</WAYPOINT>\n")
	writeElement(t, loadstarBase, "W://root/b",
		"<WAYPOINT>\n## [ADDRESS] W://root/b\n## [STATUS] S_PRG\n</WAYPOINT>\n")

	// Should not panic and list both
	listWaypoints("")
}

func TestListWaypoints_Filter(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	writeElement(t, loadstarBase, "W://root/cli/cmd_log",
		"<WAYPOINT>\n## [ADDRESS] W://root/cli/cmd_log\n## [STATUS] S_STB\n</WAYPOINT>\n")
	writeElement(t, loadstarBase, "W://root/test/test_a",
		"<WAYPOINT>\n## [ADDRESS] W://root/test/test_a\n## [STATUS] S_IDL\n</WAYPOINT>\n")

	// Filter by "cli" should only show cmd_log
	listWaypoints("cli")
}

func TestListWaypoints_Empty(t *testing.T) {
	setupCmdTest(t)
	// No waypoint files — should print "no waypoints found"
	listWaypoints("")
}
