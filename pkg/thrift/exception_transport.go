package thrift

import (
	"errors"
	"io"
)

// TTransportError kind of TTransportException.
type TTransportError byte

const (
	TTransportErrorUnknown TTransportError = iota
	TTransportErrorEOF
	TTransportErrorTimeout
)

// TTransportException a transport-level exception.
type TTransportException struct {
	kind TTransportError
	err  error
}

// NewTTransportException returns new TTransportException.
func NewTTransportException(k TTransportError, text string) *TTransportException {
	return &TTransportException{k, errors.New(text)}
}

// NewTTransportExceptionFromError returns new TTransportException.
// it will returns nil if given err is nil.
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

// Kind returns kind of TTransportException.
func (e *TTransportException) Kind() TTransportError {
	return e.kind
}

// Unwrap returns holding error.
func (e *TTransportException) Unwrap() error {
	return e.err
}

// Timeout returns true if is timeout error.
func (e *TTransportException) Timeout() bool {
	return e.kind == TTransportErrorTimeout || isTimeout(e.err)
}

// Error returns error message.
func (e *TTransportException) Error() string {
	return e.err.Error()
}

func isTimeout(err error) bool {
	var t interface{ Timeout() bool }
	return errors.As(err, &t) && t.Timeout()
}
