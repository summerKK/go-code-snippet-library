package fast

import (
	"container/list"
	"sync"

	"github.com/summerKK/go-code-snippet-library/cache"
)

// 实现并发安全的缓存切片
// 使用 LRU 淘汰算法(最少使用,最近原则)
type cacheShard struct {
	locker sync.RWMutex

	// 最大存放 entry 个数
	maxEntries int

	// 当前存储数量
	currentEntries int

	// 当一个 entry 从缓存中移除时调用该回调函数,默认 nil
	// key 是任意的可以比较类型; value 是 interface{}
	onEvicted func(key string, value interface{})

	// 双向链表
	ll *list.List

	// 存放键值对
	cache map[string]*list.Element
}

func (c *cacheShard) set(key string, value interface{}) {
	c.locker.Lock()
	defer c.locker.Unlock()

	if e, ok := c.cache[key]; ok {
		c.ll.MoveToBack(e)
		en := e.Value.(*entry)
		en.value = value
		return
	}

	en := &entry{key: key, value: value}
	e := c.ll.PushBack(en)
	c.cache[key] = e
	c.currentEntries++

	if c.maxEntries > 0 && c.len() > c.maxEntries {
		c.removeElement(c.ll.Front())
	}
}

func newCacheShard(maxEntries int, onEvicted func(key string, value interface{})) *cacheShard {
	return &cacheShard{
		maxEntries: maxEntries,
		onEvicted:  onEvicted,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

func (c *cacheShard) get(key string) interface{} {
	c.locker.RLock()
	defer c.locker.RUnlock()

	if e, ok := c.cache[key]; ok {
		// 移到链表链尾.淘汰的时候先淘汰链头的数据
		c.ll.MoveToBack(e)
		en := e.Value.(*entry)
		return en.value
	}

	return nil
}

func (c *cacheShard) del(key string) {
	if e, ok := c.cache[key]; ok {
		c.removeElement(e)
	}
}

func (c *cacheShard) len() int {
	c.locker.RLock()
	defer c.locker.RUnlock()

	return c.currentEntries
}

func (c *cacheShard) removeElement(e *list.Element) {
	c.locker.Lock()
	defer c.locker.Unlock()

	if e == nil {
		return
	}

	c.ll.Remove(e)
	en := e.Value.(*entry)
	delete(c.cache, en.key)
	c.currentEntries--

	if c.onEvicted != nil {
		c.onEvicted(en.key, en.value)
	}
}

type entry struct {
	key   string
	value interface{}
}

func (e *entry) len() int {
	return cache.CalcLen(e.value)
}
