package pay

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

type Order struct {
	ID           int       `gorm:"primaryKey;column:id;type:int(11);"`
	Uid          int       `gorm:"column:uid;type:bigint(20);not null;default:0;index:uid_key;"`                         // 用户ID
	Ymd          int       `gorm:"column:ymd;type:int(11);not null;default:0;"`                                          // 年月日
	OrderNo      string    `gorm:"column:order_no;type:varchar(32);not null;default:'';unique:order_key;index:uid_key;"` // 订单ID
	MOrderNo     string    `gorm:"column:m_order_no;type:varchar(32);not null;default:'';unique:m_order_key;"`           // 商户订单ID
	PayId        int       `gorm:"column:pay_id;type:int(11);not null;default:0;"`                                       // 支付配置id
	Account      int       `gorm:"column:account;type:int(11);not null;default:0;"`                                      // 支付金额
	Cash         int       `gorm:"column:cash;type:int(11);not null;default:0;"`                                         //cash金额
	GiftCash     int       `gorm:"column:gift_cash;type:int(11);not null;default:0;"`                                    //额外赠送cash
	Bonus        int       `gorm:"column:bonus;type:int(11);not null;default:0;"`                                        //bonus
	RequestTime  time.Time `gorm:"column:request_time;"`                                                                 //下单时间
	Email        string    `gorm:"column:email;type:varchar(32);not null;default:'';"`                                   // 邮箱地址
	Name         string    `gorm:"column:name;type:varchar(32);not null;default:'';"`                                    // 用户名
	Phone        string    `gorm:"column:phone;type:varchar(32);not null;default:'';"`                                   //手机
	RedirectTime time.Time `gorm:"column:redirect_time;"`                                                                //下单拿h5地址时间
	Status       int8      `gorm:"column:status;type:tinyint(4);not null;default:0;"`                                    // 状态 0=等待支付，1=支付完成，2=下单失败
	CompleteTime int       `gorm:"column:complete_time;type:int(11);not null;default:0;"`                                //订单完成时间
	H5Url        string    `gorm:"column:h5_url;type:varchar(128);not null;default:'';"`                                 // h5地址
	Type         int       `gorm:"column:type;type:tinyint(4);not null;default:1;"`                                      // 类型 1:充值;2:礼包
	GiftId       int       `gorm:"column:gift_id;type:int(11);not null;default:0;"`                                      //类型 1:充值;2:礼包
	Remark       string    `gorm:"column:remark;type:varchar(64);not null;default:'';"`                                  //备注
	RoomId       string    `gorm:"column:room_id;type:varchar(255);default:null;comment:房间id，房间内充值数值记录"`                 //房间id，房间内充值数值记录
	Channel      int       `gorm:"column:channel;type:int(11);not null;default:1;"`                                      //渠道
	CreatedAt    time.Time `gorm:"column:created_at;"`
	UpdatedAt    time.Time `gorm:"column:updated_at;"`
}

func (t *Order) GetTableName(ym string) string {
	return "order_" + ym
}

func (t *Order) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *Order) UpdateById(db *gorm.DB, id uint64, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *Order) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *Order) GetFirstById(db *gorm.DB, id int64) (Order Order, err error) {
	err = db.Model(t).Where("id = ?", id).First(&Order).Error
	return
}

func (t *Order) GetList(db *gorm.DB, OrderCond Order) (list []Order, err error) {
	err = db.Model(t).Where(&OrderCond).Find(&list).Error
	return
}

// 分页列表
func (t *Order) GetPageList(db *gorm.DB, BankInfoCond *Order, page, size int) (list []Order, err error) {
	now := time.Now()
	//preOneMonth := now.AddDate(0, -1, 0).Format(timeutil.MonthNumberFormat)
	//preTwoMonth := now.AddDate(0, -2, 0).Format(timeutil.MonthNumberFormat)
	err = db.Table(t.GetTableName(now.Format(timeutil.MonthNumberFormat))).
		Where(BankInfoCond).
		Scopes(model.Paginate(page, size)).
		Find(&list).
		Error
	return
}

