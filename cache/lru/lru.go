package lru

import (
	"container/list"

	"github.com/summerKK/go-code-snippet-library/cache"
)

// LRU 是一个 LRU cache,他不是并发安全的
// 根据最近最少使用原则淘汰元素
type lru struct {
	maxBytes int

	usedBytes int

	// 当元素被移除的时候触发的事件
	onEvicted func(key string, value interface{})

	list *list.List

	// 用来存键值对
	cache map[string]*list.Element
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	return &lru{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		list:      list.New(),
		cache:     make(map[string]*list.Element),
	}
}

func (l *lru) Set(key string, value interface{}) {
	// 如果 key 存在
	if e, ok := l.cache[key]; ok {
		// 移到队尾
		l.list.MoveToBack(e)
		en := e.Value.(*entry)
		en.value = value
		l.usedBytes = l.usedBytes - en.Len() + cache.CalcLen(value)

		return
	}

	en := &entry{key: key, value: value}
	e := l.list.PushBack(en)
	l.cache[key] = e
	// 添加使用容量
	l.usedBytes += en.Len()
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.DelOldest()
	}
}

func (l *lru) Get(key string) interface{} {
	if e, ok := l.cache[key]; ok {
		l.list.MoveToBack(e)
		return e.Value.(*entry).value
	}

	return nil
}

func (l *lru) Del(key string) {
	if e, ok := l.cache[key]; ok {
		l.removeElement(e)
	}
}

func (l *lru) DelOldest() {
	l.removeElement(l.list.Front())
}

func (l *lru) Len() int {
	return l.list.Len()
}

func (l *lru) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	l.list.Remove(e)
	en := e.Value.(*entry)
	delete(l.cache, en.key)

	l.usedBytes -= en.Len()

	if l.onEvicted != nil {
		l.onEvicted(en.key, en.value)
	}
}

type entry struct {
	key   string
	value interface{}
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value)
}
