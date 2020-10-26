package fast_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/summerKK/go-code-snippet-library/cache/fast"
)

const maxEntrySize = 256

func BenchmarkFastCache_Get(b *testing.B) {
	cache := fast.NewFastCache(b.N, maxEntrySize, nil)
	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0
		for pb.Next() {
			cache.Set(parallelKey(id, counter), value())
			counter += 1
		}
	})

	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			cache.Get(key(counter))
			counter = counter + 1
		}
	})
}

func key(i int) string {
	return fmt.Sprintf("key-%010d", i)
}

func parallelKey(threadID int, counter int) string {
	return fmt.Sprintf("key-%04d-%06d", threadID, counter)
}

func value() []byte {
	return make([]byte, 100)
}
