package cmap

import (
	"math"
	"sync/atomic"
)

type IConcurrentMap interface {
	// 用户返回并发量
	Concurrency() int
	// 加入一个元素.注意elem不能为nil
	// 第一个返回值代表了是否新增一个元素,如果新增一个元素会替换旧的元素
	Put(key string, elem interface{}) (bool, error)
	// 获取一个元素,如果值不存在返回nil
	Get(key string) interface{}
	Delete(key string) bool
	Len() uint64
}

type concurrentMap struct {
	concurrency int
	segments    []ISegment
	total       uint64
}

func NewConcurrentMap(concurrency int, pairRedistrubitor IPairRedistributor) (*concurrentMap, error) {
	if concurrency <= 0 {
		return nil, newIllegalParameterError("concurrency to small")
	}
	if concurrency > MAX_CONCURRENCY {
		return nil, newIllegalParameterError("concurrency to large")
	}
	cmap := &concurrentMap{}
	cmap.concurrency = concurrency
	cmap.segments = make([]ISegment, concurrency)
	for i := 0; i < concurrency; i++ {
		cmap.segments[i] = newSegment(DEFAULT_BUCKET_NUMBER, pairRedistrubitor)
	}

	return cmap, nil
}

// cmap支持的最大并发量
func (c *concurrentMap) Concurrency() int {
	return c.concurrency
}

func (c *concurrentMap) Put(key string, elem interface{}) (ok bool, err error) {
	p, err := newPair(key, elem)
	if err != nil {
		return
	}
	segment := c.findSegment(p.Hash())
	ok, err = segment.Put(p)
	if ok {
		atomic.AddUint64(&c.total, 1)
	}

	return
}

func (c *concurrentMap) Get(key string) interface{} {
	hashKey := hash(key)
	segment := c.findSegment(hashKey)
	pair := segment.GetWithHash(key, hashKey)
	if pair == nil {
		return nil
	}

	return pair.Element()
}

func (c *concurrentMap) Len() uint64 {
	return atomic.LoadUint64(&c.total)
}

func (c *concurrentMap) Delete(key string) bool {
	segment := c.findSegment(hash(key))
	if segment.Delete(key) {
		atomic.AddUint64(&c.total, ^uint64(0))
		return true
	}

	return false
}

func (c *concurrentMap) findSegment(keyHash uint64) ISegment {
	if c.concurrency == 1 {
		return c.segments[0]
	}
	var keyHash32 uint32
	// 舍弃低位,保留高位.转换成 uint32
	if keyHash > math.MaxUint32 {
		keyHash32 = uint32(keyHash >> 32)
	} else {
		keyHash32 = uint32(keyHash)
	}
	// keyHash >> 16 右移16位保留高位,通过高位取模.这样hash值可以更均匀的分布在Segment中(~_~)
	return c.segments[int(keyHash32>>16)%(c.Concurrency()-1)]
}
