package server

import (
	"context"

	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/pkg/errcode"
	"google.golang.org/grpc/metadata"
)

type Auth struct {
}

func (a *Auth) GetAppKey() string {
	return "summer"
}

func (a *Auth) GetAppSecret() string {
	return "summer"
}

func (a *Auth) Check(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)
	var appKey, appSecret string
	if v, ok := md["app_key"]; ok {
		appKey = v[0]
	}

	if v, ok := md["app_secret"]; ok {
		appSecret = v[0]
	}

	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
		return errcode.TogRPCError(errcode.Unauthorized)
	}

	return nil
}
