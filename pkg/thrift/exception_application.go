package thrift

type TApplicationError byte

const (
	TApplicationErrorUnknown TApplicationError = iota
	TApplicationErrorUnknownMethod
	TApplicationErrorInvalidMessageType
	TApplicationErrorWrongMethodName
	TApplicationErrorBadSequenceID
	TApplicationErrorMissingResult
	TApplicationErrorInternalError
	TApplicationErrorProtocolError
	TApplicationErrorInvalidTransform
	TApplicationErrorInvalidProtocol
	TApplicationErrorUnsupportedClientType
)

type TApplicationException struct {
	Message string            `thrift:"1"`
	Type    TApplicationError `thrift:"2"`
}

func (e *TApplicationException) Error() string {
	return e.Message
}
