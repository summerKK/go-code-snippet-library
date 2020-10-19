package fifo

import (
	"container/list"

	"github.com/summerKK/go-code-snippet-library/cache"
)

// FIFO 是一个 FIFO cache,他不是并发安全的
type fifo struct {
	// 缓存最大的容量,单位字节
	maxBytes int

	// 当一个 entry 从缓存中移除时调用该回调函数,默认为 nil
	onEvicted func(key string, value interface{})

	// 已使用的字节数,只包括值,key不算
	usedBytes int

	// 双向链表
	ll *list.List

	// 用来存键值对
	cache map[string]*list.Element
}

func New(maxBytes int, OnEvicted func(key string, value interface{})) cache.Cache {
	return &fifo{
		maxBytes:  maxBytes,
		onEvicted: OnEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}

// 通过 Set 方法往 Cache 尾部增加一个元素(如果已经存在,则移到尾部,并修改值)
func (f *fifo) Set(key string, value interface{}) {
	if e, ok := f.cache[key]; ok {
		f.ll.MoveToBack(e)
	}
}

func (f *fifo) Get(key string) interface{} {
	panic("implement me")
}

func (f *fifo) Del(key string) {
	panic("implement me")
}

func (f *fifo) DelOldest() {
	panic("implement me")
}

func (f *fifo) Len() int {
	panic("implement me")
}
