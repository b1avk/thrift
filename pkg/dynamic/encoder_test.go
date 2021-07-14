package dynamic_test

import (
	"reflect"
	"testing"

	"github.com/b1avk/thrift/pkg/dynamic"
	"github.com/b1avk/thrift/pkg/thrift"
)

type SimpleTestCase struct {
	name  string
	value interface{}
}

var SimpleTestCases = []SimpleTestCase{
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
}

func testBasicValue(t *testing.T, p thrift.TProtocol) {
	for _, c := range SimpleTestCases {
		t.Run(c.name, func(t *testing.T) {
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
	testBasicValue(t, thrift.NewTBinaryProtocol(thrift.NewTMemoryBuffer()))
}
