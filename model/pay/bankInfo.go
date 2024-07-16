package pay

import (
	"gm/model"
	"gorm.io/gorm"
	"time"
)

type BankInfo struct {
	ID        uint64    `gorm:"primaryKey;column:id;type:int(11);"`
	Uid       uint64    `gorm:"column:uid;type:bigint(20);not null;default:0;"`          // 用户ID
	BankCode  string    `gorm:"column:bank_code;type:varchar(32);not null;default:'';"`  // 银行类型
	BankName  string    `gorm:"column:bank_name;type:varchar(32);not null;default:'';"`  // 银行名称
	AccountNo string    `gorm:"column:account_no;type:varchar(32);not null;default:'';"` // 银行账号
	Ifsc      string    `gorm:"column:ifsc;type:varchar(32);not null;default:'';"`       // ifsc号
	Name      string    `gorm:"column:name;type:varchar(32);not null;default:'';"`       //客户姓名
	Email     string    `gorm:"column:email;type:varchar(32);not null;default:'';"`      //客户邮箱
	Phone     string    `gorm:"column:phone;type:varchar(32);not null;default:'',"`      //客户手机
	Address   string    `gorm:"column:address;type:varchar(64);not null;default:'';"`    //客户地址
	Vpa       string    `gorm:"column:vpa;type:varchar(32);not null;default:'';"`        //vpa
	Remark    string    `gorm:"column:remark;type:varchar(64);default:'';"`              //备注
	CreatedAt time.Time `gorm:"column:created_at;"`
	UpdatedAt time.Time `gorm:"column:updated_at;"`
}

func (t *BankInfo) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *BankInfo) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *BankInfo) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *BankInfo) GetFirstById(db *gorm.DB, id int64) (BankInfo BankInfo, err error) {
	err = db.Model(t).Where("id = ?", id).First(&BankInfo).Error
	return
}

func (t *BankInfo) GetFirstByUid(db *gorm.DB, uid int64) (BankInfo BankInfo, err error) {
	err = db.Model(t).Where("uid = ?", uid).First(&BankInfo).Error
	return
}

func (t *BankInfo) GetList(db *gorm.DB, BankInfoCond *BankInfo) (list []BankInfo, err error) {
	err = db.Model(t).Where(BankInfoCond).Find(&list).Error
	return
}

func (t *BankInfo) GetListByUserIds(db *gorm.DB, userIds []int) (list []BankInfo, err error) {
	err = db.Model(t).Where("uid in (?)", userIds).Find(&list).Error
	return
}

// 分页列表
func (t *BankInfo) GetPageList(db *gorm.DB, BankInfoCond *BankInfo, page, size int, channelIds []int) (total int64, list []BankInfo, err error) {
	localDb := db.Model(t).Where(BankInfoCond).Where("channel in (?)", channelIds)

	err = localDb.Count(&total).Error

	if err != nil {
		return
	}

	err = localDb.Scopes(model.Paginate(page, size)).Find(&list).Error
	return
}
