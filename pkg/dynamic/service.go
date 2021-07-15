package dynamic

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
var errorType = reflect.TypeOf((*error)(nil)).Elem()

type dynamicField struct {
	splited []string

	method       string
	hasContext   bool
	returnError  bool
	args, result *TStruct
}

func parseDynamicField(t reflect.StructField) (f dynamicField, err error) {
	splited := strings.Split(t.Tag.Get("thrift"), " ")
	f.method = splited[0]
	f.splited = splited[1:]
	ft := t.Type
	si := []reflect.StructField{}
	ni := ft.NumIn()
	var id int
	for i := 0; i < ni; i++ {
		ti := ft.In(i)
		if i == 0 && ti.AssignableTo(contextType) {
			f.hasContext = true
			continue
		}
		if id, err = f.nextIdentity(); err != nil {
			return
		}
		si = append(si, reflect.StructField{
			Name: "F" + strconv.Itoa(i),
			Type: ti,
			Tag:  reflect.StructTag(fmt.Sprintf(`thrift:"%v"`, id)),
		})
	}
	if len(si) != 0 {
		f.args = NewTStruct(reflect.StructOf(si))
	}
	no := ft.NumOut()
	if no != 0 {
		so := []reflect.StructField{}
		for i := 0; i < no; i++ {
			to := ft.Out(i)
			if i == no || to.AssignableTo(errorType) {
				f.returnError = true
				continue
			}
			if id, err = f.nextIdentity(); err != nil {
				return
			}
			so = append(so, reflect.StructField{
				Name: "F" + strconv.Itoa(i),
				Type: to,
				Tag:  reflect.StructTag(fmt.Sprintf(`thrift:"%v"`, id)),
			})
		}
		if len(so) != 0 {
			f.result = NewTStruct(reflect.StructOf(so))
		}
	}
	return
}

func (f dynamicField) newArgs(args []reflect.Value) *TStruct {
	v := f.args.Copy()
	v.New()
	for i := range f.args.encoder.fieldEncoderLists {
		v.value.Field(i).Set(args[0])
		args = args[1:]
	}
	return v
}

func (f dynamicField) newResult() *TStruct {
	return f.result.Copy()
}

func (f dynamicField) returnResult(res *TStruct, err error) (results []reflect.Value) {
	if res != nil {
		for i := range res.encoder.fieldEncoderLists {
			results = append(results, res.value.Field(i))
		}
	}
	if f.returnError {
		results = append(results, reflect.ValueOf(&err).Elem())
	}
	return
}

func (f *dynamicField) nextIdentity() (v int, err error) {
	if len(f.splited) == 0 {
		err = fmt.Errorf("no splited left")
	} else {
		v, err = strconv.Atoi(f.splited[0])
		f.splited = f.splited[1:]
	}
	return
}
