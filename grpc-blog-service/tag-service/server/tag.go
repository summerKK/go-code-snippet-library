package server

import (
	"context"
	"encoding/json"
	"log"

	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/pkg/api"
	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/pkg/errcode"
	pb "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/proto"
	"google.golang.org/grpc/metadata"
)

type TagServer struct {
}

func (t *TagServer) GetTagList(ctx context.Context, request *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	log.Printf("receive metadata from context:%v,ok:%v", md, ok)

	apiService := api.NewApi("http://127.0.0.1:8000")
	list, err := apiService.GetTagList(ctx, request.GetName())
	if err != nil {
		return nil, errcode.TogRPCError(errcode.ErrorGetTagListFail)
	}

	tagList := &pb.GetTagListReply{}
	err = json.Unmarshal(list, &tagList)
	if err != nil {
		return nil, errcode.TogRPCError(errcode.Fail)
	}

	return tagList, nil
}

func NewTagServer() *TagServer {
	return &TagServer{}
}
