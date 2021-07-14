package thrift

import (
	"context"
	"io"
)

type TTransportFactory interface {
	GetTTransport(TTransport) TTransport
}

type TTransport interface {
	io.ReadWriter
	TFlusher
}

type TFlusher interface {
	Flush(ctx context.Context) (err error)
}
