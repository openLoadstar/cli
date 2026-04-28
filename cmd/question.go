package cmd

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// qRe matches OPEN_QUESTIONS lines:
//   - [Q1], [Q1 DEFERRED], [Q1 RESOLVED some-id], [Q1 DONE some-id]
var qRe = regexp.MustCompile(`^\s*-\s*\[Q(\d+)(?:\s+(DEFERRED|RESOLVED|DONE)(?:\s+([\w.-]+))?)?\]\s*(.*)$`)

type qEntry struct {
	Address  string
	QID      string
	State    string
	Ref      string // decision file ref (for RESOLVED/DONE)
	Question string
}

func scanQuestions(loadstarBase string) []qEntry {
	wpAddrs := collectAllWaypoints(loadstarBase)
	var entries []qEntry

	for _, addr := range wpAddrs {
		wpFile := addressToFilePath(loadstarBase, addr)
		data, err := os.ReadFile(wpFile)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			m := qRe.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			stateRaw := strings.TrimSpace(m[2])
			state := "OPEN"
			switch {
			case stateRaw == "DEFERRED":
				state = "DEFERRED"
			case stateRaw == "RESOLVED":
				state = "RESOLVED"
			case stateRaw == "DONE":
				state = "DONE"
			}
			entries = append(entries, qEntry{
				Address:  addr,
				QID:      "Q" + m[1],
				State:    state,
				Ref:      strings.TrimSpace(m[3]),
				Question: strings.TrimSpace(m[4]),
			})
		}
	}
	return entries
}

// updateQState rewrites [Qn <oldState> ref] → [Qn <newState> ref] in the WP file.
func updateQState(loadstarBase, wpAddr, qid, newState string) error {
	wpFile := addressToFilePath(loadstarBase, wpAddr)
	data, err := os.ReadFile(wpFile)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", wpFile, err)
	}

	target := regexp.MustCompile(`^(\s*-\s*\[` + regexp.QuoteMeta(qid) + `)(?:\s+(?:DEFERRED|RESOLVED|DONE)(?:\s+[\w.-]+)?)?(\].*)$`)

	lines := strings.Split(string(data), "\n")
	matched := false
	for i, line := range lines {
		if m := target.FindStringSubmatch(line); m != nil {
			lines[i] = m[1] + " " + newState + m[2]
			matched = true
			break
		}
	}
	if !matched {
		return fmt.Errorf("question %s not found in %s", qid, wpAddr)
	}
	return os.WriteFile(wpFile, []byte(strings.Join(lines, "\n")), 0644)
}

// updateQStateDone transitions RESOLVED → DONE, preserving the ref.
func updateQStateDone(loadstarBase, wpAddr, qid string) error {
	wpFile := addressToFilePath(loadstarBase, wpAddr)
	data, err := os.ReadFile(wpFile)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", wpFile, err)
	}

	target := regexp.MustCompile(`^(\s*-\s*\[` + regexp.QuoteMeta(qid) + `)\s+RESOLVED(\s+[\w.-]+)?(\].*)$`)

	lines := strings.Split(string(data), "\n")
	matched := false
	for i, line := range lines {
		if m := target.FindStringSubmatch(line); m != nil {
			lines[i] = m[1] + " DONE" + m[2] + m[3]
			matched = true
			break
		}
	}
	if !matched {
		return fmt.Errorf("%s is not in RESOLVED state in %s", qid, wpAddr)
	}
	return os.WriteFile(wpFile, []byte(strings.Join(lines, "\n")), 0644)
}

var withResolved bool

