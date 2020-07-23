package limiter

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

type MethodLimiter struct {
	*Limiter
}

func NewMethodLimiter() *MethodLimiter {
	return &MethodLimiter{
		Limiter: &Limiter{
			limitBuckets: make(map[string]*ratelimit.Bucket),
		},
	}
}

func (m MethodLimiter) Key(c *gin.Context) string {
	uri := c.Request.RequestURI
	index := strings.Index(uri, "?")
	if index == -1 {
		return uri
	}

	return uri[:index]
}

func (m MethodLimiter) GetBucket(key string) (*ratelimit.Bucket, bool) {
	bucket, ok := m.limitBuckets[key]

	return bucket, ok
}

func (m MethodLimiter) AddBuckets(rules ...BucketRule) Contract {
	for _, rule := range rules {
		if _, ok := m.limitBuckets[rule.Key]; !ok {
			m.limitBuckets[rule.Key] = ratelimit.NewBucketWithQuantum(rule.FillInterval, rule.Capacity, rule.Quantum)
		}
	}

	return m
}