func (t *Order) UnionAllByYmd(req request.RechargeRecords) (string, error) {
	var (
		monthList []time.Time
		err       error
	)

	if req.Start == "" && req.End == "" {
		now := time.Now()
		//默认本月和上个月
		req.Start = now.AddDate(0, -1, 0).Format(timeutil.DateNumberFormat)
		req.End = now.Format(timeutil.DateNumberFormat)
	}

	monthList, err = timeutil.GetMonthsInRange(req.Start, req.End, timeutil.DateNumberFormat)

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

	if len(req.Ymd) > 0 {
		if req.Start == req.End {
			str.WriteString(fmt.Sprintf(" and ymd = %v ", req.Start))
		} else {
			str.WriteString(fmt.Sprintf(" and ymd >= %s and ymd <= %s ", req.Start, req.End))
		}
	}

	if req.OrderNo != "" {
		str.WriteString(fmt.Sprintf(" and order_no = '%s' ", req.OrderNo))
	}

	if req.Status != nil {
		str.WriteString(fmt.Sprintf(" and status = %v ", *req.Status))
	}

	if req.Channel == 0 {
		channels := ""

		for _, id := range req.ChannelIds {
			channels += strconv.Itoa(id) + ","
		}

		channels = strings.TrimRight(channels, ",")

		if channels != "" {
			str.WriteString(fmt.Sprintf(" and channel in (%v)", channels))
		}
	} else {
		str.WriteString(fmt.Sprintf(" and channel = %v ", req.Channel))
	}

	for i, m := range monthList {
		if i > 0 {
			sb.WriteString(" UNION ALL ")
		}

		sb.WriteString(
			fmt.Sprintf(
				"SELECT * FROM %s where 1=1 %s",
				t.GetTableName(m.Format(timeutil.MonthNumberFormat)),
				str.String(),
			),
		)
	}

	return sb.String(), nil
}

