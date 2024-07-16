package pay

import (
	"gm/model"
	"gm/request"
	"gm/utils/timeutil"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// 实际赠送比：玩家已经赠送的钱的总额度/玩家充值总额度；
// 提交赠送比：（玩家已经赠送的钱总额度+玩家提交赠送但正在审核的额度）/玩家充值总额度
// 赠送总币:玩家成功赠送的总额度
type GiveMoney struct {
	model.Base
	OrderNo          string     `gorm:"type:varchar(32);not null;default:'';unique:order_key;comment:订单号"`     //订单号
	TrdOrderNo       string     `gorm:"type:varchar(32);not null;default:'';unique:trd_order_key;comment:订单号"` //三方单号
	Uid              int        `gorm:"type:bigint(20);not null;default:0;index:uid_key;comment:用户id"`         //用户id
	NickName         string     `gorm:"type:varchar(32);not null;default:'';comment:用户昵称"`                     //用户昵称
	ChannelId        int        `gorm:"type:tinyint(4);not null;default:0;comment:渠道id"`                       //渠道id
	PayChannelId     int        `gorm:"type:tinyint(4);not null;default:0;comment:支付渠道id"`                     //支付渠道id
	Amount           int        `gorm:"type:int(11);not null;default:0;comment:赠送金币"`                          //赠送金币
	AmountTotal      int        `gorm:"type:int(11);not null;default:0;comment:赠送总币"`                          //赠送总币
	Recharge         int        `gorm:"type:int(11);not null;default:0;comment:充值总币"`                          //充值总币
	GiveRate         float64    `gorm:"type:decimal(8,3);not null;default:0.00;comment:实际赠送金币比"`               //实际赠送金币比
	CommitGiveRate   float64    `gorm:"type:decimal(8,3);not null;default:0.00;comment:提交赠送金币比"`               //提交赠送金币比
	TaxRate          float64    `gorm:"type:decimal(8,3);not null;default:0.00;comment:税收(%)"`                 //税收(%)
	Status           int        `gorm:"type:tinyint(4);not null;default:0;comment:订单状态"`                       //订单状态
	PayStatus        int        `gorm:"type:tinyint(4);not null;default:0;comment:付费状态:0:免费用户,1:付费用户,"`        //付费状态
	Auditor          string     `gorm:"type:varchar(32);not null;default:'';comment:审核人"`                      //审核人
	AuditTime        *time.Time `gorm:"type:datetime;comment:审核时间"`                                            //审核时间
	CancelTime       *time.Time `gorm:"type:datetime;comment:取消时间"`                                            //取消时间
	VoidOperator     string     `gorm:"type:varchar(32);not null;default:'';comment:作废操作人"`                    //作废操作人
	ArrivalTime      *time.Time `gorm:"type:datetime;comment:到账时间"`                                            //到账时间
	VoidTime         *time.Time `gorm:"type:datetime;comment:作废时间"`                                            //作废时间
	PayCfgId         int        `gorm:"type:int(11);not null;default:0;comment:支付配置id"`                        //支付配置id
	TpFrequency      int        `gorm:"type:int(11);not null;default:0;comment:TP总局数"`                         //TP总局数
	RmFrequency      int        `gorm:"type:int(11);not null;default:0;comment:RM总局数"`                         //RM总局数
	HundredFrequency int        `gorm:"type:int(11);not null;default:0;comment:百人总局数"`                         //百人总局数
	SlotsFrequency   int        `gorm:"type:int(11);not null;default:0;comment:Slots总局数"`                      //Slots总局数
	GiveMode         int        `gorm:"type:int(11);not null;default:0;comment:赠送金币方式"`                        //赠送金币方式
	Handle           int        `gorm:"type:tinyint(4);not null;default:2;comment:是否处理:1:是2:否"`                //是否处理
	Type             int        `gorm:"type:tinyint(4);not null;default:1;comment:类型:1:normal,2:新手嘉年华"`        //类型:1:normal,2:新手嘉年华
	GiveId           int        `gorm:"type:int(11);not null;default:0;comment:提现档位id，对应表 gave_money_config"`  //提现档位id，对应表 gave_money_config
}

const (
	GiveMoneyStatusZero               = iota
	GiveMoneyStatusToBeReviewed       //待审核
	GiveMoneyStatusUserCancel         //用户取消
	GiveMoneyStatusRepulse            //拒绝
	GiveMoneyStatusInvalid            //作废
	GiveMoneyStatusInPayment          //打款中
	GiveMoneyStatusComplete           //已到账
	GiveMoneyStatusFailed             //代付失败
	GiveMoneyStatusUpstreamRevocation //代付上游撤销
	GiveMoneyStatusUpstreamAbnormal   //异常
)

const (
	GiveMoneyTypeZero           = iota
	GiveMoneyTypeNormal         //normal
	GiveMoneyTypeNoviceCarnival //新手嘉年华提现
)

// 给活动使用，新手嘉年华...
const (
	Handle            = iota
	HandleCompleted   //已处理
	HandleUncompleted //未处理
)

func (t *GiveMoney) TableName() string {
	return "give_money"
}

func (t *GiveMoney) Save(db *gorm.DB, data *GiveMoney) (err error) {
	err = db.Model(t).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}}, // 冲突时的列
				DoUpdates: clause.AssignmentColumns([]string{
					"amount_total",
					"give_rate",
					"commit_give_rate",
					"status",
					"auditor",
					"audit_time",
					"handle",
					"pay_channel_id",
					"pay_cfg_id",
				}), // 要更新的列
			},
		).
		Create(&data).
		Error
	return
}

