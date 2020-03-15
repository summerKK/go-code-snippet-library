package session

import "sync"
import uuid "github.com/satori/go.uuid"

type MemorySessionMgr struct {
	sessionMap map[string]ISession
	sync.RWMutex
}

func (m *MemorySessionMgr) Get(sessionid string) (session ISession, err error) {
	m.RLock()
	defer m.RUnlock()
	session, ok := m.sessionMap[sessionid]
	if !ok {
		err = sessionNotExists
	}
	return
}

func (m *MemorySessionMgr) Create() (session ISession, err error) {
	u := uuid.NewV4()
	m.Lock()
	session = NewMemSession()
	m.sessionMap[u.String()] = session
	m.Unlock()
	return
}
