package svc

import (
	"github.com/summerKK/mall/service/order/api/internal/config"
	"github.com/summerKK/mall/service/order/rpc/orderclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	OrderRpc orderclient.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		OrderRpc: orderclient.NewOrder(zrpc.MustNewClient(c.OrderRpc)),
	}
}
