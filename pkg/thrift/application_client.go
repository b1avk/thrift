package thrift

import (
	"context"
	"fmt"
	"sync"
)

type TClient interface {
	Call(ctx context.Context, method string, args, result TStruct) (err error)
}

type TStandardClient struct {
	iprot, oprot TProtocol
	message      TMessageHeader
	mutex        sync.Mutex
}

func NewTStandardClient(iprot, oprot TProtocol) *TStandardClient {
	if iprot == nil || oprot == nil {
		switch {
		case iprot == nil:
			iprot = oprot
		case oprot == nil:
			oprot = iprot
		default:
			panic("thrift.NewTStandardClient: iprot or oprot must be non-nil")
		}
	}
	return &TStandardClient{
		iprot: iprot,
		oprot: oprot,
		message: TMessageHeader{
			Type: CALL,
		},
	}
}

func (p *TStandardClient) Call(ctx context.Context, method string, args, result TStruct) (err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.message.Identity++
	p.message.Name = method
	if err = p.oprot.WriteMessageBegin(p.message); err != nil {
		return
	}
	if err = args.Write(p.oprot); err != nil {
		return
	}
	if err = p.oprot.WriteMapEnd(); err != nil {
		return
	}
	if err = p.oprot.Flush(ctx); err != nil {
		return
	}
	if result == nil {
		return
	}
	var h TMessageHeader
	if h, err = p.iprot.ReadMessageBegin(); err != nil {
		return
	}
	switch {
	case h.Name != method:
		err = &TApplicationException{
			Type:    TApplicationErrorWrongMethodName,
			Message: fmt.Sprintf("%s: wrong method name", method),
		}
	case h.Identity != p.message.Identity:
		err = &TApplicationException{
			Type:    TApplicationErrorBadSequenceID,
			Message: fmt.Sprintf("%s: out of order sequence response", method),
		}
	case h.Type == EXCEPTION:
		var e TApplicationException
		if err = e.Read(p.iprot); err != nil {
			return
		}
		if err = p.iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = &e
	case h.Type == REPLY:
		if err = result.Read(p.iprot); err != nil {
			return
		}
		err = p.iprot.ReadMessageEnd()
	default:
		err = &TApplicationException{
			Type:    TApplicationErrorInvalidMessageType,
			Message: fmt.Sprintf("%s: invalid message type", method),
		}
	}
	return
}

type TPoolClient struct {
	itrans, otrans TTransportFactory
	iprot, oprot   TProtocolFactory
	sequence       int32
	mutex          sync.Mutex
	pool           sync.Pool
}

func NewTPoolClient(itrans, otrans TTransportFactory, iprot, oprot TProtocolFactory) *TPoolClient {
	if itrans == nil || otrans == nil {
		switch {
		case itrans == nil:
			itrans = otrans
		case oprot == nil:
			otrans = itrans
		default:
			panic("thrift.NewTPoolClient: itrans or otrans must be non-nil")
		}
	}
	if iprot == nil || oprot == nil {
		switch {
		case iprot == nil:
			iprot = oprot
		case oprot == nil:
			oprot = iprot
		default:
			panic("thrift.NewTPoolClient: iprot or oprot must be non-nil")
		}
	}
	return &TPoolClient{
		itrans: itrans,
		otrans: otrans,
		iprot:  iprot,
		oprot:  oprot,
	}
}

func (cp *TPoolClient) Call(ctx context.Context, method string, args, result TStruct) error {
	c, ok := cp.pool.Get().(*TStandardClient)
	if !ok {
		itrans, err := cp.itrans.GetTransport(nil)
		if err != nil {
			return err
		}
		otrans, err := cp.otrans.GetTransport(nil)
		if err != nil {
			return err
		}
		c = NewTStandardClient(cp.iprot.GetProtocol(itrans), cp.oprot.GetProtocol(otrans))
	}
	cp.mutex.Lock()
	cp.sequence++
	cp.mutex.Unlock()
	defer cp.pool.Put(c)
	c.message.Identity = cp.sequence
	return c.Call(ctx, method, args, result)
}
