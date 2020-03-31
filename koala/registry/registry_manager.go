package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	Manager = &Mgr{
		Plugins: make(map[string]IRegistry),
	}
)

// 注册中心服务管理
type Mgr struct {
	Plugins map[string]IRegistry
	sync.RWMutex
}

func (p *Mgr) add(registry IRegistry) error {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.Plugins[registry.GetName()]; ok {
		return errors.New("请不要重复注册服务")
	}

	p.Plugins[registry.GetName()] = registry

	return nil
}

func (p *Mgr) initRegistry(ctx context.Context, name string, options ...Option) (registry IRegistry, err error) {
	p.RLock()
	defer p.RUnlock()
	registry, ok := p.Plugins[name]
	if !ok {
		return nil, fmt.Errorf("%v 注册中心不存在", name)
	}
	err = registry.Init(ctx, options...)
	return
}

func RegisterRegistry(registery IRegistry) (err error) {
	return Manager.add(registery)
}

func InitRegistry(ctx context.Context, name string, options ...Option) (registry IRegistry, err error) {
	return Manager.initRegistry(ctx, name, options...)
}
