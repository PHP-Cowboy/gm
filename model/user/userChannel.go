package user

import (
	"gm/model"
	"gm/request"
	"gorm.io/gorm"
	"time"
)

type UserChannel struct {
	ID          int        `gorm:"primarykey;column:id;type:int(11);"`
	ChannelName string     `gorm:"column:channel_name;type:varchar(32);not null;default:'';"` //渠道名称
	Code        string     `gorm:"column:code;type:varchar(255);not null;default:'';"`        //code
	Remark      string     `gorm:"column:remark;type:varchar(64);not null;default:'';"`       //备注
	CreatedAt   *time.Time `gorm:"column:created_at;"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index;"`
}

func (table *UserChannel) TableName() string {
	return "user_channel"
}

func (t *UserChannel) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *UserChannel) UpdateById(db *gorm.DB, id int, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *UserChannel) DeleteById(db *gorm.DB, id int) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *UserChannel) Count(db *gorm.DB) (total int64, err error) {
	err = db.Model(t).Where(t).Count(&total).Error
	return
}

func (t *UserChannel) GetFirstById(db *gorm.DB, id int) (UserChannel UserChannel, err error) {
	err = db.Model(t).Where("id = ?", id).First(&UserChannel).Error
	return
}

func (t *UserChannel) GetList(db *gorm.DB) (list []UserChannel, err error) {
	err = db.Model(t).Where(t).Find(&list).Error
	return
}

func (t *UserChannel) GetPageList(db *gorm.DB, req request.ChannelList) (list []UserChannel, err error) {
	err = db.Model(t).Where(t).Scopes(model.Paginate(req.Page, req.Size)).Find(&list).Error
	return
}
