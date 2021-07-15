package thrift

import (
	"container/list"
	"context"
	"encoding/binary"
	"fmt"
	"math"
)

// NewTCompactProtocolFactory returns new TProtocolFactory of NewtCompactProtocol.
func NewTCompactProtocolFactory(cfg *TConfiguration) TProtocolFactory {
	return &tProtocolFactory{cfg, NewTCompactProtocol}
}

// NewTCompactProtocol returns new compact protocol.
func NewTCompactProtocol(t TTransport, cfg *TConfiguration) TProtocol {
	p := &tCompactProtocol{cfg: cfg}
	p.TExtraTransport = NewTExtraTransport(t, cfg)
	p.identityStack.Init()
	p.booleanWrite = -1
	return p
}

const (
	compactProtocolID      = 0x82
	compactVersion         = 1
	compactVersionMask     = 0x1f
	compactTypeMask        = 0xe0
	compactTypeBits        = 0x07
	compactTypeShiftAmount = 5
)

type compactType = byte

const (
	compactStop compactType = iota
	compactBooleanTrue
	compactBooleanFalse
	compactByte
	compactI16
	compactI32
	compactI64
	compactDouble
	compactBinary
	compactList
	compactSet
	compactMap
	compactStruct
)

var tTypeToCompactType = map[TType]compactType{
	STOP:   compactStop,
	BOOL:   compactByte,
	I16:    compactI16,
	U16:    compactI16,
	I32:    compactI32,
	U32:    compactI32,
	I64:    compactI64,
	U64:    compactI64,
	DOUBLE: compactDouble,
	STRING: compactBinary,
	LIST:   compactList,
	SET:    compactSet,
	MAP:    compactMap,
	STRUCT: compactStruct,
}

var compactTypeToTType = map[compactType]TType{
	compactStop:         STOP,
	compactBooleanFalse: BOOL,
	compactBooleanTrue:  BOOL,
	compactI16:          I16,
	compactI32:          I32,
	compactI64:          I64,
	compactDouble:       DOUBLE,
	compactBinary:       STRING,
	compactList:         LIST,
	compactSet:          SET,
	compactMap:          MAP,
	compactStruct:       STRUCT,
}

type tCompactProtocol struct {
	TExtraTransport
	cfg *TConfiguration
	buf [binary.MaxVarintLen64]byte

	identityStack list.List
	lastIdentity  int16

	booleanWrite int16
	booleanRead  byte
}

func (p *tCompactProtocol) WriteMessageBegin(h TMessageHeader) (err error) {
	if err = p.WriteByte(compactProtocolID); err == nil {
		if err = p.WriteByte(compactVersion | (h.Type << compactTypeShiftAmount)); err == nil {
			if err = p.WriteI32(h.Identity); err == nil {
				err = p.WriteString(h.Name)
			}
		}
	}
	return
}

func (p *tCompactProtocol) WriteMessageEnd() error {
	return nil
}

func (p *tCompactProtocol) WriteStructBegin(h TStructHeader) error {
	p.identityStack.PushBack(p.lastIdentity)
	p.lastIdentity = 0
	return nil
}

func (p *tCompactProtocol) WriteStructEnd() error {
	e := p.identityStack.Back()
	p.lastIdentity = e.Value.(int16)
	p.identityStack.Remove(e)
	return nil
}

func (p *tCompactProtocol) WriteFieldBegin(h TFieldHeader) error {
	if h.Type == BOOL {
		p.booleanWrite = h.Identity
		return nil
	} else {
		if c, ok := tTypeToCompactType[h.Type]; ok {
			return p.writeFieldHeader(c, h.Identity)
		}
		return NewTProtocolException(TProtocolErrorInvalidData, fmt.Sprintf("unexpected TType: %d", h.Type))
	}
}

func (p *tCompactProtocol) writeFieldHeader(t compactType, i int16) (err error) {
	delta := i - p.lastIdentity
	if 0 < delta && delta <= 15 {
		if err = p.WriteByte(byte((delta << 4)) | t); err != nil {
			return
		}
	} else {
		if err = p.WriteByte(t); err != nil {
			return
		}
		if err = p.WriteI16(i); err != nil {
			return
		}
	}
	p.lastIdentity = i
	return
}

func (p *tCompactProtocol) WriteFieldEnd() error {
	return nil
}

func (p *tCompactProtocol) WriteFieldStop() error {
	return p.WriteByte(compactStop)
}

func (p *tCompactProtocol) WriteMapBegin(h TMapHeader) (err error) {
	if err = p.WriteByte(h.Key); err == nil {
		if err = p.WriteByte(h.Value); err == nil {
			err = p.writeSize(h.Size)
		}
	}
	return
}

