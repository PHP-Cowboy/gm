package game

import (
	"gorm.io/gorm"
	"time"
)

type EventGiftConfig struct {
	ID        int       `gorm:"primarykey;column:id;type:int(11);" json:"id"`
	Grade     int       `gorm:"column:grade;type:int(11);not null;default:0;comment:档次" json:"grade"`                                  // 档次
	Coin      int       `gorm:"column:coin;type:int(11);not null;default:0;comment:游戏币" json:"coin"`                                   // 游戏币
	CoinType  uint8     `gorm:"column:coin_type;type:tinyint(4);not null;default:1;comment:金额类别1=可提现可下注,2=不可提现可下注" json:"coin_type"`   //金额类别(1=可提现可下注,2=不可提现可下注)
	GiftLimit int       `gorm:"column:gift_limit;type:int(11);not null;default:0;comment:赠送额度" json:"gift_limit"`                      // 赠送额度
	GiftType  uint8     `gorm:"column:gift_type;type:tinyint(4);not null;default:1;comment:金额类别1=可提现可下注,2=不可提现可下注" json:"gift_type"`   //金额类别(1=可提现可下注,2=不可提现可下注)
	Bonus     int       `gorm:"column:bonus;type:int(11);not null;default:0;comment:bonus;" json:"bonus"`                              //bonus
	BonusType uint8     `gorm:"column:bonus_type;type:tinyint(4);not null;default:1;comment:金额类别1=可提现可下注,2=不可提现可下注" json:"bonus_type"` //金额类别(1=可提现可下注,2=不可提现可下注)
	Ratio     int       `gorm:"column:ratio;type:int(11);not null;default:0;comment:折扣比例(百分比,填整数值);" json:"ratio"`                     //折扣比例
	Type      int       `gorm:"column:type;type:tinyint(4);not null;default:1;comment:类型[根据不同类型给结果]" json:"type"`                      //类型[根据不同类型给结果]
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;" json:"updated_at"`
}

func (t *EventGiftConfig) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *EventGiftConfig) UpdateById(db *gorm.DB, id int, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *EventGiftConfig) UpdateByIds(db *gorm.DB, ids []int, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id in (?)", ids).Updates(mp).Error
	return
}

func (t *EventGiftConfig) DeleteById(db *gorm.DB, id int) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *EventGiftConfig) GetFirstById(db *gorm.DB, id int) (EventGiftConfig EventGiftConfig, err error) {
	err = db.Model(t).Where("id = ?", id).First(&EventGiftConfig).Error
	return
}

func (t *EventGiftConfig) GetFirstByIsClose(db *gorm.DB, isClose int) (cfg EventGiftConfig, err error) {
	err = db.Model(t).Where("is_close = ?", isClose).First(&cfg).Error
	return
}

func (t *EventGiftConfig) GetList(db *gorm.DB) (list []EventGiftConfig, err error) {
	err = db.Model(t).Find(&list).Error
	return
}
