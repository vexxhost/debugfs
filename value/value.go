package value

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/vexxhost/debugfs/file"
)

// Decoder converts raw debugfs file contents into T.
type Decoder[T any] func([]byte) (T, error)

// ErrNilDecoder reports a nil decoder passed to Read.
var ErrNilDecoder = errors.New("nil decoder")

// Read reads file and decodes its contents with decode.
func Read[T any](f file.File, decode Decoder[T]) (T, error) {
	var zero T

	if decode == nil {
		return zero, ErrNilDecoder
	}

	raw, err := f.ReadAll()
	if err != nil {
		return zero, err
	}

	v, err := decode(raw)
	if err != nil {
		return zero, fmt.Errorf("%s: decode: %w", f.Path(), err)
	}

	return v, nil
}

// Uint64 reads a decimal unsigned integer from file.
func Uint64(f file.File) (uint64, error) {
	return Read(f, ParseUint64)
}

// ParseUint64 decodes a decimal unsigned integer.
func ParseUint64(raw []byte) (uint64, error) {
	return strconv.ParseUint(strings.TrimSpace(string(raw)), 10, 64)
}
