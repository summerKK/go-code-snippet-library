package server

import (
	"context"
	"encoding/json"

	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/pkg/api"
	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/pkg/errcode"
	pb "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/proto"
)

type TagServer struct {
}

func (t *TagServer) GetTagList(ctx context.Context, request *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	api := api.NewApi("http://127.0.0.1:8000")
	list, err := api.GetTagList(ctx, request.GetName())
	if err != nil {
		return nil, errcode.TogRPCError(errcode.ErrorGetTagListFail)
	}

	tagList := &pb.GetTagListReply{}
	err = json.Unmarshal(list, &tagList)
	if err != nil {
		return nil, err
	}

	return tagList, nil
}

func NewTagServer() *TagServer {
	return &TagServer{}
}
