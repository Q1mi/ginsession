package ginsession

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// AuthMiddleware 认证中间件
// 从请求的session Data中获取isLogin
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sd := c.MustGet("session").(Session)
		log.Printf("get session:%v from AuthMD\n", sd)
		isLogin, err := sd.Get("isLogin")
		log.Println(isLogin, err, err == nil && isLogin.(bool))
		if err == nil && isLogin.(bool) {
			// 是登录状态
			c.Next()
		} else {
			c.Abort()
			c.Redirect(http.StatusFound, "/login")
		}
	}
}
