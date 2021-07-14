package thrift_test

import (
	"testing"

	"github.com/b1avk/thrift/pkg/thrift"
)

func TestTBinaryProtocolMessageHeader(t *testing.T) {
	p := thrift.NewTBinaryProtocol(thrift.NewTMemoryBuffer(), &thrift.TConfiguration{
		StrictWrite: true,
		StrictRead:  true,
	})
	if p.WriteMessageBegin(thrift.TMessageHeader{}) != nil {
		t.Fatal("fail to write message header")
	}
	if _, err := p.ReadMessageBegin(); err != nil {
		t.Fatal("fail to read message header", err)
	}
}
