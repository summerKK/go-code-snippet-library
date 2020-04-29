package scheduler

import (
	"context"
	"fmt"
	"github.com/summerKK/go-code-snippet-library/cmap"
	"github.com/summerKK/go-code-snippet-library/webcrawler/logger"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module"
	"github.com/summerKK/go-code-snippet-library/webcrawler/toolkit/buffer"
	"net/http"
	"sync"
)

type Scheduler struct {
	// 爬取最大深度,初次为0
	maxDepth uint32
	// 可以接受的URL的主域名的字典(申明爬虫可以爬取到域名,避免爬取无关紧要的数据)
	acceptDomainMap cmap.IConcurrentMap
	// 组件注册器
	registrar *module.Registrar
	//  请求的缓冲池
	reqBufferPool buffer.IPool
	// 响应的缓冲池
	respBufferPool buffer.IPool
	// 条目的缓冲池
	itemBufferPool buffer.IPool
	// 错误的缓冲池
	errorBufferPool buffer.IPool
	// 已处理的URL的字典.(把已经处理的url放入字典.避免重复处理)
	urlMap cmap.IConcurrentMap
	// 上下文.用于感知调度器的停止
	ctx context.Context
	//  取消函数,用于停止调度器
	cancelFunc context.CancelFunc
	//  状态
	status Status
	//  专用于状态的读写锁
	statusLock sync.Mutex
	// 摘要信息
	summary ISchedSummary
}

func NewScheduler() IScheduler {
	return &Scheduler{}
}

func (s *Scheduler) Init(requestArgs RequestArgs, dataArgs DataArgs, moduleArgs ModuleArgs) (err error) {
	var oldStatus Status

	defer func() {
		s.statusLock.Lock()
		if err != nil {
			s.status = oldStatus
		} else {
			s.status = SCHED_STATUS_INITED
		}
		s.statusLock.Unlock()
	}()

	logger.Logger.Info("Check status for initialization...")
	oldStatus, err = s.checkAndSetStatus(SCHED_STATUS_INITING)
	if err != nil {
		return
	}

	logger.Logger.Info("Check request arguments...")
	if err = requestArgs.Check(); err != nil {
		return
	}
	logger.Logger.Info("request arguemtns are valid.")

	logger.Logger.Info("Check data arguments...")
	if err = dataArgs.Check(); err != nil {
		return
	}
	logger.Logger.Info("data arguments are valid.")

	logger.Logger.Info("check module arguments...")
	if err = moduleArgs.Check(); err != nil {
		return
	}
	logger.Logger.Info("module arguments are valid.")

	logger.Logger.Info("initialize scheduler's fields...")
	if s.registrar == nil {
		s.registrar = module.NewRegistrar()
	} else {
		s.registrar.Clear()
	}
	s.maxDepth = requestArgs.MaxDepth
	logger.Logger.Infof("-- max depth %d\n", s.maxDepth)
	s.acceptDomainMap, _ = cmap.NewConcurrentMap(1, nil)

	for _, domain := range requestArgs.AcceptedDomains {
		_, _ = s.acceptDomainMap.Put(domain, struct{}{})
	}
	logger.Logger.Infof("-- accepted primary domains:%v", requestArgs.AcceptedDomains)

	s.urlMap, _ = cmap.NewConcurrentMap(16, nil)
	logger.Logger.Infof("-- url map: length:%d, concurrency:%d", s.urlMap.Len(), s.urlMap.Concurrency())

	s.initBuffPool(dataArgs)
	s.resetContext()
	s.summary = newSchedSummary(requestArgs, dataArgs, moduleArgs, s)
	// 注册组件。
	logger.Logger.Info("Register modules...")
	if err = s.registerModules(moduleArgs); err != nil {
		return err
	}
	logger.Logger.Info("Scheduler has been initialized.")

	return
}

func (s *Scheduler) Start(request *http.Request) (err error) {
	panic("implement me")
}

func (s *Scheduler) Stop() (err error) {
	panic("implement me")
}

func (s *Scheduler) Status() Status {
	panic("implement me")
}

func (s *Scheduler) ErrChan() <-chan error {
	panic("implement me")
}

func (s *Scheduler) Idle() bool {
	panic("implement me")
}