func (p *tCompactProtocol) WriteMapEnd() error {
	return nil
}

func (p *tCompactProtocol) WriteSetBegin(h TSetHeader) (err error) {
	if err = p.WriteByte(h.Element); err == nil {
		err = p.writeSize(h.Size)
	}
	return
}

func (p *tCompactProtocol) WriteSetEnd() error {
	return nil
}

func (p *tCompactProtocol) WriteListBegin(h TListHeader) (err error) {
	if err = p.WriteByte(h.Element); err == nil {
		err = p.writeSize(h.Size)
	}
	return
}

func (p *tCompactProtocol) WriteListEnd() error {
	return nil
}

func (p *tCompactProtocol) WriteBool(v bool) error {
	c := compactBooleanFalse
	if v {
		c = compactBooleanTrue
	}
	if p.booleanWrite != -1 {
		err := p.writeFieldHeader(c, p.booleanWrite)
		p.booleanWrite = -1
		return err
	}
	return p.WriteByte(c)
}

func (p *tCompactProtocol) WriteByte(v byte) error {
	return NewTProtocolExceptionFromError(p.TExtraTransport.WriteByte(v))
}

func (p *tCompactProtocol) WriteDouble(v float64) error {
	buf := p.buf[:8]
	binary.LittleEndian.PutUint64(p.buf[:], math.Float64bits(v))
	_, err := p.Write(buf)
	return err
}

func (p *tCompactProtocol) WriteU16(v uint16) error {
	return p.WriteU64(uint64(v))
}

func (p *tCompactProtocol) WriteI16(v int16) error {
	return p.WriteI64(int64(v))
}

func (p *tCompactProtocol) WriteU32(v uint32) error {
	return p.WriteU64(uint64(v))
}

func (p *tCompactProtocol) WriteI32(v int32) error {
	return p.WriteI64(int64(v))
}

func (p *tCompactProtocol) WriteU64(v uint64) error {
	n := binary.PutUvarint(p.buf[:], v)
	_, err := p.Write(p.buf[:n])
	return err
}

func (p *tCompactProtocol) WriteI64(v int64) error {
	n := binary.PutVarint(p.buf[:], v)
	_, err := p.Write(p.buf[:n])
	return err
}

func (p *tCompactProtocol) WriteString(v string) error {
	return p.WriteBinary([]byte(v))
}

func (p *tCompactProtocol) WriteBinary(v []byte) (err error) {
	if err = p.writeSize(len(v)); err == nil {
		_, err = p.Write(v)
	}
	return
}

func (p *tCompactProtocol) Write(v []byte) (int, error) {
	n, err := p.TExtraTransport.Write(v)
	return n, NewTProtocolExceptionFromError(err)
}

func (p *tCompactProtocol) writeSize(v int) error {
	n := binary.PutUvarint(p.buf[:], uint64(v))
	_, err := p.Write(p.buf[:n])
	return err
}

func (p *tCompactProtocol) ReadMessageBegin() (h TMessageHeader, err error) {
	var b byte
	if b, err = p.ReadByte(); err != nil {
		return
	}
	if b != compactProtocolID {
		err = NewTProtocolException(TProtocolErrorBadVersion, "bad protocol id in message header")
		return
	}
	if b, err = p.ReadByte(); err != nil {
		return
	}
	if (b & compactVersionMask) != compactVersion {
		err = NewTProtocolException(TProtocolErrorBadVersion, "bad version in message header")
		return
	}
	h.Type = (b >> compactTypeShiftAmount) & compactTypeBits
	if h.Identity, err = p.ReadI32(); err == nil {
		h.Name, err = p.ReadString()
	}
	return
}

func (p *tCompactProtocol) ReadMessageEnd() error {
	return nil
}

func (p *tCompactProtocol) ReadStructBegin() (h TStructHeader, err error) {
	p.identityStack.PushBack(p.lastIdentity)
	p.lastIdentity = 0
	return
}

func (p *tCompactProtocol) ReadStructEnd() error {
	e := p.identityStack.Back()
	p.lastIdentity = e.Value.(int16)
	p.identityStack.Remove(e)
	return nil
}

