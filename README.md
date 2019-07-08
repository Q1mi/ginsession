[![Go Report Card](https://goreportcard.com/badge/github.com/Q1mi/gin-session)](https://goreportcard.com/report/github.com/Q1mi/gin-session)

# ginsession
A session middleware for gin framework.


# Usage

Download and install:
```bash
go get github.com/Q1mi/ginsession
```

Import it in you code:

```bash
import "github.com/Q1mi/ginsession"
```

## Examples


### Redis-based

```go

package main

import (
	"fmt"
	"log"

	"net/http"

	"github.com/Q1mi/ginsession"
	"github.com/gin-gonic/gin"
)


func main(){
	r := gin.Default()
	mgrObj, err := ginsession.CreateSessionMgr("redis", "localhost:6379")
	if err != nil {
		log.Fatalf("create manager obj failed, err:%v\n", err)
		return
	}
	sm := ginsession.SessionMiddleware(mgrObj, ginsession.Options{
		Path: "/",
		Domain: "127.0.0.1",
		MaxAge: 60,
		Secure:false,
		HttpOnly:true,
	})

	r.Use(sm)

	r.GET("/incr", func(c *gin.Context) {
		session := c.MustGet("session").(ginsession.Session)
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

	"github.com/Q1mi/ginsession"
	"github.com/gin-gonic/gin"
)


func main(){
	r := gin.Default()
	mgrObj, err := ginsession.CreateSessionMgr("memory", "")
	if err != nil {
		log.Fatalf("create manager obj failed, err:%v\n", err)
		return
	}
	sm := ginsession.SessionMiddleware(mgrObj, ginsession.Options{
		Path: "/",
		Domain: "127.0.0.1",
		MaxAge: 60,
		Secure:false,
		HttpOnly:true,
	})

	r.Use(sm)

	r.GET("/incr", func(c *gin.Context) {
		session := c.MustGet("session").(ginsession.Session)
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
