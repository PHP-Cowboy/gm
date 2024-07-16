package log

import (
	"gm/model"
	"gm/request"
	"gorm.io/gorm"
)

type TpFundsFlowLog struct {
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

type TpFundsFlowLogSqlx struct {
	model.Base
	UserName      string //用户昵称
	Uid           int    //用户ID
	FlowId        string //流水号
	RoomId        int    //房间id
	DeskId        int    //桌子id
	BeforeCash    int    //变化前cash
	BeforeWinCash int    //变化前winCash
	Cash          int    //赠送金币账户
	WinCash       int    // 充值金币账户
	Nums          int    // 变动数额
	Tax           int    //税收
	Balance       int    //账户余额
	Remark        string //备注
}

func (t *TpFundsFlowLog) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *TpFundsFlowLog) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *TpFundsFlowLog) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *TpFundsFlowLog) Count(db *gorm.DB) (total int64, err error) {
	err = db.Model(t).Where(t).Count(&total).Error
	return
}

func (t *TpFundsFlowLog) GetFirstById(db *gorm.DB, id int64) (TpFundsFlowLog TpFundsFlowLog, err error) {
	err = db.Model(t).Where("id = ?", id).First(&TpFundsFlowLog).Error
	return
}

func (t *TpFundsFlowLog) GetList(db *gorm.DB) (list []TpFundsFlowLog, err error) {
	err = db.Model(t).Where(t).Find(&list).Error
	return
}

func (t *TpFundsFlowLog) GetPageList(db *gorm.DB, req request.TpFundsFlowLog) (list []TpFundsFlowLog, err error) {
	err = db.Model(t).Where(t).Scopes(model.Paginate(req.Page, req.Size)).Find(&list).Error
	return
}
