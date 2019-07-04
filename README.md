# gin_session
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


### Redis-based

```go

package main

import (
	"fmt"
	"log"

	"net/http"

	"github.com/Q1mi/gin_session"
	"github.com/gin-gonic/gin"
)


func main(){
	r := gin.Default()
	mgrObj, err := gin_session.CreateSessionMgr("redis", "localhost:6379")
	if err != nil {
		log.Fatalf("create manager obj failed, err:%v\n", err)
		return
	}
	sm := gin_session.SessionMiddleware(mgrObj, gin_session.Options{
		Path: "/",
		Domain: "127.0.0.1",
		MaxAge: 60,
		Secure:false,
		HttpOnly:true,
	})

	r.Use(sm)

	r.GET("/incr", func(c *gin.Context) {
		session := c.MustGet("session").(gin_session.Session)
		fmt.Printf("%#v\n", session)
		var count int
		v, err := session.Get("count")
		if err != nil{
			log.Printf("get count from session failed, err:%v\n", err)
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.String(http.StatusOK, "count:%v", count)
	})

	r.Run()
}
```

### memory_based:

```go

package main

import (
	"fmt"
	"log"

	"net/http"

	"github.com/Q1mi/gin_session"
	"github.com/gin-gonic/gin"
)


func main(){
	r := gin.Default()
	mgrObj, err := gin_session.CreateSessionMgr("memory", "")
	if err != nil {
		log.Fatalf("create manager obj failed, err:%v\n", err)
		return
	}
	sm := gin_session.SessionMiddleware(mgrObj, gin_session.Options{
		Path: "/",
		Domain: "127.0.0.1",
		MaxAge: 60,
		Secure:false,
		HttpOnly:true,
	})

	r.Use(sm)

	r.GET("/incr", func(c *gin.Context) {
		session := c.MustGet("session").(gin_session.Session)
		fmt.Printf("%#v\n", session)
		var count int
		v, err := session.Get("count")
		if err != nil{
			log.Printf("get count from session failed, err:%v\n", err)
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.String(http.StatusOK, "count:%v", count)
	})

	r.Run()
}
```

## TODO

1. Add more support...