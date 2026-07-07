// Package kvm reads KVM-specific debugfs data.
//
// The package currently supports the nested_run counter exposed under:
//
//	/sys/kernel/debug/kvm/*/nested_run
//
// Each KVM debugfs directory is keyed by the PID of the userspace process that
// owns the KVM VM. The package intentionally does not assume that process is
// QEMU. Callers that need libvirt, QEMU, or OpenStack metadata should enrich
// the returned PID in a higher-level package.
//
// Example:
//
//	dfs, err := debugfs.DefaultFS()
//	if err != nil {
//		return err
//	}
//
//	runs, err := kvm.NewFS(dfs).NestedRuns(ctx)
package kvm
