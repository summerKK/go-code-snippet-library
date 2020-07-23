package limiter

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

type Contract interface {
	Key(c *gin.Context) string
	GetBucket(key string) (*ratelimit.Bucket, bool)
	AddBuckets(rules ...BucketRule) Contract
}

type Limiter struct {
	limitBuckets map[string]*ratelimit.Bucket
}

type BucketRule struct {
	Key string
	// 间隔多久时间放N个令牌
	FillInterval time.Duration
	// 令牌桶的容量
	Capacity int64
	// 每次到达间隔时间后所放的具体令牌数量
	Quantum int64
}