var questionCmd = &cobra.Command{
	Use:   "question [FILTER]",
	Short: "List OPEN/DEFERRED questions from all WayPoints",
	Long: `Scan all WayPoint files and list OPEN and DEFERRED OPEN_QUESTIONS.
Use --with-resolved to include RESOLVED and DONE items as well.

Optional FILTER matches address or question text (case-insensitive substring).

Examples:
  loadstar question
  loadstar question --with-resolved
  loadstar question M://root/maintenance`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		loadstarBase := fs.AvcsPath("")
		entries := scanQuestions(loadstarBase)

		filter := ""
		if len(args) == 1 {
			filter = strings.ToLower(args[0])
		}

		var visible []qEntry
		for _, e := range entries {
			if !withResolved && (e.State == "RESOLVED" || e.State == "DONE") {
				continue
			}
			if filter != "" &&
				!strings.Contains(strings.ToLower(e.Address), filter) &&
				!strings.Contains(strings.ToLower(e.Question), filter) {
				continue
			}
			visible = append(visible, e)
		}

		if len(visible) == 0 {
			fmt.Println("no questions found")
			return
		}

		stateOrder := map[string]int{"OPEN": 0, "DEFERRED": 1, "RESOLVED": 2, "DONE": 3}
		sort.Slice(visible, func(i, j int) bool {
			oi, oj := stateOrder[visible[i].State], stateOrder[visible[j].State]
			if oi != oj {
				return oi < oj
			}
			if visible[i].Address != visible[j].Address {
				return visible[i].Address < visible[j].Address
			}
			return visible[i].QID < visible[j].QID
		})

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ADDRESS\tQID\tSTATE\tREF\tQUESTION")
		fmt.Fprintln(w, "-------\t---\t-----\t---\t--------")
		for _, e := range visible {
			q := e.Question
			if len(q) > 70 {
				q = q[:70] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", e.Address, e.QID, e.State, e.Ref, q)
		}
		w.Flush()
		fmt.Printf("\n%d question(s)\n", len(visible))
	},
}

var questionStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show OPEN_QUESTIONS statistics",
	Run: func(cmd *cobra.Command, args []string) {
		loadstarBase := fs.AvcsPath("")
		entries := scanQuestions(loadstarBase)

		counts := map[string]int{"OPEN": 0, "DEFERRED": 0, "RESOLVED": 0, "DONE": 0}
		for _, e := range entries {
			counts[e.State]++
		}

		fmt.Printf("OPEN:     %d\n", counts["OPEN"])
		fmt.Printf("DEFERRED: %d\n", counts["DEFERRED"])
		fmt.Printf("RESOLVED: %d\n", counts["RESOLVED"])
		fmt.Printf("DONE:     %d\n", counts["DONE"])
		fmt.Printf("TOTAL:    %d\n", counts["OPEN"]+counts["DEFERRED"]+counts["RESOLVED"]+counts["DONE"])
	},
}

// questionDoneCmd: RESOLVED → DONE (AI applied decision to code)
var questionDoneCmd = &cobra.Command{
	Use:   "done <wpAddr> <qid>",
	Short: "Mark a RESOLVED question as DONE (code applied)",
	Long: `Transition [Qn RESOLVED <ref>] → [Qn DONE <ref>].
Use this after you have applied the decision to the codebase.

Example:
  loadstar question done W://root/maintenance/decisions_ui Q1`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		wpAddr := args[0]
		qid := args[1]
		loadstarBase := fs.AvcsPath("")

		if err := updateQStateDone(loadstarBase, wpAddr, qid); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ %s %s → DONE\n", wpAddr, qid)
	},
}

// questionCloseCmd: close a question directly as DONE (no decision file)
var questionCloseCmd = &cobra.Command{
	Use:   "close <wpAddr> <qid> [reason]",
	Short: "Close a question as DONE without a decision file",
	Long: `Mark any question (OPEN/DEFERRED) directly as DONE with an optional reason.
Use this for "won't fix", "no longer relevant", etc.

Example:
  loadstar question close W://root/maintenance/decisions_ui Q2 "no longer relevant"`,
	Args: cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		wpAddr := args[0]
		qid := args[1]
		reason := ""
		if len(args) == 3 {
			reason = args[2]
		}
		loadstarBase := fs.AvcsPath("")

		newState := "DONE"
		if reason != "" {
			newState = "DONE — " + reason
		}
		if err := updateQState(loadstarBase, wpAddr, qid, newState); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ %s %s closed\n", wpAddr, qid)
	},
}

func init() {
	questionCmd.Flags().BoolVar(&withResolved, "with-resolved", false, "include RESOLVED and DONE items")
	questionCmd.AddCommand(questionStatsCmd)
	questionCmd.AddCommand(questionDoneCmd)
	questionCmd.AddCommand(questionCloseCmd)
}
