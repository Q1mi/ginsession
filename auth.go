package gin_session

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


// AuthMiddleware 认证中间件
// 从请求的session Data中获取isLogin
func AuthMiddleware()gin.HandlerFunc{
	return func(c *gin.Context){
		sd := c.MustGet("session").(Session)
		log.Debugf("AuthMD:%#v\n", sd)
		isLogin , err := sd.Get("isLogin")
		log.Debug(isLogin, err,err == nil && isLogin.(bool))
		if err == nil && isLogin.(bool) {
			// 是登录状态
			c.Next()
		}else {
			c.Abort()
			c.Redirect(http.StatusFound, "/login")
		}
	}
}
