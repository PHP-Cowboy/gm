package game

import (
	"gorm.io/gorm"
	"time"
)

type Dict struct {
	Id         uint64    `json:"id"`
	TypeCode   string    `gorm:"type:varchar(20);not null;primaryKey;comment:字典类型编码"`
	Code       string    `gorm:"type:varchar(50);not null;primaryKey;comment:字典编码"` //这里改成 DictCode
	Name       string    `gorm:"type:varchar(20);not null;comment:字典名称"`
	Value      string    `gorm:"type:varchar(20);not null;comment:字典值"`
	IsEdit     int       `gorm:"type:tinyint;not null;default:0;comment:是否可编辑:0:否,1:是"`
	CreateTime time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime time.Time `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
}

func (t *Dict) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *Dict) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *Dict) UpdateValueByTCode(db *gorm.DB, typeCode, code, value string) (err error) {
	err = db.Model(t).Where("type_code = ? and code = ?", typeCode, code).Update("value", value).Error
	return
}

func (t *Dict) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *Dict) GetFirstById(db *gorm.DB, id int64) (Dict Dict, err error) {
	err = db.Model(t).Where("id = ?", id).First(&Dict).Error
	return
}

func (t *Dict) GetList(db *gorm.DB, dictCond Dict) (list []Dict, err error) {
	err = db.Model(t).Where(&dictCond).Find(&list).Error
	return
}

func (t *Dict) GetDictMpByTypeCode(db *gorm.DB, typeCode string) (mp map[string]string, err error) {
	dictList := make([]Dict, 0)

	err = db.Model(t).Where("type_code = ? ", typeCode).Find(&dictList).Error

	if err != nil {
		return
	}

	mp = make(map[string]string)

	for _, d := range dictList {
		mp[d.Code] = d.Value
	}
	return
}

func (t *Dict) GetOneByUnique(db *gorm.DB, typeCode, code string) (Dict Dict, err error) {
	err = db.Model(t).Where("type_code = ? and code = ?", typeCode, code).First(&Dict).Error
	return
}

//func (t *Dict) UpsertByUniqueTCode(db *gorm.DB, list []Dict) (err error) {
//	err = db.Model(t).Clauses(clause.OnConflict{
//		Columns:   []clause.Column{{Name: "type_code,code"}},
//		DoUpdates: clause.AssignmentColumns([]string{"value"}),
//	}).Create(&list).Error
//
//	return
//}
