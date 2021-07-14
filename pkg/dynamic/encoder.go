package dynamic

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	"github.com/b1avk/thrift/pkg/thrift"
)

type ValueEncoder struct {
	InternalEncoder
}

func ValueEncoderOf(v reflect.Type) *ValueEncoder {
	if e := getValueEncoderOf(v); e != nil {
		return e
	}
	e := &ValueEncoder{InternalEncoderOf(v)}
	cache.Store(v, e)
	return e
}

func (e *ValueEncoder) Encode(v interface{}, p thrift.TProtocol) error {
	return e.InternalEncoder.Encode(reflect.ValueOf(v), p)
}

func (e *ValueEncoder) Decode(v interface{}, p thrift.TProtocol) error {
	return e.InternalEncoder.Decode(reflect.ValueOf(v).Elem(), p)
}

type InternalEncoder interface {
	Encode(v reflect.Value, p thrift.TProtocol) error
	Decode(v reflect.Value, p thrift.TProtocol) error
	Kind() thrift.TType
}

func InternalEncoderOf(v reflect.Type) (e InternalEncoder) {
	if e := getValueEncoderOf(v); e != nil {
		return e.InternalEncoder
	}
	switch v.Kind() {
	case reflect.Bool:
		e = new(boolEncoder)
	case reflect.Uint8:
		e = new(uint8Encoder)
	case reflect.Int8:
		e = new(int8Encoder)
	case reflect.Float32, reflect.Float64:
		e = new(doubleEncoder)
	case reflect.Uint16:
		e = new(uint16Encoder)
	case reflect.Int16:
		e = new(int16Encoder)
	case reflect.Uint32:
		e = new(uint32Encoder)
	case reflect.Int32:
		e = new(int32Encoder)
	case reflect.Uint64, reflect.Uint:
		e = new(uintEncoder)
	case reflect.Int64, reflect.Int:
		e = new(intEncoder)
	case reflect.String:
		e = new(stringEncoder)
	case reflect.Map:
		keyType := v.Key()
		valueType := v.Elem()
		e = &mapEncoder{
			mapType:      v,
			keyType:      keyType,
			valueType:    valueType,
			keyEncoder:   InternalEncoderOf(keyType),
			valueEncoder: InternalEncoderOf(valueType),
		}
	case reflect.Slice:
		if v.Elem().Kind() == reflect.Uint8 {
			e = new(binaryEncoder)
		} else {
			// TODO set encoder
			e = &listEncoder{
				sliceType:      v,
				elementEncoder: InternalEncoderOf(v.Elem()),
			}
		}
	case reflect.Ptr:
		valueTpye := v.Elem()
		e = &ptrEncoder{
			valueType:       valueTpye,
			InternalEncoder: InternalEncoderOf(valueTpye),
		}
	default:
		panic(fmt.Errorf("unexpected Type: %v", v.Kind()))
	}
	cache.Store(v, e)
	return
}

func getValueEncoderOf(v reflect.Type) (e *ValueEncoder) {
	if e, ok := cache.Load(v); ok {
		if e, ok := e.(*ValueEncoder); ok {
			return e
		}
		if e, ok := e.(InternalEncoder); ok {
			e := &ValueEncoder{e}
			cache.Store(v, e)
			return e
		}
	}
	return nil
}

var cache sync.Map

type boolEncoder struct{}

func (e *boolEncoder) Encode(v reflect.Value, p thrift.TProtocol) error {
	return p.WriteBool(v.Bool())
}

func (e *boolEncoder) Decode(v reflect.Value, p thrift.TProtocol) error {
	res, err := p.ReadBool()
	v.SetBool(res)
	return err
}

func (e *boolEncoder) Kind() thrift.TType {
	return thrift.BOOL
}

type uint8Encoder struct{}

func (e *uint8Encoder) Encode(v reflect.Value, p thrift.TProtocol) error {
	return p.WriteByte(byte(v.Uint()))
}

func (e *uint8Encoder) Decode(v reflect.Value, p thrift.TProtocol) (err error) {
	mustBe(v, reflect.Uint8)
	*(*uint8)(unsafe.Pointer(v.UnsafeAddr())), err = p.ReadByte()
	return
}

func (e *uint8Encoder) Kind() thrift.TType {
	return thrift.BYTE
}

type int8Encoder struct{}

