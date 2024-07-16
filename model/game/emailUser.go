package game

import (
	"gm/global"
	"gorm.io/gorm"
	"time"
)

type EmailUser struct {
	ID          int        `gorm:"primaryKey;column:id;type:int(11);"`
	EmailId     uint32     `gorm:"column:email_id;type:int(11);not null;default:0;"`       //邮件ID
	Type        uint8      `gorm:"column:type;type:tinyint(4);not null;default:1;"`        //邮件类型(1=普通邮件,2=赠送退款邮件)
	Uid         int        `gorm:"column:uid;type:bigint(20);not null;default:0;"`         //用户ID
	Title       string     `gorm:"column:title;type:varchar(64);not null;default'';"`      //邮件标题
	Msg         string     `gorm:"column:msg;type:varchar(1024);not null;default:'';"`     //邮件内容
	Status      uint8      `gorm:"column:status;type:tinyint(4);not null;default:0;"`      //读状态(0=未读,1=已读)
	ReadTime    *time.Time `gorm:"column:read_time;"`                                      //读取时间
	GetTime     *time.Time `gorm:"column:get_time;"`                                       //兑换时间
	IsAnnex     uint8      `gorm:"column:is_annex;type:tinyint(4);not null;default:0;"`    //是否有附件
	AnnexIds    string     `gorm:"column:annex_ids;type:varchar(128);not null;default:0;"` //附件IDS(多个用,分割)
	InsertTime  uint32     `gorm:"column:insert_time;type:int(11);not null;default:0;"`    //插入时间戳
	CreatedAt   time.Time  `gorm:"column:created_at;"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;"`
	HasAttach   uint8      `gorm:"type:tinyint(4);not null;default:0;"` //是否有附件(1=有,0=否)
	Attachments string     `gorm:"type:varchar(255);"`                  //附件信息
}

func (t *EmailUser) TableName(uid int) (tableName string) {
	tableNum := uid % 10

	switch tableNum {
	case 0:
		tableName = "email_00"
	case 1:
		tableName = "email_01"
	case 2:
		tableName = "email_02"
	case 3:
		tableName = "email_03"
	case 4:
		tableName = "email_04"
	case 5:
		tableName = "email_05"
	case 6:
		tableName = "email_06"
	case 7:
		tableName = "email_07"
	case 8:
		tableName = "email_08"
	case 9:
		tableName = "email_09"
	}
	return
}

func (t *EmailUser) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

// 事务方法
func (t *EmailUser) SaveEmailByMp(tx *gorm.DB, mp map[int]EmailUser) (err error) {
	emailUserMp := make(map[int][]EmailUser)

	for _, email := range mp {
		id := email.Uid % 10

		emailUser, ok := emailUserMp[id]
		if !ok {
			emailUser = make([]EmailUser, 0, 32)
		}

		emailUser = append(emailUser, email)

		emailUserMp[id] = emailUser
	}

	for remainder, email := range emailUserMp {
		err = tx.Table(t.TableName(remainder)).CreateInBatches(&email, 100).Error
		if err != nil {
			global.Logger["err"].Errorf("EmailUser SaveEmailByMp failed,err:[%v]", err.Error())
			return
		}
	}
	return
}

// 事务方法
func (t *EmailUser) SaveByUserIds(tx *gorm.DB, userIds []uint64, email EmailUser) (err error) {
	userMp := make(map[uint64][]EmailUser)

	for _, id := range userIds {
		email.Uid = int(id)
		userMp[id] = append(userMp[id], email)
	}

	for remainder, user := range userMp {
		err = tx.Table(t.TableName(int(remainder))).CreateInBatches(&user, 100).Error
		if err != nil {
			global.Logger["err"].Errorf("EmailUser SaveByUserIds failed,err:[%v]", err.Error())
			return
		}
	}
	return
}

func (t *EmailUser) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *EmailUser) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete(t, id).Error
	return
}

func (t *EmailUser) GetFirstById(db *gorm.DB, id int64) (EmailUser EmailUser, err error) {
	err = db.Model(t).Where("id = ?", id).First(&EmailUser).Error
	return
}

func (t *EmailUser) GetList(db *gorm.DB) (list []EmailUser, err error) {
	err = db.Model(t).Find(&list).Error

	return
}
