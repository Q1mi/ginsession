package ginsession

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

const (
	// SessionCookieName key of SessionID store in Cookie
	SessionCookieName = "session_id"
	// SessionName key of session in gin.context
	SessionContextName = "session"
)

// Session stores values for a session
type Session interface {
	ID() string
	Get(string) (interface{}, error)
	Set(string, interface{})
	Del(string)
	Save()
	SetExpired(int)
}

// SessionMgr a session manager
type SessionMgr interface {
	Init(addr string, options ...string) error // init the session store
	GetSession(string) (Session, error)        // get the session by sessionID
	CreateSession() (Session)           // create a new session
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

// CreateSessionMgr create a sessionMgr by given name
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
	err = sm.Init(addr, options...) // init the SessionMgr
	return
}

// SessionMiddleware gin middleware
func SessionMiddleware(sm SessionMgr, options Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get or create a session for every request come in
		// so next handlerFunc can get the session by c.Get(SessionContextName)
		var session Session
		// try to get sessionID from cookie
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			// can't get sessionID from Cookie, need to create a new session
			log.Printf("get session_id from Cookie failed，err:%v\n", err)
			session = sm.CreateSession()
			sessionID = session.ID()
		}else{
			log.Printf("SessionID:%v\n", sessionID)
			session, err = sm.GetSession(sessionID)
			if err != nil {
				// can't get session by the sessionID
				log.Printf("get Session by %s failed，err:%v\n", sessionID, err)
				session = sm.CreateSession()
				sessionID = session.ID()
			}
		}
		session.SetExpired(options.MaxAge) // set session expired time
		c.Set(SessionContextName, session)
		// must write cookie before handlerFunc return
		c.SetCookie(SessionCookieName, sessionID, options.MaxAge, options.Path, options.Domain, options.Secure, options.HttpOnly)
		c.Next()
	}
}
