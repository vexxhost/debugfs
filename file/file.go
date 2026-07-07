package file

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const defaultMaxReadSize = 1024 * 1024

// ErrInvalidPath reports a file with no resolved path.
var ErrInvalidPath = errors.New("invalid debugfs file path")

// File is a resolved debugfs file.
type File struct {
	path string
}

// New returns a File for path.
func New(path string) File {
	return File{path: path}
}

// Path returns the resolved absolute path for f.
func (f File) Path() string {
	return f.path
}

// ReadAll reads f with a bounded read suitable for small pseudo-files.
func (f File) ReadAll() ([]byte, error) {
	if f.path == "" {
		return nil, fmt.Errorf("debugfs file path: %w", ErrInvalidPath)
	}

	return readFile(f.path)
}

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(io.LimitReader(file, defaultMaxReadSize))
}
