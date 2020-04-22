package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Routes(r *gin.Engine) *gin.Engine {
	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	//r.GET("/chapter", user.GetVCode)
	//// 路由组
	//auth := r.Group("")
	//auth.Use(Middleware.Jwt())
	//{
	//	auth.POST("/sign-out", user.Logout)
	//}
	return r
}
