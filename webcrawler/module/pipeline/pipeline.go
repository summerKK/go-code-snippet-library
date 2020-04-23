package pipeline

import (
	"fmt"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/base"
)

type Pipeline struct {
	*module.Module
	itemProcessors []base.ProcessItem
	failFast       bool
}

func New(mid base.MID, scoreCalc base.CalculateScore, itemProcessors []base.ProcessItem) (*Pipeline, error) {
	m, err := module.NewModuleInternal(mid, scoreCalc)
	if err != nil {
		return nil, err
	}
	if itemProcessors == nil {
		return nil, genParameterError("nil itemProcessors")
	}
	if len(itemProcessors) == 0 {
		return nil, genParameterError("empty itemProcessors")
	}
	// 主要是为了防止在分析器创建后外界再对解析器列表进行更改.所以赋值给新的变量
	var innerItemProcessors []base.ProcessItem
	for i, processor := range itemProcessors {
		if processor == nil {
			return nil, genParameterError(fmt.Sprintf("nil item processor [%d]", i))
		}
		innerItemProcessors = append(innerItemProcessors, processor)
	}

	return &Pipeline{
		Module:         m,
		itemProcessors: innerItemProcessors,
	}, nil
}

func (p *Pipeline) ItemProcessors() []base.ProcessItem {
	return p.itemProcessors
}

func (p *Pipeline) Send(item module.Item) (errlist []error) {
	p.Module.IncrHandlingNum()
	defer p.Module.DecrHandlingNum()
	p.Module.IncrCalledCount()
	if item == nil {
		errlist = append(errlist, genParameterError("nil item"))
		return
	}
	p.Module.IncrAcceptedCount()
	currentItem := item
	// item 经过管道把结果输出.(没上一层的处理结果作为下一层的参数.linux的管道概念)
	for _, processor := range p.itemProcessors {
		m, err := processor(currentItem)
		if err != nil {
			errlist = append(errlist, err)
			if p.failFast {
				break
			}
		}
		if m != nil {
			currentItem = m
		}
	}
	if len(errlist) == 0 {
		p.Module.IncrCompletedCount()
	}

	return
}

func (p *Pipeline) FailFast() bool {
	return p.failFast
}

func (p *Pipeline) SetFailFast(failFast bool) {
	p.failFast = failFast
}
