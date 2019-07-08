package ginsession

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"sync"
)

// memory-based session serve


// memSession session store in memory
type memSession struct {
	id      string
	data    map[string]interface{}
	expired int
	rwLock  sync.RWMutex
}

// NewMemSession constructor
func NewMemSession(id string) *memSession {
	return &memSession{
		id:   id,
		data: make(map[string]interface{}, 8),
	}
}

func (m *memSession) ID() string {
	return m.id
}

func (m *memSession) Get(key string) (value interface{}, err error) {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	value, ok := m.data[key]
	if !ok {
		err = fmt.Errorf("invalid Key")
		return
	}
	return
}

func (m *memSession) Set(key string, value interface{}) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	m.data[key] = value
}

func (m *memSession) Del(key string) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	delete(m.data, key)
}

// Save save session data to store
func (m *memSession) Save() {
	return
}

// SetExpired set expired time
func (m *memSession) SetExpired(expired int) {
	m.expired = expired
}

// MemSessionMgr a memory-based manager
type MemSessionMgr struct {
	session map[string]Session
	rwLock  sync.RWMutex
}

// NewMemSessionMgr constructor
func NewMemSessionMgr() *MemSessionMgr {
	return &MemSessionMgr{
		session: make(map[string]Session, 1024),
	}
}

// Init init...
func (m *MemSessionMgr) Init(addr string, options ...string) (err error) {
	return
}

// GetSession get the session by sessionID
func (m *MemSessionMgr) GetSession(sessionID string) (sd Session, err error) {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	sd, ok := m.session[sessionID]
	if !ok {
		err = fmt.Errorf("invalid session id")
		return
	}
	return
}

// CreateSession create a new session
func (m *MemSessionMgr) CreateSession() (sd Session) {
	// 1. generate a sessionID
	sessionID := uuid.NewV4().String()
	// 2. create a new memory-based session
	sd = NewMemSession(sessionID)
	m.session[sd.ID()] = sd
	return
}

// Clear delete the request session in sessionMgr
func (m *MemSessionMgr)Clear(sessionID string){
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	delete(m.session, sessionID)
}