func (s *Scheduler) Summary() ISchedSummary {
	panic("implement me")
}

func (s *Scheduler) checkAndSetStatus(status Status) (oldStatus Status, err error) {
	s.statusLock.Lock()
	defer s.statusLock.Unlock()
	oldStatus = s.status
	err = checkStatus(oldStatus, status, nil)
	if err == nil {
		s.status = status
	}

	return
}

// initBufferPool 用于按照给定的参数初始化缓冲池。
// 如果某个缓冲池可用且未关闭，就先关闭该缓冲池。
func (s *Scheduler) initBuffPool(args DataArgs) {
	if s.reqBufferPool != nil && !s.reqBufferPool.Closed() {
		s.reqBufferPool.Close()
	}
	s.reqBufferPool, _ = buffer.NewPool(args.ReqBufCap, args.ReqMaxBufNum)
	logger.Logger.Infof("-- request buffer pool: bufferCap:%d, maxBufferNum:%d\n", s.reqBufferPool.BufCap(), s.reqBufferPool.MaxBufNum())

	if s.respBufferPool != nil && !s.respBufferPool.Closed() {
		s.respBufferPool.Close()
	}
	s.respBufferPool, _ = buffer.NewPool(args.RespBufCap, args.RespMaxBufNum)
	logger.Logger.Infof("-- response buffer pool: bufferCap:%d, maxBufferNum:%d\n", s.respBufferPool.BufCap(), s.respBufferPool.MaxBufNum())

	if s.itemBufferPool != nil && !s.itemBufferPool.Closed() {
		s.itemBufferPool.Close()
	}
	s.itemBufferPool, _ = buffer.NewPool(args.ItemBufCap, args.ItemMaxBufNum)
	logger.Logger.Infof("-- item buffer pool: bufferCap:%d, maxBufferNum:%d\n", s.itemBufferPool.BufCap(), s.itemBufferPool.MaxBufNum())

	if s.errorBufferPool != nil && !s.errorBufferPool.Closed() {
		s.errorBufferPool.Close()
	}
	s.errorBufferPool, _ = buffer.NewPool(args.ErrBufCap, args.ErrMaxBufNum)
	logger.Logger.Infof("-- error buffer pool: bufferCap:%d, maxBufferNum:%d\n", s.errorBufferPool.BufCap(), s.errorBufferPool.MaxBufNum())
}

func (s *Scheduler) resetContext() {
	s.ctx, s.cancelFunc = context.WithCancel(context.Background())
}

func (s *Scheduler) canceled() bool {
	select {
	case <-s.ctx.Done():
		return true
	default:
		return false
	}
}

// registerModules 会注册所有给定的组件。
func (s *Scheduler) registerModules(moduleArgs ModuleArgs) error {
	for _, d := range moduleArgs.Downloaders {
		if d == nil {
			continue
		}
		ok, err := s.registrar.Register(d)
		if err != nil {
			return genErrorByError(err)
		}
		if !ok {
			errMsg := fmt.Sprintf("Couldn't register downloader instance with MID %q!", d.ID())
			return genError(errMsg)
		}
	}
	logger.Logger.Infof("All downloads have been registered. (number: %d)", len(moduleArgs.Downloaders))
	for _, a := range moduleArgs.Analyzers {
		if a == nil {
			continue
		}
		ok, err := s.registrar.Register(a)
		if err != nil {
			return genErrorByError(err)
		}
		if !ok {
			errMsg := fmt.Sprintf("Couldn't register analyzer instance with MID %q!", a.ID())
			return genError(errMsg)
		}
	}
	logger.Logger.Infof("All analyzers have been registered. (number: %d)", len(moduleArgs.Analyzers))
	for _, p := range moduleArgs.Pipelines {
		if p == nil {
			continue
		}
		ok, err := s.registrar.Register(p)
		if err != nil {
			return genErrorByError(err)
		}
		if !ok {
			errMsg := fmt.Sprintf("Couldn't register pipeline instance with MID %q!", p.ID())
			return genError(errMsg)
		}
	}
	logger.Logger.Infof("All pipelines have been registered. (number: %d)", len(moduleArgs.Pipelines))

	return nil
}
