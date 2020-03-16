package session

import (
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

type RedisSessionMgr struct {
	addr     string
	password string
	pool     *redis.Pool
	sync.RWMutex
	sessionMap map[string]ISession
}

func NewRedisSessionMgr(addr string, options ...string) (session *RedisSessionMgr, err error) {
	session = &RedisSessionMgr{
		sessionMap: make(map[string]ISession, 8),
		RWMutex:    sync.RWMutex{},
	}
	err = session.Init(addr, options...)
	return
}

func (r *RedisSessionMgr) Init(addr string, options ...string) (err error) {
	r.addr = addr
	if len(options) > 0 {
		r.password = options[0]
	}
	r.newPool()
	return
}

func (r *RedisSessionMgr) Get(sessionId string) (session ISession, err error) {
	r.RLock()
	defer r.RUnlock()
	session, ok := r.sessionMap[sessionId]
	if !ok {
		err = errSessionNotExists
	}
	return
}

func (r *RedisSessionMgr) Create() (session ISession, err error) {
	u := uuid.NewV4()
	r.Lock()
	session = NewRedisSession(u.String(), r.pool)
	r.sessionMap[u.String()] = session
	r.Unlock()
	return
}

func (r *RedisSessionMgr) newPool() {
	r.pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", r.addr)
			if _, err := c.Do("AUTH", r.password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("ping")
			return err
		},
	}
}
