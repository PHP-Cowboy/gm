package game

import (
	"gorm.io/gorm"
	"time"
)

type TaskConfig struct {
	ID        uint64    `gorm:"primarykey;column:id;type:int(11);";json:"id"`
	Group     string    `gorm:"column:group;type:varchar(32);not null;default:'';"json:"group"`          // 组
	Status    uint8     `gorm:"column:status;type:tinyint(4);not null;default:1;"json:"status"`          // 状态（1=可用，0=不可用）
	Type      uint8     `gorm:"column:type;type:tinyint(4);not null;default:1;"json:"type"`              // 类别（1=普通，2=特殊）
	GameId    uint32    `gorm:"column:game_id;type:int(11);not null;default:0;";json:"game_id"`          // 游戏ID（0=所有,>0 单个游戏ID）
	GameName  string    `gorm:"column:game_name;type:varchar(32);not null;default:'';";json:"game_name"` // 游戏名称
	MoneyType uint8     `gorm:"column:money_type;type:tinyint(4);not null;default:1;";json:"money_type"` // 货币类型(1=无限制,2=cash,3=储钱罐)
	TaskType  uint8     `gorm:"column:task_type;type:tinyint(4);not null;default:1;";json:"task_type"`   // 类型(1=局数,2=胜局,3=累计赢金额)
	Num       uint32    `gorm:"column:num;type:int(11);not null;default:1;";json:"num"`                  // 条件
	Prize     uint16    `gorm:"column:prize;type:int(11);not null;default:1;";json:"prize"`              // 奖励
	Remark    string    `gorm:"column:remark;type:varchar(128);not null;default:'';";json:"remark"`      // 备注
	CreatedAt time.Time `gorm:"column:created_at;";json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;";json:"updated_at"`
}

func (t *TaskConfig) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *TaskConfig) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *TaskConfig) GetFirstById(db *gorm.DB, id int64) (TaskConfig TaskConfig, err error) {
	err = db.Model(t).Where("id = ?", id).First(&TaskConfig).Error
	return
}

func (t *TaskConfig) GetList(db *gorm.DB) (list []TaskConfig, err error) {
	err = db.Model(t).Find(&list).Error
	return
}
