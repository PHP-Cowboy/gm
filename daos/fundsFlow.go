package daos

import (
	"errors"
	"gm/global"
	"gm/model/log"
	"gm/model/user"
	"gm/request"
	"gm/response"
	"gm/utils/timeutil"
	"time"
)

func SlotFundFlowLog(req request.SlotFundsFlowLog) (res response.SlotFundsFlowLog, err error) {
	obj := log.SlotFundsFlowLog{}

	db := global.Log

	var rechargeList []log.SlotFundsFlowLog

	rechargeList, err = obj.GetPageList(db, req)

	if err != nil {
		return
	}

	var total int64

	total, err = obj.Count(db)

	if err != nil {
		return
	}

	list := make([]response.SlotFundsFlowLogList, 0, len(rechargeList))

	for _, rl := range rechargeList {
		list = append(list, response.SlotFundsFlowLogList{
			Id:            rl.Id,
			UserName:      rl.UserName,
			Uid:           rl.Uid,
			RoomId:        rl.RoomId,
			DeskId:        rl.DeskId,
			BeforeCash:    rl.BeforeCash,
			BeforeWinCash: rl.BeforeWinCash,
			Cash:          rl.Cash,
			WinCash:       rl.WinCash,
			Nums:          rl.Nums,
			Tax:           rl.Tax,
			Balance:       rl.Balance,
			Remark:        rl.Remark,
		})
	}

	res.Total = total
	res.List = list

	return
}

func TpFundFlowLog(req request.FundsFlowLogView) (res response.RoomFundsFlowLog, err error) {

	db := global.Log
	userDb := global.User

	var (
		dataList    []log.FundsFlowLogView
		total       int64
		userIds     []int
		gameIdxList []int
		userList    []user.User
		userMp      = make(map[int]string)
	)

	obj := log.FundsFlowLogView{}

	if len(req.CreatedAt) > 0 {
		req.Start, req.End, err = timeutil.DateRangeToZeroAndLastTime(req.CreatedAt[0], req.CreatedAt[1])

		if err != nil {
			global.Logger["err"].Errorf("TpFundFlowLog timeutil.DateRangeToZeroAndLastTime failed,err:[%v]", err.Error())
			return
		}

		begin := time.Date(2024, 6, 1, 0, 0, 0, 0, time.Local)

		if req.Start.Before(begin) {
			err = errors.New("start time must after 2024-06-01 00:00:00")
			return
		}
	}

	total, dataList, err = obj.GetPageList(db, req)

	if err != nil {
		global.Logger["err"].Errorf("TpFundFlowLog obj.GetPageList failed,err:[%v]", err.Error())
		return
	}

	for _, l := range dataList {
		userIds = append(userIds, l.Uid)

		if l.GameIdX > 0 {
			gameIdxList = append(gameIdxList, l.GameIdX)
		}
	}

	userObj := new(user.User)

	userList, err = userObj.GetListByUserIds(userDb, userIds)

	if err != nil {
		global.Logger["err"].Errorf("TpFundFlowLog userObj.GetListByUserIds failed,err:[%v]", err.Error())
		return
	}

	for _, ul := range userList {
		userMp[ul.Uid] = ul.UserName
	}

	gameRecordMp := make(map[int]string, len(gameIdxList))
	if len(gameIdxList) > 0 {
		gameRecordObj := new(log.GamerecordLog)

		gameRecordMp, err = gameRecordObj.GetDetailsByGameIdx(db, req, gameIdxList)
		if err != nil {
			global.Logger["err"].Errorf("TpFundFlowLog gameRecordObj.GetDetailsByGameIdx failed,err:[%v]", err.Error())
			return
		}
	}

	list := make([]response.RoomFundsFlowLogList, 0, len(dataList))

	for _, rl := range dataList {
		userName := userMp[rl.Uid]

		details := gameRecordMp[rl.GameIdX]

		tmp := response.RoomFundsFlowLogList{
			Id:            rl.Id,
			UserName:      userName,
			Uid:           rl.Uid,
			CreatedAt:     rl.CreatedAt.Format(timeutil.TimeFormat),
			GameId:        rl.GameId,
			RoomId:        rl.RoomId,
			DeskId:        rl.DeskId,
			Pos:           rl.Pos,
			BeforeCash:    rl.Cash - rl.CashNums,
			Cash:          rl.Cash,
			CashNums:      rl.CashNums,
			BeforeWinCash: rl.WinCash - rl.WinCashNums,
			WinCash:       rl.WinCash,
			WinCashNums:   rl.WinCashNums,
			BeforeBonus:   rl.Bonus - rl.BonusNums,
			Bonus:         rl.Bonus,
			BonusNums:     rl.BonusNums,
			Nums:          rl.CashNums + rl.WinCashNums + rl.BonusNums,
			Tax:           rl.Tax,
			Balance:       rl.Cash + rl.WinCash,
			Remark:        rl.Remark,
			Details:       details,
		}

		list = append(list, tmp)
	}

	res.Total = total
	res.List = list

	return
}