func (e *int8Encoder) Encode(v reflect.Value, p thrift.TProtocol) error {
	return p.WriteByte(byte(v.Int()))
}

func (e *int8Encoder) Decode(v reflect.Value, p thrift.TProtocol) error {
	mustBe(v, reflect.Int8)
	res, err := p.ReadByte()
	*(*int8)(unsafe.Pointer(v.UnsafeAddr())) = int8(res)
	return err
}

func (e *int8Encoder) Kind() thrift.TType {
	return thrift.BYTE
}

type doubleEncoder struct{}

func (e *doubleEncoder) Encode(v reflect.Value, p thrift.TProtocol) error {
	return p.WriteDouble(v.Float())
}

func (e *doubleEncoder) Decode(v reflect.Value, p thrift.TProtocol) error {
	res, err := p.ReadDouble()
	v.SetFloat(res)
	return err
}

func (e *doubleEncoder) Kind() thrift.TType {
	return thrift.DOUBLE
}

type uint16Encoder struct{}

func (e *uint16Encoder) Encode(v reflect.Value, p thrift.TProtocol) error {
	mustBe(v, reflect.Uint16)
	return p.WriteU16(uint16(v.Uint()))
}

func (e *uint16Encoder) Decode(v reflect.Value, p thrift.TProtocol) (err error) {
	mustBe(v, reflect.Uint16)
	*(*uint16)(unsafe.Pointer(v.UnsafeAddr())), err = p.ReadU16()
	return
}

func (e *uint16Encoder) Kind() thrift.TType {
	return thrift.I16
}

type int16Encoder struct{}

func (e *int16Encoder) Encode(v reflect.Value, p thrift.TProtocol) error {
	mustBe(v, reflect.Int16)
	return p.WriteI16(int16(v.Int()))
}

func (e *int16Encoder) Decode(v reflect.Value, p thrift.TProtocol) (err error) {
	mustBe(v, reflect.Int16)
	*(*int16)(unsafe.Pointer(v.UnsafeAddr())), err = p.ReadI16()
	return
}

func (e *int16Encoder) Kind() thrift.TType {
	return thrift.I16
}

type uint32Encoder struct{}

func (e *uint32Encoder) Encode(v reflect.Value, p thrift.TProtocol) error {
	mustBe(v, reflect.Uint32)
	return p.WriteU32(uint32(v.Uint()))
}

func (e *uint32Encoder) Decode(v reflect.Value, p thrift.TProtocol) (err error) {
	mustBe(v, reflect.Uint32)
	*(*uint32)(unsafe.Pointer(v.UnsafeAddr())), err = p.ReadU32()
	return
}

func (e *uint32Encoder) Kind() thrift.TType {
	return thrift.I32
}

type int32Encoder struct{}

func (e *int32Encoder) Encode(v reflect.Value, p thrift.TProtocol) error {
	mustBe(v, reflect.Int32)
	return p.WriteI32(int32(v.Int()))
}

func (e *int32Encoder) Decode(v reflect.Value, p thrift.TProtocol) (err error) {
	mustBe(v, reflect.Int32)
	*(*int32)(unsafe.Pointer(v.UnsafeAddr())), err = p.ReadI32()
	return
}

func (e *int32Encoder) Kind() thrift.TType {
	return thrift.I32
}

type uintEncoder struct{}

func (e *uintEncoder) Encode(v reflect.Value, p thrift.TProtocol) error {
	return p.WriteU64(v.Uint())
}

func (e *uintEncoder) Decode(v reflect.Value, p thrift.TProtocol) (err error) {
	res, err := p.ReadU64()
	v.SetUint(res)
	return
}

func (e *uintEncoder) Kind() thrift.TType {
	return thrift.I64
}

type intEncoder struct{}

func (e *intEncoder) Encode(v reflect.Value, p thrift.TProtocol) error {
	return p.WriteI64(v.Int())
}

func (e *intEncoder) Decode(v reflect.Value, p thrift.TProtocol) (err error) {
	res, err := p.ReadI64()
	v.SetInt(res)
	return
}

func (e *intEncoder) Kind() thrift.TType {
	return thrift.I64
}

type stringEncoder struct{}

func (e *stringEncoder) Encode(v reflect.Value, p thrift.TProtocol) error {
	return p.WriteString(v.String())
}

func (e *stringEncoder) Decode(v reflect.Value, p thrift.TProtocol) error {
	res, err := p.ReadString()
	v.SetString(res)
	return err
}

