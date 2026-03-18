package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// allowedCreateTypes lists types that users can directly create.
// H (History) and B (BlackBox) are CLI-managed and not directly creatable.
var allowedCreateTypes = map[string]bool{
	"M": true, "W": true, "L": true, "S": true,
}

var createCmd = &cobra.Command{
	Use:   "create [TYPE] [ID]",
	Short: "Create a new LOADSTAR element (M, W, L, S)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		elemType := strings.ToUpper(args[0])
		id := args[1]
		parent, _ := cmd.Flags().GetString("parent")

		// 1. TYPE validation
		if !allowedCreateTypes[elemType] {
			fmt.Fprintf(os.Stderr, "error: invalid type %q — allowed: M, W, L, S\n", elemType)
			os.Exit(1)
		}

		// 2. Parse and validate parent address
		if parent == "" {
			fmt.Fprintln(os.Stderr, "error: --parent flag is required")
			os.Exit(1)
		}
		parentAddr, err := svc.ParseAddress(parent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid parent address %q: %v\n", parent, err)
			os.Exit(1)
		}
		loadstarBase := fs.AvcsPath("")
		parentFile := parentAddr.ToFilePath(loadstarBase)
		if !fs.Exists(parentFile) {
			fmt.Fprintf(os.Stderr, "error: parent element not found: %s\n", parentFile)
			os.Exit(1)
		}

		// 3. Build new address and check for duplicates
		newPath := strings.TrimSuffix(parentAddr.Path, "/") + "/" + id
		newAddrStr := elemType + "://" + newPath
		newAddr, _ := svc.ParseAddress(newAddrStr)
		newFile := newAddr.ToFilePath(loadstarBase)
		if fs.Exists(newFile) {
			fmt.Fprintf(os.Stderr, "error: element already exists: %s\n", newAddrStr)
			os.Exit(1)
		}

		// 4. Generate MD template and write file
		content := buildTemplate(elemType, newAddrStr, parent)
		if err := fs.Write(newFile, content); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to write element file: %v\n", err)
			os.Exit(1)
		}

		// 5. Register new address in parent's CONTAINS.ITEMS
		if err := appendToContains(parentFile, newAddrStr); err != nil {
			fmt.Fprintf(os.Stderr, "warning: element created but failed to update parent CONTAINS: %v\n", err)
		}

		fmt.Printf("created: %s\n", newAddrStr)
	},
}

var editCmd = &cobra.Command{
	Use:   "edit [ADDRESS]",
	Short: "Edit an existing element by address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addr, err := svc.ParseAddress(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		loadstarBase := fs.AvcsPath("")
		filePath := addr.ToFilePath(loadstarBase)
		if !fs.Exists(filePath) {
			fmt.Fprintf(os.Stderr, "error: element not found: %s\n", args[0])
			os.Exit(1)
		}

		// Shadow History snapshot before editing
		ts := time.Now().Format("20060102T150405")
		dotName := strings.ReplaceAll(addr.Path, "/", ".")
		histPath := filepath.Join(loadstarBase, "HISTORY", dotName+"_"+ts+".md")
		if err := fs.CopyFile(filePath, histPath); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not create history snapshot: %v\n", err)
		}

		// Get mtime before edit
		statBefore, _ := os.Stat(filePath)

		// Launch editor
		editor := resolveEditor()
		editorArgs := strings.Fields(editor)
		editorArgs = append(editorArgs, filePath)
		c := exec.Command(editorArgs[0], editorArgs[1:]...)
		c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
		if err := c.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: editor exited with error: %v\n", err)
			os.Exit(1)
		}

		// Detect changes
		statAfter, _ := os.Stat(filePath)
		if statBefore != nil && statAfter != nil && !statAfter.ModTime().After(statBefore.ModTime()) {
			fmt.Println("no changes detected")
			return
		}
		fmt.Printf("saved: %s\n", args[0])
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [ADDRESS]",
	Short: "Delete an element by address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		addr, err := svc.ParseAddress(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		loadstarBase := fs.AvcsPath("")
		filePath := addr.ToFilePath(loadstarBase)
		if !fs.Exists(filePath) {
			fmt.Fprintf(os.Stderr, "error: element not found: %s\n", args[0])
			os.Exit(1)
		}

		// Confirmation prompt
		if !force {
			fmt.Printf("delete %s? [y/N] ", args[0])
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" && answer != "yes" {
				fmt.Println("aborted")
				return
			}
		}

		// History First: backup before delete
		ts := time.Now().Format("20060102T150405")
		dotName := strings.ReplaceAll(addr.Path, "/", ".")
		histPath := filepath.Join(loadstarBase, "HISTORY", dotName+"_"+ts+"_deleted.md")
		if err := fs.CopyFile(filePath, histPath); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not create history backup: %v\n", err)
		}

		// Parse parent address from LINEAGE and remove from CONTAINS
		content, err := fs.Read(filePath)
		if err == nil {
			if parentAddr := parseLineageParent(content); parentAddr != "" {
				pa, perr := svc.ParseAddress(parentAddr)
				if perr == nil {
					parentFile := pa.ToFilePath(loadstarBase)
					if err2 := removeFromContains(parentFile, args[0]); err2 != nil {
						fmt.Fprintf(os.Stderr, "warning: could not remove from parent CONTAINS: %v\n", err2)
					}
				} else {
					fmt.Fprintf(os.Stderr, "warning: could not parse parent address %q\n", parentAddr)
				}
			} else {
				fmt.Fprintln(os.Stderr, "warning: LINEAGE.PARENT not found — parent CONTAINS not updated")
			}
		}

		if err := os.Remove(filePath); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to delete file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("deleted: %s\n", args[0])
	},
}

