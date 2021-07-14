package thrift

import (
	"errors"
)

// TProtocolError kind of TProtocolException.
type TProtocolError byte

const (
	TProtocolErrorUnknown TProtocolError = iota
	TProtocolErrorInvalidData
	TProtocolErrorNegativeSize
	TProtocolErrorSizeLimit
	TProtocolErrorBadVersion
)

// TProtocolException a protocol-level exception.
type TProtocolException struct {
	kind TProtocolError
	err  error
}

// NewTProtocolException returns new TProtocolException.
func NewTProtocolException(k TProtocolError, text string) *TProtocolException {
	return &TProtocolException{k, errors.New(text)}
}

// NewTProtocolExceptionFromError returns new TProtocolException.
// it will returns nil if given err is nil.
func NewTProtocolExceptionFromError(err error) error {
	if err == nil {
		return nil
	}
	if err, ok := err.(*TProtocolException); ok {
		return err
	}
	return &TProtocolException{TProtocolErrorUnknown, err}
}

// Kind returns kind of TProtocolException.
func (e *TProtocolException) Kind() TProtocolError {
	return e.kind
}

// Unwrap returns holding error.
func (e *TProtocolException) Unwrap() error {
	return e.err
}

// Error returns error message.
func (e *TProtocolException) Error() string {
	return e.err.Error()
}
