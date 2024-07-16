package pay

import (
	"gorm.io/gorm"
	"time"
)

type PayGift struct {
	ID           uint64    `gorm:"primaryKey;column:id;type:int(11);"`
	Name         string    `gorm:"column:name;type:varchar(32);not null;default:'';"`         // 支付包名称
	Cash         int       `gorm:"column:cash;type:int(11);not null;default:0;"`              // 到账金额
	Account      int       `gorm:"column:account;type:int(11);not null;default:0;"`           // 支付金额
	Status       int8      `gorm:"column:status;type:tinyint(4);not null;default:1;"`         // 状态（1=可用，0=不可用）
	AddMoney     int       `gorm:"column:add_money;type:int(11);not null;default:0;"`         //额外赠送无限制
	AddMoneyType int8      `gorm:"column:add_money_type;type:tinyint(4);not null;default:1;"` //无限制赠送类别（1=金额，2=比例）
	AddCash      int       `gorm:"column:add_cash;type:int(11);not null;default:0;"`          //赠送cash
	AddCashType  int8      `gorm:"column:add_cash_type;type:tinyint(4);not null;default:1;"`  //赠送cash金额（1=金额，2=比例）
	Bonus        int       `gorm:"column:bonus;type:int(11);default:0;"`                      //储钱罐
	BonusType    int8      `gorm:"column:bonus_type;type:tinyint(4);default:1;"`              //储钱罐类型（1=金额，2=比例）
	Ratio        int       `gorm:"column:ratio;type:int(11);not null;default:0;"`             //优惠比例(百分比,填整数值)
	Remark       string    `gorm:"column:remark;type:varchar(64);default:'';"`                //备注
	Type         int       `gorm:"column:_type;type:tinyint(4);not null;default:1;comment:赠送类型(1:正常，2:首提)"`
	ReplaceId    int       `gorm:"type:int;not null;default:0;comment:替换的目标id"`
	CreatedAt    time.Time `gorm:"column:created_at;"`
	UpdatedAt    time.Time `gorm:"column:updated_at;"`
}

func (t *PayGift) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *PayGift) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *PayGift) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *PayGift) GetFirstById(db *gorm.DB, id uint64) (PayGift PayGift, err error) {
	err = db.Model(t).Where("id = ?", id).First(&PayGift).Error
	return
}

func (t *PayGift) GetOneByReplaceId(db *gorm.DB, replace_id uint64) (data PayGift, err error) {
	err = db.Model(t).Where("replace_id = ?", replace_id).First(&data).Error
	return
}

func (t *PayGift) GetList(db *gorm.DB, PayGiftCond *PayGift) (list []PayGift, err error) {
	err = db.Model(t).Where(PayGiftCond).Find(&list).Error
	return
}
