package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"
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
	listWaypoints("", false)
}

func TestListWaypoints_Filter(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	writeElement(t, loadstarBase, "W://root/cli/cmd_log",
		"<WAYPOINT>\n## [ADDRESS] W://root/cli/cmd_log\n## [STATUS] S_STB\n</WAYPOINT>\n")
	writeElement(t, loadstarBase, "W://root/test/test_a",
		"<WAYPOINT>\n## [ADDRESS] W://root/test/test_a\n## [STATUS] S_IDL\n</WAYPOINT>\n")

	// Filter by "cli" should only show cmd_log
	listWaypoints("cli", false)
}

func TestListWaypoints_Empty(t *testing.T) {
	setupCmdTest(t)
	// No waypoint files — should print "no waypoints found"
	listWaypoints("", false)
}

func TestListWaypoints_Recent(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	writeElement(t, loadstarBase, "W://root/older",
		"<WAYPOINT>\n## [ADDRESS] W://root/older\n## [STATUS] S_STB\n</WAYPOINT>\n")
	// Backdate the older file by 1 hour so mtime ordering is deterministic
	olderPath := filepath.Join(loadstarBase, "WAYPOINT", "root.older.md")
	past := time.Now().Add(-1 * time.Hour)
	if err := os.Chtimes(olderPath, past, past); err != nil {
		t.Fatalf("chtimes failed: %v", err)
	}

	writeElement(t, loadstarBase, "W://root/newer",
		"<WAYPOINT>\n## [ADDRESS] W://root/newer\n## [STATUS] S_STB\n</WAYPOINT>\n")

	// Should not panic in either mode
	listWaypoints("", false)
	listWaypoints("", true)
}

func TestFormatMTime_Zero(t *testing.T) {
	if got := formatMTime(time.Time{}); got != "—" {
		t.Errorf("expected dash for zero time, got %q", got)
	}
}

func TestFormatMTime_Format(t *testing.T) {
	ts := time.Date(2026, 4, 28, 15, 47, 0, 0, time.Local)
	if got := formatMTime(ts); got != "2026-04-28 15:47" {
		t.Errorf("expected 2026-04-28 15:47, got %q", got)
	}
}
