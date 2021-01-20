package cow

import (
	"errors"
	"fmt"
	"sync/atomic"
)

// 并发安全的整数数组
type ConcurrentArray interface {
	// Set 用于设置指定索引上的元素值
	Set(index uint32, elem int) (err error)
	// Get 用户获取指定索引上的元素值
	Get(index uint32) (elem int, err error)
	// 获取数组长度
	Len() uint32
}

func NewConcurrentArray(length uint32) ConcurrentArray {
	arr := intArray{}
	arr.length = length
	arr.val.Store(make([]int, arr.length))

	return &arr
}

type intArray struct {
	length  uint32
	val     atomic.Value
	version int32
}

func (arr *intArray) Set(index uint32, elem int) (err error) {
	if err = arr.checkLen(index); err != nil {
		return
	}
	if err = arr.checkVal(); err != nil {
		return
	}

	// 不要这样做！否则会形成竞态条件！
	// oldArray := array.val.Load().([]int)
	// oldArray[index] = elem
	// array.val.Store(oldArray)

	// 下面这样做并发的时候可能覆盖其它goroutine的值
	// 因为不是原子操作.在 ② 获取到数组的值后到 ③ 的过程中其它 goroutine 可能会往数组里面写入值
	//① newArr := make([]int, arr.length)
	//② copy(newArr, arr.val.Load().([]int))
	//newArr[index] = elem
	//③ arr.val.Store(newArr)

	// 乐观自旋锁
	for {
		// 获取当前版本号(乐观锁)
		version := atomic.LoadInt32(&arr.version)
		// 用新数组存原来的元素
		newArr := make([]int, arr.length)
		copy(newArr, arr.val.Load().([]int))
		newArr[index] = elem
		// 如果之前获取的version和现在arr.version不一致,代表其它goroutine已经修改过arr.所以下一次继续尝试更新
		if atomic.CompareAndSwapInt32(&arr.version, version, version+1) {
			// 在此处代码仍然会发生中断,造成和上面一样的情况.但是大大减少了其它goroutine写操作被覆盖的问题
			arr.val.Store(newArr)
			break
		}
	}

	return nil
}

func (arr *intArray) Get(index uint32) (elem int, err error) {
	if err = arr.checkLen(index); err != nil {
		return
	}
	if err = arr.checkVal(); err != nil {
		return
	}

	elem = arr.val.Load().([]int)[index]

	return
}

func (arr *intArray) Len() uint32 {
	return arr.length
}

// 检查数组是否越界
func (arr *intArray) checkLen(index uint32) error {
	if index >= arr.length {
		return fmt.Errorf("Index out of range [0,%d]", arr.length)
	}

	return nil
}

func (arr *intArray) checkVal() error {
	v := arr.val.Load()
	if v == nil {
		return errors.New("Invalid int array")
	}

	return nil
}
