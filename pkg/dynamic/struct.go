package dynamic

import (
	"reflect"

	"github.com/b1avk/thrift/pkg/thrift"
)

type TStruct struct {
	typ     reflect.Type
	value   reflect.Value
	encoder InternalEncoder
}

func NewTStruct(v reflect.Type) *TStruct {
	mustBe(v, reflect.Struct)
	return &TStruct{
		typ:     v,
		encoder: InternalEncoderOf(v),
	}
}

func (e *TStruct) New() {
	e.value.Set(reflect.New(e.typ).Elem())
}

func (e *TStruct) Copy() *TStruct {
	return &TStruct{e.typ, reflect.Value{}, e.encoder}
}

func (e *TStruct) Write(p thrift.TProtocol) error {
	return e.encoder.Encode(e.value, p)
}

func (e *TStruct) Read(p thrift.TProtocol) error {
	e.New()
	return e.encoder.Decode(e.value, p)
}

func (e *TStruct) Value() reflect.Value {
	return e.value
}
