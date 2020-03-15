package session

import "sync"

type MemSession struct {
	id   string
	data map[string]interface{}
	sync.RWMutex
}

func NewMemSession() *MemSession {
	return &MemSession{
		data:    make(map[string]interface{}, 8),
		RWMutex: sync.RWMutex{},
	}
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
		err = keyNotExistsInSession
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
