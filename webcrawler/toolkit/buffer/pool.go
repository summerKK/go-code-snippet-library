package buffer

import (
	"fmt"
	"github.com/summerKK/go-code-snippet-library/webcrawler/errors"
	"sync"
	"sync/atomic"
)

// hello,world
type Pool struct {
	// 缓冲器的容量,当当前缓冲池空的时候没有缓冲器可以使用.会创建新的缓冲器,这时候就会用上bufCap
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

func NewPool(bufCap uint32, maxBufNum uint32) (IPool, error) {
	if bufCap == 0 {
		errMsg := fmt.Sprintf("illegal buffer cap for buffer pool:%d", bufCap)
		return nil, errors.NewIllegalParamsError(errMsg)
	}
	if maxBufNum == 0 {
		errMsg := fmt.Sprintf("illegal max buffer num for buffer pool:%d", maxBufNum)
		return nil, errors.NewIllegalParamsError(errMsg)
	}
	bufChan := make(chan IBuf, maxBufNum)
	buf, _ := NewBuffer(bufCap)
	bufChan <- buf

	return &Pool{
		bufCap:    bufCap,
		maxBufNum: maxBufNum,
		bufChan:   bufChan,
		bufNum:    1,
	}, nil
}

func (p *Pool) BufCap() uint32 {
	return p.bufCap
}

func (p *Pool) MaxBufNum() uint32 {
	return p.maxBufNum
}

func (p *Pool) BufNum() uint32 {
	// 因为缓冲器的实际数量会动态的.所以需要添加锁
	return atomic.LoadUint32(&p.bufNum)
}

func (p *Pool) Total() uint64 {
	return atomic.LoadUint64(&p.total)
}

func (p *Pool) Put(datum interface{}) (err error) {
	// 这里closed方法不用加锁.为什么会在buffer.go的Buffer.Put方法加锁是因为在Put方法里面给channel发送了值.必须确定channel没有被关闭
	if p.Closed() {
		return ErrClosedBufPool
	}
	var count uint32
	// 最大重试次数,如果对当前缓冲池所有的缓冲器操作都达到5次,就创建新的缓冲器并放入缓冲池中
	maxCount := p.BufNum() * 5
	var ok bool
	// 尝试把数据写入缓冲器中.
	for buf := range p.bufChan {
		ok, err = p.putData(buf, datum, &count, maxCount)
		if ok || err != nil {
			break
		}
	}
	return
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

	// 缓冲器已经关闭了
	if err != nil {
		return
	}

	// 缓冲器未放满,递增count.
	*count++

	// 如果尝试向缓冲器放入数据的失败次数达到阈值，
	// 并且池中缓冲器的数量未达到最大值，
	// 那么就尝试创建一个新的缓冲器，先放入数据再把它放入池。
	if *count >= maxCount && p.BufNum() < p.MaxBufNum() {
		p.rwlock.Lock()
		// 确定当前的缓冲器数量是否超过阈值
		if p.BufNum() < p.MaxBufNum() {
			if p.Closed() {
				p.rwlock.Unlock()
				return
			}
			newBuf, _ := NewBuffer(p.bufCap)
			_, err = newBuf.Put(datum)
			if err != nil {
				return
			}
			p.bufChan <- newBuf
			atomic.AddUint32(&p.bufNum, 1)
			atomic.AddUint64(&p.total, 1)
			ok = true
		}
		p.rwlock.Unlock()
		// 清理count
		*count = 0
	}

	return
}

func (p *Pool) Get() (datum interface{}, err error) {
	if p.Closed() {
		return nil, ErrClosedBufPool
	}
	var count uint32
	maxCount := p.BufNum() * 10
	for buf := range p.bufChan {
		datum, err = p.GetData(buf, &count, maxCount)
		if datum != nil || err != nil {
			break
		}
	}
	return
}

func (p *Pool) GetData(buf IBuf, count *uint32, maxCount uint32) (datum interface{}, err error) {
	if p.Closed() {
		return nil, ErrClosedBufPool
	}

	defer func() {
		// 注意这个判断失败次数的代码要放在defer里面,因为这里如果条件达成直接返回,缓冲器是不会放入缓冲池中.
		//	不会执行后面的 p.bufChan <- buf的代码
		// 失败次数超过阈值 && 当前缓冲器已经空了 && 缓冲池的缓冲器数量 > 1.那么就直接关掉当前缓冲器并不归给缓冲池
		if *count >= maxCount && buf.Len() == 0 && p.BufNum() > 1 {
			buf.Close()
			atomic.AddUint32(&p.bufNum, ^uint32(0))
			*count = 0
			return
		}

		p.rwlock.RLock()
		if p.Closed() {
			// 缓冲池已经关闭,拿出来的缓冲器无法放入池子中.这里需要把缓冲器总数量减1
			atomic.AddUint32(&p.bufNum, ^uint32(0))
			err = ErrClosedBufPool
		} else {
			p.bufChan <- buf
		}
		p.rwlock.RUnlock()
	}()

	datum, err = buf.Get()
	if datum != nil {
		atomic.AddUint64(&p.total, ^uint64(0))
		return
	}
	if err != nil {
		return
	}
	// 增加失败次数,并重试.知道达到最大重试次数(maxCount)
	*count++

	return
}

func (p *Pool) Close() bool {
	if atomic.CompareAndSwapUint32(&p.closed, 0, 1) {
		p.rwlock.Lock()
		close(p.bufChan)
		p.rwlock.Unlock()
		return true
	}

	return false
}

func (p *Pool) Closed() bool {
	if atomic.LoadUint32(&p.closed) == 1 {
		return true
	}

	return false
}
