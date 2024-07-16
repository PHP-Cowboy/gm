package pay

import (
	"gorm.io/gorm"
	"time"
)

// 用户体现记录表
type GaveMoneyUser struct {
	Uid      int       `gorm:"type:bigint(20);not null;default:0;comment:用户id"`                      //用户id
	GaveId   int       `gorm:"type:int(11);not null;default:0;comment:提现档位id，对应表 gave_money_config"` //提现档位id，对应表 gave_money_config
	Num      int       `gorm:"type:int(11);not null;default:0;comment:提现次数"`                         //提现次数
	LastTime time.Time `gorm:"column:lasttime;type:datetime;default:null;comment:上次提现的时间"`
}

func (t *GaveMoneyUser) UpdateNumByUidGiveId(db *gorm.DB, uid, giveId int, num int) (err error) {
	err = db.Model(t).Where("uid = ? and gave_id = ?", uid, giveId).Update("num", gorm.Expr("num + ?", num)).Error
	return
}
