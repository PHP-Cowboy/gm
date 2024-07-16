package log

import (
	"gm/model"
	"gm/request"
	"gorm.io/gorm"
)

type SlotFundsFlowLog struct {
	model.Base
	UserName      string `gorm:"type:varchar(64);not null;comment:用户昵称"`   //用户昵称
	Uid           int    `gorm:"type:int(11);not null;comment:用户ID"`       //用户ID
	RoomId        int    `gorm:"type:int(11);not null;comment:房间id"`       //房间id
	DeskId        int    `gorm:"type:int(11);not null;comment:桌子id"`       //桌子id
	BeforeCash    int    `gorm:"type:int(11);not null;comment:变化前cash"`    //变化前cash
	BeforeWinCash int    `gorm:"type:int(11);not null;comment:变化前winCash"` //变化前winCash
	Cash          int    `gorm:"type:int(11);not null;comment:赠送金币账户"`     //赠送金币账户
	WinCash       int    `gorm:"type:int(11);not null;comment:充值金币账户"`     // 充值金币账户
	Nums          int    `gorm:"type:int(11);not null;comment:变动数额"`       // 变动数额
	Tax           int    `gorm:"type:int(11);not null;comment:税收"`         //税收
	Balance       int    `gorm:"type:int(11);not null;comment:账户余额"`       //账户余额
	Remark        string `gorm:"type:varchar(64);not null;comment:备注"`     //备注
}

func (t *SlotFundsFlowLog) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *SlotFundsFlowLog) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *SlotFundsFlowLog) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *SlotFundsFlowLog) Count(db *gorm.DB) (total int64, err error) {
	err = db.Model(t).Where(t).Count(&total).Error
	return
}

func (t *SlotFundsFlowLog) GetFirstById(db *gorm.DB, id int64) (SlotFundsFlowLog SlotFundsFlowLog, err error) {
	err = db.Model(t).Where("id = ?", id).First(&SlotFundsFlowLog).Error
	return
}

func (t *SlotFundsFlowLog) GetList(db *gorm.DB) (list []SlotFundsFlowLog, err error) {
	err = db.Model(t).Where(t).Find(&list).Error
	return
}

func (t *SlotFundsFlowLog) GetPageList(db *gorm.DB, req request.SlotFundsFlowLog) (list []SlotFundsFlowLog, err error) {
	err = db.Model(t).Where(t).Scopes(model.Paginate(req.Page, req.Size)).Find(&list).Error
	return
}
