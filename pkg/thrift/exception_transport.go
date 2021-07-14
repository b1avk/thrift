package thrift

import (
	"errors"
	"io"
)

type TTransportError byte

const (
	TTransportErrorUnknown TTransportError = iota
	TTransportErrorEOF
	TTransportErrorTimeout
)

type TTransportException struct {
	kind TTransportError
	err  error
}

func NewTTransportException(k TTransportError, text string) *TTransportException {
	return &TTransportException{k, errors.New(text)}
}

func NewTTransportExceptionFromError(err error) error {
	if err == nil {
		return nil
	}
	if err, ok := err.(*TTransportException); ok {
		return err
	}
	e := &TTransportException{TTransportErrorUnknown, err}
	if errors.Is(e, io.EOF) || errors.Is(e, io.ErrUnexpectedEOF) {
		e.kind = TTransportErrorEOF
	}
	if isTimeout(err) {
		e.kind = TTransportErrorTimeout
	}
	return e
}

func (e *TTransportException) Kind() TTransportError {
	return e.kind
}

func (e *TTransportException) Unwrap() error {
	return e.err
}

func (e *TTransportException) Error() string {
	return e.err.Error()
}

func isTimeout(err error) bool {
	var t interface{ Timeout() bool }
	return errors.As(err, &t) && t.Timeout()
}
