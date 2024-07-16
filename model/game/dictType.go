package game

import (
	"gorm.io/gorm"
	"time"
)

type DictType struct {
	Code       string    `gorm:"type:varchar(50);primaryKey;comment:字典类型编码"` //这里改成 TypeCode
	Name       string    `gorm:"type:varchar(20);not null;comment:字典类型名称"`
	CreateTime time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime time.Time `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
}

func (t *DictType) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *DictType) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *DictType) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *DictType) GetFirstById(db *gorm.DB, id int64) (dictType DictType, err error) {
	err = db.Model(t).Where("id = ?", id).First(&dictType).Error
	return
}

func (t *DictType) GetList(db *gorm.DB) (list []DictType, err error) {
	err = db.Model(t).Find(&list).Error
	return
}
