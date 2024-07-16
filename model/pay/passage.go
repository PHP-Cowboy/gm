package pay

import (
	"gm/global"
	"gm/model"
	"gm/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 支付通道
type Passage struct {
	model.Base
	PassageId            int            `gorm:"type:int(11);not null;default:1;comment:通道id"`
	Entrance             int            `gorm:"type:tinyint(4);not null;default:1;comment:入口"`
	CollectionRate       float64        `gorm:"type:decimal(10,2);not null;default:0.00;comment:代收费率"`
	CollectionStatus     int            `gorm:"type:tinyint(4);not null;default:1;comment:代收状态:1:开启,2:关闭"`
	CollectionSort       int            `gorm:"type:int(11);not null;default:1;comment:排序"`
	CollectionChannelIds model.GormList `gorm:"type:varchar(255);comment:代收渠道id"`
	PaymentRate          float64        `gorm:"type:decimal(10,2);not null;default:0.00;comment:代付费率"`
	PaymentStatus        int            `gorm:"type:tinyint(4);not null;default:1;comment:代付状态:1:开启,2:关闭"`
	PaymentSort          int            `gorm:"type:int(11);not null;default:1;comment:排序"`
	PaymentChannelIds    model.GormList `gorm:"type:varchar(255);comment:代付渠道id"`
}

func (t *Passage) Create(db *gorm.DB, data Passage) (err error) {
	err = db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"passage_id",
				"entrance",
				"collection_rate",
				"collection_status",
				"collection_sort",
				"collection_channel_ids",
				"payment_rate",
				"payment_status",
				"payment_sort",
				"payment_channel_ids",
			}),
		}).
		Create(&data).
		Error

	return
}

func (t *Passage) UpdateById(db *gorm.DB, id int, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *Passage) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Delete(t, id).Error
	return
}

func (t *Passage) GetListPaymentSort(db *gorm.DB) (list []Passage, err error) {
	err = db.Model(t).Where("payment_status = 1").Order("payment_sort asc").Find(&list).Error
	return
}

func (t *Passage) GetList(db *gorm.DB) (list []Passage, err error) {
	err = db.Model(t).Where("payment_status = 1").Find(&list).Error
	return
}

// 分页列表
func (t *Passage) GetPageList(db *gorm.DB, req request.PassageList) (total int64, list []Passage, err error) {
	localDb := db.Model(t).Where(Passage{Entrance: req.Entrance, CollectionStatus: req.CollectionStatus, PaymentStatus: req.PaymentStatus})

	err = localDb.Count(&total).Error

	if err != nil {
		global.Logger["err"].Errorf("Passage GetPageList failed,err:[%v]", err.Error())
		return
	}

	err = localDb.Order("").Scopes(model.Paginate(req.Page, req.Size)).Find(&list).Error
	return
}

func (t *Passage) GetFirstById(db *gorm.DB, id int64) (data Passage, err error) {
	err = db.Model(t).Where("id = ?", id).First(&data).Error
	return
}
