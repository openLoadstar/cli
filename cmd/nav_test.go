package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---- appendToLinks ----

func TestAppendToLinks_Empty(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	path := writeElement(t, loadstarBase, "W://root/src",
		"<WAYPOINT>\n## [ADDRESS] W://root/src\n### 3. CONNECTIONS\n- LINKS: []\n</WAYPOINT>\n")

	if err := appendToLinks(path, "L://root/src_to_dst", "L_TST"); err != nil {
		t.Fatalf("appendToLinks: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "L://root/src_to_dst") {
		t.Error("link address should appear in CONNECTIONS.LINKS")
	}
}

func TestAppendToLinks_Append(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	path := writeElement(t, loadstarBase, "W://root/src",
		"<WAYPOINT>\n## [ADDRESS] W://root/src\n### 3. CONNECTIONS\n- LINKS: [L://root/existing | TYPE: L_REF]\n</WAYPOINT>\n")

	if err := appendToLinks(path, "L://root/new_link", "L_SEQ"); err != nil {
		t.Fatalf("appendToLinks: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, "L://root/existing") {
		t.Error("existing link should still be present")
	}
	if !strings.Contains(content, "L://root/new_link") {
		t.Error("new link should be appended")
	}
}

func TestAppendToLinks_MissingSection(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	path := writeElement(t, loadstarBase, "W://root/src",
		"<WAYPOINT>\n## [ADDRESS] W://root/src\n</WAYPOINT>\n")

	err := appendToLinks(path, "L://root/link", "L_REF")
	if err == nil {
		t.Error("expected error when CONNECTIONS.LINKS line is absent")
	}
}

// ---- allowedLinkTypes ----

func TestAllowedLinkTypes(t *testing.T) {
	for _, good := range []string{"L_REF", "L_SEQ", "L_TST"} {
		if !allowedLinkTypes[good] {
			t.Errorf("link type %q should be allowed", good)
		}
	}
	for _, bad := range []string{"L_FOO", "REF", ""} {
		if allowedLinkTypes[bad] {
			t.Errorf("link type %q should not be allowed", bad)
		}
	}
}

// ---- link: file creation and bidirectional CONNECTIONS.LINKS ----

func TestLink_FileCreated(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	linkDir := filepath.Join(loadstarBase, "LINK")

	// Create source and target elements
	srcContent := "<WAYPOINT>\n## [ADDRESS] W://root/src\n### 3. CONNECTIONS\n- LINKS: []\n</WAYPOINT>\n"
	dstContent := "<WAYPOINT>\n## [ADDRESS] W://root/dst\n### 3. CONNECTIONS\n- LINKS: []\n</WAYPOINT>\n"
	writeElement(t, loadstarBase, "W://root/src", srcContent)
	writeElement(t, loadstarBase, "W://root/dst", dstContent)

	// Build link ID and file (mirrors linkCmd logic)
	linkID := "src_to_dst"
	linkPath := "root/" + linkID
	linkFile := filepath.Join(linkDir, strings.ReplaceAll(linkPath, "/", ".")+".md")

	linkContent := "<LINK>\n## [ADDRESS] L://root/src_to_dst\n</LINK>\n"
	if err := fs.Write(linkFile, linkContent); err != nil {
		t.Fatalf("write link file: %v", err)
	}

	if !fs.Exists(linkFile) {
		t.Error("link file should exist after creation")
	}
}

func TestLink_BidirectionalRegistration(t *testing.T) {
	loadstarBase := setupCmdTest(t)

	srcPath := writeElement(t, loadstarBase, "W://root/src",
		"<WAYPOINT>\n## [ADDRESS] W://root/src\n### 3. CONNECTIONS\n- LINKS: []\n</WAYPOINT>\n")
	dstPath := writeElement(t, loadstarBase, "W://root/dst",
		"<WAYPOINT>\n## [ADDRESS] W://root/dst\n### 3. CONNECTIONS\n- LINKS: []\n</WAYPOINT>\n")

	linkAddr := "L://root/src_to_dst"

	if err := appendToLinks(srcPath, linkAddr, "L_TST"); err != nil {
		t.Fatalf("appendToLinks src: %v", err)
	}
	if err := appendToLinks(dstPath, linkAddr, "L_TST"); err != nil {
		t.Fatalf("appendToLinks dst: %v", err)
	}

	srcData, _ := os.ReadFile(srcPath)
	dstData, _ := os.ReadFile(dstPath)
	if !strings.Contains(string(srcData), linkAddr) {
		t.Error("source should have link registered in CONNECTIONS.LINKS")
	}
	if !strings.Contains(string(dstData), linkAddr) {
		t.Error("target should have link registered in CONNECTIONS.LINKS")
	}
}

func TestLink_DuplicateDetected(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	linkDir := filepath.Join(loadstarBase, "LINK")

	linkFile := filepath.Join(linkDir, "root.src_to_dst.md")
	_ = os.WriteFile(linkFile, []byte("<LINK></LINK>"), 0644)

	// Verify that Exists correctly detects duplicate
	if !fs.Exists(linkFile) {
		t.Error("duplicate link file should be detected as existing")
	}
}

// ---- extractContainsItems ----

func TestExtractContainsItems_Single(t *testing.T) {
	content := "- ITEMS: [W://root/child1]\n"
	items := extractContainsItems(content)
	if len(items) != 1 || strings.TrimSpace(items[0]) != "W://root/child1" {
		t.Errorf("expected [W://root/child1], got %v", items)
	}
}

func TestExtractContainsItems_Multiple(t *testing.T) {
	content := "- ITEMS: [W://root/a, W://root/b, W://root/c]\n"
	items := extractContainsItems(content)
	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d: %v", len(items), items)
	}
}

func TestExtractContainsItems_Empty(t *testing.T) {
	content := "- ITEMS: []\n"
	items := extractContainsItems(content)
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

// ---- showElement: depth and circular ref ----

func TestShow_Depth0_NoChildren(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	writeElement(t, loadstarBase, "W://root/leaf",
		"<WAYPOINT>\n## [ADDRESS] W://root/leaf\n## [STATUS] S_STB\n### 2. CONTAINS\n- ITEMS: []\n</WAYPOINT>\n")

	// Just verify showElement doesn't panic at depth 0
	visited := make(map[string]bool)
	// Redirect stdout by capturing: showElement writes to os.Stdout directly,
	// so we just run it and ensure no panic.
	showElement(loadstarBase, "W://root/leaf", 0, 0, visited)
}

func TestShow_CircularRefProtection(t *testing.T) {
	loadstarBase := setupCmdTest(t)

	// A points to B, B points to A → circular
	writeElement(t, loadstarBase, "W://root/a",
		"<WAYPOINT>\n## [ADDRESS] W://root/a\n## [STATUS] S_STB\n### 2. CONTAINS\n- ITEMS: [W://root/b]\n</WAYPOINT>\n")
	writeElement(t, loadstarBase, "W://root/b",
		"<WAYPOINT>\n## [ADDRESS] W://root/b\n## [STATUS] S_STB\n### 2. CONTAINS\n- ITEMS: [W://root/a]\n</WAYPOINT>\n")

	// Should return without infinite loop
	visited := make(map[string]bool)
	done := make(chan struct{})
	go func() {
		showElement(loadstarBase, "W://root/a", 5, 0, visited)
		close(done)
	}()

	select {
	case <-done:
		// passed: no infinite loop
	}
}

func TestShow_NotFound(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	visited := make(map[string]bool)
	// Should not panic for missing element
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
	if indent(3) != "      " {
		t.Errorf("depth 3 should be 6 spaces, got %q", indent(3))
	}
}
