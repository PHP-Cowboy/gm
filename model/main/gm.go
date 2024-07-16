package main

import (
	"gm/model"
	"gm/model/gm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

func Gm() {
	dsn := "root:430022@tcp(192.168.0.254)/gm?charset=utf8mb4&parseTime=True&loc=Local"

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

	//_ = db.Set(model.TableOptions, model.GetOptions("菜单表")).AutoMigrate(&gm.Menu{})
	//_ = db.Set(model.TableOptions, model.GetOptions("角色表")).AutoMigrate(&gm.Role{})
	//_ = db.Set(model.TableOptions, model.GetOptions("角色菜单表")).AutoMigrate(&gm.RoleMenu{})
	//_ = db.Set(model.TableOptions, model.GetOptions("用户表")).AutoMigrate(&gm.User{})
	//_ = db.Set(model.TableOptions, model.GetOptions("用户角色表")).AutoMigrate(&gm.UserRole{})
	//_ = db.Set(model.TableOptions, model.GetOptions("用户数据")).AutoMigrate(&gm.UserData{})
	//_ = db.Set(model.TableOptions, model.GetOptions("封禁")).AutoMigrate(&gm.Banned{})
	//_ = db.Set(model.TableOptions, model.GetOptions("首页五分钟数据")).AutoMigrate(&gm.FiveMinuteData{})
	//_ = db.Set(model.TableOptions, model.GetOptions("报表数据")).AutoMigrate(&gm.ReportData{})
	_ = db.Set(model.TableOptions, model.GetOptions("角色渠道")).AutoMigrate(&gm.RoleChannel{})

}
