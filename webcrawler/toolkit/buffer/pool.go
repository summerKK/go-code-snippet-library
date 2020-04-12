package buffer

import (
	"sync"
	"sync/atomic"
)

// hello,world
type Pool struct {
	// 缓冲器的统一容量
	bufCap uint32
	// 最大缓冲器的数量
	maxBufNum uint32
	// 缓冲器的实际数量
	bufNum uint32
	// 池中数据的总量
	total uint64
	// 存放缓冲器的channel
	bufChan chan IBuf
	// 缓冲池的关闭状态 0未关闭 1关闭
	closed uint32
	rwlock sync.RWMutex
}

func (p *Pool) BufCap() uint32 {
	return p.bufCap
}

func (p *Pool) MaxBufNum() uint32 {
	return p.maxBufNum
}

func (p *Pool) BufNum() uint32 {
	// 因为缓冲器的实际数量会动态的.所以需要添加锁
	p.rwlock.RLock()
	defer p.rwlock.RUnlock()
	return p.bufNum
}

func (p *Pool) Total() uint64 {
	// 加锁原理同上
	p.rwlock.RLock()
	defer p.rwlock.RUnlock()
	return p.total
}

func (p *Pool) Put(datum interface{}) (err error) {
	// 这里closed方法不用加锁.为什么会在buffer.go的Buffer.Put方法加锁是因为在Put方法里面给channel发送了值.必须确定channel没有被关闭
	if p.Closed() {
		return ErrClosedBufPool
	}
	var count uint32
	maxCount := p.BufNum() * 5
	var ok bool
	for buf := range p.bufChan {
		ok, err = p.putData(buf, datum, &count, maxCount)
		if ok || err != nil {
			break
		}
	}
	return
}

func (p *Pool) Get() (datum interface{}, err error) {
	panic("implement me")
}

func (p *Pool) Close() bool {
	panic("implement me")
}

func (p *Pool) Closed() bool {
	panic("implement me")
}

func (p *Pool) putData(buf IBuf, datum interface{}, count *uint32, maxCount uint32) (ok bool, err error) {
	if p.Closed() {
		return false, ErrClosedBufPool
	}
	defer func() {
		// 这里加锁操作是为了避免在channel关闭的时候往channel里面发送数据,从而造成panic
		p.rwlock.RLock()
		if p.Closed() {
			// channel已经关闭,就不归还n拿到的缓冲器.对应的bufNum也需要-1
			atomic.AddUint32(&p.bufNum, ^uint32(0))
			err = ErrClosedBufPool
		} else {
			// 归还拿到的缓冲器
			p.bufChan <- buf
		}
		p.rwlock.RUnlock()
	}()

	ok, err = buf.Put(datum)
	if ok {
		atomic.AddUint64(&p.total, 1)
		return
	}

	// 缓冲器已经满了
	if err != nil {
		return
	}

	// 执行到这里得时候说明buf缓冲器已经满了无法放入数据.
	*count++
}
