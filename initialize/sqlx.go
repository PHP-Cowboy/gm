package initialize

import (
	"git.dev666.cc/external/breezedup/goserver/core/logger"
	"gm/global"
	"strconv"
	"za.game/lib/rds"
)

func InitSqlx() {
	var err error

	//初始化游戏数据库

	info := global.ServerConfig.GameInfo

	global.GameDB, err = rds.InitSqlDB(info.User, info.Password, info.Host, strconv.Itoa(info.Port), info.Name, 30, 300, 59)
	if err != nil {
		logger.Logger.Error("InitApp: init GameDB failed! err:[%v]", err)
	} else {
		logger.Logger.Info("InitApp: init GameDB Success!")
	}

	//初始用户化数据库
	userInfo := global.ServerConfig.UserInfo
	global.NDB, err = rds.InitSqlDB(userInfo.User, userInfo.Password, userInfo.Host, strconv.Itoa(userInfo.Port), userInfo.Name, 30, 300, 59)
	if err != nil {
		logger.Logger.Error("InitApp: init NDB failed! err:[%v]", err)
	} else {
		logger.Logger.Info("InitApp: init NDB Success!")
	}

	//初始化支付数据库
	payInfo := global.ServerConfig.PayInfo
	global.PayDB, err = rds.InitSqlDB(payInfo.User, payInfo.Password, payInfo.Host, strconv.Itoa(payInfo.Port), payInfo.Name, 30, 300, 59)
	if err != nil {
		logger.Logger.Error("InitApp: init PayDB failed! err:[%v]", err)
	} else {
		logger.Logger.Info("InitApp: init PayDB Success!")
	}

	//初始化日志数据库
	logInfo := global.ServerConfig.LogInfo
	global.LogDB, err = rds.InitSqlDB(logInfo.User, logInfo.Password, logInfo.Host, strconv.Itoa(logInfo.Port), logInfo.Name, 30, 300, 59)
	if err != nil {
		logger.Logger.Error("InitApp: init PayDB failed! err:[%v]", err)
	} else {
		logger.Logger.Info("InitApp: init PayDB Success!")
	}

	//初始化redis
	redisInfo := global.ServerConfig.RedisInfo
	global.RedisPool, err = rds.InitRedis(redisInfo.Username, redisInfo.Password, redisInfo.Host, strconv.Itoa(redisInfo.Port), redisInfo.Db, 4, 50, 2000)
	if err != nil {
		logger.Logger.Error("InitApp: init RedisPool failed! err:[%v]", err)
	} else {
		logger.Logger.Info("InitApp: init RedisPool Success!")
	}
}
