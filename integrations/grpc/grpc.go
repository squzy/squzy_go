package squzy_grpc

import (
	"context"
	api "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/squzy/squzy_go/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
)

func NewClientUnaryInterceptor(app *core.Application) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		trx := app.CreateTransaction(method, api.TransactionType_TRANSACTION_TYPE_GRPC, core.GetTransactionFromContext(ctx))
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		md.Set(core.CONTEXT_KEY, trx.Id)
		err := invoker(metadata.NewOutgoingContext(ctx, md), method, req, reply, cc, opts...)
		trx.SetMeta().End()
		return err
	}
}

type streamWrap struct {
	grpc.ClientStream
	trx *core.Transaction
}

func (s streamWrap) RecvMsg(m interface{}) error {
	err := s.ClientStream.RecvMsg(m)
	if err == io.EOF {
		s.trx.End()
	}
	return err
}

func NewStreamUnaryInterceptor(app *core.Application) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		trx := app.CreateTransaction(method, api.TransactionType_TRANSACTION_TYPE_GRPC, core.GetTransactionFromContext(ctx))
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		trx.SetMeta()
		md.Set(core.CONTEXT_KEY, trx.Id)
		s, err := streamer(metadata.NewOutgoingContext(ctx, md), desc, cc, method, opts...)
		if err != nil {
			return s, err
		}
		return &streamWrap{
			ClientStream: s,
			trx: trx,
		}, nil
	}
}