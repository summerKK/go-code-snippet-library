package cmap

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type ISegment interface {
	// bool代表是否新增了新数据
	Put(p IPair) (bool, error)
	Get(key string) IPair
	GetWithHash(key string, keyHash uint64) IPair
	Delete(key string) bool
	Size() uint64
}

type segment struct {
	buckets           []IBucket
	bucketsLen        int
	pairTotal         uint64
	pairRedistributor IPairRedistributor
	lock              sync.Mutex
}

func newSegment(bucketNum int, pairRedistributor IPairRedistributor) ISegment {
	if bucketNum <= 0 {
		bucketNum = DEFAULT_BUCKET_NUMBER
	}
	if pairRedistributor == nil {
		pairRedistributor = newDefaultPairRedistributor(DEFAULT_BUCKET_LOAD_FACTOR, bucketNum)
	}
	buckets := make([]IBucket, bucketNum)
	for i := 0; i < bucketNum; i++ {
		buckets[i] = newBucket()
	}

	return &segment{
		buckets:           buckets,
		bucketsLen:        bucketNum,
		pairRedistributor: pairRedistributor,
	}
}

func (s *segment) Put(p IPair) (ok bool, err error) {
	// 这里需要原子操作.因为涉及到bucket的重新分配
	s.lock.Lock()
	defer s.lock.Unlock()
	b := s.buckets[int(p.Hash()%uint64(s.bucketsLen))]
	ok, err = b.Put(p, nil)
	if ok {
		newTotal := atomic.AddUint64(&s.pairTotal, 1)
		err = s.redistrubite(newTotal, b.Size())
	}

	return
}

func (s *segment) redistrubite(pairTotal, bucketSize uint64) (err error) {
	// 这里添加defer的原因是保证程序的健壮.因为这个`pairRedistributor`可能是外部传进来的
	// 当然也可以用默认的.因为可能是外部传入的,所以不可控.加defer.是程序更加健壮
	defer func() {
		if p := recover(); p != nil {
			if pErr, ok := p.(error); ok {
				err = newPairRedistributorError(pErr.Error())
			} else {
				err = newPairRedistributorError(fmt.Sprintf("%s", p))
			}
		}
	}()
	// 更新阀值
	s.pairRedistributor.UpdateThreshold(pairTotal, s.bucketsLen)
	// 查看bucket的状态(是否超重)
	bucketStatus := s.pairRedistributor.CheckBucketStatus(pairTotal, bucketSize)
	// 重新分配bucket
	newBuckets, changed := s.pairRedistributor.Redistrie(bucketStatus, s.buckets)
	if changed {
		s.buckets = newBuckets
		s.bucketsLen = len(s.buckets)
	}

	return
}

func (s *segment) Get(key string) IPair {
	return s.GetWithHash(key, hash(key))
}

func (s *segment) GetWithHash(key string, keyHash uint64) IPair {
	s.lock.Lock()
	b := s.buckets[int(keyHash%uint64(s.bucketsLen))]
	s.lock.Unlock()
	return b.Get(key)
}

func (s *segment) Delete(key string) bool {
	s.lock.Lock()
	b := s.buckets[int(hash(key)%uint64(s.bucketsLen))]
	ok := b.Delete(key, nil)
	if ok {
		total := atomic.AddUint64(&s.pairTotal, ^uint64(0))
		_ = s.redistrubite(total, b.Size())
	}
	s.lock.Unlock()

	return ok
}

func (s *segment) Size() uint64 {
	return atomic.LoadUint64(&s.pairTotal)
}
