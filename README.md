# debugfs

Go helpers for reading Linux debugfs data.

This module provides a small set of composable primitives for resolving files
under a debugfs mount, reading pseudo-files safely, decoding typed values, and
building subsystem-specific readers on top.

The module path is defined in [go.mod](./go.mod).

This project is licensed under the Apache License 2.0. See [LICENSE](./LICENSE).

API documentation lives in the Go package docs so it stays close to the code and
renders cleanly in `go doc` and pkg.go.dev:

https://pkg.go.dev/github.com/vexxhost/debugfs
