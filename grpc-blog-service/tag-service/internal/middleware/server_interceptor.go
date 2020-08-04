package middleware

import (
	"context"
	"log"
	"runtime/debug"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/global"
	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/metatext"
	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/pkg/errcode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ErrLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		errLog := "error log: method:%s,code:%v,message:%v,details:%v"
		status := errcode.FromError(err)
		log.Printf(errLog, info.FullMethod, status.Code(), status.Err().Error(), status.Details())
	}

	return resp, err
}

func AccessLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	requestLog := "access request log: method:%s begin_time:%d,request:%v"
	beginTime := time.Now().Local().Unix()
	log.Printf(requestLog, info.FullMethod, beginTime, req)

	resp, err = handler(ctx, req)

	responseLog := "access response log: method:%s end_time:%d,response:%v"
	endTime := time.Now().Local().Unix()
	log.Printf(responseLog, info.FullMethod, endTime, resp)

	return resp, err
}

func Recovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			recoveryLog := "recovery log: method:%s,message:%v \nstack:%s"
			log.Printf(recoveryLog, info.FullMethod, r, string(debug.Stack()[:]))
		}
	}()

	resp, err = handler(ctx, req)

	return resp, err
}

func ServerTracing(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	parentSpanContext, _ := global.Tracer.Extract(opentracing.TextMap, metatext.MetadataTextMap{MD: md})
	spanOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
		ext.SpanKindRPCServer,
		ext.RPCServerOption(parentSpanContext),
	}
	span := global.Tracer.StartSpan(info.FullMethod, spanOptions...)
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)

	return handler(ctx, req)
}
