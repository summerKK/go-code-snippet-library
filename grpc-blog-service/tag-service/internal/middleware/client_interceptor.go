package middleware

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func defaultContextTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	var cancelFunc context.CancelFunc
	// 检查是否设置了超时时间,如果没有设置超时时间添加默认超时时间
	if _, ok := ctx.Deadline(); !ok {
		defaultTimeout := time.Second * 60
		ctx, cancelFunc = context.WithTimeout(ctx, defaultTimeout)
	}

	return ctx, cancelFunc
}

func UnaryContextTimeout() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, cancelFunc := defaultContextTimeout(ctx)
		if cancelFunc != nil {
			defer log.Printf("timeout! method:%s", method)
			defer cancelFunc()
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func StreamContextTimeout() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx, cancelFunc := defaultContextTimeout(ctx)
		if cancelFunc != nil {
			defer cancelFunc()
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}
