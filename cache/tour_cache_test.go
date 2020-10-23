package cache_test

import (
	"log"
	"sync"
	"testing"

	"github.com/matryer/is"
	"github.com/summerKK/go-code-snippet-library/cache"
	"github.com/summerKK/go-code-snippet-library/cache/lru"
)

func TestTourCache_Get(t *testing.T) {
	db := map[string]string{
		"key1": "str1",
		"key2": "str2",
		"key3": "str3",
		"key4": "str4",
		"key5": "str5",
	}

	getter := cache.GetFunc(func(key string) interface{} {
		log.Println("[From DB] find key", key)

		if v, ok := db[key]; ok {
			return v
		}

		return nil
	})

	tourCache := cache.NewTourCache(getter, lru.New(0, nil))

	assertIs := is.New(t)

	var wg sync.WaitGroup

	for k, v := range db {
		wg.Add(1)
		go func(k, v string) {
			defer wg.Done()
			assertIs.Equal(v, tourCache.Get(k))
			assertIs.Equal(v, tourCache.Get(k))
		}(k, v)
	}

	wg.Wait()

	assertIs.Equal(tourCache.Get("nil"), nil)
	assertIs.Equal(tourCache.Get("nil"), nil)

	assertIs.Equal(tourCache.Stat().NGet, 12)
	assertIs.Equal(tourCache.Stat().NHit, 5)
}
