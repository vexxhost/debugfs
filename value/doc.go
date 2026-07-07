// Package value decodes typed values from resolved debugfs files.
//
// The package is intentionally independent of any specific debugfs subsystem.
// It operates on file.File values, usually resolved by the root debugfs
// package.
//
// For common decimal counters:
//
//	file, err := dfs.Open("kvm", "1234-16", "nested_run")
//	if err != nil {
//		return err
//	}
//
//	count, err := value.Uint64(file)
//
// For custom formats, provide a decoder:
//
//	v, err := value.Read(file, func(raw []byte) (MyValue, error) {
//		return parseMyValue(raw)
//	})
package value
