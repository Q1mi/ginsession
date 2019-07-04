package gin_session

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	"log"
	"strconv"
	"sync"
	"time"
)

// redis版Session服务

type redisSession struct {
	id string
	data map[string]interface{}
	loadFlag sync.Once
	loadFunc func()
	modifyFlag bool // 是否修改的标志位
	expired int // 超时时间
	rwLock sync.RWMutex
	client *redis.Client
}


// NewRedisSession 是redisSession构造函数
func NewRedisSession(id string, client *redis.Client)(session Session){
	r:= &redisSession{
		id: id,
		data: make(map[string]interface{}, 8),
		client: client,
	}
	r.loadFunc = func(){
		loadFromRedis(r)
	}
	return r
}

func (r *redisSession) ID() string {
	return r.id
}

// load session data from redis
func loadFromRedis(r *redisSession){
	data, err := r.client.Get(r.id).Result()
	if err != nil {
		r.data = make(map[string]interface{})
		return
	}
	// unmarshal
	err = json.Unmarshal([]byte(data), &r.data)
	if err != nil {
		r.data = make(map[string]interface{})
		return
	}
}

func (r *redisSession) Get(key string) (value interface{}, err error) {
	r.loadFlag.Do(r.loadFunc) // 加载一次

	r.rwLock.RLock()
	defer r.rwLock.RUnlock()
	value, ok := r.data[key]
	if !ok{
		err = fmt.Errorf("invalid Key")
		return
	}
	return

}

func (r *redisSession) Set(key string,value interface{}) {
	// 获取写锁
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	r.data[key] = value
	r.modifyFlag = true
}

func (r *redisSession) Del(key string) {
	// 删除key对应的键值对
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	delete(r.data, key)
}

func (r *redisSession)SetExpired(expired int){
	r.expired = expired
}


func (r *redisSession) Save() {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	if !r.modifyFlag {
		return
	}
	data, err := json.Marshal(r.data)
	if err != nil {
		log.Fatalf("marshal r.data failed, err:%v\n", err)
		return
	}
	r.client.Set(r.id, data, time.Second*time.Duration(r.expired))
	r.modifyFlag = false
}



type redisSessionMgr struct {
	session map[string]Session
	rwLock sync.RWMutex
	client *redis.Client
}

func NewRedisSessionMgr()*redisSessionMgr{
	return &redisSessionMgr{
		session:make(map[string]Session, 1024),
	}
}


func (r *redisSessionMgr) Init(addr string, options ...string) (err error){
	var (
		password string
		db int
	)

	if len(options) == 1{
		password = options[0]
	}
	if len(options) == 2{
		password = options[0]
		tmpDB, err := strconv.ParseInt(options[1], 10, 64)
		if err != nil {
			log.Fatalln("invalid redis DB params")
		}
		db = int(tmpDB)
	}
	r.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,  // use default DB
	})

	_, err = r.client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisSessionMgr) GetSession(sessionID string) (sd Session, err error) {
	// 取之前加锁
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()
	sd, ok := r.session[sessionID]
	if !ok {
		err = fmt.Errorf("invalid session id")
		return
	}
	return
}

// CreateSession 创建一条Session记录
func (r *redisSessionMgr) CreateSession() (sd Session,err error) {
	// 1. 造一个sessionID
	uuidObj, err := uuid.NewV4()
	if err != nil {
		return
	}
	// 2. 造一个和它对应的Session
	sd = NewRedisSession(uuidObj.String(), r.client)
	// 将创建的SessionData保存起来
	r.session[sd.ID()] = sd
	// 3. 返回SessionData
	return
}





