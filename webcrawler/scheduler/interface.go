package scheduler

import (
	"github.com/summerKK/go-code-snippet-library/webcrawler/module"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/analyzer"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/downloader"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/pipeline"
	"net/http"
)

type Status uint8

const (
	// 未初始化
	SCHED_STATUS_UINIT = iota
	// 初始化中
	SCHED_STATUS_INITING
	// 已初始化
	SCHED_STATUS_INITED
	//  启动中
	SCHED_STATUS_STARTING
	// 已启动
	SCHED_STATUS_STARTED
	// 停止中
	SCHED_STATUS_STOPPING
	// 已停止
	SCHED_STATUS_STOPPED
)

// 请求参数
type RequestArgs struct {
	// 代表可以接受的URL的主域名的列表,不在该列表的URL都会被忽略(限定爬虫爬取的范围)
	AcceptedDomains []string `json:"accepted_primary_domains"`
	// 爬取的最大深度
	MaxDepth uint32 `json:"max_depth"`
}

// 定义数据缓冲池容量
type DataArgs struct {
	ReqBufCap     uint32 `json:"req_buf_cap"`
	ReqMaxBufNum  uint32 `json:"req_max_buf_num"`
	RespBufCap    uint32 `json:"resp_buf_cap"`
	RespMaxBufNum uint32 `json:"resp_max_buf_num"`
	ItemBufCap    uint32 `json:"item_buf_cap"`
	ItemMaxBufNum uint32 `json:"item_max_buf_num"`
	ErrBufCap     uint32 `json:"err_buf_cap"`
	ErrMaxBufNum  uint32 `json:"err_max_buf_num"`
}

// module 模块
type ModuleArgs struct {
	Downloaders []downloader.IDownloader
	Analyzers   []analyzer.IAnalyzer
	Pipelines   []pipeline.IPipeline
}

// ModuleArgsSummary 代表组件相关的参数容器的摘要类型。
type ModuleArgsSummary struct {
	DownloaderListSize int `json:"downloader_list_size"`
	AnalyzerListSize   int `json:"analyzer_list_size"`
	PipelineListSize   int `json:"pipeline_list_size"`
}

// BufferPoolSummaryStruct 代表缓冲池的摘要类型。
type BufPoolSummaryStruct struct {
	BufferCap       uint32 `json:"buffer_cap"`
	MaxBufferNumber uint32 `json:"max_buffer_number"`
	BufferNumber    uint32 `json:"buffer_number"`
	Total           uint64 `json:"total"`
}

type SummaryStruct struct {
	RequestArgs RequestArgs            `json:"request_args"`
	DataArgs    DataArgs               `json:"data_args"`
	ModuleArgs  ModuleArgsSummary      `json:"module_args"`
	Status      string                 `json:"status"`
	Downloaders []module.SummaryStruct `json:"downloaders"`
	Analyzers   []module.SummaryStruct `json:"analyzers"`
	Pipelines   []module.SummaryStruct `json:"pipelines"`
	ReqBufPool  BufPoolSummaryStruct   `json:"req_buf_pool"`
	RespBufPool BufPoolSummaryStruct   `json:"resp_buf_pool"`
	ItemBufPool BufPoolSummaryStruct   `json:"item_buf_pool"`
	ErrBufPool  BufPoolSummaryStruct   `json:"err_buf_pool"`
	NumUrl      uint64                 `json:"num_url"`
}

type IArgs interface {
	// 检查参数有效性
	Check()
}

type IScheduler interface {
	Init(requestArgs RequestArgs, dataArgs DataArgs, moduleArgs ModuleArgs) (err error)
	Start(request *http.Request) (err error)
	Stop() (err error)
	// 当前运行状态
	Status() Status
	// 获取处理错误chan
	ErrChan() <-chan error
	// 判断是否空闲
	Idle() bool
	Summary() ISchedSummary
}

type ISchedSummary interface {
	Struct() SummaryStruct
	String() string
}
