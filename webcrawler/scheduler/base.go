package scheduler

import (
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/base"
	"net/http"
)

type Status uint8

// BufferPoolSummaryStruct 代表缓冲池的摘要类型。
type BufPoolSummaryStruct struct {
	BufferCap       uint32 `json:"buffer_cap"`
	MaxBufferNumber uint32 `json:"max_buffer_number"`
	BufferNumber    uint32 `json:"buffer_number"`
	Total           uint64 `json:"total"`
}

type SummaryStruct struct {
	RequestArgs RequestArgs          `json:"request_args"`
	DataArgs    DataArgs             `json:"data_args"`
	ModuleArgs  ModuleArgsSummary    `json:"module_args"`
	Status      string               `json:"status"`
	Downloaders []base.SummaryStruct `json:"downloaders"`
	Analyzers   []base.SummaryStruct `json:"analyzers"`
	Pipelines   []base.SummaryStruct `json:"pipelines"`
	ReqBufPool  BufPoolSummaryStruct `json:"req_buf_pool"`
	RespBufPool BufPoolSummaryStruct `json:"resp_buf_pool"`
	ItemBufPool BufPoolSummaryStruct `json:"item_buf_pool"`
	ErrBufPool  BufPoolSummaryStruct `json:"err_buf_pool"`
	NumUrl      uint64               `json:"num_url"`
}

type IArgs interface {
	// 检查参数有效性
	Check() error
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
