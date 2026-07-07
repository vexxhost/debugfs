package debugfs

import (
	"errors"
	"fmt"
	stdfs "io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	debugfile "github.com/vexxhost/debugfs/file"
)

const (
	// DefaultMountPoint is the conventional Linux debugfs mount point.
	DefaultMountPoint = "/sys/kernel/debug"
)

// ErrInvalidPath reports a path that cannot be resolved under an FS mount point.
var ErrInvalidPath = errors.New("invalid debugfs path")

// FS resolves debugfs files relative to a mount point.
type FS struct {
	mountPoint string
	fsys       stdfs.FS
}

// NewFS returns an FS rooted at mountPoint.
func NewFS(mountPoint string) (FS, error) {
	if mountPoint == "" {
		return FS{}, fmt.Errorf("debugfs mount point: %w", ErrInvalidPath)
	}

	abs, err := filepath.Abs(mountPoint)
	if err != nil {
		return FS{}, fmt.Errorf("debugfs mount point: %w", err)
	}

	root := filepath.Clean(abs)
	return FS{
		mountPoint: root,
		fsys:       os.DirFS(root),
	}, nil
}

// DefaultFS returns an FS rooted at DefaultMountPoint.
func DefaultFS() (FS, error) {
	return NewFS(DefaultMountPoint)
}

// MountPoint returns the debugfs mount point used by fs.
func (fs FS) MountPoint() string {
	return fs.mountPoint
}

// Open resolves elem under fs and returns the resulting debugfs file.
func (fs FS) Open(elem ...string) (debugfile.File, error) {
	name, err := fs.name(elem...)
	if err != nil {
		return debugfile.File{}, err
	}

	return debugfile.New(fs.absolutePath(name)), nil
}

// Glob returns files matching a pattern under the debugfs mount point.
func (fs FS) Glob(elem ...string) ([]debugfile.File, error) {
	pattern, err := fs.name(elem...)
	if err != nil {
		return nil, err
	}

	matches, err := stdfs.Glob(fs.fsys, pattern)
	if err != nil {
		return nil, err
	}

	files := make([]debugfile.File, 0, len(matches))
	for _, match := range matches {
		files = append(files, debugfile.New(fs.absolutePath(match)))
	}

	return files, nil
}

func (fs FS) name(elem ...string) (string, error) {
	if fs.mountPoint == "" {
		return "", fmt.Errorf("debugfs mount point: %w", ErrInvalidPath)
	}

	parts := make([]string, 0, len(elem))

	for _, e := range elem {
		if e == "" {
			continue
		}
		if path.IsAbs(e) || filepath.IsAbs(e) || containsParentSegment(e) {
			return "", fmt.Errorf("%q: %w", e, ErrInvalidPath)
		}
		parts = append(parts, e)
	}

	name := path.Join(parts...)
	if !stdfs.ValidPath(name) || name == "." {
		return "", fmt.Errorf("%q: %w", strings.Join(elem, "/"), ErrInvalidPath)
	}

	return name, nil
}

func (fs FS) absolutePath(name string) string {
	return filepath.Join(fs.mountPoint, filepath.FromSlash(name))
}

func containsParentSegment(path string) bool {
	for _, part := range strings.Split(filepath.ToSlash(path), "/") {
		if part == ".." {
			return true
		}
	}
	return false
}
