package game

import (
	"gm/model"
	"gm/request"
	"gorm.io/gorm"
)

type Roomview struct {
	Id            uint64 `gorm:"column:id"`
	RoomId        uint64 `gorm:"column:RoomId;type:int(11);"`
	SvrId         uint32 `gorm:"column:Svrid;type:int(11);not null;default:0;"`          //服务器id
	GameId        uint32 `gorm:"column:GameId;type:int(11);not null;default:0;"`         //游戏ID
	RoomIndex     uint32 `gorm:"column:RoomIndex;type:int(11);not null;default:0;"`      //房间index
	Base          uint32 `gorm:"column:Base;type:int(11);not null;default:0;"`           //底注
	MinEntry      int    `gorm:"column:MinEntry;type:int(11);not null;default:0;"`       //进入限制(下) 0代表无限制
	MaxEntry      int    `gorm:"column:MaxEntry;type:int(11);not null;default:0;"`       //进入限制(上) 0代表无限制
	RoomName      string `gorm:"column:RoomName;type:varchar(100);not null;default:'';"` //房间名称
	RoomType      uint8  `gorm:"column:RoomType;type:int(11);not null;default:0;"`       //类型 1体验大厅 2正常大厅
	RoomSwitch    int    `gorm:"column:RoomSwitch;type:tinyint(4);not null;default:0;"`  //房间开关
	RoomWelfare   int    `gorm:"column:RoomWelfare;type:tinyint(4);not null;default:0;"` //房间赠送
	Desc          string `gorm:"column:Desc;type:varchar(255);not null;default:'';"`     //房间描述
	Tax           int    `gorm:"column:Tax;type:int(11);not null;default:0;"`            //税千分比
	BonusDiscount int    `gorm:"column:BonusDiscount;type:int(11);not null;default:0;"`  //比例千分比
	AiSwitch      int    `gorm:"column:AiSwitch;type:tinyint(4);not null;default:0;"`    //ai开关
	AiLimit       int    `gorm:"column:AiLimit;type:int(11);not null;default:0;"`        //ai人数限制
	RechargeLimit int    `gorm:"column:RechargeLimit;type:int(11);not null;default:0;"`  //充值准入值
	PoolID        int    `gorm:"column:PoolID;type:int(11);not null;default:0;"`         //奖金池ID
	ExtData       string `gorm:"column:ExtData;type:varchar(500);not null;default:'';"`  //特殊配置
}

func (t *Roomview) UpdateByRoomId(db *gorm.DB, roomId uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("RoomId = ?", roomId).Updates(mp).Error
	return
}

func (t *Roomview) GetList(db *gorm.DB) (list []Roomview, err error) {
	err = db.Model(t).Find(&list).Error

	return
}

func (t *Roomview) GetPageList(db *gorm.DB, req request.RoomList) (total int64, list []Roomview, err error) {
	localDb := db.Model(t).Where(&t)

	if req.RoomName != "" {
		localDb.Where("RoomName like ?", "%"+req.RoomName+"%")
	}

	err = localDb.Count(&total).Error

	if err != nil {
		return
	}

	err = localDb.Order("GameId asc,RoomId asc").Scopes(model.Paginate(req.Page, req.Size)).Find(&list).Error

	return
}

func (t *Roomview) GetListByRoomIds(db *gorm.DB, roomIds []uint64) (list []Roomview, err error) {
	err = db.Model(t).Where("RoomId in (?)", roomIds).Find(&list).Error

	return
}
