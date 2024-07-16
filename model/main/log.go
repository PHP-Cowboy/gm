package main

import (
	"gm/model"
	log2 "gm/model/game"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

func Log() {
	dsn := "root:430022@tcp(192.168.0.254)/logdb?charset=utf8mb4&parseTime=True&loc=Local"

	logger := logger2.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger2.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			Colorful:      true,        //禁用彩色打印
			LogLevel:      logger2.Info,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "t_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: logger,
	})

	if err != nil {
		panic(err)
	}

	_ = db.Set(model.TableOptions, model.GetOptions("用户游戏战况")).AutoMigrate(&log2.GameSituation{})
}
