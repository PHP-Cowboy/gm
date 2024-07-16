package user

import (
	"gm/global"
	"gm/model"
	"gm/request"
	"gm/utils/slice"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id         int       `gorm:"primarykey;column:id;type:int(11);"`
	Uid        int       `gorm:"column:uid;type:bigint(20);not null;default:0;index:unique;"` //用户ID
	IsGuest    uint8     `gorm:"column:is_guest;type:tinyint(4);not null;default:0;"`         //是否游客
	IsSend     uint8     `gorm:"column:is_send;type:tinyint(4);not null;default:0;"`          //是否赠送
	Device     string    `gorm:"column:device;type:varchar(64);not null;default:'';"`         //设备码
	UserName   string    `gorm:"column:user_name;type:varchar(32);not null;default:'';"`      //用户名
	Icon       int8      `gorm:"column:icon;type:tinyint(4);not null;default:0;"`             //头像
	Phone      string    `gorm:"column:phone;type:varchar(16);not null;default:'';"`          //电话
	Email      string    `gorm:"column:email;type:varchar(32);not null;default:''1;"`         //邮箱
	Pwd        string    `gorm:"column:pwd;type:varchar(32);not null;default:'';"`            //密码
	Token      string    `gorm:"column:token;type:varchar(512);not null;default:'';"`         //token
	ChannelId  int       `gorm:"column:channel_id;type:int(11);not null;default:0;"`          //用户ID
	Expire     time.Time `gorm:"column:expire;"`
	TpNew      int       `gorm:"column:tp_new"`
	RegIp      string    `gorm:"column:reg_ip"`      //注册ip
	RegVersion string    `gorm:"column:reg_version"` //注册版本号
	CreatedAt  time.Time `gorm:"column:created_at;"`
	UpdatedAt  time.Time `gorm:"column:updated_at;"`
	Gpcadid    string    `gorm:"column:gpcadid"` //gpcadid
}

func (t *User) TableName() string {
	return "user"
}

func (t *User) UpdateById(db *gorm.DB, id int, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *User) GetOneByUid(db *gorm.DB, uid int) (u User, err error) {
	err = db.Table(t.TableName()).Where("uid = ?", uid).First(&u).Error
	return
}

func (t *User) GetPageList(db *gorm.DB, req request.GetGameUserList) (total int64, list []User, err error) {
	localDb := db.Table(t.TableName())

	if req.Uid > 0 {
		localDb.Where("uid = ?", req.Uid)
	} else if len(req.UserIds) > 0 {
		localDb.Where("uid in (?)", req.UserIds)
	}

	localDb.Where(User{
		ChannelId: req.ChannelId,
	})

	if req.IsGuest != nil {
		localDb.Where("is_guest = ?", *req.IsGuest)
	}

	if req.StartCreatedAt != "" {
		localDb.Where("created_at >= ?", req.StartCreatedAt)
	}

	if req.EndCreatedAt != "" {
		localDb.Where("created_at <= ?", req.EndCreatedAt)
	}

	if req.StartUpdatedAt != "" {
		localDb.Where("updated_at >= ?", req.StartUpdatedAt)
	}

	if req.EndUpdatedAt != "" {
		localDb.Where("updated_at <= ?", req.EndUpdatedAt)
	}

	localDb.Where("channel_id in (?)", req.ChannelIds)

	err = localDb.Order("id desc").Count(&total).Error

	if err != nil {
		return
	}

	err = localDb.
		Scopes(model.Paginate(req.Page, req.Size)).
		Find(&list).
		Error
	return
}

func (t *User) GetListByUserIds(db *gorm.DB, userIds []int) (list []User, err error) {
	err = db.Model(t).Where("uid in (?)", userIds).Find(&list).Error
	return
}

func (t *User) MapOfChannelByUserIds(db *gorm.DB, userIds []int) (mp map[int]int, err error) {
	var dataList []User
	err = db.Model(t).Select("uid,channel_id").Where("uid in (?)", userIds).Find(&dataList).Error

	if err != nil {
		global.Logger["err"].Infof("查询用户数据失败:" + err.Error())
		return
	}

	mp = make(map[int]int, len(dataList))

	for _, d := range dataList {
		mp[d.Uid] = d.ChannelId
	}

	return
}

func (t *User) GetMpByUserIds(db *gorm.DB, userIds []int) (mp map[int]User, err error) {
	var list []User

	err = db.Model(t).Where("uid in (?)", userIds).Find(&list).Error

	if err != nil {
		return
	}

	mp = make(map[int]User)

	for _, l := range list {
		mp[l.Uid] = l
	}
	return
}

func (t *User) GetPageListByCreateTime(db *gorm.DB, start, end string, page, size int) (dataList []User, err error) {
	err = db.Model(t).
		Where("created_at >= ? and created_at <= ?", start, end).
		Scopes(model.Paginate(page, size)).
		Find(&dataList).
		Error
	return
}

func (t *User) CountByCreateTime(db *gorm.DB, start, end string) (count int64, err error) {
	err = db.Model(t).
		Where("created_at >= ? and created_at <= ?", start, end).
		Count(&count).
		Error
	return
}

func (t *User) GetListByCreateTime(db *gorm.DB, start, end string) (userIds []int, err error) {
	dataList := make([]User, 0)

	err = db.Model(t).
		Select("uid").
		Where("created_at >= ? and created_at <= ?", start, end).
		Find(&dataList).
		Error

	if err != nil {
		return
	}

	userIds = make([]int, 0, len(dataList))

	for _, d := range dataList {
		userIds = append(userIds, d.Uid)
	}

	userIds = slice.UniqueSlice(userIds)

	return
}

func (t *User) CountChannelRegNumByCreateTime(db *gorm.DB, start, end string) (map[int]int, error) {
	dataList := make([]User, 0)

	err := db.Model(t).
		Select("uid,channel_id").
		Where("created_at >= ? and created_at <= ?", start, end).
		Find(&dataList).
		Error

	if err != nil {
		return nil, err
	}

	mp := make(map[int]int)

	for _, dl := range dataList {
		mp[dl.ChannelId]++
	}

	return mp, nil
}

func (t *User) CountChannelNumByUserIds(db *gorm.DB, userIds []uint64) (map[int]int, error) {
	dataList := make([]User, 0)

	err := db.Model(t).Where("uid in (?)", userIds).Find(&dataList).Error

	if err != nil {
		return nil, err
	}

	mp := make(map[int]int)

	for _, dl := range dataList {
		mp[dl.ChannelId]++
	}

	return mp, nil
}
