package thrift_test

import (
	"testing"

	"github.com/b1avk/thrift/pkg/thrift"
)

func TestTCompactProtocolMessageHeader(t *testing.T) {
	p := thrift.NewTCompactProtocol(thrift.NewTMemoryBuffer(), nil)
	if err := p.WriteMessageBegin(thrift.TMessageHeader{}); err != nil {
		t.Fatal("fail to write message header", err)
	}
	if _, err := p.ReadMessageBegin(); err != nil {
		t.Fatal("fail to read message header", err)
	}
}
