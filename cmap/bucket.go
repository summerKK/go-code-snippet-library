package cmap

import (
	"bytes"
	"sync"
	"sync/atomic"
)

// IBucket 代表并发安全的散列桶的接口。
type IBucket interface {
	// Put 会放入一个键-元素对。
	// 第一个返回值表示是否新增了键-元素对。
	// 若在调用此方法前已经锁定lock，则不要把lock传入！否则必须传入对应的lock！
	Put(p IPair, lock sync.Locker) (bool, error)
	// Get 会获取指定键的键-元素对。
	Get(key string) IPair
	// GetFirstPair 会返回第一个键-元素对。
	GetFirstPair() IPair
	// Delete 会删除指定的键-元素对。
	// 若在调用此方法前已经锁定lock，则不要把lock传入！否则必须传入对应的lock！
	Delete(key string, lock sync.Locker) bool
	// Clear 会清空当前散列桶。
	// 若在调用此方法前已经锁定lock，则不要把lock传入！否则必须传入对应的lock！
	Clear(lock sync.Locker)
	// Size 会返回当前散列桶的尺寸。
	Size() uint64
	// String 会返回当前散列桶的字符串表示形式。
	String() string
}

// bucket 代表并发安全的散列桶的类型。
type bucket struct {
	// firstValue 存储的是键-元素对列表的表头。
	firstValue atomic.Value
	size       uint64
}

// 占位符。
// 由于原子值不能存储nil，所以当散列桶空时用此符占位。
var placeholder IPair = &pair{}

// newBucket 会创建一个Bucket类型的实例。
func newBucket() IBucket {
	b := &bucket{}
	b.firstValue.Store(placeholder)
	return b
}

func (b *bucket) Put(p IPair, lock sync.Locker) (bool, error) {
	if p == nil {
		return false, newIllegalParameterError("nil pair")
	}
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	firstPair := b.GetFirstPair()
	// 第一个元素为空,就把当前元素插入到第一个元素
	if firstPair == nil {
		b.firstValue.Store(p)
		atomic.AddUint64(&b.size, 1)
		return true, nil
	}
	var target IPair
	key := p.Key()
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			target = v
			break
		}
	}
	if target != nil {
		// 这里的添加是原子操作
		_ = target.SetElement(p)
		// false代表更新元素
		return false, nil
	}

	// 把p插入到链表表头
	_ = p.SetNext(firstPair)
	b.firstValue.Store(p)
	atomic.AddUint64(&b.size, 1)

	return true, nil
}

func (b *bucket) Get(key string) IPair {
	firstPair := b.GetFirstPair()
	if firstPair == nil {
		return nil
	}
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			return v
		}
	}

	return nil
}

func (b *bucket) GetFirstPair() IPair {
	if v := b.firstValue.Load(); v == nil {
		return nil
	} else if p, ok := v.(IPair); !ok || p == placeholder {
		return nil
	} else {
		return p
	}
}

func (b *bucket) Delete(key string, lock sync.Locker) bool {
	var result bool
	if lock == nil {
		// 不加锁的方法
		result = b.deleteWithoutLock(key)
	} else {
		result, _ = b.deleteWithLock(key, lock)
	}

	return result
}

// 使用锁删除
func (b *bucket) deleteWithLock(key string, lock sync.Locker) (bool, error) {
	if lock == nil {
		return false, newIllegalParameterError("nil lock")
	}

	lock.Lock()
	defer lock.Unlock()

	firstPair := b.GetFirstPair()
	if firstPair == nil {
		return false, nil
	}

	// 保存删除节点的之前节点
	var prePairs []IPair
	// 目标值
	var target IPair
	// 目标值的下个节点
	var breakPoint IPair
	for v := firstPair; v != nil; v = firstPair.Next() {
		if v.Key() == key {
			target = v
			breakPoint = v.Next()
		}
		prePairs = append(prePairs, v)
	}

	// 没有找到对应的值
	if target == nil {
		return false, nil
	}

	// 删除当前节点
	newFirstpair := breakPoint
	for i := len(prePairs) - 1; i > 0; i-- {
		iPair := prePairs[i].Copy()
		_ = iPair.SetNext(newFirstpair)
		newFirstpair = iPair
	}

	if newFirstpair != nil {
		b.firstValue.Store(newFirstpair)
		atomic.AddUint64(&b.size, ^uint64(0))
	} else {
		b.firstValue.Store(placeholder)
		atomic.AddUint64(&b.size, 0)
	}

	return true, nil
}

// 不用锁删除(删除更高效)
func (b *bucket) deleteWithoutLock(key string) bool {
	firstPair := b.GetFirstPair()
	if firstPair == nil {
		return false
	}

	var prePair IPair
	var target IPair
	var nextPair IPair
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			target = v
			nextPair = v.Next()
			break
		}
		prePair = v
	}

	if target == nil {
		return false
	}

	// bucket只存在一个元素
	if prePair == nil {
		b.firstValue.Store(placeholder)
		atomic.AddUint64(&b.size, 0)
	} else {
		_ = prePair.SetNext(nextPair)
		atomic.AddUint64(&b.size, ^uint64(0))
	}

	return true
}

func (b *bucket) Clear(lock sync.Locker) {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	atomic.StoreUint64(&b.size, 0)
	b.firstValue.Store(placeholder)
}

func (b *bucket) Size() uint64 {
	return atomic.LoadUint64(&b.size)
}

func (b *bucket) String() string {
	var buf bytes.Buffer
	buf.WriteString("[ ")
	for v := b.GetFirstPair(); v != nil; v = v.Next() {
		buf.WriteString(v.String())
		buf.WriteString(" ")
	}
	buf.WriteString("]")
	return buf.String()
}
