package game

import (
	"gorm.io/gorm"
	"time"
)

type EmailAnnex struct {
	ID         uint64    `gorm:"primaryKey;column:id;type:int(11);"`
	Name       string    `gorm:"column:name;type:varchar(32);not null;default:'';"`      //附件名称
	EnName     string    `gorm:"column:enName;type:varchar(32);not null;default:'';"`    //附件英文名
	Type       uint8     `gorm:"column:type;type:tinyint(4);not null;default:1;"`        //类型(1=筹码)
	Amount     int       `gorm:"column:amount;type:int(11);not null;default:0;"`         //金额
	AmountType uint8     `gorm:"column:amount_type;type:tinyint(4);not null;default:1;"` //筹码类型(1=不可提现可下注，2=可提现下注)
	Unit       string    `gorm:"column:unit;type:varchar(16);not null;default:'';"`      //单位
	Remark     string    `gorm:"column:remark;type:varchar(64);default:'';"`             //备注
	CreatedAt  time.Time `gorm:"column:created_at;"`
	UpdatedAt  time.Time `gorm:"column:updated_at;"`
}

func (t *EmailAnnex) TableName() (tableName string) {
	tableName = "email_annex"
	return
}

func (t *EmailAnnex) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *EmailAnnex) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *EmailAnnex) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete(t, id).Error
	return
}

func (t *EmailAnnex) GetFirstById(db *gorm.DB, id int64) (EmailAnnex EmailAnnex, err error) {
	err = db.Model(t).Where("id = ?", id).First(&EmailAnnex).Error
	return
}

func (t *EmailAnnex) GetList(db *gorm.DB) (list []EmailAnnex, err error) {

	err = db.Model(t).Find(&list).Error

	return
}
