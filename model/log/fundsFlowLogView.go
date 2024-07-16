package log

import (
	"fmt"
	"gm/model"
	"gm/request"
	"gm/utils/timeutil"
	"gorm.io/gorm"
	"strings"
	"time"
)

type FundsFlowLogView struct {
	Id          int       `gorm:"primaryKey;type:int(11) unsigned AUTO_INCREMENT;comment:id"`       //id
	CreatedAt   time.Time `gorm:"autoCreateTime;type:datetime;not null;default:now();comment:创建时间"` //创建时间
	Uid         int       `gorm:"type:int(11);not null;comment:用户ID"`                               //用户ID
	GameId      int       `gorm:"type:int(11);not null;comment:游戏id"`                               //游戏id
	GameIdX     int       `gorm:"type:int;not null;comment:对局标识"`                                   //对局标识
	RoomId      int       `gorm:"type:int(11);not null;comment:房间id"`                               //房间id
	DeskId      int       `gorm:"type:int(11);not null;comment:桌子id"`                               //桌子id
	Pos         int       `gorm:"type:tinyint(4);not null;comment:座位号"`                             //座位号
	Cash        int64     `gorm:"not null;comment:当前cash余额"`                                        //当前cash余额
	CashNums    int64     `gorm:"not null;comment:cash变动数额"`                                        // cash变动数额
	WinCash     int64     `gorm:"not null;comment:充值金币账户"`                                          // 当前winCash余额
	WinCashNums int64     `gorm:"not null;comment:winCash变动数额"`                                     // winCash变动数额
	Bonus       int64     `gorm:"not null;comment:bonus"`                                           //bonus
	BonusNums   int64     `gorm:"not null;comment:bonus变动数额"`                                       //bonus变动数额
	Tax         int64     `gorm:"not null;comment:税收"`                                              //税收
	Remark      string    `gorm:"type:varchar(64);not null;comment:备注"`                             //备注
}

func (t *FundsFlowLogView) TableName(ym string) string {
	return "room_funds_flow_log_" + ym
}

func (t *FundsFlowLogView) GetList(db *gorm.DB) (list []FundsFlowLogView, err error) {
	err = db.Table(t.TableName(time.Now().Format(timeutil.MonthNumberFormat))).Where(t).Find(&list).Error
	return
}

func (t *FundsFlowLogView) GetPageList(db *gorm.DB, req request.FundsFlowLogView) (total int64, list []FundsFlowLogView, err error) {

	sql := t.JoinByTimeNew(req)

	countQuery := fmt.Sprintf("SELECT COUNT(`id`) as count FROM (%s) as tmp", sql)

	// 执行原始 SQL 查询并扫描结果
	err = db.Raw(countQuery).Scan(&total).Error
	if err != nil {
		return
	}

	// 添加LIMIT和OFFSET进行分页
	offset := (req.Page - 1) * req.Size
	query := fmt.Sprintf("%s order by created_at desc LIMIT %d OFFSET %d", sql, req.Size, offset)

	// 执行原始SQL查询
	err = db.Raw(query).Scan(&list).Error

	return
}

func (t *FundsFlowLogView) JoinByTimeNew(req request.FundsFlowLogView) string {
	var (
		start, end time.Time
	)

	if len(req.CreatedAt) > 0 {
		start = req.Start
		end = req.End
	} else {
		now := time.Now()
		start = timeutil.GetFirstDateOfMonth(now)
		end = timeutil.GetLastTime(timeutil.GetLastDateOfMonth(now))
	}

	roomLogTable := model.GetMonthsBetween("room_funds_flow_log", start, end)
	userLogTable := model.GetMonthsBetween("user_funds_flow_log", start, end)

	// 构建UNION ALL查询的SQL字符串
	var sb strings.Builder
	for i, tableName := range roomLogTable {
		if i > 0 {
			sb.WriteString(" UNION ALL ")
		}

		sb.WriteString(
			fmt.Sprintf(
				"SELECT id,created_at,uid,game_id,game_id_x,room_id,desk_id,pos,cash,cash_nums,win_cash,win_cash_nums,bonus,bonus_nums,tax,remark FROM %s",
				tableName,
			),
		)

		sb.WriteString(
			fmt.Sprintf(
				" UNION ALL SELECT `id`,`created_at`,`uid`,(`type` + 10000 ) AS `game_id`,0 AS `game_id_x`,0 AS `room_id`,0 AS `desk_id`,0 AS `pos`,`left_cash` AS `cash`,`cash` AS `cash_nums`,`left_win_cash` AS `win_cash`,`win_cash` AS `win_cash_nums`,`left_bonus` AS `bonus`,`bonus` AS `bonus_nums`,0 AS `tax`,`remark` AS `remark` FROM `%s`",
				userLogTable[i],
			),
		)
	}

	temTable := sb.String()

	var tmpSb strings.Builder

	tmpSb.WriteString(fmt.Sprintf("select * from (%s) as `t` where `t`.`created_at` >= '%s' and `t`.`created_at` <= '%s'", temTable, start.Format(timeutil.TimeFormat), end.Format(timeutil.TimeFormat)))

	if req.Uid > 0 {
		tmpSb.WriteString(fmt.Sprintf(" and uid = %d", req.Uid))
	}

	if req.GameId > 0 {
		tmpSb.WriteString(fmt.Sprintf(" and game_id = %d", req.GameId))
	}

	return tmpSb.String()
}

func (t *FundsFlowLogView) JoinByTime(req request.FundsFlowLogView) string {
	var (
		start, end time.Time
	)

	if len(req.CreatedAt) > 0 {
		start = req.Start
		end = req.End
	} else {
		now := time.Now()
		start = timeutil.GetFirstDateOfMonth(now)
		end = timeutil.GetLastTime(timeutil.GetLastDateOfMonth(now))
	}

	monthTable := model.GetMonthsBetween("room_funds_flow_log", start, end)

	// 构建UNION ALL查询的SQL字符串
	var sb strings.Builder
	for i, tableName := range monthTable {
		if i > 0 {
			sb.WriteString(" UNION ALL ")
		}
		sb.WriteString(fmt.Sprintf("SELECT id,created_at,uid,game_id,room_id,desk_id,pos,cash,cash_nums,win_cash,win_cash_nums,bonus,bonus_nums,tax,remark FROM %s", tableName))
	}

	sb.WriteString(" UNION ALL SELECT `id`,`created_at`,`uid`,(`type` + 10000 ) AS `game_id`,0 AS `room_id`,0 AS `desk_id`,0 AS `pos`,`left_cash` AS `cash`,`cash` AS `cash_nums`,`left_win_cash` AS `win_cash`,`win_cash` AS `win_cash_nums`,`left_bonus` AS `bonus`,`bonus` AS `bonus_nums`,0 AS `tax`,`remark` AS `remark` FROM `user_funds_flow_log`")

	temTable := sb.String()

	var tmpSb strings.Builder

	tmpSb.WriteString(fmt.Sprintf("select * from (%s) as `t` where `t`.`created_at` >= '%s' and `t`.`created_at` <= '%s'", temTable, start.Format(timeutil.TimeFormat), end.Format(timeutil.TimeFormat)))

	if req.Uid > 0 {
		tmpSb.WriteString(fmt.Sprintf(" and uid = %d", req.Uid))
	}

	if req.GameId > 0 {
		tmpSb.WriteString(fmt.Sprintf(" and game_id = %d", req.GameId))
	}

	return tmpSb.String()
}
