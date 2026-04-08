package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const AvcsDir = ".loadstar"

var loadstarDirs = []string{
	"MAP", "WAYPOINT", "COMMON",
	".clionly/LOG", ".clionly/MONITOR", ".clionly/TODO",
}

// FindRoot walks up from the given directory looking for an existing .loadstar folder.
// If none is found, it auto-initialises .loadstar in startDir and returns startDir.
func FindRoot(startDir string) (string, error) {
	dir := startDir
	for {
		if Exists(filepath.Join(dir, AvcsDir)) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root — auto-initialise in startDir.
			fmt.Fprintf(os.Stderr, "warning: no .loadstar/ found in %s or any parent directory.\n", startDir)
			fmt.Fprintf(os.Stderr, "         Run this command from your project directory (where .loadstar/ exists).\n")
			fmt.Fprintf(os.Stderr, "         Initialising a new .loadstar/ in current directory: %s\n", startDir)
			if err := Init(startDir); err != nil {
				return "", err
			}
			return startDir, nil
		}
		dir = parent
	}
}

// Init creates the .loadstar directory structure under projectRoot,
// and writes an initial M://root Map file.
func Init(projectRoot string) error {
	for _, d := range loadstarDirs {
		path := filepath.Join(projectRoot, AvcsDir, d)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	rootMap := filepath.Join(projectRoot, AvcsDir, "MAP", "root.md")
	if !Exists(rootMap) {
		content := "<MAP>\n## [ADDRESS] M://root\n## [STATUS] S_PRG\n\n### IDENTITY\n- SUMMARY: Project root map\n\n### WAYPOINTS\n(없음)\n\n### COMMENT\n(없음)\n</MAP>\n"
		if err := os.WriteFile(rootMap, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

// FS implements internal.Storage backed by the local filesystem.
type FS struct {
	Root string // project root (directory that contains .loadstar/)
}

func NewFS(root string) *FS {
	return &FS{Root: root}
}

// AvcsPath returns the absolute path for a relative .loadstar sub-path.
func (f *FS) AvcsPath(rel string) string {
	return filepath.Join(f.Root, AvcsDir, rel)
}

func (f *FS) Read(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (f *FS) Write(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func (f *FS) Exists(path string) bool {
	return Exists(path)
}

func (f *FS) CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// ListByPrefix returns files in dir whose base name starts with prefix.
func (f *FS) ListByPrefix(dir, prefix string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), prefix) {
			result = append(result, filepath.Join(dir, e.Name()))
		}
	}
	return result, nil
}

// Package-level helpers (kept for backward compatibility)

func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func WriteFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
