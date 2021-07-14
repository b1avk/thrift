package thrift

import (
	"context"
	"io"
)

type TTransportFactory interface {
	GetTTransport(TTransport) (TTransport, error)
}

type TTransport interface {
	io.ReadWriter
	TFlusher
}

type TFlusher interface {
	Flush(ctx context.Context) (err error)
}
