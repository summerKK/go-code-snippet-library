package cache

import (
	"log"
	"sync"
)

const DefaultMaxBytes = 1 << 29

type Cache interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Del(key string)
	DelOldest()
	Len() int
}

type safeCache struct {
	m     sync.RWMutex
	cache Cache

	nhit, nget int
}

func NewSafeCache(cache Cache) *safeCache {
	return &safeCache{
		cache: cache,
	}
}

func (s *safeCache) Set(key string, value interface{}) {
	s.m.Lock()
	defer s.m.Unlock()

	s.cache.Set(key, value)
}

func (s *safeCache) Get(key string) interface{} {
	s.m.RLock()
	defer s.m.RUnlock()

	// 获取
	s.nget++

	if s.cache == nil {
		return nil
	}

	value := s.cache.Get(key)
	// 命中缓存
	if value != nil {
		log.Println("[Cache] hit")
		s.nhit++
	}

	return value
}

func (s *safeCache) stat() *Stat {
	s.m.RLock()
	defer s.m.RUnlock()

	return &Stat{
		NHit: s.nhit,
		NGet: s.nget,
	}
}

type Stat struct {
	NHit, NGet int
}
