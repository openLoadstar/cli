package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"path/filepath"

	"github.com/spf13/cobra"
)

var allowedLinkTypes = map[string]bool{
	"L_REF": true, "L_SEQ": true, "L_TST": true,
}

var linkCmd = &cobra.Command{
	Use:   "link [SOURCE] [TARGET]",
	Short: "Create a logical link between two elements",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		linkType, _ := cmd.Flags().GetString("type")
		linkType = strings.ToUpper(linkType)

		// Validate link type
		if !allowedLinkTypes[linkType] {
			fmt.Fprintf(os.Stderr, "error: invalid link type %q — allowed: L_REF, L_SEQ, L_TST\n", linkType)
			os.Exit(1)
		}

		srcAddr, err := svc.ParseAddress(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid source address: %v\n", err)
			os.Exit(1)
		}
		dstAddr, err := svc.ParseAddress(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid target address: %v\n", err)
			os.Exit(1)
		}

		loadstarBase := fs.AvcsPath("")
		srcFile := srcAddr.ToFilePath(loadstarBase)
		dstFile := dstAddr.ToFilePath(loadstarBase)

		if !fs.Exists(srcFile) {
			fmt.Fprintf(os.Stderr, "error: source element not found: %s\n", args[0])
			os.Exit(1)
		}
		if !fs.Exists(dstFile) {
			fmt.Fprintf(os.Stderr, "error: target element not found: %s\n", args[1])
			os.Exit(1)
		}

		// Build Link ID and address
		linkID := srcAddr.ID + "_to_" + dstAddr.ID
		// Use source path prefix to derive link path
		pathPrefix := strings.Join(strings.Split(srcAddr.Path, "/")[:len(strings.Split(srcAddr.Path, "/"))-1], "/")
		linkPath := pathPrefix + "/" + linkID
		linkAddrStr := "L://" + linkPath
		now := time.Now().Format("2006-01-02")

		// Check for duplicate
		linkFile := filepath.Join(loadstarBase, "LINK", strings.ReplaceAll(linkPath, "/", ".")+".md")
		if fs.Exists(linkFile) {
			fmt.Fprintf(os.Stderr, "error: link already exists: %s\n", linkAddrStr)
			os.Exit(1)
		}

		// Create Link md file
		linkContent := fmt.Sprintf("<LINK>\n## [ADDRESS] %s\n## [STATUS] S_STB\n\n### 1. IDENTITY\n- SUMMARY: %s → %s (%s)\n- METADATA: [Created: %s]\n- SYNCED_AT: %s\n\n### 2. CONTAINS\n- ITEMS: []\n- PAYLOAD:\n  - SOURCE: %s\n  - TARGET: %s\n  - TYPE: %s\n\n### 3. CONNECTIONS\n- LINEAGE: [PARENT: M://%s, CHILDREN: []]\n- LINKS: []\n\n### 4. RESOURCES\n- SAVEPOINTS: []\n</LINK>\n",
			linkAddrStr, args[0], args[1], linkType, now, now,
			args[0], args[1], linkType,
			pathPrefix,
		)
		if err := fs.Write(linkFile, linkContent); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to create link file: %v\n", err)
			os.Exit(1)
		}

		// Register in SOURCE CONNECTIONS.LINKS
		if err := appendToLinks(srcFile, linkAddrStr, linkType); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not update source CONNECTIONS.LINKS: %v\n", err)
		}
		// Register in TARGET CONNECTIONS.LINKS (reverse reference)
		if err := appendToLinks(dstFile, linkAddrStr, linkType); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not update target CONNECTIONS.LINKS: %v\n", err)
		}

		fmt.Printf("linked: %s -> %s [%s]\n", args[0], args[1], linkType)
	},
}

var showCmd = &cobra.Command{
	Use:   "show [ADDRESS]",
	Short: "Display element metadata in terminal",
	Args:  cobra.ExactArgs(1),
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

	// Extract CONTAINS.ITEMS addresses
	children := extractContainsItems(content)
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

func extractContainsItems(content string) []string {
	re := regexp.MustCompile(`-\s*ITEMS:\s*\[([^\]]*)\]`)
	m := re.FindStringSubmatch(content)
	if len(m) < 2 || strings.TrimSpace(m[1]) == "" {
		return nil
	}
	parts := strings.Split(m[1], ",")
	return parts
}

// appendToLinks adds a link entry to the CONNECTIONS.LINKS list.
func appendToLinks(filePath, linkAddr, linkType string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	linksRe := regexp.MustCompile(`^(\s*-\s*LINKS:\s*\[)(.*?)(\].*)$`)
	entry := fmt.Sprintf("%s | TYPE: %s", linkAddr, linkType)
	for i, line := range lines {
		m := linksRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		inner := strings.TrimSpace(m[2])
		if inner == "" {
			lines[i] = m[1] + entry + m[3]
		} else {
			lines[i] = m[1] + inner + ",\n    " + entry + m[3]
		}
		return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
	}
	return fmt.Errorf("CONNECTIONS.LINKS line not found in %s", filePath)
}

func init() {
	linkCmd.Flags().String("type", "", "Link type: L_REF | L_SEQ | L_TST")
	linkCmd.MarkFlagRequired("type")
	showCmd.Flags().Int("depth", 0, "Depth of child elements to display (default: 0)")
}
