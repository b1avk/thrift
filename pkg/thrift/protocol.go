package thrift

import (
	"fmt"
)

type TProtocolFactory interface {
	GetProtocol(TTransport) TProtocol
}

type TProtocol interface {
	TFlusher

	WriteMessageBegin(h TMessageHeader) (err error)

	WriteMessageEnd() (err error)

	WriteStructBegin(h TStructHeader) (err error)

	WriteStructEnd() (err error)

	WriteFieldBegin(h TFieldHeader) (err error)

	WriteFieldEnd() (err error)

	WriteFieldStop() (err error)

	WriteMapBegin(h TMapHeader) (err error)

	WriteMapEnd() (err error)

	WriteSetBegin(h TSetHeader) (err error)

	WriteSetEnd() (err error)

	WriteListBegin(h TListHeader) (err error)

	WriteListEnd() (err error)

	WriteBool(v bool) (err error)

	WriteByte(v byte) (err error)

	WriteDouble(v float64) (err error)

	WriteU16(v uint16) (err error)

	WriteI16(v int16) (err error)

	WriteU32(v uint32) (err error)

	WriteI32(v int32) (err error)

	WriteU64(v uint64) (err error)

	WriteI64(v int64) (err error)

	WriteString(v string) (err error)

	WriteBinary(v []byte) (err error)

	ReadMessageBegin() (h TMessageHeader, err error)

	ReadMessageEnd() (err error)

	ReadStructBegin() (h TStructHeader, err error)

	ReadStructEnd() (err error)

	ReadFieldBegin() (h TFieldHeader, err error)

	ReadFieldEnd() (err error)

	ReadMapBegin() (h TMapHeader, err error)

	ReadMapEnd() (err error)

	ReadSetBegin() (h TSetHeader, err error)

	ReadSetEnd() (err error)

	ReadListBegin() (h TListHeader, err error)

	ReadListEnd() (err error)

	ReadBool() (v bool, err error)

	ReadByte() (v byte, err error)

	ReadDouble() (v float64, err error)

	ReadU16() (v uint16, err error)

	ReadI16() (v int16, err error)

	ReadU32() (v uint32, err error)

	ReadI32() (v int32, err error)

	ReadU64() (v uint64, err error)

	ReadI64() (v int64, err error)

	ReadString() (v string, err error)

	ReadBinary() (v []byte, err error)

	Skip(v TType) (err error)
}

func Skip(v TType, p TProtocol) (err error) {
	switch v {
	case BOOL:
		_, err = p.ReadBool()
	case BYTE:
		_, err = p.ReadByte()
	case DOUBLE:
		_, err = p.ReadDouble()
	case U16:
		_, err = p.ReadU16()
	case I16:
		_, err = p.ReadI16()
	case U32:
		_, err = p.ReadU32()
	case I32:
		_, err = p.ReadI32()
	case U64:
		_, err = p.ReadU64()
	case I64:
		_, err = p.ReadI64()
	case STRING:
		_, err = p.ReadString()
	case STRUCT:
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
			if err = p.Skip(h.Type); err != nil {
				return
			}
			if err = p.ReadFieldEnd(); err != nil {
				return
			}
		}
		err = p.ReadStructEnd()
	case MAP:
		var h TMapHeader
		if h, err = p.ReadMapBegin(); err != nil {
			return
		}
		for i := 0; i < h.Size; i++ {
			if err = p.Skip(h.Key); err != nil {
				return
			}
			if err = p.Skip(h.Value); err != nil {
				return
			}
		}
		err = p.ReadMapEnd()
	case SET:
		var h TSetHeader
		if h, err = p.ReadSetBegin(); err != nil {
			return
		}
		for i := 0; i < h.Size; i++ {
			if err = p.Skip(h.Element); err != nil {
				return
			}
		}
		err = p.ReadSetEnd()
	case LIST:
		var h TListHeader
		if h, err = p.ReadListBegin(); err != nil {
			return
		}
		for i := 0; i < h.Size; i++ {
			if err = p.Skip(h.Element); err != nil {
				return
			}
		}
		err = p.ReadListEnd()
	default:
		err = NewTProtocolException(TProtocolErrorInvalidData, fmt.Sprintf("unexpected TType: %d", v))
	}
	return
}
