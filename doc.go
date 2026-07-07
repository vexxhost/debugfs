// Package debugfs resolves files under a Linux debugfs mount.
//
// The package intentionally stays close to filesystem primitives. FS uses the
// standard library io/fs path model for names and globbing, resolves them under
// a configured debugfs mount point, and returns file.File values. Bounded
// pseudo-file reads live in package file. Typed decoding lives in package
// value, and subsystem-specific readers live in packages such as kvm.
//
// A direct read starts with DefaultFS or NewFS, then uses the returned
// file.File:
//
//	dfs, err := debugfs.DefaultFS()
//	if err != nil {
//		return err
//	}
//
//	file, err := dfs.Open("kvm", "1234-16", "nested_run")
//	if err != nil {
//		return err
//	}
//
//	raw, err := file.ReadAll()
//
// For typed values, use package value. For KVM-specific data, use package kvm.
package debugfs
