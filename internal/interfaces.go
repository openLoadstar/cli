package internal

// Storage defines file I/O operations on the .loadstar directory.
type Storage interface {
	Read(path string) (string, error)
	Write(path, content string) error
	Exists(path string) bool
	CopyFile(src, dst string) error
	ListByPrefix(dir, prefix string) ([]string, error)
}

// GitClient defines git operations required by LOADSTAR.
type GitClient interface {
	Commit(message string) (string, error)
	LatestHash() (string, error)
}
