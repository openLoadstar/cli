package internal

// Storage defines file I/O operations on the .loadstar directory.
type Storage interface {
	Read(path string) (string, error)
	Write(path, content string) error
	Exists(path string) bool
	CopyFile(src, dst string) error
	ListByPrefix(dir, prefix string) ([]string, error)
}

