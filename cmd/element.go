package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// allowedCreateTypes lists types that users can directly create.
// B (BlackBox) is CLI-managed and not directly creatable.
var allowedCreateTypes = map[string]bool{
	"M": true, "W": true,
}

var createCmd = &cobra.Command{
	Use:   "create [TYPE] [ID]",
	Short: "Create a new LOADSTAR element (M, W)",
	Long: `Create a new LOADSTAR element under the specified parent.

Types:
  M   Map       — index that groups WayPoints
  W   WayPoint  — unit of work / intent

  Note: B (BlackBox) is auto-managed by the CLI.

Examples:
  loadstar create M cli --parent M://root
  loadstar create W cmd_log --parent M://root/cli`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		elemType := strings.ToUpper(args[0])
		id := args[1]
		parent, _ := cmd.Flags().GetString("parent")

		// 1. TYPE validation
		if !allowedCreateTypes[elemType] {
			fmt.Fprintf(os.Stderr, "error: invalid type %q — allowed: M, W\n", elemType)
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

		// 5. Register new address in parent's child list
		parentContent, _ := os.ReadFile(parentFile)
		parentStr := string(parentContent)
		if strings.Contains(parentStr, "<MAP>") {
			if err := appendToWaypoints(parentFile, newAddrStr); err != nil {
				fmt.Fprintf(os.Stderr, "warning: element created but failed to update parent WAYPOINTS: %v\n", err)
			}
		} else {
			if err := appendToChildren(parentFile, newAddrStr); err != nil {
				fmt.Fprintf(os.Stderr, "warning: element created but failed to update parent CHILDREN: %v\n", err)
			}
		}

		// 6. Auto-create BlackBox when creating a WayPoint
		if elemType == "W" {
			bbAddrStr := "B://" + strings.TrimPrefix(newAddrStr, "W://")
			bbAddr, _ := svc.ParseAddress(bbAddrStr)
			bbFile := bbAddr.ToFilePath(loadstarBase)
			if !fs.Exists(bbFile) {
				bbContent := buildTemplate("B", bbAddrStr, "")
				if err := fs.Write(bbFile, bbContent); err != nil {
					fmt.Fprintf(os.Stderr, "warning: WayPoint created but failed to create BlackBox: %v\n", err)
				} else {
					fmt.Printf("created: %s\n", bbAddrStr)
				}
			}
		}

		fmt.Printf("created: %s\n", newAddrStr)
	},
}

