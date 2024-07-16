package initialize

import (
	"github.com/gin-gonic/gin"
	"gm/api"
	"gm/middlewares"
	"gm/router"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	//跨域
	r.Use(middlewares.Cors())

	group := r.Group("/v1/")
	group.POST("login", api.Login)
	//无需鉴权
	router.Router(group)
	//需鉴权路由
	router.AuthRouter(group)

	return r
}
