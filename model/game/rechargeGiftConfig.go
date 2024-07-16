package game

import (
	"gorm.io/gorm"
	"time"
)

// 充值赠送礼包配置
type RechargeGiftConfig struct {
	ID           int       `gorm:"primarykey;column:id;type:int(11);" json:"id"`
	BasicRewards int       `gorm:"column:basic_rewards;type:int(11);not null;default:0;comment:基础奖励;" json:"basic_rewards"`               //基础奖励
	BasicType    uint8     `gorm:"column:basic_type;type:tinyint(4);not null;default:1;comment:金额类别1=可提现可下注,2=不可提现可下注" json:"basic_type"` //金额类别(1=可提现可下注,2=不可提现可下注)
	GiftRewards  int       `gorm:"column:gift_rewards;type:int(11);not null;default:0;comment:赠送额度" json:"gift_rewards"`                  // 赠送额度
	GiftType     uint8     `gorm:"column:gift_type;type:tinyint(4);not null;default:1;comment:金额类别1=可提现可下注,2=不可提现可下注" json:"gift_type"`   //金额类别(1=可提现可下注,2=不可提现可下注)
	Gift2Rewards int       `gorm:"column:gift_rewards;type:int(11);not null;default:0;comment:赠送额度" json:"gift2_rewards"`                 // 赠送额度
	Gift2Type    uint8     `gorm:"column:gift_type;type:tinyint(4);not null;default:1;comment:金额类别1=可提现可下注,2=不可提现可下注" json:"gift2_type"`  //金额类别(1=可提现可下注,2=不可提现可下注)
	Total        int       `gorm:"column:total;type:int(11);not null;default:0;comment:总额度;" json:"total"`                                //总额度
	Price        int       `gorm:"column:price;type:int(11);not null;default:0;comment:价格;" json:"price"`                                 //价格
	Times        int       `gorm:"column:times;type:int(11);not null;default:0;comment:出现次数;" json:"times"`                               //出现次数
	Interval     int       `gorm:"column:interval;type:int(11);not null;default:0;comment:间隔时长(s);" json:"interval"`                      //间隔时长(s)
	Ratio        int       `gorm:"column:ratio;type:int(11);not null;default:0;comment:折扣比例(百分比,填整数值);" json:"ratio"`                     //折扣比例
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;" json:"updated_at"`
}

func (t *RechargeGiftConfig) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *RechargeGiftConfig) UpdateById(db *gorm.DB, id int, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *RechargeGiftConfig) UpdateByIds(db *gorm.DB, ids []int, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id in (?)", ids).Updates(mp).Error
	return
}

func (t *RechargeGiftConfig) DeleteById(db *gorm.DB, id int) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *RechargeGiftConfig) GetFirstById(db *gorm.DB, id int) (RechargeGiftConfig RechargeGiftConfig, err error) {
	err = db.Model(t).Where("id = ?", id).First(&RechargeGiftConfig).Error
	return
}

func (t *RechargeGiftConfig) GetFirstByIsClose(db *gorm.DB, isClose int) (cfg RechargeGiftConfig, err error) {
	err = db.Model(t).Where("is_close = ?", isClose).First(&cfg).Error
	return
}

func (t *RechargeGiftConfig) GetList(db *gorm.DB) (list []RechargeGiftConfig, err error) {
	err = db.Model(t).Find(&list).Error
	return
}
