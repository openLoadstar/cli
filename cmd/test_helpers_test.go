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

	for _, d := range []string{"MAP", "WAYPOINT", "COMMON",
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
	typeMap := map[string]string{"M": "MAP", "W": "WAYPOINT"}
	dir := typeMap[parts[0]]
	if dir == "" {
		t.Fatalf("unsupported type in address: %s", addrStr)
	}
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
	return "<WAYPOINT>\n## [ADDRESS] " + addr + "\n## [STATUS] S_STB\n\n### CONNECTIONS\n- PARENT: " + parent + "\n- CHILDREN: []\n- REFERENCE: []\n\n### TODO\n(없음)\n</WAYPOINT>\n"
}