func (t *GiveMoney) UpdateById(db *gorm.DB, id int, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error

	return
}

func (t *GiveMoney) Create(db *gorm.DB, data *GiveMoney) (err error) {
	err = db.Model(t).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}}, // 冲突时的列
				DoUpdates: clause.AssignmentColumns([]string{
					"trd_order_no",
					"status",
					"arrival_time",
				}), // 要更新的列
			},
		).
		Create(&data).
		Error
	return
}

func (t *GiveMoney) CreateInBatches(db *gorm.DB, list []GiveMoney) (err error) {
	err = db.Model(t).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}}, // 冲突时的列
				DoUpdates: clause.AssignmentColumns([]string{
					"amount_total",
					"give_rate",
					"commit_give_rate",
					"status",
					"auditor",
					"audit_time",
					"handle",
				}), // 要更新的列
			},
		).
		CreateInBatches(&list, model.BatchSize).
		Error
	return
}

func (t *GiveMoney) CreateInBatchesPayConfigId(db *gorm.DB, list []GiveMoney) (err error) {
	err = db.Model(t).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}}, // 冲突时的列
				DoUpdates: clause.AssignmentColumns([]string{
					"pay_channel_id",
					"pay_cfg_id",
					"status",
				}), // 要更新的列
			},
		).
		CreateInBatches(&list, model.BatchSize).
		Error
	return
}

func (t *GiveMoney) GetOneByOrderNo(db *gorm.DB, orderNo string) (data GiveMoney, err error) {
	err = db.Model(t).Where("order_no = ?", orderNo).First(&data).Error

	return
}

func (t *GiveMoney) GetListByIds(db *gorm.DB, ids []int) (list []GiveMoney, err error) {
	err = db.Model(t).Where("id in (?)", ids).Find(&list).Error

	return
}

func (t *GiveMoney) GetCompleteListByIds(db *gorm.DB, ids []int) (list []GiveMoney, err error) {
	err = db.Model(t).Where("id in (?) and status = ?", ids, GiveMoneyStatusComplete).Find(&list).Error

	return
}

