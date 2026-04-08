package address

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Address represents a parsed LOADSTAR URI (e.g. W://root/dev/auth)
type Address struct {
	Type string // M, W, B
	Path string // root/dev/auth
	ID   string // auth
}

func Parse(raw string) (*Address, error) {
	// e.g. W://root/dev/auth
	parts := strings.SplitN(raw, "://", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid address format: %s", raw)
	}
	segments := strings.Split(parts[1], "/")
	return &Address{
		Type: parts[0],
		Path: parts[1],
		ID:   segments[len(segments)-1],
	}, nil
}

// ToFilePath converts a logical address to a physical file path.
// e.g. W://root/cli/cmd_create -> <baseDir>/WAYPOINT/root.cli.cmd_create.md
func (a *Address) ToFilePath(baseDir string) string {
	typeDir := typeDirMap[a.Type]
	dotName := strings.ReplaceAll(a.Path, "/", ".")
	return filepath.Join(baseDir, typeDir, dotName+".md")
}

var typeDirMap = map[string]string{
	"M": "MAP",
	"W": "WAYPOINT",
}
