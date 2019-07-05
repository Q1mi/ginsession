package ginsession

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

const (
	// SessionCookieName 保存在Cookie中的Session标识
	SessionCookieName = "session_id"
	// SessionName 保存在请求上下文中的Session标识
	SessionName       = "session"
)


// Session session
type Session interface {
	ID() string
	Get(string) (interface{}, error)
	Set(string, interface{})
	Del(string)
	Save()
	SetExpired(int)
}

// SessionMgr 全局的Session管理者
type SessionMgr interface {
	Init(addr string, options ...string) error // 初始化对应的Session存储
	GetSession(string) (Session, error)        // 根据SessionID获取对应的Session
	CreateSession() (Session, error)           // 创建一个新的Session记录
}

// Options Cookie Options
type Options struct {
	Path   string
	Domain string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

// CreateSessionMgr 用于初始化一个SessionMgr
func CreateSessionMgr(name string, addr string, options ...string) (sm SessionMgr, err error) {

	switch name {
	case "memory":
		sm = NewMemSessionMgr()
	case "redis":
		sm = NewRedisSessionMgr()
	default:
		err = fmt.Errorf("unsupport %s\n", name)
		return
	}
	err = sm.Init(addr, options...) // 初始化SessionMgr
	return
}

// SessionMiddleware gin middleware
func SessionMiddleware(sm SessionMgr, options Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 0. 请求进来之后给每个请求分配个session
		// 后续处理函数只需通过c.Get("session")即可操作该请求对应的Session
		var session Session
		// 1. 先从请求的Cookie中获取session_id
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			// 取不到session_id，创建一份新的Session
			log.Printf("get session_id from Cookie failed，err:%v\n", err)
			session, _ = sm.CreateSession()
			sessionID = session.ID()
		}
		log.Printf("SessionID:%v\n", sessionID)
		session, err = sm.GetSession(sessionID)
		if err != nil {
			// 根据sessionID取不到Session数据
			log.Printf("get Session by %s failed，err:%v\n", sessionID, err)
			session, _ = sm.CreateSession()
			sessionID = session.ID()
		}
		session.SetExpired(options.MaxAge)

		c.Set(SessionName, session)
		// 回写Cookie要在handlerFunc返回前
		c.SetCookie(SessionCookieName, sessionID, options.MaxAge, options.Path, options.Domain, options.Secure, options.HttpOnly)
		c.Next()
	}
}
