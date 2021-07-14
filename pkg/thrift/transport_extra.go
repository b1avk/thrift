package thrift

import (
	"io"
)

// TExtraTransport enhanced version of TTransport
type TExtraTransport interface {
	TTransport
	io.ByteReader
	io.ByteWriter
}

// NewTExtraTransport wraps t to TExtraTransport.
func NewTExtraTransport(t TTransport, cfg *TConfiguration) TExtraTransport {
	cfg.Propagate(t)
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

func (t *tExtraTransport) SetTConfiguration(cfg *TConfiguration) {
	cfg.Propagate(t.TTransport)
}
