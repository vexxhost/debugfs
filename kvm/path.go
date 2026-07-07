package kvm

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vexxhost/debugfs/file"
)

func pidFromFile(f file.File) (int, error) {
	path := f.Path()
	dir := filepath.Base(filepath.Dir(path))
	pidText, _, found := strings.Cut(dir, "-")
	if !found {
		return 0, fmt.Errorf("%s: parse KVM debugfs directory: missing dash", path)
	}

	pid, err := strconv.Atoi(pidText)
	if err != nil {
		return 0, fmt.Errorf("%s: parse KVM debugfs PID: %w", path, err)
	}

	return pid, nil
}
