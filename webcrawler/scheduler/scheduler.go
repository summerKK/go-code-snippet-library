package scheduler

import (
	"context"
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
	registrar module.Registrar
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
	logger.Logger.Info("Check status for initialization...")
	var oldStatus Status
	oldStatus, err = s.checkAndSetStatus(SCHED_STATUS_INITING)
	if err != nil {
		return
	}
	defer func() {
		s.statusLock.Lock()
		if err != nil {
			s.status = oldStatus
		} else {
			s.status = SCHED_STATUS_INITED
		}
	}()
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
