package thrift

import (
	"bytes"
	"context"
)

type TMemoryBuffer struct {
	*bytes.Buffer
}

func NewTMemoryBuffer(cfg *TConfiguration) *TMemoryBuffer {
	return &TMemoryBuffer{bytes.NewBuffer(make([]byte, cfg.GetMaxBufferSize()))}
}

func (*TMemoryBuffer) Flush(ctx context.Context) error {
	return nil
}
