package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [ADDRESS]",
	Short: "Display element metadata in terminal",
	Long: `Print the content and status of a LOADSTAR element.
Use --depth to recursively show child elements as a tree.

Examples:
  loadstar show W://root/cli/cmd_log
  loadstar show M://root/cli
  loadstar show M://root --depth 2    # show root and 2 levels of children
  loadstar show M://root --depth 1    # show only direct children`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		depth, _ := cmd.Flags().GetInt("depth")
		addr, err := svc.ParseAddress(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		loadstarBase := fs.AvcsPath("")
		visited := make(map[string]bool)
		showElement(loadstarBase, addr.Type+"://"+addr.Path, depth, 0, visited)
	},
}

func showElement(loadstarBase, addrStr string, maxDepth, currentDepth int, visited map[string]bool) {
	if visited[addrStr] {
		fmt.Printf("%s[CIRCULAR: %s]\n", indent(currentDepth), addrStr)
		return
	}
	visited[addrStr] = true

	addr, err := svc.ParseAddress(addrStr)
	if err != nil {
		fmt.Printf("%s[INVALID: %s]\n", indent(currentDepth), addrStr)
		return
	}
	filePath := addr.ToFilePath(loadstarBase)
	if !fs.Exists(filePath) {
		fmt.Printf("%s[NOT FOUND: %s]\n", indent(currentDepth), addrStr)
		return
	}

	content, _ := fs.Read(filePath)

	// Extract STATUS
	status := extractField(content, `## \[STATUS\]\s+(\S+)`)
	fmt.Printf("%s%s  [%s]\n", indent(currentDepth), addrStr, status)

	if currentDepth >= maxDepth {
		return
	}

	// Extract children: WAYPOINTS (Map) or CHILDREN (WayPoint)
	children := extractChildren(content)
	for _, child := range children {
		child = strings.TrimSpace(child)
		if child == "" {
			continue
		}
		showElement(loadstarBase, child, maxDepth, currentDepth+1, visited)
	}
}

func indent(depth int) string {
	return strings.Repeat("  ", depth)
}

func extractField(content, pattern string) string {
	re := regexp.MustCompile(pattern)
	m := re.FindStringSubmatch(content)
	if len(m) < 2 {
		return "?"
	}
	return m[1]
}

// extractChildren extracts child addresses from either MAP (WAYPOINTS list)
// or WAYPOINT (CHILDREN: [...]) format.
func extractChildren(content string) []string {
	// Try CHILDREN: [...] format (WayPoint)
	childrenRe := regexp.MustCompile(`-\s*CHILDREN:\s*\[([^\]]*)\]`)
	m := childrenRe.FindStringSubmatch(content)
	if m != nil && strings.TrimSpace(m[1]) != "" {
		return strings.Split(m[1], ",")
	}

	// Try WAYPOINTS list format (Map) — lines starting with "- W://" or "- M://"
	var result []string
	lines := strings.Split(content, "\n")
	inWaypoints := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "### WAYPOINTS" {
			inWaypoints = true
			continue
		}
		if inWaypoints {
			if strings.HasPrefix(trimmed, "###") || trimmed == "" && len(result) > 0 {
				break
			}
			if strings.HasPrefix(trimmed, "- ") && strings.Contains(trimmed, "://") {
				addr := strings.TrimPrefix(trimmed, "- ")
				result = append(result, strings.TrimSpace(addr))
			}
			if trimmed == "(없음)" {
				break
			}
		}
	}
	return result
}

func init() {
	showCmd.Flags().Int("depth", 0, "Depth of child elements to display (default: 0)")
}
