package file

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestReadAll(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "counter")
	if err := os.WriteFile(path, []byte("42\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := New(path).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "42\n" {
		t.Fatalf("ReadAll() = %q, want %q", got, "42\n")
	}
}

func TestZeroFileReturnsInvalidPath(t *testing.T) {
	_, err := (File{}).ReadAll()
	if !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("ReadAll() error = %v, want ErrInvalidPath", err)
	}
}
