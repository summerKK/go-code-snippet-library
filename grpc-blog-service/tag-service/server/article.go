package server

import (
	"context"
	"encoding/json"

	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/pkg/api"
	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/pkg/errcode"
	pb "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/proto"
)

type ArticleServer struct {
}

func NewArticleServer() *ArticleServer {
	return &ArticleServer{}
}

func (t *ArticleServer) GetArticleList(ctx context.Context, request *pb.GetArticleListRequest) (*pb.GetArticleListReply, error) {
	apiService := api.NewApi("http://127.0.0.1:8000")
	list, err := apiService.GetArticleList(ctx, request.TagId)
	if err != nil {
		return nil, errcode.TogRPCError(errcode.ErrorGetArticleListFail)
	}

	articleList := &pb.GetArticleListReply{}
	err = json.Unmarshal(list, &articleList)
	if err != nil {
		return nil, errcode.TogRPCError(errcode.Fail)
	}

	return articleList, nil
}
