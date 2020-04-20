package cmap

import "sync"

type ISegment interface {
	// bool代表是否新增了新数据
	Put(p IPair) (bool, error)
	Get(key string) IPair
	GetWithHash(key string, keyHash uint64) IPair
	Delete(key string) bool
	Szie() uint64
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

func (s *segment) Put(p IPair) (bool, error) {
	panic("implement me")
}

func (s *segment) Get(key string) IPair {
	panic("implement me")
}

func (s *segment) GetWithHash(key string, keyHash uint64) IPair {
	panic("implement me")
}

func (s *segment) Delete(key string) bool {
	panic("implement me")
}

func (s *segment) Szie() uint64 {
	panic("implement me")
}
