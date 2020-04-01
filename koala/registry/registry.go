package registry

import "context"

// 注册中心
type IRegistry interface {
	GetName() (name string)
	Register(ctx context.Context, service *Service) (err error)
	UnRegister(ctx context.Context, service *Service) (err error)
	Init(ctx context.Context, options ...Option) (err error)
	GetService(ctx context.Context, name string) (service *Service, err error)
}
