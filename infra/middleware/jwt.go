package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"one/infra/base_c"
	"one/infra/helper"
	"one/models"
	"strings"
)

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		//redis 黑名单
		str, err := base_c.Redis().Get(tokenString).Result()
		if err != nil && err.Error() != "redis: nil" || "" != str {
			// 验证不通过，不再调用后续的函数处理
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"code":    401,
				"message": "Token timeout",
			})
			return
		}
		func() {
			claim, err := helper.ParseToken(tokenString)
			//logrus.Info(claim)
			if claim == nil {
				// 验证不通过，不再调用后续的函数处理
				c.Abort()
				msg := fmt.Sprint("Couldn't handle this token:", err)
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  "error",
					"code":    401,
					"message": msg,
				})
				return
			}
			c.Set(models.JwtClaimsKey, claim)
			// 验证通过，会继续访问下一个中间件
			c.Next()
		}()
	}
}