func (p *tCompactProtocol) ReadFieldBegin() (h TFieldHeader, err error) {
	var b byte
	if b, err = p.ReadByte(); err != nil {
		return
	}
	if (b & 0x0f) == compactStop {
		return
	}
	delta := int16((b & 0xf0) >> 4)
	if delta == 0 {
		if h.Identity, err = p.ReadI16(); err != nil {
			return
		}
	} else {
		h.Identity = p.lastIdentity + delta
	}
	h.Type = b & 0x0f
	switch h.Type {
	case compactBooleanTrue:
		p.booleanRead = 1
	case compactBooleanFalse:
		p.booleanRead = 2
	}
	if b, ok := compactTypeToTType[h.Type]; ok {
		h.Type = b
	} else {
		err = NewTProtocolException(TProtocolErrorInvalidData, fmt.Sprintf("unexpected compact Type: %d", h.Type))
	}
	p.lastIdentity = h.Identity
	return
}

func (p *tCompactProtocol) ReadFieldEnd() error {
	return nil
}

func (p *tCompactProtocol) ReadMapBegin() (h TMapHeader, err error) {
	if h.Key, err = p.ReadByte(); err == nil {
		if h.Value, err = p.ReadByte(); err == nil {
			h.Size, err = p.readSize()
		}
	}
	return
}

func (p *tCompactProtocol) ReadMapEnd() error {
	return nil
}

func (p *tCompactProtocol) ReadSetBegin() (h TSetHeader, err error) {
	if h.Element, err = p.ReadByte(); err == nil {
		h.Size, err = p.readSize()
	}
	return
}

func (p *tCompactProtocol) ReadSetEnd() error {
	return nil
}

func (p *tCompactProtocol) ReadListBegin() (h TListHeader, err error) {
	if h.Element, err = p.ReadByte(); err == nil {
		h.Size, err = p.readSize()
	}
	return
}

func (p *tCompactProtocol) ReadListEnd() error {
	return nil
}

func (p *tCompactProtocol) ReadBool() (bool, error) {
	if p.booleanRead != 0 {
		v := p.booleanRead == 1
		p.booleanRead = 0
		return v, nil
	}
	v, err := p.ReadByte()
	return v == compactBooleanTrue, err
}

func (p *tCompactProtocol) ReadByte() (byte, error) {
	v, err := p.TExtraTransport.ReadByte()
	return v, NewTProtocolExceptionFromError(err)
}

func (p *tCompactProtocol) ReadDouble() (float64, error) {
	buf := p.buf[:8]
	_, err := p.Read(buf)
	return math.Float64frombits(binary.LittleEndian.Uint64(buf)), err
}

func (p *tCompactProtocol) ReadU16() (uint16, error) {
	v, err := p.ReadU64()
	return uint16(v), err
}

func (p *tCompactProtocol) ReadI16() (int16, error) {
	v, err := p.ReadI64()
	return int16(v), err
}

func (p *tCompactProtocol) ReadU32() (uint32, error) {
	v, err := p.ReadU64()
	return uint32(v), err
}

func (p *tCompactProtocol) ReadI32() (int32, error) {
	v, err := p.ReadI64()
	return int32(v), err
}

func (p *tCompactProtocol) ReadU64() (uint64, error) {
	v, err := binary.ReadUvarint(p.TExtraTransport)
	return v, NewTProtocolExceptionFromError(err)
}

func (p *tCompactProtocol) ReadI64() (int64, error) {
	v, err := binary.ReadVarint(p.TExtraTransport)
	return v, NewTProtocolExceptionFromError(err)
}

func (p *tCompactProtocol) ReadString() (v string, err error) {
	var n int
	if n, err = p.readSize(); err == nil {
		v, err = p.readStringBody(n)
	}
	return
}

func (p *tCompactProtocol) ReadBinary() (v []byte, err error) {
	var n int
	if n, err = p.readSize(); err == nil {
		if err = p.cfg.CheckSizeForProtocol(n); err != nil {
			return
		}
		v = make([]byte, n)
		_, err = p.Read(v)
	}
	return
}

func (p *tCompactProtocol) Read(v []byte) (int, error) {
	n, err := p.TExtraTransport.Read(v)
	return n, NewTProtocolExceptionFromError(err)
}

func (p *tCompactProtocol) readSize() (int, error) {
	v, err := binary.ReadUvarint(p.TExtraTransport)
	return int(v), NewTProtocolExceptionFromError(err)
}

func (p *tCompactProtocol) readStringBody(n int) (string, error) {
	if err := p.cfg.CheckSizeForProtocol(n); err != nil {
		return "", err
	}
	v := make([]byte, n)
	_, err := p.Read(v)
	return string(v), err
}

func (p *tCompactProtocol) Skip(v TType) error {
	return Skip(v, p)
}

func (p *tCompactProtocol) Flush(ctx context.Context) error {
	return NewTProtocolExceptionFromError(p.TExtraTransport.Flush(ctx))
}
