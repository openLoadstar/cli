package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func wpWithQuestions(addr, parent string, questions []string) string {
	qBlock := ""
	if len(questions) > 0 {
		qBlock = "\n### ISSUE\n- OPEN_QUESTIONS:\n"
		for _, q := range questions {
			qBlock += "  " + q + "\n"
		}
	}
	return "<WAYPOINT>\n## [ADDRESS] " + addr + "\n## [STATUS] S_IDL\n\n### IDENTITY\n- SUMMARY: test\n\n### CONNECTIONS\n- PARENT: " + parent + "\n- CHILDREN: []\n- REFERENCE: []" + qBlock + "</WAYPOINT>\n"
}

func setupQuestionTest(t *testing.T) string {
	t.Helper()
	loadstarBase := setupCmdTest(t)

	writeElement(t, loadstarBase, "M://root",
		"<MAP>\n## [ADDRESS] M://root\n## [STATUS] S_STB\n\n### IDENTITY\n- SUMMARY: root\n\n### WAYPOINTS\n- W://root/wp_a\n- W://root/wp_b\n- W://root/wp_c\n\n### COMMENT\n(없음)\n</MAP>\n")

	writeElement(t, loadstarBase, "W://root/wp_a",
		wpWithQuestions("W://root/wp_a", "M://root", []string{
			"- [Q1] Should we use cascade deletion?",
			"- [Q2 DEFERRED] What about orphan maps?",
		}))

	writeElement(t, loadstarBase, "W://root/wp_b",
		wpWithQuestions("W://root/wp_b", "M://root", []string{
			"- [Q1 RESOLVED wp_b.2026-04-28.001] Chose option C.",
		}))

	writeElement(t, loadstarBase, "W://root/wp_c",
		wpWithQuestions("W://root/wp_c", "M://root", []string{
			"- [Q1 DONE wp_c.2026-04-28.001] Applied to code.",
		}))

	return filepath.Join(loadstarBase)
}

// ---- scanQuestions ----

func TestScanQuestions_ReturnsAllStates(t *testing.T) {
	loadstarBase := setupQuestionTest(t)
	entries := scanQuestions(loadstarBase)

	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}
}

func TestScanQuestions_StatesCorrect(t *testing.T) {
	loadstarBase := setupQuestionTest(t)
	entries := scanQuestions(loadstarBase)

	stateMap := map[string]string{}
	for _, e := range entries {
		stateMap[e.Address+"/"+e.QID] = e.State
	}

	if stateMap["W://root/wp_a/Q1"] != "OPEN" {
		t.Errorf("Q1 of wp_a should be OPEN, got %s", stateMap["W://root/wp_a/Q1"])
	}
	if stateMap["W://root/wp_a/Q2"] != "DEFERRED" {
		t.Errorf("Q2 of wp_a should be DEFERRED, got %s", stateMap["W://root/wp_a/Q2"])
	}
	if stateMap["W://root/wp_b/Q1"] != "RESOLVED" {
		t.Errorf("Q1 of wp_b should be RESOLVED, got %s", stateMap["W://root/wp_b/Q1"])
	}
	if stateMap["W://root/wp_c/Q1"] != "DONE" {
		t.Errorf("Q1 of wp_c should be DONE, got %s", stateMap["W://root/wp_c/Q1"])
	}
}

func TestScanQuestions_RefParsed(t *testing.T) {
	loadstarBase := setupQuestionTest(t)
	entries := scanQuestions(loadstarBase)

	for _, e := range entries {
		if e.Address == "W://root/wp_b" && e.QID == "Q1" {
			if e.Ref != "wp_b.2026-04-28.001" {
				t.Errorf("expected ref wp_b.2026-04-28.001, got %q", e.Ref)
			}
			return
		}
	}
	t.Error("wp_b Q1 not found")
}

