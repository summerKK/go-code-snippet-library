package logic

import (
	"context"
	"github.com/summerKK/mall/service/user/rpc/userclient"

	"github.com/summerKK/mall/service/user/api/internal/svc"
	"github.com/summerKK/mall/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) RegisterLogic {
	return RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	res, err := l.svcCtx.UserRpc.Register(l.ctx, &userclient.RegisterRequest{
		Name:     req.Name,
		Gender:   req.Gender,
		Mobile:   req.Mobile,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &types.RegisterResponse{
		Id:     res.Id,
		Name:   res.Name,
		Gender: res.Gender,
		Mobile: res.Mobile,
	}, nil
}
