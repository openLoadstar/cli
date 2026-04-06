package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bono/loadstar/internal/core"
	"github.com/bono/loadstar/internal/storage"
)

// setupCmdTest initialises a temporary .loadstar directory and wires the package-level
// globals (fs, svc) so that cmd functions can be called directly in tests.
func setupCmdTest(t *testing.T) (loadstarBase string) {
	t.Helper()
	root := t.TempDir()

	for _, d := range []string{"MAP", "WAYPOINT", "BLACKBOX", "COMMON",
		".clionly/LOG", ".clionly/MONITOR", ".clionly/TODO"} {
		if err := os.MkdirAll(filepath.Join(root, ".loadstar", d), 0755); err != nil {
			t.Fatalf("mkdir %s: %v", d, err)
		}
	}

	fsInst := storage.NewFS(root)
	fs = fsInst
	svc = core.NewElementService(fsInst)

	return filepath.Join(root, ".loadstar")
}

// writeElement writes raw content to the appropriate type-dir file for the given address.
func writeElement(t *testing.T, loadstarBase, addrStr, content string) string {
	t.Helper()
	parts := strings.SplitN(addrStr, "://", 2)
	if len(parts) != 2 {
		t.Fatalf("bad address: %s", addrStr)
	}
	typeMap := map[string]string{"M": "MAP", "W": "WAYPOINT", "B": "BLACKBOX"}
	dir := typeMap[parts[0]]
	dotName := strings.ReplaceAll(parts[1], "/", ".")
	path := filepath.Join(loadstarBase, dir, dotName+".md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeElement: %v", err)
	}
	return path
}

// parentMapContent returns minimal MAP element content for use as a parent.
func parentMapContent(addr string) string {
	return "<MAP>\n## [ADDRESS] " + addr + "\n## [STATUS] S_STB\n\n### IDENTITY\n- SUMMARY: test\n\n### WAYPOINTS\n(없음)\n\n### COMMENT\n(없음)\n</MAP>\n"
}

// parentWPContent returns minimal WayPoint element content for use as a parent.
func parentWPContent(addr, parent string) string {
	return "<WAYPOINT>\n## [ADDRESS] " + addr + "\n## [STATUS] S_STB\n\n### CONNECTIONS\n- PARENT: " + parent + "\n- CHILDREN: []\n- REFERENCE: []\n- BLACKBOX: B://" + strings.TrimPrefix(addr, "W://") + "\n\n### TODO\n(없음)\n</WAYPOINT>\n"
}

// ---- TestAddress_ToFilePath ----

func TestAddress_ToFilePath(t *testing.T) {
	setupCmdTest(t)

	cases := []struct {
		raw  string
		want string
	}{
		{"W://root/cli/cmd_create", filepath.Join("WAYPOINT", "root.cli.cmd_create.md")},
		{"M://root", filepath.Join("MAP", "root.md")},
		{"B://root/cli/cmd_create", filepath.Join("BLACKBOX", "root.cli.cmd_create.md")},
	}

	for _, c := range cases {
		addr, err := svc.ParseAddress(c.raw)
		if err != nil {
			t.Fatalf("ParseAddress(%s): %v", c.raw, err)
		}
		got := addr.ToFilePath("/base")
		wantFull := filepath.Join("/base", c.want)
		if got != wantFull {
			t.Errorf("ToFilePath(%s) = %q, want %q", c.raw, got, wantFull)
		}
	}
}

// ---- buildTemplate ----

func TestBuildTemplate_ContainsAddress(t *testing.T) {
	for _, typ := range []string{"M", "W"} {
		content := buildTemplate(typ, typ+"://root/x", "M://root")
		if !strings.Contains(content, typ+"://root/x") {
			t.Errorf("buildTemplate(%s) missing address", typ)
		}
	}
}

func TestBuildTemplate_UnknownType(t *testing.T) {
	result := buildTemplate("X", "X://foo", "M://root")
	if result != "" {
		t.Errorf("expected empty string for unknown type, got %q", result)
	}
}