func TestScanQuestions_Empty(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	writeElement(t, loadstarBase, "M://root",
		"<MAP>\n## [ADDRESS] M://root\n## [STATUS] S_STB\n\n### IDENTITY\n- SUMMARY: root\n\n### WAYPOINTS\n- W://root/wp_empty\n\n### COMMENT\n(없음)\n</MAP>\n")
	writeElement(t, loadstarBase, "W://root/wp_empty",
		wpWithQuestions("W://root/wp_empty", "M://root", nil))

	entries := scanQuestions(loadstarBase)
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

// ---- qRe regex ----

func TestQRe_Open(t *testing.T) {
	line := "  - [Q3] Is this the right approach?"
	m := qRe.FindStringSubmatch(line)
	if m == nil {
		t.Fatal("expected match")
	}
	// m[1]=digit, m[2]=state, m[3]=ref, m[4]=text
	if m[1] != "3" || m[2] != "" || m[4] != "Is this the right approach?" {
		t.Errorf("unexpected match: %v", m)
	}
}

func TestQRe_Deferred(t *testing.T) {
	line := "  - [Q2 DEFERRED] Postponed for v2"
	m := qRe.FindStringSubmatch(line)
	if m == nil {
		t.Fatal("expected match")
	}
	if m[2] != "DEFERRED" {
		t.Errorf("expected DEFERRED, got %q", m[2])
	}
	if m[3] != "" {
		t.Errorf("expected no ref, got %q", m[3])
	}
}

func TestQRe_Resolved(t *testing.T) {
	line := "  - [Q1 RESOLVED wp_b.2026-04-28.001] Chose option C."
	m := qRe.FindStringSubmatch(line)
	if m == nil {
		t.Fatal("expected match")
	}
	if m[2] != "RESOLVED" {
		t.Errorf("expected RESOLVED, got %q", m[2])
	}
	if m[3] != "wp_b.2026-04-28.001" {
		t.Errorf("expected ref, got %q", m[3])
	}
}

func TestQRe_Done(t *testing.T) {
	line := "  - [Q1 DONE wp_c.2026-04-28.001] Applied to code."
	m := qRe.FindStringSubmatch(line)
	if m == nil {
		t.Fatal("expected match")
	}
	if m[2] != "DONE" {
		t.Errorf("expected DONE, got %q", m[2])
	}
	if m[3] != "wp_c.2026-04-28.001" {
		t.Errorf("expected ref, got %q", m[3])
	}
}

// ---- updateQStateDone ----

func TestUpdateQStateDone_ResolvedToDone(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	writeElement(t, loadstarBase, "M://root",
		"<MAP>\n## [ADDRESS] M://root\n## [STATUS] S_STB\n\n### IDENTITY\n- SUMMARY: root\n\n### WAYPOINTS\n- W://root/wp_x\n\n### COMMENT\n(없음)\n</MAP>\n")
	writeElement(t, loadstarBase, "W://root/wp_x",
		wpWithQuestions("W://root/wp_x", "M://root", []string{
			"- [Q1 RESOLVED wp_x.2026-04-28.001] some decision",
		}))

	if err := updateQStateDone(loadstarBase, "W://root/wp_x", "Q1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(addressToFilePath(loadstarBase, "W://root/wp_x"))
	if !strings.Contains(string(data), "[Q1 DONE wp_x.2026-04-28.001]") {
		t.Errorf("expected DONE in file, got:\n%s", string(data))
	}
}

func TestUpdateQStateDone_NotResolved_Fails(t *testing.T) {
	loadstarBase := setupCmdTest(t)
	writeElement(t, loadstarBase, "M://root",
		"<MAP>\n## [ADDRESS] M://root\n## [STATUS] S_STB\n\n### IDENTITY\n- SUMMARY: root\n\n### WAYPOINTS\n- W://root/wp_y\n\n### COMMENT\n(없음)\n</MAP>\n")
	writeElement(t, loadstarBase, "W://root/wp_y",
		wpWithQuestions("W://root/wp_y", "M://root", []string{
			"- [Q1] open question",
		}))

	if err := updateQStateDone(loadstarBase, "W://root/wp_y", "Q1"); err == nil {
		t.Error("expected error for non-RESOLVED question, got nil")
	}
}