var editCmd = &cobra.Command{
	Use:   "edit [ADDRESS]",
	Short: "Edit an existing element by address",
	Long: `Open an element's markdown file in the system editor.

The editor is resolved in this order: $LOADSTAR_EDITOR, $EDITOR, notepad (Windows) / vi (Unix).

Examples:
  loadstar edit W://root/cli/cmd_log
  loadstar edit M://root/cli
  loadstar edit B://root/cli/cmd_create`,
	Args: cobra.ExactArgs(1),
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
	Long: `Delete a LOADSTAR element and remove it from its parent's child list.

You will be prompted for confirmation unless --force is set.

Note: child elements referencing the deleted element
are NOT automatically removed — verify manually after deletion.

Examples:
  loadstar delete W://root/cli/cmd_log
  loadstar delete M://root/old_feature
  loadstar delete W://root/cli/cmd_log --force`,
	Args: cobra.ExactArgs(1),
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

		// Parse parent address from CONNECTIONS and remove from parent's child list
		content, err := fs.Read(filePath)
		if err == nil {
			if parentAddr := parseParent(content); parentAddr != "" {
				pa, perr := svc.ParseAddress(parentAddr)
				if perr == nil {
					parentFile := pa.ToFilePath(loadstarBase)
					parentContent, _ := os.ReadFile(parentFile)
					parentStr := string(parentContent)
					if strings.Contains(parentStr, "<MAP>") {
						if err2 := removeFromWaypoints(parentFile, args[0]); err2 != nil {
							fmt.Fprintf(os.Stderr, "warning: could not remove from parent WAYPOINTS: %v\n", err2)
						}
					} else {
						if err2 := removeFromChildren(parentFile, args[0]); err2 != nil {
							fmt.Fprintf(os.Stderr, "warning: could not remove from parent CHILDREN: %v\n", err2)
						}
					}
				}
			} else {
				fmt.Fprintln(os.Stderr, "warning: CONNECTIONS.PARENT not found — parent not updated")
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
		return fmt.Sprintf("<MAP>\n## [ADDRESS] %s\n## [STATUS] S_IDL\n\n### IDENTITY\n- SUMMARY:\n\n### WAYPOINTS\n(없음)\n\n### COMMENT\n(없음)\n</MAP>\n", address)
	case "W":
		return fmt.Sprintf("<WAYPOINT>\n## [ADDRESS] %s\n## [STATUS] S_IDL\n\n### IDENTITY\n- SUMMARY:\n- METADATA: [Ver: 1.0, Created: %s]\n\n### CONNECTIONS\n- PARENT: %s\n- CHILDREN: []\n- REFERENCE: []\n- BLACKBOX: B://%s\n\n### TODO\n(없음)\n\n### ISSUE\n(없음)\n\n### COMMENT\n(없음)\n</WAYPOINT>\n", address, now, parent, strings.TrimPrefix(address, "W://"))
	case "B":
		return fmt.Sprintf("<BLACKBOX>\n## [ADDRESS] %s\n## [STATUS] S_IDL\n\n### DESCRIPTION\n- SUMMARY:\n- LINKED_WP: W://%s\n\n### CODE_MAP\n(미작성)\n\n### TODO\n(없음)\n\n### ISSUE\n(없음)\n\n### COMMENT\n(없음)\n</BLACKBOX>\n", address, strings.TrimPrefix(address, "B://"))
	default:
		return ""
	}
}

// appendToWaypoints adds newAddr to a MAP's WAYPOINTS section.
func appendToWaypoints(filePath, newAddr string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "### WAYPOINTS" {
			// Insert after the header. If next line is "(없음)", replace it.
			ins := i + 1
			if ins < len(lines) && strings.TrimSpace(lines[ins]) == "(없음)" {
				lines[ins] = "- " + newAddr
			} else {
				lines = append(lines[:ins], append([]string{"- " + newAddr}, lines[ins:]...)...)
			}
			return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
		}
	}
	return fmt.Errorf("WAYPOINTS section not found in %s", filePath)
}

// removeFromWaypoints removes targetAddr from a MAP's WAYPOINTS section.
func removeFromWaypoints(filePath, targetAddr string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "- "+targetAddr {
			lines = append(lines[:i], lines[i+1:]...)
			return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
		}
	}
	return nil // Not found — acceptable
}

// appendToChildren adds newAddr to a WAYPOINT's CHILDREN list.
func appendToChildren(filePath, newAddr string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	childrenRe := regexp.MustCompile(`^(\s*-\s*CHILDREN:\s*\[)(.*?)(\].*)$`)
	for i, line := range lines {
		m := childrenRe.FindStringSubmatch(line)
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
	return fmt.Errorf("CHILDREN line not found in %s", filePath)
}

// removeFromChildren removes targetAddr from a WAYPOINT's CHILDREN list.
func removeFromChildren(filePath, targetAddr string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	childrenRe := regexp.MustCompile(`^(\s*-\s*CHILDREN:\s*\[)(.*?)(\].*)$`)
	for i, line := range lines {
		m := childrenRe.FindStringSubmatch(line)
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
	return fmt.Errorf("CHILDREN line not found in %s", filePath)
}

// parseParent extracts the PARENT address from CONNECTIONS section.
func parseParent(content string) string {
	re := regexp.MustCompile(`(?m)^-\s*PARENT:\s*(.+)$`)
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
