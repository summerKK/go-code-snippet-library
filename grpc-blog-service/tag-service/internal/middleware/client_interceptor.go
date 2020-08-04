package middleware

import (
	"context"
	"log"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/global"
	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/metatext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func ClientTracing() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var parentCtx opentracing.SpanContext
		var spanOpts []opentracing.StartSpanOption
		parentSpan := opentracing.SpanFromContext(ctx)
		// 检查是否包含上级的跨度信息
		if parentSpan != nil {
			parentCtx = parentSpan.Context()
			spanOpts = append(spanOpts, opentracing.ChildOf(parentCtx))
		}
		spanOpts = append(spanOpts, []opentracing.StartSpanOption{
			opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
			ext.SpanKindRPCClient,
		}...)

		span := global.Tracer.StartSpan(method, spanOpts...)
		defer span.Finish()

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		_ = global.Tracer.Inject(span.Context(), opentracing.TextMap, metatext.MetadataTextMap{MD: md})
		newCtx := opentracing.ContextWithSpan(metadata.NewOutgoingContext(ctx, md), span)

		return invoker(newCtx, method, req, reply, cc, opts...)
	}
}