func (t *GiveMoney) GetPageList(db *gorm.DB, req request.GiveList) (total int64, list []GiveMoney, err error) {
	localDb := db.Table(t.TableName()).Where(GiveMoney{Uid: req.Uid, ChannelId: req.ChannelId, PayChannelId: req.PayChannelId})

	if req.GiveRateMin > 0 {
		localDb.Where("give_rate >= ?", req.GiveRateMin)
	}

	if req.GiveRateMax > 0 {
		localDb.Where("give_rate <= ?", req.GiveRateMax)
	}

	if req.CommitGiveRateMin > 0 {
		localDb.Where("commit_give_rate <= ?", req.CommitGiveRateMin)
	}

	if req.CommitGiveRateMax > 0 {
		localDb.Where("commit_give_rate >= ?", req.CommitGiveRateMax)
	}

	if req.PayStatus != nil {
		localDb.Where("pay_status = ?", *req.PayStatus)
	}

	if req.Status != nil {
		localDb.Where("status = ?", *req.Status)
	}

	if len(req.CreatedAt) > 0 {
		localDb.Where("created_at >= ? and created_at <= ?", req.StartCreatedAt, req.EndCreatedAt)
	}

	localDb.Where("channel_id in (?)", req.ChannelIds)

	err = localDb.Count(&total).Error
	if err != nil {
		return
	}

	err = localDb.Order("id desc").Scopes(model.Paginate(req.Page, req.Size)).Find(&list).Error

	return
}

func (t *GiveMoney) BatchHandle(db *gorm.DB, req request.CheckIds, mp map[string]interface{}) (err error) {
	err = db.Table(t.TableName()).Where("id in (?)", req.Ids).Updates(mp).Error
	return
}

func (t *GiveMoney) UpdateStatusByIds(db *gorm.DB, req request.CheckIds, status int) (err error) {
	err = db.Table(t.TableName()).Where("id in (?)", req.Ids).Update("status", status).Error
	return
}

func (t *GiveMoney) GetListByHandle(db *gorm.DB, handle int) (dataList []GiveMoney, err error) {
	err = db.Model(t).Where("handle = ?", handle).Find(&dataList).Error
	return
}

func (t *GiveMoney) GetAuditorList(db *gorm.DB) (dataList []GiveMoney, err error) {
	err = db.Model(t).Where("auditor != ''").Find(&dataList).Error
	return
}

func (t *GiveMoney) GetYesterdayTodayListByUid(db *gorm.DB, uid int) (mp map[string]int, err error) {
	var (
		ymds     []string
		list     []GiveMoney
		todayMax int
	)

	mp = make(map[string]int)
	mp["today"] = 0
	mp["yesterday"] = 0

	now := time.Now()

	today := now.Format(timeutil.DateNumberFormat)

	yesterday := now.AddDate(0, 0, -1).Format(timeutil.DateNumberFormat)

	ymds = append(ymds, today, yesterday)

	//根据到账时间统计
	err = db.Model(t).
		Where("uid = ? ", uid).
		Where("arrival_time >= ?", timeutil.GetZeroTime(now.AddDate(0, 0, -1))).
		Where("arrival_time <= ?", timeutil.GetLastTime(now)).
		Find(&list).Error

	for _, g := range list {
		switch g.CreatedAt.Format(timeutil.DateNumberFormat) {
		case yesterday:
			mp["yesterday"] += g.Amount / 100

		case today:
			if g.Amount > todayMax {
				todayMax = g.Amount / 100
			}
			mp["today"] += g.Amount / 100
		}
	}

	mp["today_withdraw_max"] = todayMax

	return
}

type PeoplesAndAmount struct {
	PeopleNum int
	SumAmount int
}

// key: channel,value sum(amount)
func (t *GiveMoney) MapChannelAmountByArrivalTime(db *gorm.DB, dateTime time.Time) (mp map[int]PeoplesAndAmount, err error) {
	var dataList []GiveMoney

	err = db.Table(t.TableName()).
		Select("uid,channel_id,amount").
		Where(
			"arrival_time >= ? and arrival_time <= ?",
			timeutil.GetZeroTime(dateTime).Format(timeutil.TimeFormat),
			timeutil.GetLastTime(dateTime).Format(timeutil.TimeFormat),
		).Find(&dataList).
		Error

	if err != nil {
		return
	}

	mp = make(map[int]PeoplesAndAmount)

	for _, d := range dataList {
		v, ok := mp[d.ChannelId]

		if !ok {
			v = PeoplesAndAmount{}
		}

		v.SumAmount += d.Amount
		v.PeopleNum++

		mp[d.ChannelId] = v
	}

	return
}

