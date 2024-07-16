package user

import (
	"gm/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 封禁信息
type Banned struct {
	model.Base
	Cate int    `gorm:"type:tinyint(4);not null;uniqueIndex:cateInfo;comment:1:IP,2:device"`
	Info string `gorm:"type:varchar(64);not null;uniqueIndex:cateInfo;comment:类别信息"`
	Uid  int    `gorm:"type:int(11);not null;comment:管理后台操作人id"`
}

type BannedCateType int

const (
	BannedCate       BannedCateType = iota
	BannedCateIP     BannedCateType = iota
	BannedCateDevice BannedCateType = iota
)

func (t *Banned) TableName() string {
	return "banned"
}

func (t *Banned) Upsert(db *gorm.DB, cate BannedCateType, info string, uid int) (err error) {
	err = db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "cate,info"}},
			DoUpdates: clause.AssignmentColumns([]string{"cate", "info", "uid"}),
		}).
		Create(&Banned{
			Cate: int(cate),
			Info: info,
			Uid:  uid,
		}).
		Error

	return
}

func (t *Banned) DeleteByUserCateInfo(db *gorm.DB, cate BannedCateType, info string) (err error) {
	err = db.Model(t).Where("cate = ? and info = ?", cate, info).Delete(&Banned{}).Error
	return
}
