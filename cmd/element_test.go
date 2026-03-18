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
// Returns the loadstarBase path.
func setupCmdTest(t *testing.T) (loadstarBase string) {
	t.Helper()
	root := t.TempDir()

	// Create required sub-directories
	for _, d := range []string{"MAP", "WAYPOINT", "LINK", "SAVEPOINT", "BLACKBOX", "HISTORY", "COMMON"} {
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
	typeMap := map[string]string{"M": "MAP", "W": "WAYPOINT", "L": "LINK", "S": "SAVEPOINT"}
	dir := typeMap[parts[0]]
	dotName := strings.ReplaceAll(parts[1], "/", ".")
	path := filepath.Join(loadstarBase, dir, dotName+".md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeElement: %v", err)
	}
	return path
}

// parentContent returns minimal MAP element content for use as a parent.
func parentContent(addr string) string {
	return "<MAP>\n## [ADDRESS] " + addr + "\n## [STATUS] S_STB\n\n### 2. CONTAINS\n- ITEMS: []\n\n### 3. CONNECTIONS\n- LINEAGE: [PARENT: -, CHILDREN: []]\n- LINKS: []\n</MAP>\n"
}

// ---- TestAddress_ToFilePath ----

func TestAddress_ToFilePath(t *testing.T) {
	setupCmdTest(t)

	cases := []struct {
		raw  string
		want string // relative to loadstarBase, using OS separator
	}{
		{"W://root/cli/cmd_create", filepath.Join("WAYPOINT", "root.cli.cmd_create.md")},
		{"M://root", filepath.Join("MAP", "root.md")},
		{"S://root/sp1", filepath.Join("SAVEPOINT", "root.sp1.md")},
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
	cases := []string{"M", "W", "L", "S"}
	for _, typ := range cases {
		content := buildTemplate(typ, typ+"://root/x", "M://root")
		if !strings.Contains(content, typ+"://root/x") {
			t.Errorf("buildTemplate(%s) missing address", typ)
		}
		if !strings.Contains(content, "M://root") {
			t.Errorf("buildTemplate(%s) missing parent", typ)
		}
	}
}

func TestBuildTemplate_UnknownType(t *testing.T) {
	result := buildTemplate("X", "X://foo", "M://root")
	if result != "" {
		t.Errorf("expected empty string for unknown type, got %q", result)
	}
}

// ---- appendToContains / removeFromContains ----

func TestAppendToContains(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	path := writeElement(t, loadstarBase, "M://root", parentContent("M://root"))

	if err := appendToContains(path, "W://root/child1"); err != nil {
		t.Fatalf("appendToContains: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "W://root/child1") {
		t.Error("child not found in ITEMS after appendToContains")
	}
}

func TestAppendToContains_SecondChild(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	path := writeElement(t, loadstarBase, "M://root", parentContent("M://root"))

	_ = appendToContains(path, "W://root/child1")
	_ = appendToContains(path, "W://root/child2")

	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, "W://root/child1") || !strings.Contains(content, "W://root/child2") {
		t.Error("both children should be in ITEMS")
	}
}

func TestRemoveFromContains(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	path := writeElement(t, loadstarBase, "M://root", parentContent("M://root"))

	_ = appendToContains(path, "W://root/child1")
	_ = appendToContains(path, "W://root/child2")

	if err := removeFromContains(path, "W://root/child1"); err != nil {
		t.Fatalf("removeFromContains: %v", err)
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

// ---- parseLineageParent ----

func TestParseLineageParent(t *testing.T) {
	content := "- LINEAGE: [PARENT: M://root/cli, CHILDREN: []]\n"
	got := parseLineageParent(content)
	if got != "M://root/cli" {
		t.Errorf("got %q, want %q", got, "M://root/cli")
	}
}

func TestParseLineageParent_Missing(t *testing.T) {
	got := parseLineageParent("no lineage here")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

// ---- Create: TYPE validation ----

func TestCreate_InvalidType(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	writeElement(t, loadstarBase, "M://root", parentContent("M://root"))

	// allowedCreateTypes excludes "H" and "B"
	for _, bad := range []string{"H", "B", "X", "z"} {
		if allowedCreateTypes[bad] {
			t.Errorf("type %q should not be in allowedCreateTypes", bad)
		}
	}
	for _, good := range []string{"M", "W", "L", "S"} {
		if !allowedCreateTypes[good] {
			t.Errorf("type %q should be in allowedCreateTypes", good)
		}
	}
}

// ---- Create: duplicate detection via file presence ----

func TestCreate_DuplicateDetected(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	writeElement(t, loadstarBase, "M://root", parentContent("M://root"))

	// Pre-create child
	childContent := buildTemplate("W", "W://root/dup", "M://root")
	writeElement(t, loadstarBase, "W://root/dup", childContent)

	// Verify Exists returns true for the duplicate
	addr, _ := svc.ParseAddress("W://root/dup")
	if !fs.Exists(addr.ToFilePath(loadstarBase)) {
		t.Error("expected element to be detected as existing")
	}
}

// ---- Create: parent CONTAINS updated ----

func TestCreate_ParentContains(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	parentPath := writeElement(t, loadstarBase, "M://root", parentContent("M://root"))

	newAddrStr := "W://root/newchild"
	content := buildTemplate("W", newAddrStr, "M://root")
	newAddr, _ := svc.ParseAddress(newAddrStr)
	newFile := newAddr.ToFilePath(loadstarBase)

	if err := fs.Write(newFile, content); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := appendToContains(parentPath, newAddrStr); err != nil {
		t.Fatalf("appendToContains: %v", err)
	}

	data, _ := os.ReadFile(parentPath)
	if !strings.Contains(string(data), newAddrStr) {
		t.Error("parent CONTAINS.ITEMS should contain new child address")
	}
}

// ---- Edit: Shadow History snapshot ----

func TestEdit_ShadowHistory(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	elemPath := writeElement(t, loadstarBase, "W://root/elem", buildTemplate("W", "W://root/elem", "M://root"))

	histDir := filepath.Join(loadstarBase, "HISTORY")
	// Simulate Shadow History creation (as done in editCmd)
	histPath := filepath.Join(histDir, "root.elem_20260101T000000.md")
	if err := fs.CopyFile(elemPath, histPath); err != nil {
		t.Fatalf("CopyFile: %v", err)
	}

	if _, err := os.Stat(histPath); os.IsNotExist(err) {
		t.Error("history snapshot file should exist after shadow copy")
	}
}

// ---- Delete: History First backup ----

func TestDelete_HistoryBackup(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	elemPath := writeElement(t, loadstarBase, "W://root/victim", buildTemplate("W", "W://root/victim", "M://root"))

	histDir := filepath.Join(loadstarBase, "HISTORY")
	histPath := filepath.Join(histDir, "root.victim_20260101T000000_deleted.md")

	// Simulate History First backup
	if err := fs.CopyFile(elemPath, histPath); err != nil {
		t.Fatalf("CopyFile: %v", err)
	}
	if err := os.Remove(elemPath); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	if _, err := os.Stat(histPath); os.IsNotExist(err) {
		t.Error("history backup should exist after delete")
	}
	if _, err := os.Stat(elemPath); !os.IsNotExist(err) {
		t.Error("original file should be gone after delete")
	}
}

// ---- Delete: parent CONTAINS updated ----

func TestDelete_ParentContainsRemoved(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	parentPath := writeElement(t, loadstarBase, "M://root", parentContent("M://root"))
	elemPath := writeElement(t, loadstarBase, "W://root/child", buildTemplate("W", "W://root/child", "M://root"))

	_ = appendToContains(parentPath, "W://root/child")

	// Simulate delete: remove from parent CONTAINS, then delete file
	if err := removeFromContains(parentPath, "W://root/child"); err != nil {
		t.Fatalf("removeFromContains: %v", err)
	}
	_ = os.Remove(elemPath)

	data, _ := os.ReadFile(parentPath)
	if strings.Contains(string(data), "W://root/child") {
		t.Error("deleted child should be removed from parent CONTAINS")
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

// ---- Additional behaviors ----

func TestAppendToContains_NoItemsLine(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	path := writeElement(t, loadstarBase, "M://root", "<MAP>\n## [ADDRESS] M://root\n## [STATUS] S_STB\n\n### 2. CONTAINS\n- PAYLOAD:\n\n### 3. CONNECTIONS\n- LINEAGE: [PARENT: -, CHILDREN: []]\n- LINKS: []\n</MAP>\n")
	if err := appendToContains(path, "W://root/child1"); err == nil {
		t.Fatal("expected error when ITEMS line is missing")
	}
}

func TestRemoveFromContains_MissingTargetNoPanic(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	path := writeElement(t, loadstarBase, "M://root", parentContent("M://root"))
	// Removing an address that doesn't exist should succeed and keep line unchanged
	original, _ := os.ReadFile(path)
	if err := removeFromContains(path, "W://root/absent"); err != nil {
		t.Fatalf("removeFromContains unexpected error: %v", err)
	}
	after, _ := os.ReadFile(path)
	if string(after) != string(original) {
		t.Error("content should remain unchanged when removing non-existent child")
	}
}

func TestBuildTemplate_WaypointFieldsPresent(t *testing.T) {
	content := buildTemplate("W", "W://root/id", "M://root")
	mandatory := []string{"<WAYPOINT>", "## [ADDRESS] W://root/id", "- EXECUTOR:", "- RESPONSE_STATUS:", "</WAYPOINT>"}
	for _, m := range mandatory {
		if !strings.Contains(content, m) {
			t.Errorf("waypoint template missing %q", m)
		}
	}
}

func TestParseLineageParent_ExtraSpaces(t *testing.T) {
	line := "- LINEAGE: [PARENT:   M://root/cli  , CHILDREN: []]"
	got := parseLineageParent(line)
	if got != "M://root/cli" {
		t.Errorf("got %q, want %q", got, "M://root/cli")
	}
}

func TestAddress_ToFilePath_RelativeBaseHandling(t *testing.T) {
	setupCmdTest(t)
	addr, err := svc.ParseAddress("L://root/edge/case")
	if err != nil {
		t.Fatalf("ParseAddress: %v", err)
	}
	got := addr.ToFilePath("/base")
	want := filepath.Join("/base", "LINK", "root.edge.case.md")
	if got != want {
		t.Errorf("ToFilePath = %q, want %q", got, want)
	}
}
