package thrift

import (
	"bytes"
	"context"
)

// TMemoryBuffer memory buffer-based implementation for TTransport.
type TMemoryBuffer struct {
	*bytes.Buffer
}

// NewTMemoryBuffer returns new empty TMemoryBuffer.
func NewTMemoryBuffer() *TMemoryBuffer {
	return &TMemoryBuffer{new(bytes.Buffer)}
}

// Flush flushing is no-op; always returns nil.
func (*TMemoryBuffer) Flush(ctx context.Context) error {
	return nil
}
