package logic

import (
	"context"
	"github.com/summerKK/mall/service/order/model"
	"github.com/summerKK/mall/service/order/rpc/internal/svc"
	"github.com/summerKK/mall/service/order/rpc/order"
	"github.com/summerKK/mall/service/product/rpc/product"
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

func (l *CreateLogic) Create(in *order.CreateRequest) (*order.CreateResponse, error) {
	// 检验用户是否存在
	_, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoRequest{Id: in.Uid})
	if err != nil {
		return nil, err
	}

	// 查询产品是否存在
	productRes, err := l.svcCtx.ProductRpc.Detail(l.ctx, &product.DetailRequest{Id: in.Pid})
	if err != nil {
		return nil, err
	}
	if productRes.Stock <= 0 {
		return nil, status.Error(100, "产品库存不足")
	}

	o := model.Order{
		Uid:    in.Uid,
		Pid:    in.Pid,
		Amount: in.Amount,
		Status: 0,
	}
	res, err := l.svcCtx.OrderModel.Insert(&o)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	o.Id, err = res.LastInsertId()
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	// 更新库存
	_, err = l.svcCtx.ProductRpc.Update(l.ctx, &product.UpdateRequest{
		Id:     productRes.Id,
		Name:   productRes.Name,
		Desc:   productRes.Desc,
		Stock:  productRes.Stock - 1,
		Amount: productRes.Amount,
		Status: productRes.Status,
	})
	if err != nil {
		return nil, err
	}

	return &order.CreateResponse{Id: o.Id}, nil
}
