# gin-session
A session middleware for gin framework.


# Usage

Download and install:
```bash
go get github.com/Q1mi/gin-session
```

Import it in you code:

```bash
import "github.com/Q1mi/gin_session"
```

## Examples


Redis-based

```go
package main

import (
	"github.com/Q1mi/gin_session"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

var log = logrus.New()

func main(){
	log.SetReportCaller(true)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	mgrObj, err := ice_cube2.CreateSessionMgr("redis", "localhost:6379")
	if err != nil {
		log.Errorf("create manager obj failed, err:%v\n", err)
		return
	}
	sm := go_session.SessionMiddleware(mgrObj, ice_cube2.Options{
		Path: "/",
		Domain: "127.0.0.1",
		MaxAge: 300,
		Secure:false,
		HttpOnly:true,
	})
	r.Use(sm)

    r.GET("/incr", func(c *gin.Context) {
            session := c.MustGet("session").(Session)
            var count int
            v := session.Get("count")
            if v == nil {
                count = 0
            } else {
                count = v.(int)
                count++
            }
            session.Set("count", count)
            session.Save()
            c.JSON(200, gin.H{"count": count})
        })

	r.Run()
}
```

memory_based:

```go
package main

import (
	"github.com/Q1mi/gin_session"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

var log = logrus.New()

func main(){
	log.SetReportCaller(true)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	mgrObj, err := go_session.CreateSessionMgr("memory", "")
	if err != nil {
		log.Errorf("create manager obj failed, err:%v\n", err)
		return
	}
	sm := go_session.SessionMiddleware(mgrObj, ice_cube2.Options{
		Path: "/",
		Domain: "127.0.0.1",
		MaxAge: 300,
		Secure:false,
		HttpOnly:true,
	})
	r.Use(sm)

	r.GET("/incr", func(c *gin.Context) {
    		session := c.MustGet("session").(Session)
    		var count int
    		v := session.Get("count")
    		if v == nil {
    			count = 0
    		} else {
    			count = v.(int)
    			count++
    		}
    		session.Set("count", count)
    		session.Save()
    		c.JSON(200, gin.H{"count": count})
    	})


	r.Run()
}
```