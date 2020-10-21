package fifo

import (
	"testing"

	"github.com/matryer/is"
)

func TestFifo_Set(t *testing.T) {
	assertIs := is.New(t)

	cache := New(24, nil)
	cache.DelOldest()
	cache.Set("k1", 1)
	v := cache.Get("k1")
	assertIs.Equal(v, 1)

	cache.Del("k1")
	assertIs.Equal(0, cache.Len())
}

func TestFifo_OnEvicted(t *testing.T) {
	assertIs := is.New(t)

	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}

	cache := New(24, onEvicted)
	cache.Set("k1", 1)
	cache.Set("k2", 2)
	cache.Set("k3", 3)

	cache.Del("k1")
	cache.Del("k2")

	assertIs.Equal(keys, []string{"k1", "k2"})
	assertIs.Equal(1, cache.Len())
}
