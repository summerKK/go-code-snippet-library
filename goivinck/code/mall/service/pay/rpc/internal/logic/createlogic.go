package logic

import (
	"context"
	"github.com/summerKK/mall/service/order/rpc/order"
	"github.com/summerKK/mall/service/pay/model"
	"github.com/summerKK/mall/service/pay/rpc/internal/svc"
	"github.com/summerKK/mall/service/pay/rpc/pay"
	"github.com/summerKK/mall/service/user/rpc/user"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *pay.CreateRequest) (*pay.CreateResponse, error) {
	// 检查用户是否存在
	_, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoRequest{Id: in.Uid})
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.OrderRpc.Detail(l.ctx, &order.DetailRequest{Id: in.Oid})
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.PayModel.FindOneByOid(in.Oid)
	if err == nil {
		return nil, status.Error(100, "订单已存在")
	}

	newPay := model.Pay{
		Uid:    in.Uid,
		Oid:    in.Oid,
		Amount: in.Amount,
		Source: 0,
		Status: 0,
	}

	res, err := l.svcCtx.PayModel.Insert(&newPay)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	newPay.Id, err = res.LastInsertId()
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	return &pay.CreateResponse{Id: newPay.Id}, nil
}
