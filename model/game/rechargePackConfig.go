package game

import (
	"gorm.io/gorm"
	"time"
)

// 充值礼包配置
type RechargePackConfig struct {
	ID           int       `gorm:"primarykey;column:id;type:int(11);" json:"id"`
	GameId       int       `gorm:"column:game_id;type:int(11);not null;default:0;comment:游戏id;" json:"game_id"`                           //游戏id
	RoomId       int       `gorm:"column:room_id;type:int(11);not null;default:0;comment:游戏id;" json:"room_id"`                           //房间id
	BasicRewards int       `gorm:"column:basic_rewards;type:int(11);not null;default:0;comment:基础奖励;" json:"basic_rewards"`               //基础奖励
	BasicType    uint8     `gorm:"column:basic_type;type:tinyint(4);not null;default:1;comment:金额类别1=可提现可下注,2=不可提现可下注" json:"basic_type"` //金额类别(1=可提现可下注,2=不可提现可下注)
	GiftRewards  int       `gorm:"column:gift_rewards;type:int(11);not null;default:0;comment:赠送额度" json:"gift_rewards"`                  // 赠送额度
	GiftType     uint8     `gorm:"column:gift_type;type:tinyint(4);not null;default:1;comment:金额类别1=可提现可下注,2=不可提现可下注" json:"gift_type"`   //金额类别(1=可提现可下注,2=不可提现可下注)
	Bonus        int       `gorm:"column:bonus;type:int(11);not null;default:0;comment:bonus;" json:"bonus"`                              //bonus
	BonusType    uint8     `gorm:"column:bonus_type;type:tinyint(4);not null;default:1;comment:金额类别1=可提现可下注,2=不可提现可下注" json:"bonus_type"` //金额类别(1=可提现可下注,2=不可提现可下注)
	Total        int       `gorm:"column:total;type:int(11);not null;default:0;comment:总额度;" json:"total"`                                //总额度
	Price        int       `gorm:"column:price;type:int(11);not null;default:0;comment:价格;" json:"price"`                                 //价格
	Times        int       `gorm:"column:times;type:int(11);not null;default:0;comment:出现次数;" json:"times"`                               //出现次数
	Interval     int       `gorm:"column:interval;type:int(11);not null;default:0;comment:间隔时长(s);" json:"interval"`                      //间隔时长(s)
	Ratio        int       `gorm:"column:ratio;type:int(11);not null;default:0;comment:折扣比例(百分比,填整数值);" json:"ratio"`                     //折扣比例
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;" json:"updated_at"`
}

func (t *RechargePackConfig) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *RechargePackConfig) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *RechargePackConfig) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *RechargePackConfig) GetFirstById(db *gorm.DB, id uint64) (RechargePackConfig RechargePackConfig, err error) {
	err = db.Model(t).Where("id = ?", id).First(&RechargePackConfig).Error
	return
}

func (t *RechargePackConfig) GetList(db *gorm.DB) (list []RechargePackConfig, err error) {
	err = db.Model(t).Find(&list).Error
	return
}
