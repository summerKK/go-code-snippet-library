package session

import "sync"
import uuid "github.com/satori/go.uuid"

type Manager struct {
	sessionMap map[string]ISession
	sync.RWMutex
}

func (m *Manager) Get(sessionid string) (session ISession, err error) {
	m.RLock()
	defer m.RUnlock()
	session, ok := m.sessionMap[sessionid]
	if !ok {
		err = sessionNotExists
	}
	return
}

func (m *Manager) Create() (session ISession, err error) {
	u := uuid.NewV4()
	m.Lock()
	session = NewMemSession()
	m.sessionMap[u.String()] = session
	m.Unlock()
	return
}
