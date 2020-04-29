package scheduler

import (
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
