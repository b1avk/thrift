package thrift

import (
	"context"
	"encoding/binary"
	"math"
)

func NewTBinaryProtocol(t TTransport, cfg *TConfiguration) TProtocol {
	p := &tBinaryProtocol{cfg: cfg.NonNil()}
	p.TExtraTransport = NewTExtraTransport(t, p.cfg)
	return p
}

type tBinaryProtocol struct {
	TExtraTransport
	cfg *TConfiguration
	buf [8]byte
}

func (p *tBinaryProtocol) SetTConfiguration(cfg *TConfiguration) {
	p.cfg = cfg.NonNil()
	p.cfg.Propagate(p.TExtraTransport)
}

func (p *tBinaryProtocol) WriteMessageBegin(h TMessageHeader) (err error) {
	// TODO implement.
	return
}

func (p *tBinaryProtocol) WriteMessageEnd() error {
	return nil
}

func (p *tBinaryProtocol) WriteStructBegin(h TStructHeader) error {
	return nil
}

func (p *tBinaryProtocol) WriteStructEnd() error {
	return nil
}

func (p *tBinaryProtocol) WriteFieldBegin(h TFieldHeader) (err error) {
	if err = p.WriteByte(h.Type); err == nil {
		err = p.WriteI16(h.Identity)
	}
	return
}

func (p *tBinaryProtocol) WriteFieldEnd() error {
	return nil
}

func (p *tBinaryProtocol) WriteFieldStop() error {
	return p.WriteByte(STOP)
}

func (p *tBinaryProtocol) WriteMapBegin(h TMapHeader) (err error) {
	if err = p.WriteByte(h.Key); err == nil {
		if err = p.WriteByte(h.Value); err == nil {
			err = p.writeSize(h.Size)
		}
	}
	return
}

func (p *tBinaryProtocol) WriteMapEnd() error {
	return nil
}

func (p *tBinaryProtocol) WriteSetBegin(h TSetHeader) (err error) {
	if err = p.WriteByte(h.Element); err == nil {
		err = p.writeSize(h.Size)
	}
	return
}

func (p *tBinaryProtocol) WriteSetEnd() error {
	return nil
}

func (p *tBinaryProtocol) WriteListBegin(h TListHeader) (err error) {
	if err = p.WriteByte(h.Element); err == nil {
		err = p.writeSize(h.Size)
	}
	return
}

func (p *tBinaryProtocol) WriteListEnd() error {
	return nil
}

func (p *tBinaryProtocol) WriteBool(v bool) error {
	if v {
		return p.WriteByte(1)
	}
	return p.WriteByte(0)
}

func (p *tBinaryProtocol) WriteByte(v byte) error {
	return NewTProtocolExceptionFromError(p.TExtraTransport.WriteByte(v))
}

func (p *tBinaryProtocol) WriteDouble(v float64) error {
	return p.WriteU64(math.Float64bits(v))
}

func (p *tBinaryProtocol) WriteU16(v uint16) (err error) {
	buf := p.buf[:2]
	binary.BigEndian.PutUint16(buf, v)
	_, err = p.Write(buf)
	return
}

func (p *tBinaryProtocol) WriteI16(v int16) error {
	return p.WriteU16(uint16(v))
}

func (p *tBinaryProtocol) WriteU32(v uint32) (err error) {
	buf := p.buf[:4]
	binary.BigEndian.PutUint32(buf, v)
	_, err = p.Write(buf)
	return
}

func (p *tBinaryProtocol) WriteI32(v int32) error {
	return p.WriteU32(uint32(v))
}

func (p *tBinaryProtocol) WriteU64(v uint64) (err error) {
	buf := p.buf[:8]
	binary.BigEndian.PutUint64(buf, v)
	_, err = p.Write(buf)
	return
}

func (p *tBinaryProtocol) WriteI64(v int64) error {
	return p.WriteU64(uint64(v))
}

func (p *tBinaryProtocol) WriteString(v string) error {
	return p.WriteBinary([]byte(v))
}

func (p *tBinaryProtocol) WriteBinary(v []byte) (err error) {
	if err = p.writeSize(len(v)); err == nil {
		_, err = p.Write(v)
	}
	return
}

func (p *tBinaryProtocol) Write(v []byte) (int, error) {
	n, err := p.TExtraTransport.Write(v)
	return n, NewTProtocolExceptionFromError(err)
}

