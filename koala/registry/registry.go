package registry

import "context"

// 注册中心
type IRegistry interface {
	GetName() string
	Register(ctx context.Context, service *Service)
	UnRegister(ctx context.Context, service *Service)
	Init(ctx context.Context, name string, options ...Option) error
}
