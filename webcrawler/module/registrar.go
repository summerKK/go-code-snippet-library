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

func NewRegistrar() *Registrar {
	return &Registrar{
		ModulesMap: make(map[MType]map[MID]IModule, 4),
	}
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

// Get 用户获取一个指定类型的组件实例
// 本函数会通过负载均衡的方式返回¬
func (r *Registrar) Get(mType MType) (selectedModule IModule, err error) {
	moduels, err := r.GetAllTypeBy(mType)
	if err != nil {
		return
	}
	minScore := uint64(0)
	for _, module := range moduels {
		// 计算score分值
		SetScore(module)
		score := module.Score()
		if score < minScore || minScore == 0 {
			minScore = score
			selectedModule = module
		}
	}

	return
}

func (r *Registrar) GetAllTypeBy(mType MType) (result map[MID]IModule, err error) {
	if !Legalletter(mType) {
		errMsg := fmt.Sprintf("incorrect module type:%s", mType)
		err = errors.NewIllegalParamsError(errMsg)
		return
	}

	r.locker.RLock()
	defer r.locker.RUnlock()
	modules := r.ModulesMap[mType]
	if modules == nil {
		err = ErrModuleNotFoundInstance
		return
	}
	// 这里要给变量重新赋值,因为modules也是一个map.避免引发数据竞争
	result = map[MID]IModule{}
	for i, i2 := range modules {
		result[i] = i2
	}

	return
}

func (r *Registrar) GetAll() (result map[MID]IModule) {
	result = map[MID]IModule{}
	r.locker.RLock()
	defer r.locker.RUnlock()
	for _, moduels := range r.ModulesMap {
		for mid, module := range moduels {
			result[mid] = module
		}
	}

	return
}

func (r *Registrar) Clear() {
	r.locker.Lock()
	defer r.locker.Unlock()
	r.ModulesMap = map[MType]map[MID]IModule{}
}
