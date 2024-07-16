package log

import (
	"fmt"
	"gm/global"
	"gm/model"
	"gm/request"
	"gm/utils/timeutil"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type UserLoginLog struct {
	model.Base
	Uid                int       `gorm:"type:int;not null;comment:用户id"`
	Nickname           string    `gorm:"type:varchar(64);not null;default:'';comment:用户昵称"`
	ChannelId          int       `gorm:"type:int;not null;comment:渠道id"`
	Channel            string    `gorm:"type:varchar(32);not null;default:'';comment:渠道名称"`
	Assets             int       `gorm:"type:int;not null;default:0;comment:总资产cash+winCash"`
	ReferralCommission int       `gorm:"type:int;not null;default:0;comment:推荐佣金"`
	Ip                 string    `gorm:"type:varchar(16);not null;default:'';comment:登录ip"`
	Device             string    `gorm:"type:varchar(64);not null;default:'';comment:登录设备号"`
	Version            string    `gorm:"type:varchar(64);not null;default:'';comment:登录版本号"`
	LoginMode          int       `gorm:"type:tinyint(4);not null;default:0;comment:登录方式(1:游客,2:手机号)"`
	RegTime            time.Time `gorm:"type:datetime;not null;comment:注册时间"`
}

func (t *UserLoginLog) TableName(ym string) string {
	if ym == "" {
		ym = time.Now().Format(timeutil.MonthNumberFormat)
	}
	return "user_login_log_" + ym
}

func (t *UserLoginLog) UnionAllByYmd(req request.GetLoginLogList) (string, error) {
	var (
		monthList []time.Time
		err       error
	)

	if req.StartTime == "" && req.EndTime == "" {
		now := time.Now()
		//默认本月和上个月
		req.StartTime = now.AddDate(0, -1, 0).Format(timeutil.TimeFormat)
		req.EndTime = now.Format(timeutil.TimeFormat)
	}

	monthList, err = timeutil.GetMonthsInRange(req.StartTime, req.EndTime, timeutil.TimeFormat)

	if err != nil {
		global.Logger["err"].Errorf("(t *Order) UnionAllByYmd timeutil.GetMonthsInRange failed,err:[%v]", err.Error())
		return "", err
	}

	var (
		sb  strings.Builder
		str strings.Builder
	)

	if req.Uid > 0 {
		str.WriteString(fmt.Sprintf(" and uid = %v ", req.Uid))
	}

	if req.LoginMode > 0 {
		str.WriteString(fmt.Sprintf(" and login_mode = %v ", req.LoginMode))
	}

	if len(req.CreatedAt) > 0 {
		str.WriteString(fmt.Sprintf(" and created_at >= '%s' and created_at <= '%s' ", req.StartTime, req.EndTime))
	}

	if req.ChannelId == 0 {
		channels := ""

		for _, id := range req.ChannelIds {
			channels += strconv.Itoa(id) + ","
		}

		channels = strings.TrimRight(channels, ",")

		if channels != "" {
			str.WriteString(fmt.Sprintf(" and channel_id in (%v)", channels))
		}
	} else {
		str.WriteString(fmt.Sprintf(" and channel_id = %v ", req.ChannelId))
	}

	for i, m := range monthList {
		if i > 0 {
			sb.WriteString(" UNION ALL ")
		}

		sb.WriteString(
			fmt.Sprintf(
				"SELECT * FROM %s where 1=1 %s",
				t.TableName(m.Format(timeutil.MonthNumberFormat)),
				str.String(),
			),
		)
	}

	return sb.String(), nil
}

func (t *UserLoginLog) GetPageList(db *gorm.DB, req request.GetLoginLogList) (total int64, list []UserLoginLog, err error) {
	var sql string
	sql, err = t.UnionAllByYmd(req)
	if err != nil {
		global.Logger["err"].Errorf("(t *UserLoginLog) GetPageList t.UnionAllByYmd failed,err:[%v]", err.Error())
		return 0, nil, err
	}

	err = db.Raw(fmt.Sprintf(`select count(id) from (%s) as t`, sql)).Scan(&total).Error

	if err != nil {
		global.Logger["err"].Errorf("Order db.Raw count failed,err:[%v]", err.Error())
		return
	}

	limit := (req.Page - 1) * req.Size

	if limit < 0 {
		limit = 0
	}

	err = db.Raw(fmt.Sprintf(`select * from (%s) as t order by t.created_at desc limit %v,%v`, sql, limit, req.Size)).Scan(&list).Error

	if err != nil {
		global.Logger["err"].Errorf("Order Find failed,err:[%v]", err.Error())
		return
	}

	return
}

func (t *UserLoginLog) GetLastLoginInfoByUid(db *gorm.DB, uid int) (data UserLoginLog, err error) {
	err = db.Table(t.TableName("")).Order("id desc").Where("uid = ?", uid).First(&data).Error
	return
}

func (t *UserLoginLog) GetListByUserIds(db *gorm.DB, ids []int, dateTime time.Time) (mp map[int]struct{}, err error) {
	var (
		dataList   []UserLoginLog
		start, end string
	)

	start = timeutil.GetZeroTime(dateTime).Format(timeutil.TimeFormat)
	end = timeutil.GetLastTime(dateTime).Format(timeutil.TimeFormat)

	err = db.Table(t.TableName(dateTime.Format(timeutil.MonthNumberFormat))).
		Select("uid").
		Where("uid in (?) and created_at >= ? and created_at <= ?", ids, start, end).
		Find(&dataList).
		Error

	if err != nil {
		return
	}

	mp = make(map[int]struct{})

	for _, d := range dataList {
		mp[d.Uid] = struct{}{}
	}

	return
}

func (t *UserLoginLog) GetListByUserIdsAndCreatedAt(db *gorm.DB, ids []int, dateTime time.Time) (mp map[int]struct{}, err error) {
	var dataList []UserLoginLog

	err = db.Table(t.TableName(dateTime.Format(timeutil.MonthNumberFormat))).
		Select("uid").
		Where("uid in (?)", ids).
		Where("created_at >= ? and created_at <= ?", timeutil.GetZeroTime(dateTime), timeutil.GetLastTime(dateTime)).
		Find(&dataList).
		Error

	if err != nil {
		return
	}

	mp = make(map[int]struct{})

	for _, d := range dataList {
		mp[d.Uid] = struct{}{}
	}

	return
}

func (t *UserLoginLog) GetListByCreatedAt(db *gorm.DB, start, end, ym string) (mp, userChannelMp map[int]int, err error) {
	var dataList []UserLoginLog

	err = db.Table(t.TableName(ym)).
		Select("uid,channel_id").
		Where("created_at >= ? and created_at <= ?", start, end).
		Find(&dataList).
		Error

	if err != nil {
		return
	}

	mp = make(map[int]int)
	userChannelMp = make(map[int]int)

	uidMp := make(map[int]struct{})

	for _, d := range dataList {
		userChannelMp[d.Uid] = d.ChannelId

		_, ok := uidMp[d.Uid]

		//同一用户多次登录去重
		if ok {
			continue
		}

		//渠道登录用户数
		mp[d.ChannelId]++
		uidMp[d.Uid] = struct{}{}
	}

	return
}

func (t *UserLoginLog) CountChannelLoginNumByCreateAt(db *gorm.DB, start, end, ym string) (mp map[int]int, err error) {
	var dataList []UserLoginLog

	err = db.Table(t.TableName(ym)).
		Select("uid,channel_id").
		Where("created_at >= ? and created_at <= ?", start, end).
		Find(&dataList).
		Error

	if err != nil {
		return
	}

	mp = make(map[int]int)

	uidMp := make(map[int]struct{})

	for _, d := range dataList {

		_, ok := uidMp[d.Uid]

		//同一用户多次登录去重
		if ok {
			continue
		}

		//渠道登录用户数
		mp[d.ChannelId]++
		uidMp[d.Uid] = struct{}{}
	}

	return
}

func (t *UserLoginLog) CountByCreatedAt(db *gorm.DB, start, end string) (count int64, err error) {

	//err = db.Table(t.TableName("")).
	//	Select("distinct(uid) as count").
	//	Where("created_at >= ? and created_at <= ?", start, end).
	//	Count(&count).
	//	Error

	sql := fmt.Sprintf("SELECT COUNT(DISTINCT(uid)) FROM %s where created_at >= ? and created_at <= ?", t.TableName(""))

	err = db.Raw(sql, start, end).Scan(&count).Error

	return
}

func (t *UserLoginLog) CountByUserIdsCreatedAt(db *gorm.DB, userIds []int, start, end string) (count int64, err error) {

	//err = db.Table(t.TableName("")).
	//	Select("distinct(uid) as count").
	//	Where("created_at >= ? and created_at <= ?", start, end).
	//	Count(&count).
	//	Error

	sql := fmt.Sprintf("SELECT COUNT(DISTINCT(uid)) FROM %s where uid in (?) and created_at >= ? and created_at <= ?", t.TableName(""))

	err = db.Raw(sql, userIds, start, end).Scan(&count).Error

	return
}
