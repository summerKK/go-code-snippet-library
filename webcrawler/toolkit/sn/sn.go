package sn

import (
	"math"
	"sync"
	"sync/atomic"
)

type Generator struct {
	// start 代表序列号的最小值。
	start uint64
	//  代表最大值,超过这个值就从0开始.当然cycleCount++
	max        uint64
	next       uint64
	cycleCount uint64
	locker     sync.RWMutex
}

// NewSNGenertor 会创建一个序列号生成器。
// 参数start用于指定第一个序列号的值。
// 参数max用于指定序列号的最大值。
func NewGenerator(start uint64, max uint64) *Generator {
	if max == 0 {
		max = math.MaxUint64
	}

	return &Generator{
		start: start,
		max:   max,
		next:  start,
	}
}

func (g *Generator) Min() uint64 {
	return g.start
}

func (g *Generator) Max() uint64 {
	return g.max
}

func (g *Generator) Next() uint64 {
	return atomic.LoadUint64(&g.next)
}

func (g *Generator) CycleCount() uint64 {
	return atomic.LoadUint64(&g.cycleCount)
}

func (g *Generator) Get() uint64 {
	g.locker.Lock()
	defer g.locker.Unlock()
	id := g.next
	if id == g.max {
		g.next = g.start
		g.cycleCount++
	} else {
		g.next++
	}
	return id
}
