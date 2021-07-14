package thrift

import (
	"errors"
)

type TProtocolError byte

const (
	TProtocolErrorUnknown TProtocolError = iota
	TProtocolErrorInvalidData
	TProtocolErrorNegativeSize
	TProtocolErrorSizeLimit
	TProtocolErrorBadVersion
)

type TProtocolException struct {
	kind TProtocolError
	err  error
}

func NewTProtocolException(k TProtocolError, text string) *TProtocolException {
	return &TProtocolException{k, errors.New(text)}
}

func NewTProtocolExceptionFromError(err error) error {
	if err == nil {
		return nil
	}
	if err, ok := err.(*TProtocolException); ok {
		return err
	}
	return &TProtocolException{TProtocolErrorUnknown, err}
}

func (e *TProtocolException) Kind() TProtocolError {
	return e.kind
}

func (e *TProtocolException) Unwrap() error {
	return e.err
}

func (e *TProtocolException) Error() string {
	return e.err.Error()
}
