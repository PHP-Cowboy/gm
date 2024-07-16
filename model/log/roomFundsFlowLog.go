package log

import (
	"gm/model"
	"gm/request"
	"gm/utils/timeutil"
	"gorm.io/gorm"
	"time"
)

// 游戏资金流水日志
type RoomFundsFlowLog struct {
	Id           int       `gorm:"primaryKey;type:int(11) unsigned AUTO_INCREMENT;comment:id"`       //id
	CreatedAt    time.Time `gorm:"autoCreateTime;type:datetime;not null;default:now();comment:创建时间"` //创建时间
	Uid          int       `gorm:"type:int(11);not null;comment:用户ID"`                               //用户ID
	GameId       int       `gorm:"type:int(11);not null;comment:游戏id"`                               //游戏id
	GameIdX      int64     `gorm:"not null;comment:对局标识;default:0"`                                  //对局标识
	RoomId       int       `gorm:"type:int(11);not null;comment:房间id"`                               //房间id
	DeskId       int       `gorm:"type:int(11);not null;comment:桌子id"`                               //桌子id
	Pos          int       `gorm:"type:tinyint(4);not null;comment:座位号"`                             //座位号
	Cash         int64     `gorm:"not null;comment:当前cash余额"`                                        //当前cash余额
	CashNums     int64     `gorm:"not null;comment:cash变动数额"`                                        // cash变动数额
	WinCash      int64     `gorm:"not null;comment:充值金币账户"`                                          // 当前winCash余额
	WinCashNums  int64     `gorm:"not null;comment:winCash变动数额"`                                     // winCash变动数额
	Withdraw     int64     `gorm:"not null;comment:赠送储钱罐提现额度"`                                       //赠送储钱罐提现额度
	WithdrawNums int64     `gorm:"not null;comment:变动的储钱罐提现额度"`                                      //变动的储钱罐提现额度
	ExtraGift    int64     `gorm:"not null;comment:额外赠送"`                                            //额外赠送
	Tax          int64     `gorm:"not null;comment:税收"`                                              //税收
	Remark       string    `gorm:"type:varchar(64);not null;comment:备注"`                             //备注
}

func (t *RoomFundsFlowLog) TableName(ym string) string {
	return "room_funds_flow_log_" + ym
}

func (t *RoomFundsFlowLog) GetList(db *gorm.DB) (list []RoomFundsFlowLog, err error) {
	err = db.Table(t.TableName(time.Now().Format(timeutil.MonthNumberFormat))).Where(t).Find(&list).Error
	return
}

func (t *RoomFundsFlowLog) GetPageList(db *gorm.DB, req request.RoomFundsFlowLog) (total int64, list []RoomFundsFlowLog, err error) {
	localDb := db.Table(t.TableName(time.Now().Format(timeutil.MonthNumberFormat))).Where(RoomFundsFlowLog{Uid: req.Uid, GameId: req.GameId})

	if len(req.CreatedAt) > 0 {
		localDb.Where("created_at >= ? and created_at <= ?", req.StartCreatedAt, req.EndCreatedAt)
	}

	err = localDb.Count(&total).Error

	if err != nil {
		return
	}

	err = localDb.Order("id desc").Scopes(model.Paginate(req.Page, req.Size)).Find(&list).Error
	return
}

func (t *RoomFundsFlowLog) GetPlayMpByUserIds(db *gorm.DB, ids []int, dateTime time.Time) (mp map[int]struct{}, err error) {
	var dataList []RoomFundsFlowLog

	err = db.Table(t.TableName(dateTime.Format(timeutil.MonthNumberFormat))).Select("uid").Where("uid in (?)", ids).Find(&dataList).Error

	if err != nil {
		return
	}

	mp = make(map[int]struct{})

	for _, d := range dataList {
		mp[d.Uid] = struct{}{}
	}

	return
}
