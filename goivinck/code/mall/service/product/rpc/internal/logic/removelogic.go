package logic

import (
	"context"
	"github.com/summerKK/mall/service/product/model"
	"google.golang.org/grpc/status"

	"github.com/summerKK/mall/service/product/rpc/internal/svc"
	"github.com/summerKK/mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveLogic {
	return &RemoveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveLogic) Remove(in *product.RemoveRequest) (*product.RemoveResponse, error) {
	_, err := l.svcCtx.ProductModel.FindOne(in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(100, "产品不存在")
		}

		return nil, status.Error(500, err.Error())
	}

	err = l.svcCtx.ProductModel.Delete(in.Id)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	return &product.RemoveResponse{}, nil
}
