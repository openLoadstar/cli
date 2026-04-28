package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var showRecent bool

var showCmd = &cobra.Command{
	Use:   "show [FILTER]",
	Short: "List all waypoints with status",
	Long: `List all WayPoint elements with their address, status, and last-modified time.
Optionally filter by keyword (case-insensitive match against address).

Examples:
  loadstar show                # list all waypoints (sorted by address)
  loadstar show cli            # filter: addresses containing "cli"
  loadstar show --recent       # sort by most recently modified first
  loadstar show cli --recent   # combine filter + recent sort`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filter := ""
		if len(args) > 0 {
			filter = args[0]
		}
		listWaypoints(filter, showRecent)
	},
}

func init() {
	showCmd.Flags().BoolVar(&showRecent, "recent", false, "sort by most recently modified file first")
}

func listWaypoints(filter string, recent bool) {
	loadstarBase := fs.AvcsPath("")
	wpDir := filepath.Join(loadstarBase, "WAYPOINT")

	files, err := os.ReadDir(wpDir)
	if err != nil {
		fmt.Println("no waypoints found")
		return
	}

	type wpInfo struct {
		Address      string
		Status       string
		LastModified time.Time
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

		fullPath := filepath.Join(wpDir, f.Name())
		content, err := fs.Read(fullPath)
		if err != nil {
			continue
		}
		status := extractField(content, `## \[STATUS\]\s+(\S+)`)

		var mtime time.Time
		if info, err := f.Info(); err == nil {
			mtime = info.ModTime()
		}

		items = append(items, wpInfo{Address: addr, Status: status, LastModified: mtime})
	}

	if len(items) == 0 {
		fmt.Println("no waypoints found")
		return
	}

	if recent {
		sort.Slice(items, func(i, j int) bool {
			return items[i].LastModified.After(items[j].LastModified)
		})
	} else {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Address < items[j].Address
		})
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ADDRESS\tSTATUS\tLAST_MODIFIED")
	fmt.Fprintln(w, "-------\t------\t-------------")
	for _, item := range items {
		fmt.Fprintf(w, "%s\t%s\t%s\n", item.Address, item.Status, formatMTime(item.LastModified))
	}
	w.Flush()

	fmt.Printf("\n%d waypoint(s)\n", len(items))
}

func formatMTime(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	return t.Format("2006-01-02 15:04")
}

func extractField(content, pattern string) string {
	re := regexp.MustCompile(pattern)
	m := re.FindStringSubmatch(content)
	if len(m) < 2 {
		return "?"
	}
	return m[1]
}
