package logic

import (
	"context"
	"github.com/summerKK/mall/common/cryptx"
	"github.com/summerKK/mall/service/user/model"
	"github.com/summerKK/mall/service/user/rpc/user"
	"google.golang.org/grpc/status"

	"github.com/summerKK/mall/service/user/rpc/internal/svc"
	"github.com/summerKK/mall/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *userclient.RegisterRequest) (*userclient.RegisterResponse, error) {
	// 判断手机是否存在
	_, err := l.svcCtx.UserModel.FindOneByMobile(in.Mobile)
	if err == nil {
		return nil, status.Error(100, "该用户已存在")
	}
	if err == model.ErrNotFound {
		newUser := model.User{
			Name:     in.Name,
			Gender:   in.Gender,
			Mobile:   in.Mobile,
			Password: cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, in.Password),
		}
		res, err := l.svcCtx.UserModel.Insert(&newUser)
		if err != nil {
			return nil, status.Error(500, err.Error())
		}
		newUser.Id, err = res.LastInsertId()
		if err != nil {
			return nil, status.Error(500, err.Error())
		}

		return &user.RegisterResponse{
			Id:     newUser.Id,
			Name:   newUser.Name,
			Gender: newUser.Gender,
			Mobile: newUser.Mobile,
		}, nil
	}

	return nil, status.Error(500, err.Error())
}
