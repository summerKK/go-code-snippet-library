package svc

import (
	"github.com/summerKK/mall/service/product/api/internal/config"
	"github.com/summerKK/mall/service/product/rpc/productclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	ProductRpc productclient.Product
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		ProductRpc: productclient.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
	}
}
