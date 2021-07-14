package thrift

// TApplicationError kind of TApplicationException.
type TApplicationError = int32

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

// TApplicationException an application-level exception.
type TApplicationException struct {
	Message string            `thrift:"1"`
	Type    TApplicationError `thrift:"2"`
}

// Write writes values to p.
func (e *TApplicationException) Write(p TProtocol) (err error) {
	if err = p.WriteStructBegin(TStructHeader{"TApplicationException"}); err != nil {
		return
	}
	if len(e.Message) > 0 {
		if err = p.WriteFieldBegin(TFieldHeader{"message", STRING, 1}); err != nil {
			return
		}
		if err = p.WriteString(e.Message); err != nil {
			return
		}
		if err = p.WriteFieldEnd(); err != nil {
			return
		}
	}
	if err = p.WriteFieldBegin(TFieldHeader{"type", I32, 2}); err != nil {
		return
	}
	if err = p.WriteI32(e.Type); err != nil {
		return
	}
	if err = p.WriteFieldEnd(); err != nil {
		return
	}
	if err = p.WriteFieldStop(); err != nil {
		return
	}
	err = p.WriteStructEnd()
	return
}

// Read reads values from p.
func (e *TApplicationException) Read(p TProtocol) (err error) {
	if _, err = p.ReadStructBegin(); err != nil {
		return
	}
	var h TFieldHeader
	for {
		if h, err = p.ReadFieldBegin(); err != nil {
			return
		}
		if h.Type == STOP {
			break
		}
		switch h.Identity {
		case 1:
			if h.Type == STRING {
				e.Message, err = p.ReadString()
			} else {
				err = p.Skip(h.Type)
			}
			if err != nil {
				return
			}
		case 2:
			if h.Type == I32 {
				e.Type, err = p.ReadI32()
			} else {
				err = p.Skip(h.Type)
			}
			if err != nil {
				return
			}
		default:
			if err = p.Skip(h.Type); err != nil {
				return
			}
		}
		if err = p.ReadFieldEnd(); err != nil {
			return
		}
	}
	err = p.ReadStructEnd()
	return
}

// Error returns a error message.
func (e *TApplicationException) Error() string {
	return e.Message
}
