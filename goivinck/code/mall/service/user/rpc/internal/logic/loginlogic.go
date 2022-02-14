package logic

import (
	"context"
	"github.com/summerKK/mall/common/cryptx"
	"github.com/summerKK/mall/service/user/model"
	"google.golang.org/grpc/status"

	"github.com/summerKK/mall/service/user/rpc/internal/svc"
	"github.com/summerKK/mall/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *userclient.LoginRequest) (*userclient.LoginResponse, error) {
	// 查询用户是否存在
	user, err := l.svcCtx.UserModel.FindOneByMobile(in.Mobile)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(100, "用户不存在")
		}
		return nil, status.Error(500, err.Error())
	}

	// 判断密码是否正确
	pwd := cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, in.Password)
	if pwd != user.Password {
		return nil, status.Error(100, "密码错误")
	}

	return &userclient.LoginResponse{
		Id:     user.Id,
		Name:   user.Name,
		Gender: user.Gender,
		Mobile: user.Mobile,
	}, nil
}
