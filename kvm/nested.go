package kvm

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/vexxhost/debugfs/value"
)

type NestedRun struct {
	// PID is the userspace process ID associated with the KVM VM.
	PID int

	// Path is the debugfs file that provided Count.
	Path string

	// Count is the value of the KVM nested_run counter.
	Count uint64
}

// NestedRuns reads all KVM nested_run counters exposed by debugfs.
func (fs FS) NestedRuns(ctx context.Context) ([]NestedRun, error) {
	files, err := fs.debugfs.Glob("kvm", "[0-9]*-*", "nested_run")
	if err != nil {
		return nil, err
	}

	runs := make([]NestedRun, 0, len(files))

	for _, file := range files {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		pid, err := pidFromFile(file)
		if err != nil {
			return nil, err
		}

		count, err := value.Uint64(file)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, fmt.Errorf("read nested_run for pid %d: %w", pid, err)
		}

		runs = append(runs, NestedRun{
			PID:   pid,
			Path:  file.Path(),
			Count: count,
		})
	}

	sort.Slice(runs, func(i, j int) bool {
		if runs[i].PID == runs[j].PID {
			return runs[i].Path < runs[j].Path
		}
		return runs[i].PID < runs[j].PID
	})

	return runs, nil
}
