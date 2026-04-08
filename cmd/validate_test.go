package cmd

import (
	"testing"
)

func TestValidate_AllValid(t *testing.T) {
	loadstarBase := setupCmdTest(t)

	writeElement(t, loadstarBase, "M://root",
		"<MAP>\n## [ADDRESS] M://root\n## [STATUS] S_STB\n\n### WAYPOINTS\n- W://root/a\n\n### COMMENT\n</MAP>\n")
	writeElement(t, loadstarBase, "W://root/a",
		"<WAYPOINT>\n## [ADDRESS] W://root/a\n## [STATUS] S_STB\n\n### CONNECTIONS\n- PARENT: M://root\n- CHILDREN: []\n- REFERENCE: []\n\n### TODO\n(없음)\n</WAYPOINT>\n")

	// Should not panic — all references valid
	validateCmd.Run(validateCmd, nil)
}

func TestValidate_BrokenParent(t *testing.T) {
	loadstarBase := setupCmdTest(t)

	writeElement(t, loadstarBase, "W://root/orphan",
		"<WAYPOINT>\n## [ADDRESS] W://root/orphan\n## [STATUS] S_STB\n\n### CONNECTIONS\n- PARENT: M://root/nonexistent\n- CHILDREN: []\n- REFERENCE: []\n\n### TODO\n(없음)\n</WAYPOINT>\n")

	// Should report broken PARENT reference
	validateCmd.Run(validateCmd, nil)
}

func TestValidate_BrokenChild(t *testing.T) {
	loadstarBase := setupCmdTest(t)

	writeElement(t, loadstarBase, "W://root/parent",
		"<WAYPOINT>\n## [ADDRESS] W://root/parent\n## [STATUS] S_STB\n\n### CONNECTIONS\n- PARENT: M://root\n- CHILDREN: [W://root/parent/missing]\n- REFERENCE: []\n\n### TODO\n(없음)\n</WAYPOINT>\n")
	writeElement(t, loadstarBase, "M://root",
		"<MAP>\n## [ADDRESS] M://root\n## [STATUS] S_STB\n\n### WAYPOINTS\n- W://root/parent\n\n### COMMENT\n</MAP>\n")

	// Should report broken CHILDREN reference
	validateCmd.Run(validateCmd, nil)
}

func TestValidate_BrokenMapWaypoint(t *testing.T) {
	loadstarBase := setupCmdTest(t)

	writeElement(t, loadstarBase, "M://root",
		"<MAP>\n## [ADDRESS] M://root\n## [STATUS] S_STB\n\n### WAYPOINTS\n- W://root/exists\n- W://root/ghost\n\n### COMMENT\n</MAP>\n")
	writeElement(t, loadstarBase, "W://root/exists",
		"<WAYPOINT>\n## [ADDRESS] W://root/exists\n## [STATUS] S_STB\n\n### CONNECTIONS\n- PARENT: M://root\n- CHILDREN: []\n- REFERENCE: []\n\n### TODO\n(없음)\n</WAYPOINT>\n")

	// Should report W://root/ghost as broken WAYPOINTS reference
	validateCmd.Run(validateCmd, nil)
}

func TestValidate_EmptyFields(t *testing.T) {
	loadstarBase := setupCmdTest(t)

	writeElement(t, loadstarBase, "M://root",
		"<MAP>\n## [ADDRESS] M://root\n## [STATUS] S_STB\n\n### WAYPOINTS\n- W://root/a\n\n### COMMENT\n</MAP>\n")
	writeElement(t, loadstarBase, "W://root/a",
		"<WAYPOINT>\n## [ADDRESS] W://root/a\n## [STATUS] S_STB\n\n### CONNECTIONS\n- PARENT: M://root\n- CHILDREN: []\n- REFERENCE: []\n\n### TODO\n(없음)\n</WAYPOINT>\n")

	// Empty CHILDREN/REFERENCE should not report errors
	validateCmd.Run(validateCmd, nil)
}
