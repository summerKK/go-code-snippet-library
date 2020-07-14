package cmap

import "sync/atomic"

// BucketStatus 代表散列桶状态的类型。
type BucketStatus uint8

const (
	// BUCKET_STATUS_NORMAL 代表散列桶正常。
	BUCKET_STATUS_NORMAL BucketStatus = 0
	// BUCKET_STATUS_UNDERWEIGHT 代表散列桶过轻。
	BUCKET_STATUS_UNDERWEIGHT BucketStatus = 1
	// BUCKET_STATUS_OVERWEIGHT 代表散列桶过重。
	BUCKET_STATUS_OVERWEIGHT BucketStatus = 2
)

type IPairRedistributor interface {
	//  根据键-元素对总数和散列桶总数计算并更新阈值
	UpdateThreshold(pairTotal uint64, bucketNum int)
	CheckBucketStatus(pairTotal uint64, bucketSize uint64) (bucketStatus BucketStatus)
	Redistrie(bucketStatus BucketStatus, buckets []IBucket) (newBuckets []IBucket, changed bool)
}

type pariRedistributor struct {
	// loadFactor 代表装载因子。
	loadFactor float64
	// upperThreshold 代表散列桶重量的上阈限。
	// 当某个散列桶的尺寸增至此值时会触发再散列。
	upperThreshold uint64
	// overweightBucketCount 代表过重的散列桶的计数。
	overweightBucketCount uint64
	// emptyBucketCount 代表空的散列桶的计数。
	emptyBucketCount uint64
}

func (p *pariRedistributor) UpdateThreshold(pairTotal uint64, bucketNum int) {
	var average float64
	average = float64(pairTotal / uint64(bucketNum))
	if average < 100 {
		average = 100
	}
	// defer func() {
	// 	fmt.Printf(bucketCountTemplate,
	// 		pairTotal,
	// 		bucketNumber,
	// 		average,
	// 		atomic.LoadUint64(&pr.upperThreshold),
	// 		atomic.LoadUint64(&pr.emptyBucketCount))
	// }()
	atomic.StoreUint64(&p.upperThreshold, uint64(average*p.loadFactor))
}

// bucketStatusTemplate 代表调试用散列桶状态信息模板。
var bucketStatusTemplate = `Check bucket status: 
    pairTotal: %d
    bucketSize: %d
    upperThreshold: %d
    overweightBucketCount: %d
    emptyBucketCount: %d
    bucketStatus: %d
	
`

func (p *pariRedistributor) CheckBucketStatus(pairTotal uint64, bucketSize uint64) (bucketStatus BucketStatus) {
	// defer func() {
	// 	fmt.Printf(bucketStatusTemplate,
	// 		pairTotal,
	// 		bucketSize,
	// 		atomic.LoadUint64(&pr.upperThreshold),
	// 		atomic.LoadUint64(&pr.overweightBucketCount),
	// 		atomic.LoadUint64(&pr.emptyBucketCount),
	// 		bucketStatus)
	// }()
	if bucketSize > DEFAULT_BUCKET_MAX_SIZE || bucketSize >= atomic.LoadUint64(&p.upperThreshold) {
		atomic.AddUint64(&p.overweightBucketCount, 1)
		bucketStatus = BUCKET_STATUS_OVERWEIGHT
		return
	}
	if bucketSize == 0 {
		atomic.AddUint64(&p.emptyBucketCount, 1)
	}
	return
}

// redistributionTemplate 代表重新分配信息模板。
var redistributionTemplate = `Redistributing: 
    bucketStatus: %d
    currentNumber: %d
    newNumber: %d

`

func (p *pariRedistributor) Redistrie(bucketStatus BucketStatus, buckets []IBucket) (newBuckets []IBucket, changed bool) {
	// 当前桶的总数量
	currentNumber := uint64(len(buckets))
	newNumber := currentNumber
	// defer func() {
	// 	fmt.Printf(redistributionTemplate,
	// 		bucketStatus,
	// 		currentNumber,
	// 		newNumber)
	// }()
	switch bucketStatus {
	case BUCKET_STATUS_OVERWEIGHT:
		// 当过载的bucket < 4 可以忽略
		if atomic.LoadUint64(&p.overweightBucketCount)*4 < currentNumber {
			return nil, false
		}
		// buckets扩大一倍
		newNumber = currentNumber << 1
	case BUCKET_STATUS_UNDERWEIGHT:
		if currentNumber < 100 || atomic.LoadUint64(&p.emptyBucketCount)*4 < currentNumber {
			return nil, false
		}
		// buckets缩小1倍
		newNumber = currentNumber >> 1
		if newNumber < 2 {
			newNumber = 2
		}
	default:
		return nil, false
	}
	// 没有扩容
	if newNumber == currentNumber {
		// 参数置零
		atomic.StoreUint64(&p.overweightBucketCount, 0)
		atomic.StoreUint64(&p.emptyBucketCount, 0)
		return nil, false
	}
	var pairs []IPair
	// 取出所有bucket的pair.重新分配bucket
	for _, b := range buckets {
		for e := b.GetFirstPair(); e != nil; e = e.Next() {
			pairs = append(pairs, e)
		}
	}
	// 需要新增bucket
	if newNumber > currentNumber {
		// 新增bucket的时候要情况所有bucket里面的pair,因为要重新分配
		for i := uint64(0); i < currentNumber; i++ {
			buckets[i].Clear(nil)
		}
		// 创建新增的bucket
		for j := newNumber - currentNumber; j > 0; j-- {
			buckets = append(buckets, newBucket())
		}
	} else {
		// 需要缩小bucket
		buckets = make([]IBucket, newNumber)
		for i := uint64(0); i < newNumber; i++ {
			buckets[i] = newBucket()
		}
	}
	// 重新把元素分配到bucket里面去
	var count int
	for _, p := range pairs {
		index := int(p.Hash() % newNumber)
		b := buckets[index]
		_, _ = b.Put(p, nil)
		count++
	}
	atomic.StoreUint64(&p.overweightBucketCount, 0)
	atomic.StoreUint64(&p.emptyBucketCount, 0)
	return buckets, true
}

// newDefaultPairRedistributor 会创建一个PairRedistributor类型的实例。
// 参数loadFactor代表散列桶的负载因子。
// 参数bucketNumber代表散列桶的数量。
func newDefaultPairRedistributor(loadFactor float64, bucketNumber int) IPairRedistributor {
	if loadFactor <= 0 {
		loadFactor = DEFAULT_BUCKET_LOAD_FACTOR
	}
	pr := &pariRedistributor{}
	pr.loadFactor = loadFactor
	pr.UpdateThreshold(0, bucketNumber)
	return pr
}
