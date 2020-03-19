package session

import "sync"

type MemSession struct {
	id   string
	data map[string]interface{}
	sync.RWMutex
	status int
}

func NewMemSession(id string) *MemSession {
	return &MemSession{
		id:      id,
		data:    make(map[string]interface{}, 8),
		RWMutex: sync.RWMutex{},
	}
}

func (m *MemSession) IsModify() bool {
	if m.status == status_modify {
		return true
	}

	return false
}

func (m *MemSession) Id() string {
	return m.id
}

func (m *MemSession) Set(key string, value interface{}) (err error) {
	m.Lock()
	defer m.Unlock()
	m.data[key] = value
	return
}

func (m *MemSession) Get(key string) (value interface{}, err error) {
	m.RLock()
	defer m.RUnlock()
	value, ok := m.data[key]
	if !ok {
		err = errKeyNotExistsInSession
	}
	return
}

func (m *MemSession) Del(key string) (err error) {
	m.Lock()
	defer m.Unlock()
	delete(m.data, key)
	return
}

func (m *MemSession) Save() (err error) {
	return
}
