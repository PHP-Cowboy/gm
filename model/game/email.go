package game

import (
	"gm/model"
	"gorm.io/gorm"
	"time"
)

type Email struct {
	ID          uint64    `gorm:"primaryKey;column:id;type:int(11);"`
	Type        uint8     `gorm:"type:tinyint(4);not null;default:1;"`     //邮件类型(1=系统邮件,2=奖励邮件)
	Title       string    `gorm:"type:varchar(64);not null;default:'';"`   //邮件标题
	Msg         string    `gorm:"type:varchar(1024);not null;default:'';"` //邮件内容
	Status      uint8     `gorm:"type:tinyint(4);not null;default:0;"`     //状态(0=未发送,1=已发送,2=发送中)
	IsAnnex     uint8     `gorm:"type:tinyint(4);not null;default:0;"`     //是否有附件(1=有,0=否) todo del
	AnnexIds    string    `gorm:"type:varchar(128);not null;default:'';"`  //附件IDS todo del
	SendType    uint8     `gorm:"type:tinyint(4);not null;default:1;"`     //发送类型(1=全服发送,2=在线玩家发送,3=指定玩家)
	UserIds     string    `gorm:"type:varchar(2048);not null;default:'';"` //接受玩家IDS
	StartTime   int64     `gorm:"column:start_time;"`                      //开始时间
	EndTime     int64     `gorm:"column:end_time;"`                        //结束时间
	EventId     uint64    `gorm:"type:int(11);default:0;"`                 //关联事件id
	Condition   uint64    `gorm:"type:int(11);default:0;"`                 //发放条件
	CreatedAt   time.Time `gorm:"column:created_at;"`
	UpdatedAt   time.Time `gorm:"column:updated_at;"`
	HasAttach   uint8     `gorm:"type:tinyint(4);not null;default:0;"` //是否有附件(1=有,0=否)
	Attachments string    `gorm:"type:varchar(255);"`                  //附件信息
}

func (t *Email) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *Email) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *Email) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete(t, id).Error
	return
}

func (t *Email) GetFirstById(db *gorm.DB, id int64) (Email Email, err error) {
	err = db.Model(t).Where("id = ?", id).First(&Email).Error
	return
}

func (t *Email) GetList(db *gorm.DB) (list []Email, err error) {
	err = db.Model(t).Find(&list).Error

	return
}

func (t *Email) GetPageList(db *gorm.DB, page, size int) (list []Email, err error) {
	err = db.Model(t).Where(&t).Scopes(model.Paginate(page, size)).Find(&list).Error

	return
}
