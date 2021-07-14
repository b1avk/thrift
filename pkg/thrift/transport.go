package thrift

import (
	"context"
	"io"
)

// TTransportFactory a factory of TTransport.
type TTransportFactory interface {
	// GetTransport returns new TTransport.
	GetTransport(t TTransport) (TTransport, error)
}

// TTransport interface that groups io.ReadWriter and Flusher.
// it should be used by TProtocol
type TTransport interface {
	io.ReadWriter
	TFlusher
}

// TFlusher interface that wraps Flush method which
// allows to flush underlying buffer.
// it implemented by TTransport and TProtocol.
type TFlusher interface {
	Flush(ctx context.Context) (err error)
}
