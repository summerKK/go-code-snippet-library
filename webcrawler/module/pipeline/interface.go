package pipeline

import (
	"github.com/summerKK/go-code-snippet-library/webcrawler/module"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/data"
)

type ProcessItem func(item data.Item) (result data.Item, err error)

type IPipeline interface {
	module.IModule
	// 返回当前处理函数的列表
	ItemProcessors() []ProcessItem
	// 向条目处理管道发送条目,条目需要依次经过若干个条目处理函数
	Send(item data.IData) []error
	// 快速失败,如果在条目处理函数中有一个失败,就立即返回错误忽略后面的操作
	FailFast() bool
	// 设置是否快速失败
	SetFailFast(failFast bool)
}
