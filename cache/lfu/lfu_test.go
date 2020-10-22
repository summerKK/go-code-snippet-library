package lfu

import (
	"testing"

	"github.com/matryer/is"
)

func TestLfu_Set(t *testing.T) {
	assertIs := is.New(t)

	cache := New(24, nil)
	cache.DelOldest()

	cache.Set("k1", 1)
	v := cache.Get("k1")
	assertIs.Equal(v, 1)

	cache.Del("k1")
	assertIs.Equal(cache.Len(), 0)
}

func TestLfu_DelOldest(t *testing.T) {
	assertIs := is.New(t)

	cache := New(24, nil)
	cache.DelOldest()

	cache.Set("k1", 1)
	cache.Set("k2", 2)

	cache.Get("k1")
	cache.Get("k1")
	cache.Get("k2")

	cache.DelOldest()

	assertIs.Equal(cache.Get("k2"), nil)
}
