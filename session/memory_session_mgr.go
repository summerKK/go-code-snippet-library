package session

import "sync"
import uuid "github.com/satori/go.uuid"

type MemorySessionMgr struct {
	sessionMap map[string]ISession
	sync.RWMutex
}

func NewMemSessionMgr(addr string, options ...string) (session *MemorySessionMgr, err error) {
	session = &MemorySessionMgr{
		sessionMap: make(map[string]ISession, 8),
		RWMutex:    sync.RWMutex{},
	}
	err = session.Init(addr, options...)
	return
}

func (m *MemorySessionMgr) Init(addr string, options ...string) (err error) {
	return
}

func (m *MemorySessionMgr) Get(sessionid string) (session ISession, err error) {
	m.RLock()
	defer m.RUnlock()
	session, ok := m.sessionMap[sessionid]
	if !ok {
		err = errSessionNotExists
	}
	return
}

func (m *MemorySessionMgr) Create() (session ISession, err error) {
	u := uuid.NewV4()
	m.Lock()
	session = NewMemSession(u.String())
	m.sessionMap[u.String()] = session
	m.Unlock()
	return
}
