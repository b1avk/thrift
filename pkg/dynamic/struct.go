package dynamic

import (
	"reflect"

	"github.com/b1avk/thrift/pkg/thrift"
)

// TStruct dynamic implementation for thrift.TStruct.
type TStruct struct {
	typ     reflect.Type
	value   reflect.Value
	encoder *structEncoder
}

// NewTStruct returns new TStruct for v.
func NewTStruct(v reflect.Type) *TStruct {
	mustBe(v, reflect.Struct)
	return &TStruct{
		typ:     v,
		encoder: InternalEncoderOf(v).(*structEncoder),
	}
}

// New initial value of e.
func (e *TStruct) New() {
	if !(e == nil || e.typ == nil) {
		e.value.Set(reflect.New(e.typ).Elem())
	}
}

// Copy returns new TStruct with same type and encoder.
func (e *TStruct) Copy() *TStruct {
	if e == nil || e.typ == nil {
		return nil
	}
	return &TStruct{e.typ, reflect.New(e.typ).Elem(), e.encoder}
}

// Write writes e.value to p.
func (e *TStruct) Write(p thrift.TProtocol) error {
	return e.encoder.Encode(e.value, p)
}

// Read reads e.value from p.
func (e *TStruct) Read(p thrift.TProtocol) error {
	e.New()
	return e.encoder.Decode(e.value, p)
}

// Value returns e.value.
func (e *TStruct) Value() reflect.Value {
	return e.value
}
