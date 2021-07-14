package thrift

import (
	"context"
	"io"
)

type TTransport interface {
	io.ReadWriter
	TFlusher
}

type TFlusher interface {
	Flush(ctx context.Context) (err error)
}