func (e *stringEncoder) Kind() thrift.TType {
	return thrift.STRING
}

type binaryEncoder struct{}

func (e *binaryEncoder) Encode(v reflect.Value, p thrift.TProtocol) error {
	return p.WriteBinary(v.Bytes())
}

func (e *binaryEncoder) Decode(v reflect.Value, p thrift.TProtocol) error {
	res, err := p.ReadBinary()
	v.SetBytes(res)
	return err
}

func (e *binaryEncoder) Kind() thrift.TType {
	return thrift.STRING
}

type mapEncoder struct {
	mapType, keyType, valueType reflect.Type
	keyEncoder, valueEncoder    InternalEncoder
}

func (e *mapEncoder) Encode(v reflect.Value, p thrift.TProtocol) (err error) {
	l := v.Len()
	h := thrift.TMapHeader{
		Key:   e.keyEncoder.Kind(),
		Value: e.valueEncoder.Kind(),
		Size:  l,
	}
	if err = p.WriteMapBegin(h); err == nil {
		for iter := v.MapRange(); iter.Next(); {
			if err = e.keyEncoder.Encode(iter.Key(), p); err != nil {
				return
			}
			if err = e.valueEncoder.Encode(iter.Value(), p); err != nil {
				return
			}
		}
		err = p.WriteListEnd()
	}
	return
}

func (e *mapEncoder) Decode(v reflect.Value, p thrift.TProtocol) (err error) {
	var h thrift.TMapHeader
	if h, err = p.ReadMapBegin(); err == nil {
		if h.Key != e.keyEncoder.Kind() || h.Value != e.valueEncoder.Kind() {
			for i := 0; i < h.Size; i++ {
				if err = p.Skip(h.Key); err != nil {
					return
				}
				if err = p.Skip(h.Value); err != nil {
					return
				}
			}
		} else {
			if v.IsNil() {
				v.Set(reflect.MakeMap(e.mapType))
			}
			for i := 0; i < h.Size; i++ {
				key := reflect.New(e.keyType).Elem()
				if err = e.keyEncoder.Decode(key, p); err != nil {
					return
				}
				value := reflect.New(e.valueType).Elem()
				if err = e.valueEncoder.Decode(value, p); err != nil {
					return
				}
				v.SetMapIndex(key, value)
			}
		}
		err = p.ReadMapEnd()
	}
	return
}

func (e *mapEncoder) Kind() thrift.TType {
	return thrift.MAP
}

type listEncoder struct {
	sliceType      reflect.Type
	elementEncoder InternalEncoder
}

func (e *listEncoder) Encode(v reflect.Value, p thrift.TProtocol) (err error) {
	l := v.Len()
	h := thrift.TListHeader{
		Element: e.elementEncoder.Kind(),
		Size:    l,
	}
	if err = p.WriteListBegin(h); err == nil {
		for i := 0; i < l; i++ {
			if err = e.elementEncoder.Encode(v.Index(i), p); err != nil {
				return
			}
		}
		err = p.WriteListEnd()
	}
	return
}

func (e *listEncoder) Decode(v reflect.Value, p thrift.TProtocol) (err error) {
	var h thrift.TListHeader
	if h, err = p.ReadListBegin(); err == nil {
		if h.Element != e.elementEncoder.Kind() {
			for i := 0; i < h.Size; i++ {
				if err = p.Skip(h.Element); err != nil {
					return
				}
			}
		} else {
			if h.Size > v.Len() || v.IsNil() {
				v.Set(reflect.MakeSlice(e.sliceType, h.Size, h.Size))
			}
			for i := 0; i < h.Size; i++ {
				if err = e.elementEncoder.Decode(v.Index(i), p); err != nil {
					return
				}
			}
		}
		err = p.ReadListEnd()
	}
	return
}

func (e *listEncoder) Kind() thrift.TType {
	return thrift.LIST
}

type ptrEncoder struct {
	valueType reflect.Type
	InternalEncoder
}

func (e *ptrEncoder) Decode(v reflect.Value, p thrift.TProtocol) error {
	v.Set(reflect.New(e.valueType))
	return e.InternalEncoder.Decode(v.Elem(), p)
}

func mustBe(v reflect.Value, k reflect.Kind) {
	if v.Kind() != k {
		panic(fmt.Sprintf("reflection: value must be %v not %v", k, v.Kind()))
	}
}
