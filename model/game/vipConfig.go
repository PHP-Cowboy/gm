package game

import (
	"gorm.io/gorm"
	"time"
)

type VipConfig struct {
	ID             int       `gorm:"primarykey;column:id;type:int(11);" json:"id"`
	Level          uint8     `gorm:"column:level;type:tinyint(4);not null;default:0;comment:等级" json:"level"`                           // 等级
	NeedExp        uint64    `gorm:"column:need_exp;type:int(11);not null;default:0;comment:升级所需经验" json:"need_exp"`                    // 升级所需经验
	WithdrawNums   uint8     `gorm:"column:withdraw_nums;type:tinyint(4);not null;default:0;comment:每日提现次数" json:"withdraw_nums"`       // 每日提现次数
	WithdrawMoney  uint64    `gorm:"column:withdraw_money;type:int(11);not null;default:0;comment:每日提现金额" json:"withdraw_money"`        // 每日提现金额
	DayPrizeType   uint8     `gorm:"column:day_prize_type;type:tinyint(4);not null;default:0;comment:每日奖品类型" json:"day_prize_type"`     // 每日奖品类型
	DayPrizeNums   uint64    `gorm:"column:day_prize_nums;type:int(11);not null;default:0;comment:每日奖品数量" json:"day_prize_nums"`        //每日奖品数量
	WeekPrizeType  uint8     `gorm:"column:week_prize_type;type:tinyint(4);not null;default:0;comment:每周奖品类型" json:"week_prize_type"`   //每周奖品类型
	WeekPrizeNums  uint64    `gorm:"column:week_prize_nums;type:int(11);not null;default:0;comment:每周奖品数量" json:"week_prize_nums"`      //每周奖品数量
	MonthPrizeType uint8     `gorm:"column:month_prize_type;type:tinyint(4);not null;default:0;comment:每月奖品类型" json:"month_prize_type"` //每月奖品类型
	MonthPrizeNums uint64    `gorm:"column:month_prize_nums;type:int(11);not null;default:0;comment:每月奖品数量" json:"month_prize_nums"`    //每月奖品数量
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime;" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;" json:"updated_at"`
}

func (t *VipConfig) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *VipConfig) UpdateById(db *gorm.DB, id int, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *VipConfig) DeleteById(db *gorm.DB, id int) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *VipConfig) GetFirstById(db *gorm.DB, id int) (vipConfig VipConfig, err error) {
	err = db.Model(t).Where("id = ?", id).First(&vipConfig).Error
	return
}

func (t *VipConfig) GetList(db *gorm.DB) (list []VipConfig, err error) {
	err = db.Model(t).Find(&list).Error
	return
}
