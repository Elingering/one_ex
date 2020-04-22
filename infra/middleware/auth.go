package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"one/models"
	"regexp"
)

//权限检查中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//根据上下文获取载荷claims 从claims获得role
		claims := c.MustGet(models.JwtClaimsKey).(*models.Claim)
		role := claims.Role
		e := models.Casbin()
		pat := "/v[0-9]+" //正则
		re,_ := regexp.Compile(pat)
		path := re.ReplaceAllString(c.Request.URL.Path,"")
		//log.Info(role+" | "+path+" | "+c.Request.Method)
		//检查权限
		ok, err := e.EnforceSafe(role, path, c.Request.Method)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": false,
				"msg":    "错误消息" + err.Error(),
			})
			c.Abort()
			return
		}
		if ok {
			c.Next()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": false,
				"msg":    "很抱歉您没有此权限",
			})
			c.Abort()
			return
		}
	}
}
