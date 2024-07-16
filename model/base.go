package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

const (
	BatchSize = 100
)

const (
	GmStatusZero = iota
	GmStatusNormal
	GmStatusDisable
)

type Base struct {
	Id        int       `gorm:"primaryKey;type:int(11) unsigned AUTO_INCREMENT;comment:id"`
	CreatedAt time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;comment:更新时间"`
}

type Creator struct {
	CreatorId int    `gorm:"type:int(11) unsigned;comment:操作人id"`
	Creator   string `gorm:"type:varchar(32);comment:操作人昵称"`
}

const TableOptions string = "gorm:table_options"

func GetOptions(tableName string) string {
	return "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci comment '" + tableName + "'"
}

func AutoIncrementOptions(tableName string) string {
	return "ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci comment '" + tableName + "'"
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 9999:
			pageSize = 9999
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

type GormList []string

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), g)
}

// 根据开始时间和结束时间获取所有月份
func GetMonthsBetween(table string, start, end time.Time) []string {
	var months []string
	current := start

	// 确保开始时间是月初
	start = time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, start.Location())

	// 循环直到超过结束时间
	for !current.After(end) {
		// 构建表名，例如 "logs_2023_03"
		tableName := fmt.Sprintf("%s_%d%02d", table, current.Year(), current.Month())
		months = append(months, tableName)

		// 移至下一个月的第一天
		current = current.AddDate(0, 1, 0)
	}

	return months
}
