package game

import (
	"gorm.io/gorm"
	"time"
)

type GameBase struct {
	ID        int            `gorm:"primarykey;column:id;type:int(11);"`
	Uuid      uint32         `gorm:"column:uuid;type:int(11);not null;default:0;index:unique;"`                              //游戏ID
	GameName  string         `db:"game_name" gorm:"column:game_name;type:varchar(32);not null;default:'';";json:"game_name"` //游戏名称
	EnName    string         `db:"en_name" gorm:"column:en_name;type:varchar(32);not null;default:'';";json:"en_name"`       //英文名
	Type      uint8          `db:"type" gorm:"column:type;type:tinyint(4);not null;default:1;";json:"type"`                  //类型(1=tp)
	Icon      string         `db:"icon" gorm:"column:icon;type:varchar(256);not null;default:'';";json:"icon"`               //图片
	ApiUrl    string         `db:"api_url" gorm:"column:api_url;type:varchar(512);not null;default:'';";json:"api_url"`      //api地址
	Status    uint8          `db:"status" gorm:"column:status;type:tinyint(4);not null;default:0;";json:"status"`            //状态（0=等待上线,1=上线,2=维护,3=下线）
	Sort      uint32         `db:"sort" gorm:"column:sort;type:int(11);not null;default:0;";json:"sort"`                     //排序（大的靠前）
	Remark    string         `db:"remark" gorm:"column:remark;type:varchar(64);default:'';";json:"remark"`                   //备注
	CreatedAt time.Time      `db:"created_at" gorm:"column:created_at;";json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" gorm:"column:updated_at;";json:"updated_at"`
	DeletedAt gorm.DeletedAt `db:"deleted_at" gorm:"column:deleted_at;index;"`
}

func (t *GameBase) TableName() (tableName string) {
	tableName = "game_base"
	return
}

func (t *GameBase) GetList(db *gorm.DB) (list []GameBase, err error) {
	err = db.Model(t).Find(&list).Error
	return
}

func (t *GameBase) GetListByIds(db *gorm.DB, ids []int) (list []GameBase, err error) {
	err = db.Model(t).Where("id in (?)", ids).Find(&list).Error
	return
}
