package squzy_grpc

import (
	"context"
	"fmt"
	api "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/squzy/squzy_go/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
)

func NewClientUnaryInterceptor(app *core.Application) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		trx := app.CreateTransaction(fmt.Sprintf("%s/%s",method, cc.Target()), api.TransactionType_TRANSACTION_TYPE_GRPC, core.GetTransactionFromContext(ctx))
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		md.Set(app.GetTracingHeader(), trx.GetId())
		err := invoker(metadata.NewOutgoingContext(ctx, md), method, req, reply, cc, opts...)
		trx.SetMeta(&core.TransactionMeta{
			Host:   cc.Target(),
			Path:   method,
			Method: "RPC",
		}).End(err)
		return err
	}
}

type clientStreamWrap struct {
	grpc.ClientStream
	trx *core.Transaction
}

type serverStreamWrap struct {
	grpc.ServerStream
	trx *core.Transaction
}

func (s serverStreamWrap) Context() context.Context {
	ctx := s.ServerStream.Context()
	return core.ContextWithTransaction(ctx, s.trx)
}

func (s clientStreamWrap) RecvMsg(m interface{}) error {
	err := s.ClientStream.RecvMsg(m)
	var e error = nil
	if err != io.EOF {
		e = err
	}
	s.trx.End(e)
	return err
}

func NewClientStreamUnaryInterceptor(app *core.Application) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		trx := app.CreateTransaction(fmt.Sprintf("%s/%s",method, cc.Target()), api.TransactionType_TRANSACTION_TYPE_GRPC, core.GetTransactionFromContext(ctx))
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		md.Set(app.GetTracingHeader(), trx.GetId())
		s, err := streamer(metadata.NewOutgoingContext(ctx, md), desc, cc, method, opts...)
		if err != nil {
			return s, err
		}
		trx.SetMeta(&core.TransactionMeta{
			Host:   cc.Target(),
			Path:   method,
			Method: "RPC",
		}).End(err)
		return &clientStreamWrap{
			ClientStream: s,
			trx:          trx,
		}, nil
	}
}

func NewServerStreamInterceptor(app *core.Application) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromOutgoingContext(ss.Context())
		var parent *core.Transaction
		if ok {
			id := md.Get(app.GetTracingHeader())[0]
			if id != "" {
				parent = &core.Transaction{
					Id: id,
				}
			}
		}
		trx := app.CreateTransaction(info.FullMethod, api.TransactionType_TRANSACTION_TYPE_ROUTER, parent)
		err := handler(srv, &serverStreamWrap{
			ss,
			trx,
		})
		trx.SetMeta(&core.TransactionMeta{
			Path:   info.FullMethod,
			Method: "RPC",
			Host:   app.GetHost(),
		}).End(err)
		return err
	}
}

func NewServerUnaryInterceptor(app *core.Application) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromOutgoingContext(ctx)
		var parent *core.Transaction
		if ok {
			id := md.Get(app.GetTracingHeader())[0]
			if id != "" {
				parent = &core.Transaction{
					Id: id,
				}
			}
		}
		trx := app.CreateTransaction(info.FullMethod, api.TransactionType_TRANSACTION_TYPE_ROUTER, parent)

		trx.SetMeta(&core.TransactionMeta{
			Path:   info.FullMethod,
			Method: "RPC",
			Host:   app.GetHost(),
		})
		resp, err := handler(core.ContextWithTransaction(ctx, trx), req)

		trx.End(err)
		return resp, err
	}
}
