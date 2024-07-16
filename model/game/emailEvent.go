package game

import (
	"gm/model"
	"gm/request"
	"gorm.io/gorm"
)

// 邮件关联事件
type EmailEvent struct {
	Id   uint64 `gorm:"primaryKey;type:int(11) unsigned AUTO_INCREMENT;comment:id"`
	Name string `gorm:"type:varchar(32);not null;default:'';comment:事件名称"`
}

func (t *EmailEvent) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *EmailEvent) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *EmailEvent) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete(t, id).Error
	return
}

func (t *EmailEvent) GetFirstById(db *gorm.DB, id int64) (EmailEvent EmailEvent, err error) {
	err = db.Model(t).Where("id = ?", id).First(&EmailEvent).Error
	return
}

func (t *EmailEvent) GetList(db *gorm.DB) (list []EmailEvent, err error) {
	err = db.Model(t).Find(&list).Error

	return
}

func (t *EmailEvent) GetPageList(db *gorm.DB, req request.GetEmailEventList) (total int64, list []EmailEvent, err error) {
	localDb := db.Model(t)

	err = localDb.Count(&total).Error

	if err != nil {
		return
	}
	err = localDb.Scopes(model.Paginate(req.Page, req.Size)).Find(&list).Error

	return
}