func (p *tBinaryProtocol) writeSize(v int) error {
	return p.WriteU32(uint32(v))
}

func (p *tBinaryProtocol) ReadMessageBegin() (h TMessageHeader, err error) {
	// TODO implement
	return
}

func (p *tBinaryProtocol) ReadMessageEnd() error {
	return nil
}

func (p *tBinaryProtocol) ReadStructBegin() (h TStructHeader, err error) {
	return
}

func (p *tBinaryProtocol) ReadStructEnd() error {
	return nil
}

func (p *tBinaryProtocol) ReadFieldBegin() (h TFieldHeader, err error) {
	if h.Type, err = p.ReadByte(); err == nil && h.Type != STOP {
		h.Identity, err = p.ReadI16()
	}
	return
}

func (p *tBinaryProtocol) ReadFieldEnd() error {
	return nil
}

func (p *tBinaryProtocol) ReadMapBegin() (h TMapHeader, err error) {
	if h.Key, err = p.ReadByte(); err == nil {
		if h.Value, err = p.ReadByte(); err == nil {
			h.Size, err = p.readSize()
		}
	}
	return
}

func (p *tBinaryProtocol) ReadMapEnd() error {
	return nil
}

func (p *tBinaryProtocol) ReadSetBegin() (h TSetHeader, err error) {
	if h.Element, err = p.ReadByte(); err == nil {
		h.Size, err = p.readSize()
	}
	return
}

func (p *tBinaryProtocol) ReadSetEnd() error {
	return nil
}

func (p *tBinaryProtocol) ReadListBegin() (h TListHeader, err error) {
	if h.Element, err = p.ReadByte(); err == nil {
		h.Size, err = p.readSize()
	}
	return
}

func (p *tBinaryProtocol) ReadListEnd() error {
	return nil
}

func (p *tBinaryProtocol) ReadBool() (bool, error) {
	v, err := p.ReadByte()
	return v == 1, err
}

func (p *tBinaryProtocol) ReadByte() (byte, error) {
	v, err := p.TExtraTransport.ReadByte()
	return v, NewTProtocolExceptionFromError(err)
}

func (p *tBinaryProtocol) ReadDouble() (float64, error) {
	v, err := p.ReadU64()
	return math.Float64frombits(v), err
}

func (p *tBinaryProtocol) ReadU16() (uint16, error) {
	buf := p.buf[:2]
	_, err := p.Read(buf)
	return binary.BigEndian.Uint16(buf), err
}

func (p *tBinaryProtocol) ReadI16() (int16, error) {
	v, err := p.ReadU16()
	return int16(v), err
}

func (p *tBinaryProtocol) ReadU32() (uint32, error) {
	buf := p.buf[:4]
	_, err := p.Read(buf)
	return binary.BigEndian.Uint32(buf), err
}

func (p *tBinaryProtocol) ReadI32() (int32, error) {
	v, err := p.ReadU32()
	return int32(v), err
}

func (p *tBinaryProtocol) ReadU64() (uint64, error) {
	buf := p.buf[:8]
	_, err := p.Read(buf)
	return binary.BigEndian.Uint64(buf), err
}

func (p *tBinaryProtocol) ReadI64() (int64, error) {
	v, err := p.ReadU64()
	return int64(v), err
}

func (p *tBinaryProtocol) ReadString() (v string, err error) {
	var n int
	if n, err = p.readSize(); err == nil {
		v, err = p.readStringBody(n)
	}
	return
}

func (p *tBinaryProtocol) ReadBinary() (v []byte, err error) {
	var n int
	if n, err = p.readSize(); err == nil {
		v = make([]byte, n)
		_, err = p.Read(v)
	}
	return
}

func (p *tBinaryProtocol) Read(v []byte) (int, error) {
	n, err := p.TExtraTransport.Read(v)
	return n, NewTProtocolExceptionFromError(err)
}

func (p *tBinaryProtocol) readSize() (int, error) {
	v, err := p.ReadU32()
	return int(v), err
}

func (p *tBinaryProtocol) readStringBody(n int) (string, error) {
	v := make([]byte, n)
	_, err := p.Read(v)
	return string(v), err
}

func (p *tBinaryProtocol) Skip(v TType) error {
	return Skip(v, p)
}

func (p *tBinaryProtocol) Flush(ctx context.Context) error {
	return NewTProtocolExceptionFromError(p.TExtraTransport.Flush(ctx))
}
