package main

import (
	"context"
	"log"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/internal/middleware"
	pb "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func main() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	md := metadata.New(map[string]string{"go": "hello,world"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	options := []grpc.DialOption{
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				// 设置超时时间
				middleware.UnaryContextTimeout(),
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
