package fifo

import (
	"container/list"

	"github.com/summerKK/go-code-snippet-library/cache"
)

// FIFO 是一个 FIFO cache,他不是并发安全的
// 通过 map 和 双向链表实现查询o(1),插入o(1)
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
		en := e.Value.(*entry)
		// key 对应的元素移动到了链尾,减去原来元素的大小,加上新增的的元素的大小
		f.usedBytes = f.usedBytes - cache.CalcLen(en.value) + cache.CalcLen(value)
		en.value = value
		return
	}

	// 构建一个entry
	en := &entry{key, value}
	// 加入到链表尾部,并返回Element类型数据
	e := f.ll.PushBack(en)
	// 加入到map里面
	f.cache[key] = e

	f.usedBytes += en.Len()
	if f.maxBytes > 0 && f.usedBytes > f.maxBytes {
		f.DelOldest()
	}
}

func (f *fifo) Get(key string) interface{} {
	if e, ok := f.cache[key]; ok {
		return e.Value.(*entry).value
	}

	return nil
}

// 删除知道 key 的数据
func (f *fifo) Del(key string) {
	if e, ok := f.cache[key]; ok {
		f.removeElement(e)
	}
}

// 删除链头的数据(先进先出)
func (f *fifo) DelOldest() {
	f.removeElement(f.ll.Front())
}

func (f *fifo) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	// 从链表删除元素
	f.ll.Remove(e)
	en := e.Value.(*entry)
	// 从使用容量中删除删除元素的容量
	f.usedBytes -= en.Len()
	// 从 map 中删除指定元素
	delete(f.cache, en.key)

	// 触发删除事件,如果存在的话
	if f.onEvicted != nil {
		f.onEvicted(en.key, en.value)
	}
}

// 获取链表元素个数
func (f *fifo) Len() int {
	return f.ll.Len()
}

// 封装一下链表的 value
type entry struct {
	key   string
	value interface{}
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value)
}
