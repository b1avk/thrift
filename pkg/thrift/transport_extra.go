package thrift

import (
	"io"
)

type TExtraTransport interface {
	TTransport
	io.ByteReader
	io.ByteWriter
}

func NewTExtraTransport(t TTransport) TExtraTransport {
	if t, ok := t.(TExtraTransport); ok {
		return t
	}
	return &tExtraTransport{TTransport: t}
}

type tExtraTransport struct {
	TTransport
	cache [1]byte
}

func (t *tExtraTransport) WriteByte(b byte) error {
	t.cache[0] = b
	_, err := t.Write(t.cache[:])
	return NewTTransportExceptionFromError(err)
}

func (t *tExtraTransport) ReadByte() (byte, error) {
	_, err := t.Read(t.cache[:])
	return t.cache[0], NewTTransportExceptionFromError(err)
}
