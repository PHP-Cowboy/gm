package main

import (
	"fmt"
	"gm/corn"
	"gm/global"
	"gm/initialize"
	"gm/msgcenter"
	"os"
	"os/signal"
	"syscall"
	"za.game/lib/dbconn"
)

func main() {
	initialize.InitLogger()

	initialize.InitConfig()

	initialize.InitMysql()

	initialize.InitSqlx()

	initialize.InitRedis()

	//za.game/lib/dbconn
	dbconn.InitApp()

	//防止没有dbconn配置信息，用来过渡用
	if dbconn.NDB == nil {
		dbconn.NDB = global.NDB
	}
	if dbconn.GameDB == nil {
		dbconn.GameDB = global.GameDB
	}
	if dbconn.LogDB == nil {
		dbconn.LogDB = global.LogDB
	}
	if dbconn.PayDB == nil {
		dbconn.PayDB = global.PayDB
	}
	if dbconn.RedisPool == nil {
		dbconn.RedisPool = global.RedisPool
	}

	if global.ServerConfig.Mode != "local" {
		go corn.Consumer()
	}

	msgcenter.InitMct(&global.ServerConfig.Mct)

	g := initialize.InitRouter()

	serverConfig := global.ServerConfig

	fmt.Println("服务启动中,端口:", serverConfig.Port)

	go func() {
		err := g.Run(fmt.Sprintf(":%d", serverConfig.Port))
		if err != nil {
			panic("启动失败:" + err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	msgcenter.Shutdown()
}
