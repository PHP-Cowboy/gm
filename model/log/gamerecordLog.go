package log

import (
	"fmt"
	"gm/model"
	"gm/request"
	"gm/utils/timeutil"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type GamerecordLog struct {
	Id      int    `gorm:"primaryKey;type:int(11) unsigned AUTO_INCREMENT;comment:id"`
	GameIdx int    `gorm:"type:int;not null;default:0;comment:对局标识"`
	Details string `gorm:"type:varchar;NOT NULL;comment:详情"`
}

func (t *GamerecordLog) JoinByTime(req request.FundsFlowLogView, gameldxList []int) string {
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

	tableList := model.GetMonthsBetween("GameRecord_log", start, end)

	gameIdxs := ""

	for _, id := range gameldxList {
		gameIdxs += strconv.Itoa(id) + ","
	}

	gameIdxs = strings.TrimRight(gameIdxs, ",")

	// 构建UNION ALL查询的SQL字符串
	var sb strings.Builder
	for i, tableName := range tableList {
		if i > 0 {
			sb.WriteString(" UNION ALL ")
		}

		sb.WriteString("SELECT game_idx,details FROM ")

		sb.WriteString(tableName)

		if gameIdxs != "" {
			sb.WriteString(fmt.Sprintf(" where game_idx in (%v)", gameIdxs))
		}
	}

	sql := fmt.Sprintf("select * from (%s) as `t`", sb.String())

	return sql
}

func (t *GamerecordLog) GetDetailsByGameIdx(db *gorm.DB, req request.FundsFlowLogView, gameIdx []int) (mp map[int]string, err error) {
	mp = make(map[int]string, len(gameIdx))

	var dataList []GamerecordLog

	sql := t.JoinByTime(req, gameIdx)

	err = db.Raw(sql).Scan(&dataList).Error

	if err != nil {
		return
	}

	for _, log := range dataList {
		mp[log.GameIdx] = log.Details
	}

	return
}
