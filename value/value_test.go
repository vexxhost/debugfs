package value

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/vexxhost/debugfs"
	"github.com/vexxhost/debugfs/file"
)

func TestUint64(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "counter")
	if err := os.WriteFile(path, []byte("42\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	fs, err := debugfs.NewFS(root)
	if err != nil {
		t.Fatal(err)
	}

	file, err := fs.Open("counter")
	if err != nil {
		t.Fatal(err)
	}

	got, err := Uint64(file)
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Fatalf("Uint64() = %d, want 42", got)
	}
}

func TestReadWrapsDecoderErrorWithPath(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "counter"), []byte("not-a-number\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	fs, err := debugfs.NewFS(root)
	if err != nil {
		t.Fatal(err)
	}

	file, err := fs.Open("counter")
	if err != nil {
		t.Fatal(err)
	}

	_, err = Uint64(file)
	if err == nil {
		t.Fatal("Uint64() error = nil, want error")
	}
	var numErr *strconv.NumError
	if !errors.As(err, &numErr) {
		t.Fatalf("Uint64() error = %v, want strconv.NumError", err)
	}
	if got := err.Error(); got == "" || !containsAll(got, file.Path(), "decode") {
		t.Fatalf("Uint64() error = %q, want path and decode context", got)
	}
}

func TestReadRejectsNilDecoder(t *testing.T) {
	_, err := Read(file.File{}, Decoder[uint64](nil))
	if !errors.Is(err, ErrNilDecoder) {
		t.Fatalf("Read() error = %v, want ErrNilDecoder", err)
	}
}

func containsAll(s string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(s, part) {
			return false
		}
	}
	return true
}
