package main

import (
	"context"
	"log"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/global"
	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/internal/middleware"
	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/pkg/tracer"
	pb "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type Auth struct {
	AppKey    string
	AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"app_key":    a.AppKey,
		"app_secret": a.AppSecret,
	}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return false
}

func init() {
	err := setUpTracer()
	if err != nil {
		log.Printf("setUpTracer error:%v", err)
	}
}

func main() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	md := metadata.New(map[string]string{"go": "hello,world"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	auth := &Auth{
		AppKey:    "summer",
		AppSecret: "summer",
	}

	options := []grpc.DialOption{
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				// 设置超时时间
				middleware.UnaryContextTimeout(),
				// 链路追踪
				middleware.ClientTracing(),
				// 设置重试
				grpc_retry.UnaryClientInterceptor(
					grpc_retry.WithMax(2),
					grpc_retry.WithCodes(
						codes.Internal,
						codes.Unknown,
						codes.DeadlineExceeded,
					),
				),
			),
		),
		// grpc认证
		grpc.WithPerRPCCredentials(auth),
	}

	clientConn, err := GetClientConn(ctx, ":8001", options)

	if err != nil {
		panic(err)
	}
	defer clientConn.Close()

	tagServiceClient := pb.NewTagServiceClient(clientConn)
	list, err := tagServiceClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Go"})
	if err != nil {
		log.Fatalf("tagServiceClient.GetTagList err:%v", err)
	}
	log.Printf("getTagList:%+v", list)
}

func GetClientConn(ctx context.Context, target string, options []grpc.DialOption) (*grpc.ClientConn, error) {
	options = append(options, grpc.WithInsecure())
	return grpc.DialContext(ctx, target, options...)
}

func setUpTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer("rpc-blog-service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer

	return nil
}
