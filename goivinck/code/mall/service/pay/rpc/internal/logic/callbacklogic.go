package logic

import (
	"context"
	"github.com/summerKK/mall/service/order/rpc/order"
	"github.com/summerKK/mall/service/pay/model"
	"github.com/summerKK/mall/service/user/rpc/user"
	"google.golang.org/grpc/status"

	"github.com/summerKK/mall/service/pay/rpc/internal/svc"
	"github.com/summerKK/mall/service/pay/rpc/pay"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackLogic {
	return &CallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CallbackLogic) Callback(in *pay.CallbackRequest) (*pay.CallbackResponse, error) {
	// 查询用户是否存在
	_, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoRequest{Id: in.Uid})
	if err != nil {
		return nil, err
	}

	// 查询订单是否存在
	_, err = l.svcCtx.OrderRpc.Detail(l.ctx, &order.DetailRequest{Id: in.Oid})
	if err != nil {
		return nil, err
	}

	res, err := l.svcCtx.PayModel.FindOne(in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(100, "订单不存在")
		}
		return nil, status.Error(500, err.Error())
	}
	// 支付金额与订单金额不符
	if in.Amount != res.Amount {
		return nil, status.Error(100, "支付金额与订单金额不符")
	}

	res.Source = in.Source
	res.Status = in.Status

	// 更新支付信息
	err = l.svcCtx.PayModel.Update(res)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	// 更新订单状态
	_, err = l.svcCtx.OrderRpc.Paid(l.ctx, &order.PaidRequest{Id: in.Oid})
	if err != nil {
		return nil, err
	}

	return &pay.CallbackResponse{}, nil
}
