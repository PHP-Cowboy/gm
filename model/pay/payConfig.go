package pay

import (
	"gm/model"
	"gm/request"
	"gorm.io/gorm"
	"time"
)

type PayConfig struct {
	ID             int       `gorm:"primaryKey;column:id;type:int(11);"`
	Name           string    `gorm:"column:name;type:varchar(32);not null;default:'';"`              // 支付渠道名称
	Icon           string    `gorm:"column:icon;type:varchar(128);not null;default:'';"`             // 支付图标
	Url            string    `gorm:"column:url;type:varchar(128);not null;default:'';"`              // 支付host
	BackUrl        string    `gorm:"column:back_url;type:varchar(128);not null;default:'';"`         // 代收回调地址
	PaymentBackUrl string    `gorm:"column:payment_back_url;type:varchar(128);not null;default:'';"` // 代付回调地址
	AppId          string    `gorm:"column:app_id;type:varchar(64);not null;default:'';"`            //appID
	Secret         string    `gorm:"column:secret;type:varchar(128);not null;default:'';"`           //secret 秘钥
	Merchant       string    `gorm:"column:merchant;type:varchar(64);not null;default:'',"`          //商户号
	Status         int8      `gorm:"column:status;type:tinyint(4);not null;default:0;"`              //状态 (1=可用，0=不可用)
	Remark         string    `gorm:"column:remark;type:varchar(64);default:'';"`                     //备注
	Markers        string    `gorm:"column:markers;type:varchar(64);default:'';"`                    //调用接口标记
	CreatedAt      time.Time `gorm:"column:created_at;"`
	UpdatedAt      time.Time `gorm:"column:updated_at;"`
}

func (t *PayConfig) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *PayConfig) UpdateById(db *gorm.DB, id int, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *PayConfig) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *PayConfig) GetFirstById(db *gorm.DB, id int) (PayConfig PayConfig, err error) {
	err = db.Model(t).Where("id = ?", id).First(&PayConfig).Error
	return
}

func (t *PayConfig) GetList(db *gorm.DB) (list []PayConfig, err error) {
	err = db.Model(t).Find(&list).Error
	return
}

func (t *PayConfig) GetListByIds(db *gorm.DB, ids []int) (list []PayConfig, err error) {
	err = db.Model(t).Where("id in (?)", ids).Find(&list).Error
	return
}

func (t *PayConfig) GetPageList(db *gorm.DB, req request.ConfigList) (total int64, list []PayConfig, err error) {
	localDb := db.Model(t).Where(PayConfig{Name: req.Name})

	err = localDb.Count(&total).Error

	if err != nil {
		return
	}
	err = localDb.Scopes(model.Paginate(req.Page, req.Size)).Find(&list).Error
	return
}
