package dynamic_test

import (
	"reflect"
	"testing"

	"github.com/b1avk/thrift/pkg/dynamic"
	"github.com/b1avk/thrift/pkg/thrift"
)

type BasicTestCase struct {
	name  string
	value interface{}
}

var BasicTestCases = []BasicTestCase{
	{
		name:  "BooleanTrue",
		value: true,
	},
	{
		name:  "BooleanFalse",
		value: false,
	},
	{
		name:  "Uint8",
		value: uint8(255),
	},
	{
		name:  "Int8",
		value: int8(-128),
	},
	{
		name:  "Float32",
		value: float32(0.123),
	},
	{
		name:  "Float64",
		value: float64(0.321),
	},
	{
		name:  "Uint16",
		value: uint16(255),
	},
	{
		name:  "Int16",
		value: int16(-128),
	},
	{
		name:  "Uint32",
		value: uint32(255),
	},
	{
		name:  "Int32",
		value: int32(-128),
	},
	{
		name:  "Uint64",
		value: uint64(255),
	},
	{
		name:  "Int64",
		value: int64(-128),
	},
	{
		name:  "Uint",
		value: uint(255),
	},
	{
		name:  "Int",
		value: int(-128),
	},
	{
		name:  "String",
		value: "Hello World",
	},
	{
		name:  "Binary",
		value: []byte("Hello World"),
	},
	{
		name:  "Slice",
		value: []string{"Is", "This", "World", "Or", "Mars", "?"},
	},
	{
		name:  "Map",
		value: map[string]int{"Hello": 1, "Hi": 2},
	},
	{
		name:  "BooleanTruePtr",
		value: toPTR(true).(*bool),
	},
	{
		name:  "BooleanFalsePtr",
		value: toPTR(false).(*bool),
	},
	{
		name:  "StringPtr",
		value: toPTR("Hello World").(*string),
	},
	{
		name: "Struct",
		value: BasicStruct{
			BooleanTrue:  true,
			BooleanFalse: false,
			Double:       0.123,
			String:       "Hello Mars",
			Enum:         BasicEnum(123),
			Nested: &BasicStruct{
				BooleanTrue:  true,
				BooleanFalse: false,
				Double:       0.321,
				String:       "Hello World",
				Nested: &BasicStruct{
					Enum: BasicEnum(321),
				},
			},
		},
	},
}

type BasicEnum int32

type BasicStruct struct {
	BooleanTrue  bool         `thrift:"0"`
	BooleanFalse bool         `thrift:"1"`
	OptionaBool  bool         `thrift:"2,optional"`
	Double       float64      `thrift:"3"`
	String       string       `thrift:"4"`
	Set          []string     `thrift:"6,optional,set"`
	Enum         BasicEnum    `thrift:"7"`
	Nested       *BasicStruct `thrift:"8"`
}

type GetProtocol func() thrift.TProtocol

func TestSetEncoder(t *testing.T) {
	e := dynamic.InternalEncoderOf(reflect.TypeOf((*BasicStruct)(nil)).Elem())
	if e, ok := e.(interface {
		FieldHeader() map[int]thrift.TFieldHeader
	}); ok {
		if !(e.FieldHeader()[5].Type == thrift.SET) {
			t.Fatal("field index 5 must be SET")
		}
	} else {
		t.Fatal("invalid encoder")
	}
}

func testBasicValue(t *testing.T, getProtocol GetProtocol) {
	for _, c := range BasicTestCases {
		t.Run(c.name, func(t *testing.T) {
			p := getProtocol()
			vt := reflect.TypeOf(c.value)
			e := dynamic.ValueEncoderOf(vt)
			if err := e.Encode(c.value, p); err != nil {
				t.Fatal(err)
			}
			rv := reflect.New(vt)
			if err := e.Decode(rv.Interface(), p); err != nil {
				t.Fatal(err)
			}
			r := rv.Elem().Interface()
			if !reflect.DeepEqual(c.value, r) {
				t.Fatal("value obtained for encode and decode mismatch")
			}
		})
	}
}

func TestBasicValueBinaryProtocol(t *testing.T) {
	testBasicValue(t, func() thrift.TProtocol {
		return thrift.NewTBinaryProtocol(thrift.NewTMemoryBuffer(), nil)
	})
}

func TestBasicValueCompactProtocol(t *testing.T) {
	testBasicValue(t, func() thrift.TProtocol {
		return thrift.NewTCompactProtocol(thrift.NewTMemoryBuffer(), nil)
	})
}

func toPTR(s interface{}) interface{} {
	r := reflect.New(reflect.TypeOf(s))
	r.Elem().Set(reflect.ValueOf(s))
	return r.Interface()
}
