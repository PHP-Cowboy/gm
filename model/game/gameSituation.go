package game

import (
	"fmt"
	"gm/model"
	"gorm.io/gorm"
)

// 用户游戏战况
type GameSituation struct {
	model.Base
	Uid       int    `gorm:"type:int(11);uniqueIndex:userGameType;not null;comment:用户id"`
	GameName  string `gorm:"type:varchar(32);uniqueIndex:userGameType;not null;default:'';comment:游戏名称"`
	RoomType  string `gorm:"type:varchar(32);uniqueIndex:userGameType;not null;default:'';comment:房间类型"`
	Total     int    `gorm:"type:int(11);not null;comment:总局数"`
	WinNum    int    `gorm:"type:int(11);not null;comment:赢局数"`
	WinMoney  int    `gorm:"type:int(11);not null;comment:总赢额度"`
	LossMoney int    `gorm:"type:int(11);not null;comment:总输额度"`
}

func (t *GameSituation) TableName(uid int) string {
	return "game_situation_0" + fmt.Sprintf("%v", uid%10)
}

func (t *GameSituation) GetList(db *gorm.DB, uid int) (list []GameSituation, err error) {
	err = db.Table(t.TableName(uid)).Where("uid = ?", uid).Find(&list).Error
	return
}