func init() {
	createCmd.Flags().String("parent", "", "Parent element address (e.g. M://root/cli)")
	deleteCmd.Flags().Bool("force", false, "Skip confirmation prompt")
}

// buildTemplate generates an MD template string for the given element type.
func buildTemplate(elemType, address, parent string) string {
	now := time.Now().Format("2006-01-02")
	switch elemType {
	case "M":
		return fmt.Sprintf("<MAP>\n## [ADDRESS] %s\n## [STATUS] S_IDL\n\n### 1. IDENTITY\n- SUMMARY:\n- METADATA: [Ver: 1.0, Created: %s]\n- SYNCED_AT: %s\n\n### 2. CONTAINS\n- ITEMS: []\n- PAYLOAD:\n\n### 3. CONNECTIONS\n- LINEAGE: [PARENT: %s, CHILDREN: []]\n- LINKS: []\n\n### 4. RESOURCES\n- SAVEPOINTS: []\n\n### 5. TODO\n- REQUESTER: %s\n- RESPONSE_STATUS: PENDING\n- TECH_SPEC:\n- EXECUTION_HISTORY: []\n</MAP>\n", address, now, now, parent, parent)
	case "W":
		return fmt.Sprintf("<WAYPOINT>\n## [ADDRESS] %s\n## [STATUS] S_IDL\n\n### 1. IDENTITY\n- SUMMARY:\n- METADATA: [Ver: 1.0, Created: %s, Priority: P2]\n- SYNCED_AT: %s\n\n### 2. CONTAINS\n- ITEMS: []\n- PAYLOAD:\n\n### 3. CONNECTIONS\n- LINEAGE: [PARENT: %s, CHILDREN: []]\n- LINKS: []\n\n### 4. RESOURCES\n- SAVEPOINTS: []\n\n### 5. TODO\n- REQUESTER: %s\n- EXECUTOR: %s\n- RESPONSE_STATUS: PENDING\n- TECH_SPEC:\n- EXECUTION_HISTORY: []\n</WAYPOINT>\n", address, now, now, parent, parent, address)
	case "L":
		return fmt.Sprintf("<LINK>\n## [ADDRESS] %s\n## [STATUS] S_IDL\n\n### 1. IDENTITY\n- SUMMARY:\n- METADATA: [Created: %s]\n- SYNCED_AT: %s\n\n### 2. CONTAINS\n- ITEMS: []\n- PAYLOAD:\n  - SOURCE:\n  - TARGET:\n  - TYPE:\n\n### 3. CONNECTIONS\n- LINEAGE: [PARENT: %s, CHILDREN: []]\n- LINKS: []\n\n### 4. RESOURCES\n- SAVEPOINTS: []\n</LINK>\n", address, now, now, parent)
	case "S":
		return fmt.Sprintf("<SAVEPOINT>\n## [ADDRESS] %s\n## [STATUS] S_IDL\n\n### 1. IDENTITY\n- SUMMARY:\n- METADATA: [Created: %s]\n- SYNCED_AT: %s\n\n### 2. CONTAINS\n- ITEMS: []\n- PAYLOAD:\n\n### 3. CONNECTIONS\n- LINEAGE: [PARENT: %s, CHILDREN: []]\n- LINKS: []\n\n### 4. RESOURCES\n- SAVEPOINTS: []\n</SAVEPOINT>\n", address, now, now, parent)
	default:
		return ""
	}
}

// appendToContains adds newAddr to the CONTAINS.ITEMS list using line scanning.
func appendToContains(filePath, newAddr string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	itemsRe := regexp.MustCompile(`^(\s*-\s*ITEMS:\s*\[)(.*?)(\].*)$`)
	for i, line := range lines {
		m := itemsRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		inner := strings.TrimSpace(m[2])
		if inner == "" {
			lines[i] = m[1] + newAddr + m[3]
		} else {
			lines[i] = m[1] + inner + ", " + newAddr + m[3]
		}
		return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
	}
	return fmt.Errorf("CONTAINS.ITEMS line not found in %s", filePath)
}

// removeFromContains removes targetAddr from the CONTAINS.ITEMS list using line scanning.
func removeFromContains(filePath, targetAddr string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	itemsRe := regexp.MustCompile(`^(\s*-\s*ITEMS:\s*\[)(.*?)(\].*)$`)
	for i, line := range lines {
		m := itemsRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		parts := strings.Split(m[2], ",")
		var kept []string
		for _, p := range parts {
			if strings.TrimSpace(p) != targetAddr {
				kept = append(kept, p)
			}
		}
		lines[i] = m[1] + strings.Join(kept, ",") + m[3]
		return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
	}
	return fmt.Errorf("CONTAINS.ITEMS line not found in %s", filePath)
}

// parseLineageParent extracts the PARENT address from a LINEAGE line.
func parseLineageParent(content string) string {
	re := regexp.MustCompile(`LINEAGE:\s*\[PARENT:\s*([^,\]]+)`)
	m := re.FindStringSubmatch(content)
	if len(m) < 2 {
		return ""
	}
	return strings.TrimSpace(m[1])
}

// resolveEditor returns the editor command to use, respecting env vars.
func resolveEditor() string {
	if e := os.Getenv("LOADSTAR_EDITOR"); e != "" {
		return e
	}
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	if runtime.GOOS == "windows" {
		return "notepad"
	}
	return "vi"
}
