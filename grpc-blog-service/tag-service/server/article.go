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

type ArticleServer struct {
	auth *Auth
}

func NewArticleServer() *ArticleServer {
	return &ArticleServer{}
}

func (t *ArticleServer) GetArticleList(ctx context.Context, request *pb.GetArticleListRequest) (*pb.GetArticleListReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	log.Printf("receive metadata from context:%v,ok:%v", md, ok)

	if err := t.auth.Check(ctx); err != nil {
		return nil, err
	}

	apiService := api.NewApi("http://grpc-summer.cc:8000")
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
