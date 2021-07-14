package thrift

import (
	"context"
	"fmt"
)

type TClient interface {
	Call(ctx context.Context, method string, args, result TStruct) (err error)
}

type TStandardClient struct {
	iprot, oprot TProtocol
	message      TMessageHeader
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
	p.message.Identity++
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
