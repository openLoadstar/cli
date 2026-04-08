package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [FILTER]",
	Short: "List all waypoints with status",
	Long: `List all WayPoint elements with their address and status.
Optionally filter by keyword (case-insensitive match against address).

Examples:
  loadstar show                # list all waypoints
  loadstar show cli            # filter: addresses containing "cli"
  loadstar show test           # filter: addresses containing "test"`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filter := ""
		if len(args) > 0 {
			filter = args[0]
		}
		listWaypoints(filter)
	},
}

func listWaypoints(filter string) {
	loadstarBase := fs.AvcsPath("")
	wpDir := filepath.Join(loadstarBase, "WAYPOINT")

	files, err := os.ReadDir(wpDir)
	if err != nil {
		fmt.Println("no waypoints found")
		return
	}

	type wpInfo struct {
		Address string
		Status  string
	}

	var items []wpInfo
	lowerFilter := strings.ToLower(filter)

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}
		dotName := strings.TrimSuffix(f.Name(), ".md")
		path := strings.ReplaceAll(dotName, ".", "/")
		addr := "W://" + path

		if filter != "" && !strings.Contains(strings.ToLower(addr), lowerFilter) {
			continue
		}

		content, err := fs.Read(filepath.Join(wpDir, f.Name()))
		if err != nil {
			continue
		}
		status := extractField(content, `## \[STATUS\]\s+(\S+)`)

		items = append(items, wpInfo{Address: addr, Status: status})
	}

	if len(items) == 0 {
		fmt.Println("no waypoints found")
		return
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Address < items[j].Address
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ADDRESS\tSTATUS")
	fmt.Fprintln(w, "-------\t------")
	for _, item := range items {
		fmt.Fprintf(w, "%s\t%s\n", item.Address, item.Status)
	}
	w.Flush()

	fmt.Printf("\n%d waypoint(s)\n", len(items))
}

func extractField(content, pattern string) string {
	re := regexp.MustCompile(pattern)
	m := re.FindStringSubmatch(content)
	if len(m) < 2 {
		return "?"
	}
	return m[1]
}
