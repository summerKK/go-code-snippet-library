package module

type MID string
type Type string
type CalculateScore func(counts Counts) uint64

type SummaryStruct struct {
	ID        MID         `json:"id"`
	Called    uint64      `json:"called"`
	Accepted  uint64      `json:"accepted"`
	Completed uint64      `json:"completed"`
	Handing   uint64      `json:"handing"`
	Extra     interface{} `json:"extra,omitempty"`
}

// Counts 代表用于汇集组件内部计数的类型。
type Counts struct {
	// CalledCount 代表调用计数。
	CalledCount uint64
	// AcceptedCount 代表接受计数。
	AcceptedCount uint64
	// CompletedCount 代表成功完成计数。
	CompletedCount uint64
	// HandlingNumber 代表实时处理数。
	HandlingNumber uint64
}

const (
	TYPE_DOWNLOADER Type = "downloader"
	TYPE_ANALYZER   Type = "analyzer"
	TYPE_PIPELINE   Type = "pipeline"
)

var legalletterMap = map[Type]string{
	TYPE_DOWNLOADER: "D",
	TYPE_ANALYZER:   "A",
	TYPE_PIPELINE:   "P",
}

type IModule interface {
	// 当前组件的ID
	ID() MID
	// 当前组件的网络地址
	Addr() string
	Score() uint64
	SetScore(score uint64)
	// 获取评分计算器
	ScoreCalculator() CalculateScore
	// 获取组件被调用次数
	CalledCount() uint64
	// 获取组件可以被调用的次数,组件一般会因为超负荷h或参数错误而拒绝调用
	AcceptedCount() uint64
	// 获取组件已经完成的调用次数
	CompletedCount() uint64
	// 获取组件正在处理的调用次数
	HandlingNum() uint64
	// 一次性获取所有计数
	Counts() Counts
	Summary() SummaryStruct
}

type IRegistrar interface {
	Register(module IModule) (bool, error)
	UnRegister(mid MID) (bool, error)
	Get(moduleType Type) (IModule, error)
	GetAllTypeBy(moduleType Type) (map[MID]IModule, error)
	GetAll() map[MID]IModule
	// 清除所有组件
	Clear()
}

type IModuleInternal interface {
	IModule
	IncrCalledCount()
	IncrAcceptedCount()
	IncrCompletedCount()
	IncrHandlingNum()
	DecrHandlingNum()
	// 清空所有计数
	Clear()
}
