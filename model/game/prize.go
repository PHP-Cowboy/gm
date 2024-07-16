package game

import (
	"gm/model"
	"gm/request"
	"gorm.io/gorm"
	"time"
)

type Prize struct {
	ID        uint64    `gorm:"primaryKey;column:id;type:int(11);" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(32);not null;default:'';comment:奖励名称" json:"name"`                                     //奖励名称
	EnName    string    `gorm:"column:en_name;type:varchar(32);not null;default:'';comment:奖励英文名称" json:"en_name"`                             //附件英文名
	Type      uint8     `gorm:"column:type;type:tinyint(4);not null;default:1;comment:类型(1:筹码)" json:"type"`                                   //类型(1:筹码)
	GoodsNum  int       `gorm:"column:goods_num;type:int(11);not null;default:0;comment:物品数量" json:"goods_num"`                                //物品数量
	GoodsType uint8     `gorm:"column:goods_type;type:tinyint(4);not null;default:0;comment:物品类别(1:winCash,2:cash,3:bonus)" json:"goods_type"` //物品类别(1:winCash,2:cash,3:bonus)
	Unit      string    `gorm:"column:unit;type:varchar(16);not null;default:'';comment:单位" json:"unit"`                                       //单位
	Remark    string    `gorm:"column:remark;type:varchar(64);default:'';comment:备注" json:"remark"`                                            //备注
	CreatedAt time.Time `gorm:"column:created_at;"`
	UpdatedAt time.Time `gorm:"column:updated_at;"`
}

func (t *Prize) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *Prize) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *Prize) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete(t, id).Error
	return
}

func (t *Prize) GetFirstById(db *gorm.DB, id int64) (Prize Prize, err error) {
	err = db.Model(t).Where("id = ?", id).First(&Prize).Error
	return
}

func (t *Prize) GetList(db *gorm.DB) (list []Prize, err error) {
	err = db.Model(t).Find(&list).Error

	return
}

func (t *Prize) GetPageList(db *gorm.DB, req request.GetSingPrizeList) (total int64, list []Prize, err error) {

	localDb := db.Model(t)

	err = localDb.Count(&total).Error

	if err != nil {
		return
	}

	err = localDb.Scopes(model.Paginate(req.Page, req.Size)).Find(&list).Error

	return
}
