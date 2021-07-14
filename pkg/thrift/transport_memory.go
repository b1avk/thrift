package thrift

import (
	"bytes"
	"context"
)

type TMemoryBuffer struct {
	*bytes.Buffer
}

func NewTMemoryBuffer() *TMemoryBuffer {
	return &TMemoryBuffer{new(bytes.Buffer)}
}

func (*TMemoryBuffer) Flush(ctx context.Context) error {
	return nil
}
