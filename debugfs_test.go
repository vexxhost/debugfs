package debugfs

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestFileReadAll(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "kvm", "123-4"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "kvm", "123-4", "nested_run"), []byte("42\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	fs, err := NewFS(root)
	if err != nil {
		t.Fatal(err)
	}

	file, err := fs.Open("kvm", "123-4", "nested_run")
	if err != nil {
		t.Fatal(err)
	}

	got, err := file.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "42\n" {
		t.Fatalf("ReadAll() = %q, want %q", got, "42\n")
	}
	if file.Path() != filepath.Join(root, "kvm", "123-4", "nested_run") {
		t.Fatalf("Path() = %q", file.Path())
	}
}

func TestGlob(t *testing.T) {
	root := t.TempDir()
	paths := []string{
		filepath.Join(root, "kvm", "123-4", "nested_run"),
		filepath.Join(root, "kvm", "456-7", "nested_run"),
	}
	for _, path := range paths {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("0\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	fs, err := NewFS(root)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.Glob("kvm", "[0-9]*-*", "nested_run")
	if err != nil {
		t.Fatal(err)
	}

	gotPaths := make([]string, 0, len(got))
	for _, file := range got {
		gotPaths = append(gotPaths, file.Path())
	}

	if !slices.Equal(gotPaths, paths) {
		t.Fatalf("Glob() paths = %#v, want %#v", gotPaths, paths)
	}
}

func TestRejectsEscapingPaths(t *testing.T) {
	fs, err := NewFS(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	_, err = fs.Open("..", "counter")
	if !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("Open() error = %v, want ErrInvalidPath", err)
	}

	_, err = fs.Open(filepath.Dir(fs.MountPoint()))
	if !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("Open() error = %v, want ErrInvalidPath", err)
	}
}

func TestNewFSRejectsEmptyMountPoint(t *testing.T) {
	_, err := NewFS("")
	if !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("NewFS() error = %v, want ErrInvalidPath", err)
	}
}
