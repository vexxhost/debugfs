package kvm

import "github.com/vexxhost/debugfs"

// FS reads KVM-specific debugfs data.
type FS struct {
	debugfs debugfs.FS
}

// NewFS returns a KVM debugfs reader backed by dfs.
func NewFS(dfs debugfs.FS) FS {
	return FS{debugfs: dfs}
}