func (t *Order) RechargeRecords(db *gorm.DB, req request.RechargeRecords) (total int64, list []Order, err error) {

	var sql string

	sql, err = t.UnionAllByYmd(req)

	if err != nil {
		global.Logger["err"].Errorf("(t *Order) RechargeRecords UnionAllByYmd failed,err:[%v]", err.Error())
		return
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

func (t *Order) GetYesterdayTodayListByUid(db *gorm.DB, uid int) (mp map[string]int, err error) {
	var (
		ymds     []string
		list     []Order
		todayMax int
	)

	mp = make(map[string]int)
	mp["today"] = 0
	mp["yesterday"] = 0

	now := time.Now()

	today := now.Format(timeutil.DateNumberFormat)

	yesterday := now.AddDate(0, 0, -1).Format(timeutil.DateNumberFormat)

	ymds = append(ymds, today, yesterday)

	err = db.Table(t.GetTableName(now.Format(timeutil.MonthNumberFormat))).Where("uid = ? and status = 3", uid).Where("ymd in (?)", ymds).Find(&list).Error

	for _, o := range list {
		switch strconv.Itoa(o.Ymd) {
		case yesterday:
			mp["yesterday"] += o.Account

		case today:
			if o.Account > todayMax {
				todayMax = o.Account
			}
			mp["today"] += o.Account
		}
	}

	mp["today_recharge_max"] = todayMax

	return
}

func (t *Order) GetListByUserIdsAndYmd(db *gorm.DB, ids []int, dateTime time.Time) (mp map[int]int, err error) {
	var dataList []Order

	err = db.Table(t.GetTableName(dateTime.Format(timeutil.MonthNumberFormat))).
		Select("uid,account").
		Where("uid in (?)", ids).
		Where("ymd = ? ", dateTime.Format(timeutil.DateNumberFormat)).
		Find(&dataList).
		Error

	if err != nil {
		global.Logger["err"].Errorf("Order GetListByUserIdsAndYmd failed,err:[%v]", err.Error())
		return
	}

	mp = make(map[int]int)

	for _, d := range dataList {
		v, ok := mp[d.Uid]

		if !ok {
			v = 0
		}

		v += d.Account

		mp[d.Uid] = v
	}

	return
}

type AmountAndCash struct {
	Amount int
	Cash   int
}

func (t *Order) GetAmountAndCoinByUserIdsAndYmd(db *gorm.DB, ids []int, dateTime time.Time) (mp map[int]AmountAndCash, err error) {
	var dataList []Order

	err = db.Table(t.GetTableName(dateTime.Format(timeutil.MonthNumberFormat))).
		Select("uid,account,cash,gift_cash").
		Where("uid in (?) and status = 3", ids).
		Where("ymd = ? ", dateTime.Format(timeutil.DateNumberFormat)).
		Find(&dataList).
		Error

	if err != nil {
		global.Logger["err"].Errorf("Order GetAmountAndCoinByUserIdsAndYmd failed,err:[%v]", err.Error())
		return
	}

	mp = make(map[int]AmountAndCash)

	for _, d := range dataList {
		v, ok := mp[d.Uid]

		if !ok {
			v = AmountAndCash{}
		}

		v.Amount += d.Account
		v.Cash += (d.Cash + d.GiftCash)

		mp[d.Uid] = v
	}

	return
}

func (t *Order) GetListByYmd(db *gorm.DB, dateTime time.Time) (mp map[int]int, err error) {
	var dataList []Order

	err = db.Table(t.GetTableName(dateTime.Format(timeutil.MonthNumberFormat))).
		Select("uid,account").
		Where("ymd = ? and status = 3", dateTime.Format(timeutil.DateNumberFormat)).
		Find(&dataList).
		Error

	if err != nil {
		global.Logger["err"].Errorf("Order GetListByYmd failed,err:[%v]", err.Error())
		return
	}

	mp = make(map[int]int)

	for _, d := range dataList {
		v, ok := mp[d.Uid]

		if !ok {
			v = 0
		}

		v += d.Account

		mp[d.Uid] = v
	}

	return
}

func (t *Order) GetUidChannelListByYmd(db *gorm.DB, dateTime time.Time) (mp map[int]int, err error) {
	var dataList []Order

	err = db.Table(t.GetTableName(dateTime.Format(timeutil.MonthNumberFormat))).
		Select("uid,channel").
		Where("ymd = ? and status = 3", dateTime.Format(timeutil.DateNumberFormat)).
		Find(&dataList).
		Error

	if err != nil {
		global.Logger["err"].Errorf("Order GetListByYmd failed,err:[%v]", err.Error())
		return
	}

	mp = make(map[int]int)

	for _, d := range dataList {
		mp[d.Uid] = d.Channel
	}

	return
}

func (t *Order) CountRechargeByTime(db *gorm.DB, start, end string) (peopleNum, totalAmount int, err error) {
	var dataList []Order

	err = db.Table(t.GetTableName(time.Now().Format(timeutil.MonthNumberFormat))).
		Select("uid,account").
		Where("created_at >= ? and created_at <= ? and status = 3", start, end).
		Find(&dataList).
		Error

	if err != nil {
		global.Logger["err"].Errorf("Order CountRechargeByTime failed,err:[%v]", err.Error())
		return
	}

	mp := make(map[int]int)

	for _, d := range dataList {
		v, ok := mp[d.Uid]

		if !ok {
			v = 0
		}

		v += d.Account

		mp[d.Uid] = v
	}

	for _, v := range mp {
		peopleNum++
		totalAmount += v
	}

	return
}

type ChannelPeopleNumAndTotalAmount struct {
	PeopleNum   int
	TotalAmount int
}

func (t *Order) CountChannelRechargeByTime(db *gorm.DB, start, end string) (map[int]ChannelPeopleNumAndTotalAmount, error) {
	var dataList []Order

	err := db.Table(t.GetTableName(time.Now().Format(timeutil.MonthNumberFormat))).
		Select("uid,account,channel").
		Where("created_at >= ? and created_at <= ? and status = 3", start, end).
		Find(&dataList).
		Error

	if err != nil {
		global.Logger["err"].Errorf("Order CountRechargeByTime failed,err:[%v]", err.Error())
		return nil, err
	}

	mp := make(map[int]ChannelPeopleNumAndTotalAmount)

	uidMp := make(map[int]struct{})

	for _, d := range dataList {
		data, ok := mp[d.Channel]

		if !ok {
			data.TotalAmount = 0
			data.PeopleNum = 0
		}

		data.TotalAmount += d.Account

		_, uOk := uidMp[d.Uid]

		if !uOk {
			data.PeopleNum++
		}

		uidMp[d.Uid] = struct{}{}

		mp[d.Channel] = data
	}

	return mp, nil
}

func (t *Order) CountRechargeByUserIdsTime(db *gorm.DB, userIds []int, start, end string) (peopleNum, totalAmount int, err error) {
	var dataList []Order

	err = db.Table(t.GetTableName(time.Now().Format(timeutil.MonthNumberFormat))).
		Select("uid,account").
		Where("uid in (?) and created_at >= ? and created_at <= ? and status = 3", userIds, start, end).
		Find(&dataList).
		Error

	if err != nil {
		global.Logger["err"].Errorf("Order CountRechargeByUserIdsTime failed,err:[%v]", err.Error())
		return
	}

	mp := make(map[int]int)

	for _, d := range dataList {
		v, ok := mp[d.Uid]

		if !ok {
			v = 0
		}

		v += d.Account

		mp[d.Uid] = v
	}

	for _, v := range mp {
		peopleNum++
		totalAmount += v
	}

	return
}
