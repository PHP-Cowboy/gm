package game

import (
	"gm/model"
	"time"
)

// Dove Cash海外试玩平台
type DoveCash struct {
	model.Base
	Campaign string    `gorm:"type:varchar(32);not null;comment:广告期数"`
	Start    time.Time `gorm:"type:datetime;not null;comment:开始时间"`
	End      time.Time `gorm:"type:datetime;not null;comment:结束时间"`
}
