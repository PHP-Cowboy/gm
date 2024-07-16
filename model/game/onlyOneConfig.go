package game

import (
	"gorm.io/gorm"
	"time"
)

// 三选一礼包配置表
type OnlyOneConfig struct {
	ID        int       `gorm:"primarykey;column:id;type:int(11);" json:"id"`
	Grade     int       `gorm:"column:grade;type:int;not null;default:0;comment:档位" json:"grade"`                                      // 档位
	FirstDay  int       `gorm:"column:first_day;type:int;not null;default:0;comment:第一天领取" json:"first_day"`                           // 第一天领取
	FirstType uint8     `gorm:"column:first_type;type:tinyint(4);not null;default:1;comment:金额类别1=可提现可下注,2=不可提现可下注" json:"first_type"` //金额类别(1=可提现可下注,2=不可提现可下注)
	NextDay   int       `gorm:"column:next_day;type:int;not null;default:0;comment:第二天领取" json:"next_day"`                             // 第二天领取
	NextType  uint8     `gorm:"column:next_type;type:tinyint(4);not null;default:1;comment:金额类别1=可提现可下注,2=不可提现可下注" json:"next_type"`   //金额类别(1=可提现可下注,2=不可提现可下注)
	ThirdDay  int       `gorm:"column:third_day;type:int;not null;default:0;comment:第三天领取" json:"third_day"`                           // 第三天领取
	ThirdType uint8     `gorm:"column:third_type;type:tinyint(4);not null;default:1;comment:金额类别1=可提现可下注,2=不可提现可下注" json:"third_type"` //金额类别(1=可提现可下注,2=不可提现可下注)
	Ratio     int       `gorm:"column:ratio;type:int(11);not null;default:0;comment:折扣比例(百分比,填整数值);" json:"ratio"`                     //折扣比例
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;" json:"updated_at"`
}

func (t *OnlyOneConfig) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *OnlyOneConfig) UpdateById(db *gorm.DB, id int, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *OnlyOneConfig) UpdateByIds(db *gorm.DB, ids []int, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id in (?)", ids).Updates(mp).Error
	return
}

func (t *OnlyOneConfig) DeleteById(db *gorm.DB, id int) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *OnlyOneConfig) GetFirstById(db *gorm.DB, id int) (OnlyOneConfig OnlyOneConfig, err error) {
	err = db.Model(t).Where("id = ?", id).First(&OnlyOneConfig).Error
	return
}

func (t *OnlyOneConfig) GetFirstByIsClose(db *gorm.DB, isClose int) (cfg OnlyOneConfig, err error) {
	err = db.Model(t).Where("is_close = ?", isClose).First(&cfg).Error
	return
}

func (t *OnlyOneConfig) GetList(db *gorm.DB) (list []OnlyOneConfig, err error) {
	err = db.Model(t).Find(&list).Error
	return
}
