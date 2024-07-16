package game

import (
	"gm/model"
	"gorm.io/gorm"
	"time"
)

type SignConfig struct {
	ID           uint64         `gorm:"primarykey;column:id;type:int(11);" json:"id"`
	Name         string         `gorm:"column:name;type:varchar(32);not null;default:'';" json:"name"`            //签到名称
	SignNum      uint8          `gorm:"column:sign_num;type:tinyint(4);not null;default:0;" json:"sign_num"`      //累计签到次数
	PrizeIds     model.GormList `gorm:"column:prize_ids;type:varchar(255);not null;default:'';" json:"prize_ids"` //奖励ID
	PrizeCashId  int            `gorm:"column:prize_cash_id;type:int;not null;default:0;" json:"prize_cash_id"`   //cash
	PrizeBonusId int            `gorm:"column:prize_bonus_id;type:int;not null;default:0;" json:"prize_bonus_id"` //bonus
	Unit         string         `gorm:"column:unit;type:varchar(16);not null;default:'';" json:"unit"`            //单位
	Remark       string         `gorm:"column:remark;type:varchar(64);default:'';" json:"remark"`                 //备注
	CreatedAt    time.Time      `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;" json:"updated_at"`
}

func (t *SignConfig) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *SignConfig) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *SignConfig) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete(t, id).Error
	return
}

func (t *SignConfig) GetFirstById(db *gorm.DB, id int64) (SignConfig SignConfig, err error) {
	err = db.Model(t).Where("id = ?", id).First(&SignConfig).Error
	return
}

func (t *SignConfig) GetList(db *gorm.DB) (list []SignConfig, err error) {
	err = db.Model(t).Find(&list).Error

	return
}
