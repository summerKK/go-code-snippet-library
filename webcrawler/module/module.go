package module

import (
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/base"
	"sync/atomic"
)

type Module struct {
	mid  base.MID
	addr string
	// 评分
	score uint64
	// 评分计算器
	scoreCalc base.CalculateScore
	// 调用次数
	calledCount uint64
	// j接收次数
	acceptedCount uint64
	// 已完成次数
	completedCount uint64
	// 实时处理次数
	handlingNum uint64
	counts      uint64
}

func NewModuleInternal(mid base.MID, scoreCalc base.CalculateScore) (*Module, error) {

	return &Module{mid: mid, scoreCalc: scoreCalc}, nil
}

func (m *Module) ID() base.MID {
	return m.mid
}

func (m *Module) Addr() string {
	return m.addr
}

func (m *Module) Score() uint64 {
	return atomic.LoadUint64(&m.score)
}

func (m *Module) SetScore(score uint64) {
	atomic.CompareAndSwapUint64(&m.score, m.score, score)
}

func (m *Module) ScoreCalculator() base.CalculateScore {
	return m.scoreCalc
}

func (m *Module) CalledCount() uint64 {
	return atomic.LoadUint64(&m.calledCount)
}

func (m *Module) AcceptedCount() uint64 {
	return atomic.LoadUint64(&m.acceptedCount)
}

func (m *Module) CompletedCount() uint64 {
	return atomic.LoadUint64(&m.completedCount)
}

func (m *Module) HandlingNum() uint64 {
	return atomic.LoadUint64(&m.handlingNum)
}

func (m *Module) Counts() base.Counts {
	return base.Counts{
		CalledCount:    atomic.LoadUint64(&m.calledCount),
		AcceptedCount:  atomic.LoadUint64(&m.acceptedCount),
		CompletedCount: atomic.LoadUint64(&m.completedCount),
		HandlingNumber: atomic.LoadUint64(&m.handlingNum),
	}
}

func (m *Module) Summary() base.SummaryStruct {
	counts := m.Counts()
	return base.SummaryStruct{
		ID:        m.ID(),
		Called:    counts.CalledCount,
		Accepted:  counts.AcceptedCount,
		Completed: counts.CompletedCount,
		Handing:   counts.HandlingNumber,
		Extra:     nil,
	}
}

func (m *Module) IncrCalledCount() {
	atomic.AddUint64(&m.calledCount, 1)
}

func (m *Module) IncrAcceptedCount() {
	atomic.AddUint64(&m.acceptedCount, 1)
}

func (m *Module) IncrCompletedCount() {
	atomic.AddUint64(&m.completedCount, 1)
}

func (m *Module) IncrHandlingNum() {
	atomic.AddUint64(&m.handlingNum, 1)
}

func (m *Module) DecrHandlingNum() {
	atomic.AddUint64(&m.handlingNum, ^uint64(0))
}

func (m *Module) Clear() {
	atomic.StoreUint64(&m.calledCount, 0)
	atomic.StoreUint64(&m.acceptedCount, 0)
	atomic.StoreUint64(&m.completedCount, 0)
	atomic.StoreUint64(&m.handlingNum, 0)
}
