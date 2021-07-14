package thrift

import (
	"context"
	"encoding/binary"
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
	return p
}

type tCompactProtocol struct {
	TExtraTransport
	cfg *TConfiguration
	buf [binary.MaxVarintLen64]byte
}

func (p *tCompactProtocol) WriteMessageBegin(h TMessageHeader) (err error) {
	// TODO
	return
}

func (p *tCompactProtocol) WriteMessageEnd() error {
	return nil
}

func (p *tCompactProtocol) WriteStructBegin(h TStructHeader) error {
	return nil
}

func (p *tCompactProtocol) WriteStructEnd() error {
	return nil
}

func (p *tCompactProtocol) WriteFieldBegin(h TFieldHeader) error {
	// TODO
	return nil
}

func (p *tCompactProtocol) WriteFieldEnd() error {
	return nil
}

func (p *tCompactProtocol) WriteFieldStop() error {
	return p.WriteByte(STOP)
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
	// TODO
	return nil
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
	return p.WriteI64(int64(v))
}

func (p *tCompactProtocol) ReadMessageBegin() (h TMessageHeader, err error) {
	// TODO
	return
}

func (p *tCompactProtocol) ReadMessageEnd() error {
	return nil
}

func (p *tCompactProtocol) ReadStructBegin() (h TStructHeader, err error) {
	return
}

func (p *tCompactProtocol) ReadStructEnd() error {
	return nil
}

func (p *tCompactProtocol) ReadFieldBegin() (h TFieldHeader, err error) {
	// TODO
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
	// TODO
	return false, nil
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
	v, err := p.ReadU64()
	return int16(v), err
}

func (p *tCompactProtocol) ReadU32() (uint32, error) {
	v, err := p.ReadU64()
	return uint32(v), err
}

func (p *tCompactProtocol) ReadI32() (int32, error) {
	v, err := p.ReadU64()
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
	v, err := p.ReadU32()
	return int(v), err
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