func TestBuildTemplate_WaypointFields(t *testing.T) {
	content := buildTemplate("W", "W://root/id", "M://root")
	mandatory := []string{"<WAYPOINT>", "## [ADDRESS] W://root/id", "CONNECTIONS", "PARENT:", "CHILDREN:", "BLACKBOX:", "</WAYPOINT>"}
	for _, m := range mandatory {
		if !strings.Contains(content, m) {
			t.Errorf("waypoint template missing %q", m)
		}
	}
}

// ---- appendToWaypoints / removeFromWaypoints (MAP) ----

func TestAppendToWaypoints(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	path := writeElement(t, loadstarBase, "M://root", parentMapContent("M://root"))

	if err := appendToWaypoints(path, "W://root/child1"); err != nil {
		t.Fatalf("appendToWaypoints: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "W://root/child1") {
		t.Error("child not found in WAYPOINTS after append")
	}
}

func TestRemoveFromWaypoints(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	path := writeElement(t, loadstarBase, "M://root", parentMapContent("M://root"))

	_ = appendToWaypoints(path, "W://root/child1")
	_ = appendToWaypoints(path, "W://root/child2")

	if err := removeFromWaypoints(path, "W://root/child1"); err != nil {
		t.Fatalf("removeFromWaypoints: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if strings.Contains(content, "W://root/child1") {
		t.Error("child1 should have been removed")
	}
	if !strings.Contains(content, "W://root/child2") {
		t.Error("child2 should still be present")
	}
}

// ---- appendToChildren / removeFromChildren (WayPoint) ----

func TestAppendToChildren(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	path := writeElement(t, loadstarBase, "W://root/parent", parentWPContent("W://root/parent", "M://root"))

	if err := appendToChildren(path, "W://root/parent/child1"); err != nil {
		t.Fatalf("appendToChildren: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "W://root/parent/child1") {
		t.Error("child not found in CHILDREN after append")
	}
}

func TestRemoveFromChildren(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	path := writeElement(t, loadstarBase, "W://root/parent", parentWPContent("W://root/parent", "M://root"))

	_ = appendToChildren(path, "W://root/parent/child1")
	_ = appendToChildren(path, "W://root/parent/child2")

	if err := removeFromChildren(path, "W://root/parent/child1"); err != nil {
		t.Fatalf("removeFromChildren: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if strings.Contains(content, "W://root/parent/child1") {
		t.Error("child1 should have been removed")
	}
	if !strings.Contains(content, "W://root/parent/child2") {
		t.Error("child2 should still be present")
	}
}

// ---- parseParent ----

func TestParseParent(t *testing.T) {
	content := "- PARENT: W://root/cli\n- CHILDREN: []\n"
	got := parseParent(content)
	if got != "W://root/cli" {
		t.Errorf("got %q, want %q", got, "W://root/cli")
	}
}

func TestParseParent_Missing(t *testing.T) {
	got := parseParent("no parent here")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

// ---- Create: TYPE validation ----

func TestCreate_InvalidType(t *testing.T) {
	for _, bad := range []string{"H", "B", "L", "S", "X"} {
		if allowedCreateTypes[bad] {
			t.Errorf("type %q should not be in allowedCreateTypes", bad)
		}
	}
	for _, good := range []string{"M", "W"} {
		if !allowedCreateTypes[good] {
			t.Errorf("type %q should be in allowedCreateTypes", good)
		}
	}
}

// ---- Create: duplicate detection ----

func TestCreate_DuplicateDetected(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	writeElement(t, loadstarBase, "M://root", parentMapContent("M://root"))

	childContent := buildTemplate("W", "W://root/dup", "M://root")
	writeElement(t, loadstarBase, "W://root/dup", childContent)

	addr, _ := svc.ParseAddress("W://root/dup")
	if !fs.Exists(addr.ToFilePath(loadstarBase)) {
		t.Error("expected element to be detected as existing")
	}
}

// ---- resolveEditor ----

func TestResolveEditor_EnvOverride(t *testing.T) {
	t.Setenv("LOADSTAR_EDITOR", "myeditor")
	got := resolveEditor()
	if got != "myeditor" {
		t.Errorf("expected myeditor, got %q", got)
	}
}

func TestResolveEditor_FallbackEditor(t *testing.T) {
	t.Setenv("LOADSTAR_EDITOR", "")
	t.Setenv("EDITOR", "nano")
	got := resolveEditor()
	if got != "nano" {
		t.Errorf("expected nano, got %q", got)
	}
}
