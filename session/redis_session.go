package session

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"sync"
)

func NewRedisSession(id string, pool *redis.Pool) *RedisSession {
	return &RedisSession{
		id:      id,
		buf:     make(map[string]interface{}, 8),
		RWMutex: sync.RWMutex{},
		pool:    pool,
		status:  status_init,
	}
}

type RedisSession struct {
	id  string
	buf map[string]interface{}
	sync.RWMutex
	pool   *redis.Pool
	status int
}

func (r *RedisSession) IsModify() bool {
	if r.status == status_modify {
		return true
	}

	return false
}

func (r *RedisSession) Id() string {
	return r.id
}

func (r *RedisSession) Set(key string, value interface{}) (err error) {
	r.Lock()
	defer r.Unlock()
	r.buf[key] = value
	r.status = status_modify
	return
}

func (r *RedisSession) Get(key string) (value interface{}, err error) {
	r.RLock()
	defer r.RUnlock()
	if r.status == status_init {
		r.loadData()
		r.status = status_modify
	}
	value, ok := r.buf[key]
	if !ok {
		err = errKeyNotExistsInSession
	}
	return
}

func (r *RedisSession) Del(key string) (err error) {
	r.Lock()
	defer r.Unlock()
	delete(r.buf, key)
	r.status = status_modify
	return
}

func (r *RedisSession) Save() (err error) {
	r.RLock()
	defer r.RUnlock()
	if r.status == status_init {
		return
	}
	bytes, err := json.Marshal(r.buf)
	if err != nil {
		return
	}
	conn := r.pool.Get()
	_, err = conn.Do("set", r.id, string(bytes))
	return
}

func (r *RedisSession) loadData() (err error) {
	conn := r.pool.Get()
	replay, err := conn.Do("get", r.id)
	if err != nil {
		return
	}
	data, err := redis.String(replay, err)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(data), &r.buf)
	return
}
