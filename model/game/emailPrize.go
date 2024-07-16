package game

import (
	"gorm.io/gorm"
	"time"
)

type EmailPrize struct {
	ID        uint64    `gorm:"primaryKey;column:id;type:int(11);"`
	Uid       uint64    `gorm:"column:uid;type:bigint(11);not null;default:0;"`   //用户ID
	EmailId   uint32    `gorm:"column:email_id;type:int(11);not null;default:0;"` //邮件ID
	Money     int       `gorm:"column:money;type:int(11);not null;default:0;"`    //无限制
	Cash      int       `gorm:"column:cash;type:int(11);not null;default:0;"`     //cash
	Store     int       `gorm:"column:store;type:int(11);not null;default:0;"`    //储钱罐
	Status    uint8     `gorm:"column:status;type:tinyint(4);not null;default:0;"`
	CreatedAt time.Time `gorm:"column:created_at;"`
}

func (t *EmailPrize) TableName() (tableName string) {
	tableY := time.Now().Local().Format("200601")
	tableName = "email_prize_" + tableY
	return
}

func (t *EmailPrize) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *EmailPrize) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *EmailPrize) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete(t, id).Error
	return
}

func (t *EmailPrize) GetFirstById(db *gorm.DB, id int64) (EmailPrize EmailPrize, err error) {
	err = db.Model(t).Where("id = ?", id).First(&EmailPrize).Error
	return
}

func (t *EmailPrize) GetList(db *gorm.DB) (list []EmailPrize, err error) {
	err = db.Model(t).Find(&list).Error

	return
}
