package gin_session

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"sync"
)

// 内存版Session服务
// 接口本来就是引用无需再获取指针

// memSession支持的操作

// memSession 表示用户的Session数据
type memSession struct {
	id      string
	data    map[string]interface{}
	expired int
	rwLock  sync.RWMutex // 读写锁，锁的是上面的Data
	// 过期时间
}

// NewSessionData 构造函数
func NewMemSession(id string) *memSession {
	return &memSession{
		id:   id,
		data: make(map[string]interface{}, 8),
	}
}

func (m *memSession) ID() string {
	return m.id
}

// Get 根据key获取值
func (m *memSession) Get(key string) (value interface{}, err error) {
	// 获取读锁
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	value, ok := m.data[key]
	if !ok {
		err = fmt.Errorf("invalid Key")
		return
	}
	return
}

// Set 根据key获取值
func (m *memSession) Set(key string, value interface{}) {
	// 获取写锁
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	m.data[key] = value
}

// Del 删除Key对应的键值对
func (m *memSession) Del(key string) {
	// 删除key对应的键值对
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	delete(m.data, key)
}

// Save 保存session数据的键值对
func (m *memSession) Save() {
	return
}

func (m *memSession) SetExpired(expired int) {
	m.expired = expired
}

// MemSessionMgr 是一个全局的Session 管理
type MemSessionMgr struct {
	session map[string]Session
	rwLock  sync.RWMutex
}

// NewMemSessionMgr 构造函数
func NewMemSessionMgr() *MemSessionMgr {
	return &MemSessionMgr{
		session: make(map[string]Session, 1024),
	}
}

func (m *MemSessionMgr) Init(addr string, options ...string) (err error) {
	return
}

// GetSessionData 根据传进来的SessionID找到对应的SessionData
func (m *MemSessionMgr) GetSession(sessionID string) (sd Session, err error) {
	// 取之前加锁
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	sd, ok := m.session[sessionID]
	if !ok {
		err = fmt.Errorf("invalid session id")
		return
	}
	return
}

// CreateSession 创建一条Session记录
func (m *MemSessionMgr) CreateSession() (sd Session, err error) {
	// 1. 造一个sessionID
	uuidObj, err := uuid.NewV4()
	if err != nil {
		return
	}
	// 2. 造一个和它对应的Session
	sd = NewMemSession(uuidObj.String())
	// 将创建的SessionData保存起来
	m.session[sd.ID()] = sd
	// 3. 返回SessionData
	return
}
