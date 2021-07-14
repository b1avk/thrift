package dynamic_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/b1avk/thrift/pkg/dynamic"
	"github.com/b1avk/thrift/pkg/thrift"
)

type GreetArgs struct {
	Name string `thrift:"1"`
}

type GreetResult struct {
	Text string `thrift:"0"`
}

type GreeterService struct {
	Greet          func(name string) string                               `thrift:"greet 1 0"`
	GreetRetErr    func(name string) (string, error)                      `thrift:"greet 1 0"`
	GreetCtx       func(ctx context.Context, name string) string          `thrift:"greet 1 0"`
	GreetCtxRetErr func(ctx context.Context, name string) (string, error) `thrift:"greet 1 0"`
}

type FakeClient struct{}

func (*FakeClient) Call(ctx context.Context, method string, args, result thrift.TStruct) (err error) {
	b := thrift.NewTMemoryBuffer()
	p := thrift.NewTBinaryProtocol(b, nil)
	if err = args.Write(p); err != nil {
		return
	}
	av := &GreetArgs{}
	ae := dynamic.ValueEncoderOf(reflect.TypeOf(av))
	if err = ae.Decode(&av, p); err != nil {
		return
	}
	rv := &GreetResult{
		Text: fmt.Sprintf("Hello %s !", av.Name),
	}
	re := dynamic.ValueEncoderOf(reflect.TypeOf(rv))
	if err = re.Encode(rv, p); err != nil {
		return
	}
	err = result.Read(p)
	return
}

func TestDynamicGreeterServiceGreet(t *testing.T) {
	s := dynamic.WrapServiceClient(new(GreeterService), new(FakeClient)).(*GreeterService)
	if s.Greet("World") != "Hello World !" {
		t.Fatal(`Greet("World") must returns "Hello World !"`)
	}
}

func TestDynamicGreeterServiceGreetRetErr(t *testing.T) {
	s := dynamic.WrapServiceClient(new(GreeterService), new(FakeClient)).(*GreeterService)
	if res, err := s.GreetRetErr("World"); !(res == "Hello World !" && err == nil) {
		t.Fatal(`GreetRetErr("World") must returns ("Hello World !", nil)`)
	}
}

func TestDynamicGreeterServiceGreetCtx(t *testing.T) {
	s := dynamic.WrapServiceClient(new(GreeterService), new(FakeClient)).(*GreeterService)
	if s.GreetCtx(context.Background(), "World") != "Hello World !" {
		t.Fatal(`GreetCtx(ctx, "World") must returns "Hello World !"`)
	}
}

func TestDynamicGreeterServiceGreetCtxRetErr(t *testing.T) {
	s := dynamic.WrapServiceClient(new(GreeterService), new(FakeClient)).(*GreeterService)
	if res, err := s.GreetCtxRetErr(context.Background(), "World"); !(res == "Hello World !" && err == nil) {
		t.Fatal(`GreetCtxRetErr(ctx, "World") must returns ("Hello World !", nil)`)
	}
}
