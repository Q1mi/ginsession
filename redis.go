package ginsession

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	"log"
	"strconv"
	"sync"
	"time"
)

// redis-based session serve

// redisSession redis-based session
type redisSession struct {
	id         string
	data       map[string]interface{}
	modifyFlag bool
	expired    int
	rwLock     sync.RWMutex
	client     *redis.Client
}

// NewRedisSession redisSession constructor
func NewRedisSession(id string, client *redis.Client) (session Session) {
	r := &redisSession{
		id:     id,
		data:   make(map[string]interface{}, 8),
		client: client,
	}
	return r
}

func (r *redisSession) ID() string {
	return r.id
}

// load session data from redis
func (r *redisSession)Load()(err error) {
	data, err := r.client.Get(r.id).Bytes()
	if err != nil {
		log.Printf("get session data from redis by %s failed, err:%v\n", r.id, err)
		return
	}
	// unmarshal
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err = dec.Decode(&r.data)
	if err != nil {
		log.Printf("gob decode session data failed, err:%v\n", err)
		return
	}
	return
}


func (r *redisSession) Get(key string) (value interface{}, err error) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()
	value, ok := r.data[key]
	if !ok {
		err = fmt.Errorf("invalid Key")
		return
	}
	return

}

func (r *redisSession) Set(key string, value interface{}) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	r.data[key] = value
	r.modifyFlag = true
}

func (r *redisSession) Del(key string) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	delete(r.data, key)
	r.modifyFlag = true
}

func (r *redisSession) SetExpired(expired int) {
	r.expired = expired
}

func (r *redisSession) Save() {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	if !r.modifyFlag {
		return
	}
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(r.data)
	data, err := json.Marshal(r.data)
	if err != nil {
		log.Fatalf("gob encode r.data failed, err:%v\n", err)
		return
	}
	r.client.Set(r.id, data, time.Second*time.Duration(r.expired))
	log.Printf("set data: %v to redis.\n", data)
	r.modifyFlag = false
}

type redisSessionMgr struct {
	session map[string]Session
	rwLock  sync.RWMutex
	client  *redis.Client
}

// NewRedisSessionMgr redis-based sessionMgr constructor
func NewRedisSessionMgr() *redisSessionMgr {
	return &redisSessionMgr{
		session: make(map[string]Session, 1024),
	}
}

func (r *redisSessionMgr) Init(addr string, options ...string) (err error) {
	var (
		password string
		db       int
	)

	if len(options) == 1 {
		password = options[0]
	}
	if len(options) == 2 {
		password = options[0]
		db, err = strconv.Atoi(options[1])
		if err != nil {
			log.Fatalln("invalid redis DB param")
		}
	}
	r.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	_, err = r.client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

// GetSession load session data and add to sessionMgr
func (r *redisSessionMgr) GetSession(sessionID string) (sd Session, err error) {
	sd = NewRedisSession(sessionID, r.client)
	err = sd.Load()
	if err != nil {
		return
	}
	r.rwLock.RLock()
	r.session[sessionID] = sd
	r.rwLock.RUnlock()
	return
}

// CreateSession create a new session
func (r *redisSessionMgr) CreateSession() (sd Session) {
	// 1. generate a sessionID by uuid
	sessionID := uuid.NewV4().String()
	// 2. create a session use sessionID
	sd = NewRedisSession(sessionID, r.client)
	// 3. save the session in mgr
	r.session[sd.ID()] = sd
	return
}

// Clear delete the request session in sessionMgr
func (r *redisSessionMgr) Clear(sessionID string){
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	delete(r.session, sessionID)
}
