package module

import (
	"fmt"
	"github.com/summerKK/go-code-snippet-library/webcrawler/errors"
	"sync"
)

// 组件注册器
type Registrar struct {
	ModulesMap map[MType]map[MID]IModule
	locker     sync.RWMutex
}

func (r *Registrar) Register(module IModule) (succ bool, err error) {
	if module == nil {
		err = errors.NewIllegalParamsError("nil module instance")
		return
	}
	mid := module.ID()
	parts, err := SplitMid(mid)
	if err != nil {
		return
	}
	mType := legalletterTypeMap[parts[0]]
	if !CheckType(mType, module) {
		errMsg := fmt.Sprintf("incorrect module type:%s", mType)
		err = errors.NewIllegalParamsError(errMsg)
		return
	}

	r.locker.Lock()
	defer r.locker.Unlock()
	modules := r.ModulesMap[mType]
	if modules == nil {
		modules = map[MID]IModule{}
	}
	// 已经注册过
	if _, ok := modules[mid]; ok {
		succ = false
		return
	}

	modules[mid] = module
	// 这个步骤是需要的,虽然 modules是指针.如果在上面modules不为nil的时候.下面这步骤是不需要的.
	// but,modules是nil的时候.下面步骤是必须的
	r.ModulesMap[mType] = modules
	succ = true

	return
}

func (r *Registrar) UnRegister(mid MID) (succ bool, err error) {
	if mid == "" {
		err = errors.NewIllegalParamsError(
			fmt.Sprintf("illegal mid:%s", mid),
		)
		return
	}
	parts, err := SplitMid(mid)
	if err != nil {
		return
	}
	mType := legalletterTypeMap[parts[0]]
	if !Legalletter(mType) {
		errMsg := fmt.Sprintf("incorrect module type:%s", mType)
		err = errors.NewIllegalParamsError(errMsg)
		return
	}
	r.locker.Lock()
	defer r.locker.Unlock()

	if modules, ok := r.ModulesMap[mType]; ok {
		if _, ok := modules[mid]; ok {
			delete(modules, mid)
			succ = true
			return
		}
	}

	return
}

func (r *Registrar) Get(moduleType MType) (IModule, error) {
	panic("implement me")
}

func (r *Registrar) GetAllTypeBy(moduleType MType) (map[MID]IModule, error) {
	panic("implement me")
}

func (r *Registrar) GetAll() map[MID]IModule {
	panic("implement me")
}

func (r *Registrar) Clear() {
	panic("implement me")
}
