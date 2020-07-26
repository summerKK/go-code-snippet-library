package main

import (
	"context"
	"log"

	pb "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/proto"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	clientConn, err := GetClientConn(ctx, ":8001", nil)
	if err != nil {
		panic(err)
	}
	defer clientConn.Close()

	tagServiceClient := pb.NewTagServiceClient(clientConn)
	list, err := tagServiceClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Go"})
	if err != nil {
		panic(err)
	}
	log.Printf("getTagList:%+v", list)
}

func GetClientConn(ctx context.Context, target string, options []grpc.DialOption) (*grpc.ClientConn, error) {
	options = append(options, grpc.WithInsecure())

	return grpc.DialContext(ctx, target, options...)
}
