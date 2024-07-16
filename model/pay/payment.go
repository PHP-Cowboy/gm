package pay

import (
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	ID           uint64    `gorm:"primaryKey;column:id;type:int(11);"`
	Uid          uint64    `gorm:"column:uid;type:bigint(20);not null;default:0;index:uid_key;"`                         // 用户ID
	Ymd          int       `gorm:"column:ymd;type:int(11);not null;default:0;"`                                          // 年月日
	PayId        int       `gorm:"column:pay_id;type:int(11);not null;default:0;"`                                       // 支付配置id
	OrderNo      string    `gorm:"column:order_no;type:varchar(32);not null;default:'';unique:order_key;index:uid_key;"` // 订单ID
	Account      int       `gorm:"column:account;type:int(11);not null;default:0;"`                                      // 支付金额
	TransferId   string    `gorm:"column:transfer_id;type:varchar(32);not null;default:'';unique:tran_key;"`             //支付渠道订单ID
	BankCode     string    `gorm:"column:bank_code;type:varchar(32);not null;default:'';"`                               //代付类型
	BankName     string    `gorm:"column:bank_name;type:varchar(32);not null;default:'';"`                               //银行名称
	AccountNo    string    `gorm:"column:account_no;type:varchar(32);not null;default:'';"`                              //银行卡号
	Ifsc         string    `gorm:"column:ifsc;type:varchar(32);not null;default:'';"`                                    // ifsc
	Name         string    `gorm:"column:name;type:varchar(32);not null;default:'';"`                                    // 客户姓名
	Email        string    `gorm:"column:email;type:varchar(32);not null;default:'';"`                                   //客户邮箱
	Phone        string    `gorm:"column:phone;type:varchar(32);not null;default:'';"`                                   //手机
	Address      string    `gorm:"column:address;type:varchar(32);not null;default:'';"`                                 //客户地址
	Vpa          string    `gorm:"column:vpa;type:varchar(32);not null;default:'';"`                                     //upi 账号
	Status       int8      `gorm:"column:status;type:tinyint(4);not null;default:0;"`                                    // 状态 0=等待支付，1=支付完成，2=下单失败
	RequestTime  time.Time `gorm:"column:request_time;"`                                                                 //下单时间
	OrderTime    time.Time `gorm:"column:order_time;"`                                                                   //下单返回时间
	CompleteTime int       `gorm:"column:complete_time;type:int(11);not null;default:0;"`                                //订单完成时间
	Remark       string    `gorm:"column:remark;type:varchar(64);not null;default:'';"`                                  //备注
	CreatedAt    time.Time `gorm:"column:created_at;"`
	UpdatedAt    time.Time `gorm:"column:updated_at;"`
}

func (t *Payment) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *Payment) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *Payment) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *Payment) GetFirstById(db *gorm.DB, id int64) (Payment Payment, err error) {
	err = db.Model(t).Where("id = ?", id).First(&Payment).Error
	return
}

func (t *Payment) GetList(db *gorm.DB, PaymentCond Payment) (list []Payment, err error) {
	err = db.Model(t).Where(&PaymentCond).Find(&list).Error
	return
}
