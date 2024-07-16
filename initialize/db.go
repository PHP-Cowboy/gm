package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"gm/global"
)

// 定义自己的Writer
type SqlWriter struct {
	sqlLog *zap.SugaredLogger
}

// 实现gorm/logger.Writer接口
func (m *SqlWriter) Printf(format string, v ...interface{}) {
	//记录日志
	m.sqlLog.Info(fmt.Sprintf(format, v...))
}

func NewSqlWriter() *SqlWriter {
	l, ok := global.Logger["sql"]

	if !ok {
		panic("sql日志加载失败")
	}

	return &SqlWriter{sqlLog: l}
}

func InitMysql() {
	info := global.ServerConfig.MysqlInfo

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		info.User,
		info.Password,
		info.Host,
		info.Port,
		info.Name,
	)

	logger := logger2.New(
		NewSqlWriter(),
		logger2.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			Colorful:      true,        //禁用彩色打印
			LogLevel:      logger2.Info,
		},
	)

	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "t_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}

	userInfo := global.ServerConfig.UserInfo

	userDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		userInfo.User,
		userInfo.Password,
		userInfo.Host,
		userInfo.Port,
		userInfo.Name,
	)

	global.User, err = gorm.Open(mysql.Open(userDsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "t_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}

	logInfo := global.ServerConfig.LogInfo

	logDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		logInfo.User,
		logInfo.Password,
		logInfo.Host,
		logInfo.Port,
		logInfo.Name,
	)

	global.Log, err = gorm.Open(mysql.Open(logDsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "t_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}

	gameInfo := global.ServerConfig.GameInfo

	gameDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		gameInfo.User,
		gameInfo.Password,
		gameInfo.Host,
		gameInfo.Port,
		gameInfo.Name,
	)

	global.Game, err = gorm.Open(mysql.Open(gameDsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "t_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}

	payInfo := global.ServerConfig.PayInfo

	payDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		payInfo.User,
		payInfo.Password,
		payInfo.Host,
		payInfo.Port,
		payInfo.Name,
	)

	global.Pay, err = gorm.Open(mysql.Open(payDsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "t_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}
}