func (t *GiveMoney) GetListByUserIdsAndArrivalTime(db *gorm.DB, ids []int, dateTime time.Time) (mp map[int]int, err error) {
	var dataList []GiveMoney

	err = db.Table(t.TableName()).
		Select("uid,amount").
		Where("uid in (?)", ids).
		Where(
			"arrival_time >= ? and arrival_time <= ?",
			timeutil.GetZeroTime(dateTime).Format(timeutil.TimeFormat),
			timeutil.GetLastTime(dateTime).Format(timeutil.TimeFormat),
		).Find(&dataList).
		Error

	if err != nil {
		return
	}

	mp = make(map[int]int)

	for _, d := range dataList {
		v, ok := mp[d.Uid]

		if !ok {
			v = 0
		}

		v += d.Amount

		mp[d.Uid] = v
	}

	return
}

func (t *GiveMoney) GetListByCreatedAt(db *gorm.DB, dateTime time.Time) (mp map[int]int, err error) {
	var dataList []GiveMoney

	err = db.Table(t.TableName()).
		Select("uid,amount").
		Where("created_at >= ? and created_at <= ?", timeutil.GetZeroTime(dateTime), timeutil.GetLastTime(dateTime)).
		Find(&dataList).
		Error

	if err != nil {
		return
	}

	mp = make(map[int]int)

	for _, d := range dataList {
		v, ok := mp[d.Uid]

		if !ok {
			v = 0
		}

		v += d.Amount

		mp[d.Uid] = v
	}

	return
}

func (t *GiveMoney) SumAmountByCreatedAt(db *gorm.DB, start, end string) (totalAmount int, err error) {
	var dataList []GiveMoney

	err = db.Table(t.TableName()).
		Select("amount").
		Where("arrival_time >= ? and arrival_time <= ?", start, end).
		Find(&dataList).
		Error

	if err != nil {
		return
	}

	for _, d := range dataList {
		totalAmount += d.Amount
	}

	return
}
func (t *GiveMoney) SumChannelAmountByCreatedAt(db *gorm.DB, start, end string) (map[int]int, error) {
	var dataList []GiveMoney

	err := db.Table(t.TableName()).
		Select("amount,channel_id").
		Where("arrival_time >= ? and arrival_time <= ?", start, end).
		Find(&dataList).
		Error

	if err != nil {
		return nil, err
	}

	mp := make(map[int]int)

	for _, d := range dataList {
		mp[d.ChannelId] += d.Amount
	}

	return mp, nil
}

// 赠送统计，根据到账时间
func (t *GiveMoney) SumAmountByUserIdsArrivalTime(db *gorm.DB, userIds []int, start, end string) (totalAmount int, err error) {
	var dataList []GiveMoney

	err = db.Table(t.TableName()).
		Select("amount").
		Where("uid in (?) and arrival_time >= ? and arrival_time <= ? ", userIds, start, end).
		Find(&dataList).
		Error

	if err != nil {
		return
	}

	for _, d := range dataList {
		totalAmount += d.Amount
	}

	return
}

// 根据用户id和日期查询用户申请赠送额度
func (t *GiveMoney) SumAmountByUidAndCreatedAt(db *gorm.DB, uid int, dateTime time.Time) (sumAmount int, err error) {
	var dataList []GiveMoney

	err = db.Table(t.TableName()).
		Select("amount").
		Where("uid = ? and created_at >= ? and created_at <= ?", uid, timeutil.GetZeroTime(dateTime), timeutil.GetLastTime(dateTime)).
		Find(&dataList).
		Error

	if err != nil {
		return
	}

	for _, d := range dataList {
		sumAmount += d.Amount
	}

	return
}
