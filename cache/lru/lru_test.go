package lru

import (
	"testing"

	"github.com/matryer/is"
)

func TestLru_Set(t *testing.T) {
	assertIs := is.New(t)

	cache := New(32, nil)
	cache.Set("k1", 1)
	cache.Set("k2", 2)
	cache.Set("k3", 2)

	assertIs.Equal(3, cache.Len())

	cache.Get("k1")
	cache.DelOldest()
	cache.DelOldest()

	assertIs.Equal(nil, cache.Get("k2"))
	assertIs.Equal(nil, cache.Get("k3"))
	assertIs.Equal(1, cache.Len())
	assertIs.Equal(1, cache.Get("k1"))
}
