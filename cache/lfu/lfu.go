package lfu

import (
	"container/heap"

	"github.com/summerKK/go-code-snippet-library/cache"
)

// LFU 是一个LFU cache,它不是并发安全的
// 淘汰策略是淘汰最少使用的元素
// queue 实现了最小堆(heap)
type lfu struct {
	// 缓存最大的容量,单位字节
	maxBytes int

	// 当一个元素被删除时,触发该函数,默认为 nil
	onEvicted func(key string, value interface{})

	queue *queue

	// 已使用的容量
	usedBytes int

	// 用来存键值对
	cache map[string]*entry
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	q := make(queue, 0, 1024)
	return &lfu{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		queue:     &q,
		cache:     make(map[string]*entry),
	}
}

func (l *lfu) Set(key string, value interface{}) {
	// key已经存在
	if e, ok := l.cache[key]; ok {
		l.usedBytes = l.usedBytes - l.cache[key].Len() + cache.CalcLen(value)
		l.queue.update(e, value, e.weight+1)
		return
	}

	en := &entry{key: key, value: value}
	heap.Push(l.queue, en)
	l.cache[key] = en

	l.usedBytes += en.Len()
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.removeElement(heap.Pop(l.queue))
	}
}

// 获取元素
func (l *lfu) Get(key string) interface{} {
	if e, ok := l.cache[key]; ok {
		l.queue.update(e, e.value, e.weight+1)
		return e.value
	}

	return nil
}

// 从 cache 中删除 key 对应的元素
func (l *lfu) Del(key string) {
	if e, ok := l.cache[key]; ok {
		heap.Remove(l.queue, e.index)
		l.removeElement(e)
	}
}

// 删除使用最小的元素
func (l *lfu) DelOldest() {
	if l.queue.Len() == 0 {
		return
	}

	l.removeElement(heap.Pop(l.queue))
}

// 缓存元素个数
func (l *lfu) Len() int {
	return l.queue.Len()
}

// 删除元素的一些执行操作
func (l *lfu) removeElement(x interface{}) {
	if x == nil {
		return
	}

	en := x.(*entry)
	delete(l.cache, en.key)
	l.usedBytes -= en.Len()

	if l.onEvicted != nil {
		l.onEvicted(en.key, en.value)
	}
}
