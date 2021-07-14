package dynamic

import (
	"context"
	"reflect"

	"github.com/b1avk/thrift/pkg/thrift"
)

// WrapServiceClient wraps s into c.
func WrapServiceClient(s interface{}, c thrift.TClient) interface{} {
	sv := reflect.ValueOf(s)
	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	st := sv.Type()
	if st.Kind() != reflect.Struct {
		panic("dynamic.WrapServiceClient: service must be struct")
	}
	n := st.NumField()
	for i := 0; i < n; i++ {
		if m, ok := makeClientMethod(c, st.Field(i)); ok {
			sv.Field(i).Set(m)
		}
	}
	return s
}

func makeClientMethod(c thrift.TClient, v reflect.StructField) (reflect.Value, bool) {
	f, err := parseDynamicField(v)
	if err != nil {
		return reflect.Value{}, false
	}
	ctx := context.Background()
	return reflect.MakeFunc(v.Type, func(args []reflect.Value) (results []reflect.Value) {
		if f.hasContext {
			ctx = args[0].Interface().(context.Context)
			args = args[1:]
		}
		res := f.newResult()
		err := c.Call(ctx, f.method, f.newArgs(args), res)
		return f.returnResult(res, err)
	}), true
}
