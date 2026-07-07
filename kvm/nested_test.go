package kvm

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/vexxhost/debugfs"
)

func TestNestedRuns(t *testing.T) {
	dfs, err := debugfs.NewFS(filepath.Join("testdata", "debugfs"))
	if err != nil {
		t.Fatal(err)
	}

	got, err := NewFS(dfs).NestedRuns(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	want := []NestedRun{
		{PID: 1001, Path: filepath.Join(dfs.MountPoint(), "kvm", "1001-15", "nested_run"), Count: 0},
		{PID: 1002, Path: filepath.Join(dfs.MountPoint(), "kvm", "1002-16", "nested_run"), Count: 0},
		{PID: 2001, Path: filepath.Join(dfs.MountPoint(), "kvm", "2001-16", "nested_run"), Count: 42},
		{PID: 2002, Path: filepath.Join(dfs.MountPoint(), "kvm", "2002-101", "nested_run"), Count: 0},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("NestedRuns() = %#v, want %#v", got, want)
	}
}

func TestNestedRunsHonorsContext(t *testing.T) {
	dfs, err := debugfs.NewFS(filepath.Join("testdata", "debugfs"))
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = NewFS(dfs).NestedRuns(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("NestedRuns() error = %v, want context.Canceled", err)
	}
}

func TestPIDFromNestedRunPath(t *testing.T) {
	dfs, err := debugfs.NewFS(filepath.Join("testdata", "debugfs"))
	if err != nil {
		t.Fatal(err)
	}

	file, err := dfs.Open("kvm", "2001-16", "nested_run")
	if err != nil {
		t.Fatal(err)
	}

	got, err := pidFromFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if got != 2001 {
		t.Fatalf("pidFromFile() = %d, want 2001", got)
	}
}
