package base

import (
	"github.com/gin-gonic/gin"
	"one/infra"
	"one/infra/middleware"
	"one/routes"
)

var ginApplication *gin.Engine

func Gin() *gin.Engine {
	return ginApplication
}

type GinStarter struct {
	infra.BaseStarter
}

func (i *GinStarter) Init(ctx infra.StarterContext) {
	// 创建gin application实例
	ginApplication = initGin(ctx)
}

func (i *GinStarter) Start(ctx infra.StarterContext) {
	// 读取路由信息
	g := routes.Routes(ginApplication)
	// 启动 gin
	port := ctx.Props().GetDefault("app.server.port", "80")
	// 监听并在 0.0.0.0 上启动服务
	g.Run(":" + port)
}

func (i *GinStarter) StartBlocking() bool {
	return true
}

func initGin(ctx infra.StarterContext) *gin.Engine {
	// 设置gin模式: ""/debug, release, test
	gin.SetMode("")
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(middleware.LoggerToFile(ctx))
	app.Use(middleware.Cors())
	return app
}
