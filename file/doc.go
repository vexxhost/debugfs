// Package file provides bounded reads for resolved debugfs pseudo-files.
//
// Most callers should obtain File values from debugfs.FS rather than calling
// New directly. debugfs.FS resolves paths under a configured debugfs mount
// point; this package only knows how to read the file once a path has been
// resolved.
package file
