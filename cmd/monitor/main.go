package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/bono/loadstar/internal/git"
	"github.com/bono/loadstar/internal/storage"
)

const (
	defaultInterval = 5 * time.Minute
	flagFileName    = "checkpoint_needed.flag"
)

// monitoredDirs are the .loadstar/ subdirectories to watch for changes.
var monitoredDirs = []string{
	"WAYPOINT", "BLACKBOX", "MAP",
}

func main() {
	// Determine project root: first arg or cwd
	var startDir string
	if len(os.Args) > 1 {
		startDir = os.Args[1]
	} else {
		var err error
		startDir, err = os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: cannot determine working directory:", err)
			os.Exit(1)
		}
	}

	root, err := storage.FindRoot(startDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: .loadstar/ not found in %s or any parent\n", startDir)
		os.Exit(1)
	}

	loadstarBase := filepath.Join(root, ".loadstar")
	monitorDir := filepath.Join(loadstarBase, ".clionly", "MONITOR")
	_ = os.MkdirAll(monitorDir, 0755)

	flagPath := filepath.Join(monitorDir, flagFileName)
	gc := git.NewClient(root)

	fmt.Printf("loadstar_monitor started (interval: %s)\n", defaultInterval)
	fmt.Printf("project root: %s\n", root)

	for {
		checkAndFlag(gc, flagPath)
		time.Sleep(defaultInterval)
	}
}

// checkAndFlag checks for uncommitted .loadstar/ changes and creates a flag file if needed.
func checkAndFlag(gc *git.Client, flagPath string) {
	changed, err := gc.ChangedLoadstarFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[%s] warning: git status check failed: %v\n",
			time.Now().Format("15:04:05"), err)
		return
	}

	// Filter to monitored element directories only
	var elements []string
	for _, f := range changed {
		parts := strings.SplitN(f, "/", 3)
		if len(parts) < 3 {
			continue
		}
		dir := parts[1]
		for _, m := range monitoredDirs {
			if dir == m {
				elements = append(elements, f)
				break
			}
		}
	}

	if len(elements) == 0 {
		fmt.Printf("[%s] no uncommitted changes\n", time.Now().Format("15:04:05"))
		return
	}

	// Skip if flag already exists
	if _, err := os.Stat(flagPath); err == nil {
		fmt.Printf("[%s] checkpoint_needed.flag already exists (%d changed files)\n",
			time.Now().Format("15:04:05"), len(elements))
		return
	}

	// Create flag file
	sort.Strings(elements)
	now := time.Now()
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## [DETECTED_AT] %s\n", now.Format("2006-01-02T15:04:05")))
	sb.WriteString("## [CHANGED_FILES]\n")
	for _, e := range elements {
		sb.WriteString("- ")
		sb.WriteString(e)
		sb.WriteString("\n")
	}

	if err := os.WriteFile(flagPath, []byte(sb.String()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "[%s] error: failed to write flag file: %v\n",
			now.Format("15:04:05"), err)
		return
	}

	fmt.Printf("[%s] checkpoint needed: %d changed files → flag created\n",
		now.Format("15:04:05"), len(elements))
}
