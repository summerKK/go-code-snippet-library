package scheduler

import (
	"context"
	"errors"
	"fmt"
	"github.com/summerKK/go-code-snippet-library/cmap"
	"github.com/summerKK/go-code-snippet-library/webcrawler/logger"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/base"
	"github.com/summerKK/go-code-snippet-library/webcrawler/toolkit/buffer"
	"net/http"
	"strings"
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
	defer func() {
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("fatal scheduler error:%v", p)
			logger.Logger.Fatal(errMsg)
			err = genError(errMsg)
		}
	}()

	logger.Logger.Info("start scheduler...")
	logger.Logger.Info("check status for start...")

	var oldStatus Status
	oldStatus, err = s.checkAndSetStatus(SCHED_STATUS_STARTING)
	defer func() {
		s.statusLock.Lock()
		if err != nil {
			s.status = oldStatus
		} else {
			s.status = SCHED_STATUS_STARTED
		}
		s.statusLock.Unlock()
	}()

	if err != nil {
		return
	}

	logger.Logger.Info("check http request...")
	if request == nil {
		err = genError("nil http request")
		return
	}
	logger.Logger.Info("the http request is valid.")

	logger.Logger.Info("get the primary domain...")
	var primaryDomain string
	primaryDomain, err = getPrimaryDomain(request.Host)
	if err != nil {
		return
	}
	logger.Logger.Infof("-- primary domain:%s", primaryDomain)
	_, err = s.acceptDomainMap.Put(primaryDomain, struct{}{})
	if err != nil {
		return
	}

	// 检查缓冲池是否已经初始化好
	if err := s.checkBufferPoolForStart(); err != nil {
		return err
	}

	return
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

// checkBufferPoolForStart 会检查缓冲池是否已为调度器的启动准备就绪。
// 如果某个缓冲池不可用，就直接返回错误值报告此情况。
// 如果某个缓冲池已关闭，就按照原先的参数重新初始化它。
func (s *Scheduler) checkBufferPoolForStart() (err error) {
	if s.reqBufferPool == nil {
		err = genError("nul request buffer pool")
		return
	}
	if s.reqBufferPool != nil && s.reqBufferPool.Closed() {
		s.reqBufferPool, err = buffer.NewPool(s.reqBufferPool.BufCap(), s.reqBufferPool.MaxBufNum())
		if err != nil {
			return
		}
	}

	if s.respBufferPool == nil {
		err = genError("nul response buffer pool")
		return
	}
	if s.respBufferPool != nil && s.respBufferPool.Closed() {
		s.respBufferPool, err = buffer.NewPool(s.respBufferPool.BufCap(), s.respBufferPool.MaxBufNum())
		if err != nil {
			return
		}
	}

	if s.itemBufferPool == nil {
		err = genError("nul item buffer pool")
		return
	}
	if s.itemBufferPool != nil && s.itemBufferPool.Closed() {
		s.itemBufferPool, err = buffer.NewPool(s.itemBufferPool.BufCap(), s.itemBufferPool.MaxBufNum())
		if err != nil {
			return
		}
	}

	if s.errorBufferPool == nil {
		err = genError("nul error buffer pool")
		return
	}
	if s.errorBufferPool != nil && s.errorBufferPool.Closed() {
		s.errorBufferPool, err = buffer.NewPool(s.errorBufferPool.BufCap(), s.errorBufferPool.MaxBufNum())
		if err != nil {
			return
		}
	}

	return
}

func (s *Scheduler) download() {
	for {
		if s.canceled() {
			break
		}
		datum, err := s.reqBufferPool.Get()
		if err != nil {
			logger.Logger.Warn("the request buffer pool was closed.break request reception")
			break
		}
		request, ok := datum.(*module.Request)
		if !ok {
			errMsg := fmt.Sprintf("incorrect request type:%T", datum)
			sendError(errors.New(errMsg), "", s.errorBufferPool)
		}

		s.downloadOne(request)
	}
}

func (s *Scheduler) downloadOne(request *module.Request) {
	if request == nil {
		return
	}
	if s.canceled() {
		return
	}
	m, err := s.registrar.Get(base.TYPE_ANALYZER)
	if err != nil || m == nil {
		errMsg := fmt.Sprintf("couldn't get a downloader:%s", err)
		sendError(errors.New(errMsg), "", s.errorBufferPool)
		s.sendReq(request)
		return
	}
	downloader, ok := m.(base.IDownloader)
	if !ok {
		errMsg := fmt.Sprintf("incorrect downloader type: %T (MID: %s)", m, m.ID())
		sendError(errors.New(errMsg), m.ID(), s.errorBufferPool)
		s.sendReq(request)
		return
	}
	response, err := downloader.Download(request)
	if response != nil {
		s.sendResp(response)
	}
	if err != nil {
		sendError(err, m.ID(), s.errorBufferPool)
	}
}

func (s *Scheduler) sendReq(request *module.Request) bool {
	if request == nil {
		return false
	}
	if s.canceled() {
		return false
	}
	httpReq := request.Req()
	if httpReq == nil {
		logger.Logger.Warn("ignore the request! its http request is invalid!")
		return false
	}
	url := httpReq.URL
	if url == nil {
		logger.Logger.Warn("ignore the rquest! its http request url is invalid!")
		return false
	}
	scheme := strings.ToLower(url.Scheme)
	if scheme != "http" && scheme != "https" {
		logger.Logger.Warnf("Ignore the request! Its URL scheme is %q, but should be %q or %q. (URL: %s)\n", scheme, "http", "https", url)
		return false
	}
	if v := s.urlMap.Get(url.String()); v != nil {
		logger.Logger.Warnf("Ignore the request! Its URL is repeated. (URL: %s)\n", url)
		return false
	}
	pd, _ := getPrimaryDomain(url.Host)
	if s.acceptDomainMap.Get(pd) == nil {
		if pd == "bing.net" {
			panic(httpReq.URL)
		}
		logger.Logger.Warnf("Ignore the request! Its host %q is not in accepted primary domain map. (URL: %s)\n", httpReq.Host, url)
		return false
	}
	if request.Depth() > s.maxDepth {
		logger.Logger.Warnf("Ignore the request! Its depth %d is greater than %d. (URL: %s)\n", request.Depth(), s.maxDepth, url)
		return false
	}

	go func(req *module.Request) {
		if err := s.reqBufferPool.Put(req); err != nil {
			logger.Logger.Warnln("The request buffer pool was closed. Ignore request sending.")
		}
	}(request)

	_, _ = s.urlMap.Put(url.String(), struct{}{})

	return true
}

func (s *Scheduler) sendResp(resp *module.Response) bool {
	if resp == nil {
		return false
	}
	if s.respBufferPool.Closed() || s.canceled() {
		return false
	}

	go func(resp *module.Response) {
		if err := s.respBufferPool.Put(resp); err != nil {
			logger.Logger.Warnln("The response buffer pool was closed. Ignore response sending.")
		}
	}(resp)

	return true
}
