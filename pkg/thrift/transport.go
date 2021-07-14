package thrift

import (
	"context"
	"io"
)

type TTransportFactory interface {
	GetTransport(TTransport) (TTransport, error)
}

type TTransport interface {
	io.ReadWriter
	TFlusher
}

type TFlusher interface {
	Flush(ctx context.Context) (err error)
}
