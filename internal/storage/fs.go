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
// writes an initial M://root Map file, and sets up .claude/ hooks.
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

	// Create .claude/ hooks for meta-sync reminders
	if err := initClaudeHooks(projectRoot); err != nil {
		return err
	}

	return nil
}

// initClaudeHooks creates .claude/settings.json and hooks/loadstar-drift-check.sh
// if they don't already exist.
func initClaudeHooks(projectRoot string) error {
	hooksDir := filepath.Join(projectRoot, ".claude", "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return err
	}

	settingsPath := filepath.Join(projectRoot, ".claude", "settings.json")
	if !Exists(settingsPath) {
		settings := `{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [
          {
            "type": "command",
            "command": "bash \"$CLAUDE_PROJECT_DIR/.claude/hooks/loadstar-drift-check.sh\""
          }
        ]
      }
    ]
  }
}
`
		if err := os.WriteFile(settingsPath, []byte(settings), 0644); err != nil {
			return err
		}
	}

	hookPath := filepath.Join(hooksDir, "loadstar-drift-check.sh")
	if !Exists(hookPath) {
		hook := `#!/bin/bash
# loadstar-drift-check.sh
# PostToolUse hook: 소스코드 수정 시 LOADSTAR 메타데이터 갱신 리마인더

INPUT=$(cat)
FILE_PATH=$(echo "$INPUT" | jq -r '.tool_input.file_path // empty')

if [[ -z "$FILE_PATH" ]]; then
  exit 0
fi

if [[ "$FILE_PATH" == *".loadstar"* || "$FILE_PATH" == *".claude"* ]]; then
  exit 0
fi

BASENAME=$(basename "$FILE_PATH")
case "$BASENAME" in
  go.mod|go.sum|pom.xml|package.json|package-lock.json|*.json|*.yaml|*.yml|*.toml|*.md|*.txt|*.css|LICENSE|.gitignore)
    exit 0
    ;;
esac

echo "[LOADSTAR] 소스 파일 수정됨: $FILE_PATH"
echo "[LOADSTAR] 작업 착수 전 대상 WayPoint TECH_SPEC에 작업 항목을 [ ]로 등록했는지 확인하세요."
echo "[LOADSTAR] 작업 완료 후 [x] YYYY-MM-DD로 체크하고, 필요 시 STATUS를 갱신하세요."

exit 0
`
		if err := os.WriteFile(hookPath, []byte(hook), 0755); err != nil {
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
