package game

import (
	"gorm.io/gorm"
	"time"
)

type HallConfig struct {
	ID        uint64    `gorm:"primarykey;column:id;type:int(11);";json:"id"`
	Group     string    `gorm:"column:group;type:varchar(64);not null;default''" json:"group"`           //相同分组的按照sort顺序依次弹出
	Type      uint8     `gorm:"column:type;type:tinyint(4);not null;default:0;";json:"type"`             //弹出类型（1=签到，2=礼包，3=每日任务）
	ModelNum  uint8     `gorm:"column:model_num;type:tinyint(4);not null;default:1;";json:"model_num"`   //弹出频率
	ModelType uint8     `gorm:"column:model_type;type:tinyint(4);not null;default:0;";json:"model_type"` //频率类型（1=每天总次数，2=间隔分钟数，3=每次进入大厅）
	Sort      uint8     `gorm:"column:sort;type:tinyint(4);not null;default:1;";json:"sort"`             //排序（大的优先弹出）
	CreatedAt time.Time `gorm:"column:created_at;";json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;";json:"updated_at"`
}

func (t *HallConfig) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *HallConfig) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *HallConfig) GetFirstById(db *gorm.DB, id int64) (HallConfig HallConfig, err error) {
	err = db.Model(t).Where("id = ?", id).First(&HallConfig).Error
	return
}

func (t *HallConfig) GetList(db *gorm.DB) (list []HallConfig, err error) {
	err = db.Model(t).Find(&list).Error
	return
}
