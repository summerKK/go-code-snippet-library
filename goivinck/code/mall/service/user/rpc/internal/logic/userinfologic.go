package logic

import (
	"context"
	"github.com/summerKK/mall/service/user/model"
	"google.golang.org/grpc/status"

	"github.com/summerKK/mall/service/user/rpc/internal/svc"
	"github.com/summerKK/mall/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *userclient.UserInfoRequest) (*userclient.UserInfoResponse, error) {
	//  查询用户是否存在
	user, err := l.svcCtx.UserModel.FindOne(in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(100, "用户不存在")
		}
		return nil, status.Error(500, err.Error())
	}

	return &userclient.UserInfoResponse{
		Id:     user.Id,
		Name:   user.Name,
		Gender: user.Gender,
		Mobile: user.Mobile,
	}, nil
}
