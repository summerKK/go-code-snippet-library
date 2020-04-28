package buffer

import (
	"fmt"
	"github.com/summerKK/go-code-snippet-library/webcrawler/errors"
	"sync"
	"sync/atomic"
)

type Buffer struct {
	ch chan interface{}
	// 0 未关闭 1关闭
	closed      uint32
	closingLock sync.RWMutex
}

func NewBuffer(size uint32) (*Buffer, error) {
	if size == 0 {
		errMsg := fmt.Sprintf("illegal size for buffer:%d", size)
		return nil, errors.NewIllegalParamsError(errMsg)
	}
	return &Buffer{
		ch: make(chan interface{}, size),
	}, nil
}

func (b *Buffer) Cap() uint32 {
	return uint32(cap(b.ch))
}

func (b *Buffer) Len() uint32 {
	return uint32(len(b.ch))
}

func (b *Buffer) Put(datum interface{}) (ok bool, err error) {
	//  这里加锁是因为需要通过 Closed()方法确定通道是否关闭.给一个已经关闭的channel发送数据会造成panic
	// 在 Close()方法中,在关闭通道之前试图拿到锁,然后跟这里的锁刚好互斥.这样就不会造成在写入channel的时候别的goroutine刚好把channel关闭了
	// 这时候如果put拿到锁,关闭通道的goroutine就会阻塞等到put结束.
	b.closingLock.RLock()
	defer b.closingLock.RUnlock()
	if b.Closed() {
		return false, ErrClosedBuf
	}
	select {
	case b.ch <- datum:
		ok = true
	default:
		ok = false
	}

	return
}

func (b *Buffer) Get() (interface{}, error) {
	select {
	case datum, ok := <-b.ch:
		if !ok {
			return nil, ErrClosedBuf
		}
		return datum, nil
	default:
		return nil, nil
	}
}

func (b *Buffer) Close() bool {
	// 原子操作,把closed置位1
	if atomic.CompareAndSwapUint32(&b.closed, 0, 1) {
		// 需要使用锁来关闭channel,避免重复关闭channel而造成的panic
		b.closingLock.Lock()
		close(b.ch)
		b.closingLock.Unlock()
		return true
	}
	return false
}

func (b *Buffer) Closed() bool {
	if atomic.LoadUint32(&b.closed) == 0 {
		return false
	}
	return true
}
