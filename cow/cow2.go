package cow

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

type ConcurrentIntArray interface {
	Set(index int, elem int) (old int, err error)
	Get(index int) (elem int, err error)
	Len() int
}

type segment struct {
	val    atomic.Value
	length int
	// 0：可读可写；1：只读。
	status uint32
}

func (seg *segment) init(length int) {
	seg.length = length
	seg.val.Store(make([]int, length))
}

func (seg *segment) checkIndex(index int) error {
	if index < 0 || index >= seg.length {
		return fmt.Errorf("index out of range [0, %d) in segment", seg.length)
	}
	return nil
}

func (seg *segment) set(index int, elem int) (old int, err error) {
	if err = seg.checkIndex(index); err != nil {
		return
	}
	// 修改成功后把值改为0
	defer atomic.StoreUint32(&seg.status, 0)

	point := 10
	count := 0
	// 这里实现的是悲观自旋锁.只有拿到status = 0的值的时候才有权更新seg
	for {
		count++
		// 当status != 0的时候.代表当前有其它goroutine在更新值.这里等待下次再尝试更新
		if !atomic.CompareAndSwapUint32(&seg.status, 0, 1) {
			// 多次尝试拿不到status.就让出CPU时间给其它goroutine
			if count%point == 0 {
				// 优化 让出CPU时间.
				runtime.Gosched()
			}
			continue
		}

		newArr := make([]int, seg.length)
		copy(newArr, seg.val.Load().([]int))
		// 旧值
		old = newArr[index]
		newArr[index] = elem
		// 原子操作
		seg.val.Store(newArr)

		return
	}
}

func (seg *segment) get(index int) (elem int, err error) {
	if err = seg.checkIndex(index); err != nil {
		return
	}
	elem = seg.val.Load().([]int)[index]

	return
}

// myIntArray 代表 ConcurrentIntArray 接口的实现类型。
type myIntArray struct {
	length    int        // 元素总数量。
	segLenStd int        // 单个内部段的标准长度。
	segments  []*segment // 内部段列表。
}

// NewConcurrentIntArray 会创建一个 ConcurrentIntArray 类型值。
func NewConcurrentIntArray(length int) ConcurrentIntArray {
	if length < 0 {
		length = 0
	}
	array := new(myIntArray)
	array.init(length)
	return array
}

func (array *myIntArray) init(length int) {
	array.length = length
	// 每个segment的size
	array.segLenStd = 10 //TODO 此处是一个优化点，可以根据参数值调整。
	// 多少个segment
	segNum := length / array.segLenStd
	// 最后一个segment的size
	segLenTail := length % array.segLenStd
	if segLenTail > 0 {
		segNum = segNum + 1
	}
	array.segments = make([]*segment, segNum)
	for i := 0; i < segNum; i++ {
		seg := segment{}
		// 最后一个segment
		if i == segNum-1 && segLenTail > 0 {
			seg.init(segLenTail)
		} else {
			seg.init(array.segLenStd)
		}
		array.segments[i] = &seg
	}
}

func (array *myIntArray) Set(index int, elem int) (old int, err error) {
	if err = array.checkIndex(index); err != nil {
		return
	}
	// 获取对应的segment
	seg := array.segments[index/array.segLenStd]
	// index 34 % 10 = 4
	return seg.set(index%array.segLenStd, elem)
}

func (array *myIntArray) Get(index int) (elem int, err error) {
	if err = array.checkIndex(index); err != nil {
		return
	}
	seg := array.segments[index/array.segLenStd]
	return seg.get(index % array.segLenStd)
}

func (array *myIntArray) Len() int {
	return array.length
}

// checkIndex 用于检查索引的有效性。
func (array *myIntArray) checkIndex(index int) error {
	if index < 0 || index >= array.length {
		return fmt.Errorf("index out of range [0, %d)", array.length)
	}
	return nil
}
