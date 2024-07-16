package daos

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gm/common/constant"
	"gm/daos/cache"
	"gm/daos/rds"
	"gm/global"
	"gm/model/game"
	"gm/model/gm"
	"gm/model/log"
	"gm/model/pay"
	"gm/model/user"
	"gm/request"
	"gm/response"
	"gm/utils/ecode"
	"gm/utils/errUtil"
	"gm/utils/slice"
	"gm/utils/timeutil"
	"gorm.io/gorm"
	"net/url"
	"strconv"
	"time"
	"za.game/lib/account"
	accUser "za.game/lib/account/user"
	"za.game/lib/consts"
	"za.game/lib/dbconn"
)

func GetLoginLogList(req request.GetLoginLogList) (res response.LoginLogRsp, err error) {
	db := global.Log

	obj := new(log.UserLoginLog)

	var (
		total     int64
		pageList  []log.UserLoginLog
		channelMp = make(map[int]accUser.UserChannel)
	)

	if len(req.CreatedAt) > 0 {
		req.StartTime, req.EndTime, err = timeutil.DateRangeToZeroAndLastTimeFormat(req.CreatedAt[0], req.CreatedAt[1], timeutil.TimeFormat)

		if err != nil {
			global.Logger["err"].Errorf("GetLoginLogList timeutil.DateRangeToZeroAndLastTimeFormat failed,err:[%v]", err.Error())
			return
		}
	}

	total, pageList, err = obj.GetPageList(db, req)

	if err != nil {
		global.Logger["err"].Errorf("GetLoginLogList GetPageList failed,err:[%v]", err.Error())
		return
	}

	list := make([]response.LoginLog, 0, len(pageList))

	channelMp, err = rds.GetAllChannelMap()

	if err != nil {
		global.Logger["err"].Infof(fmt.Sprintf("cache.GetAllChannel failed, err:%s", err.Error()))
		return
	}

	for _, p := range pageList {

		userCh, cNameOk := channelMp[p.ChannelId]

		if !cNameOk {
			userCh = accUser.UserChannel{}
		}

		list = append(list, response.LoginLog{
			Id:                 p.Id,
			Uid:                p.Uid,
			Nickname:           p.Nickname,
			ChannelId:          p.ChannelId,
			Channel:            userCh.ChannelName,
			Assets:             p.Assets,
			ReferralCommission: p.ReferralCommission,
			Ip:                 p.Ip,
			Device:             p.Device,
			Version:            p.Version,
			LoginMode:          p.LoginMode,
			LoginTime:          p.CreatedAt.Format(timeutil.TimeFormat),
			RegTime:            p.RegTime.Format(timeutil.TimeFormat),
		})
	}

	res.List = list
	res.Total = total

	return
}

func GetGameUserList(req request.GetGameUserList) (res response.GameUserRsp, err error) {
	db := global.User

	obj := new(user.User)

	var (
		total    int64
		pageList []user.User
	)

	if len(req.CreatedAt) > 0 {
		req.StartCreatedAt, req.EndCreatedAt, err = timeutil.DateRangeToZeroAndLastTimeFormat(req.CreatedAt[0], req.CreatedAt[1], timeutil.TimeFormat)

		if err != nil {
			global.Logger["err"].Errorf("GetGameUserList timeutil.DateRangeToZeroAndLastTimeFormat req.CreatedAt failed,err:[%v]", err.Error())
			return
		}
	}

	if len(req.UpdatedAt) > 0 {
		req.StartUpdatedAt, req.EndUpdatedAt, err = timeutil.DateRangeToZeroAndLastTimeFormat(req.UpdatedAt[0], req.UpdatedAt[1], timeutil.TimeFormat)
		if err != nil {
			global.Logger["err"].Errorf("GetGameUserList timeutil.DateRangeToZeroAndLastTimeFormat req.UpdatedAt failed,err:[%v]", err.Error())
			return
		}
	}

	userInfoObj := new(user.UserInfo)

	var userInfos []user.UserInfo

	userInfos, err = userInfoObj.GetListUnionAll(db, req)

	if err != nil {
		global.Logger["err"].Errorf("GetGameUserList userInfoObj.GetListUnionAll failed,err:[%v]", err.Error())
		return
	}

	var uids []int

	for _, ui := range userInfos {
		uids = append(uids, ui.Uid)
	}

	req.UserIds = uids

	total, pageList, err = obj.GetPageList(db, req)

	if err != nil {
		global.Logger["err"].Errorf("GetGameUserList obj.GetPageList failed,err:[%v]", err.Error())
		return
	}

	uMp := make(map[int][]int, 0)
	udMp := make(map[int][]int, 0)

	//用户id根据分表规则，拆分进map中
	for _, pl := range pageList {
		k := pl.Uid % 5
		ids, ok := uMp[k]

		if !ok {
			ids = make([]int, 0)
		}
		ids = append(ids, pl.Uid)

		uMp[k] = ids

		udK := pl.Uid % 10

		udIds, udOk := udMp[udK]

		if !udOk {
			udIds = make([]int, 0)
		}
		udIds = append(udIds, pl.Uid)

		udMp[udK] = udIds
	}

	userInfos, err = userInfoObj.GetListByUserIds(db, uMp)

	if err != nil {
		global.Logger["err"].Errorf("GetGameUserList userInfoObj.GetListByUserIds failed,err:[%v]", err.Error())
		return
	}

	uInfoMp := make(map[int]user.UserInfo, 0)
	for _, u := range userInfos {
		uInfoMp[u.Uid] = u
	}

	list := make([]response.GameUser, 0, len(pageList))

	for _, p := range pageList {
		uinfo, _ := uInfoMp[p.Uid]

		payStatus := 0

		if uinfo.Recharge > 0 {
			payStatus = 1
		}

		list = append(list, response.GameUser{
			Id:         p.Id,
			Uid:        p.Uid,
			IsGuest:    p.IsGuest,
			IsSend:     p.IsSend,
			Device:     p.Device,
			UserName:   p.UserName,
			Icon:       p.Icon,
			Phone:      p.Phone,
			Email:      p.Email,
			ChannelId:  p.ChannelId,
			TpNew:      p.TpNew,
			RegIp:      p.RegIp,
			RegVersion: p.RegVersion,
			Asset:      uinfo.Cash + uinfo.WinCash,
			Recharge:   uinfo.Recharge,
			PayStatus:  payStatus,
			GiftCash:   uinfo.WithdrawedMoney,
			CreatedAt:  p.CreatedAt.Format(timeutil.TimeFormat),
			UpdatedAt:  p.UpdatedAt.Format(timeutil.TimeFormat),
		})
	}

	res.List = list
	res.Total = total

	return
}

func GiveList(req request.GiveList) (res response.GameGiveRsp, err error) {
	db := global.Pay

	obj := new(pay.GiveMoney)

	var (
		givePageList []pay.GiveMoney
		total        int64
	)

	if len(req.CreatedAt) > 0 {
		req.StartCreatedAt, req.EndCreatedAt, err = timeutil.DateRangeToZeroAndLastTimeFormat(req.CreatedAt[0], req.CreatedAt[1], timeutil.TimeFormat)
		if err != nil {
			global.Logger["err"].Errorf("GiveList timeutil.DateRangeToZeroAndLastTimeFormat req.UpdatedAt failed,err:[%v]", err.Error())
			return
		}
	}

	total, givePageList, err = obj.GetPageList(db, req)

	if err != nil {
		global.Logger["err"].Errorf("GiveList GetPageList failed,err:[%v]", err.Error())
		return
	}

	var payCfgMap = make(map[int]pay.PayConfig)

	payCfgMap, err = cache.GetAllPayCfgMap()
	if err != nil {
		global.Logger["err"].Errorf("GiveList GetAllPassageMap failed,err:[%v]", err.Error())
		return
	}

	var channelMp = make(map[int]accUser.UserChannel)

	channelMp, err = rds.GetAllChannelMap()

	if err != nil {
		global.Logger["err"].Errorf("GiveList GetAllChannelMap failed,err:[%v]", err.Error())
		return
	}

	list := make([]response.GiveMoney, 0, len(givePageList))

	for _, l := range givePageList {
		payCfg, _ := payCfgMap[l.PayCfgId]

		channel, _ := channelMp[l.ChannelId]

		list = append(list, response.GiveMoney{
			Id:               l.Id,
			Uid:              l.Uid,
			UserName:         l.NickName,
			PayChannelName:   payCfg.Name,
			Amount:           l.Amount,
			AmountTotal:      l.AmountTotal,
			Recharge:         l.Recharge,
			GiveRate:         l.GiveRate,
			CommitGiveRate:   l.CommitGiveRate,
			Status:           l.Status,
			Auditor:          l.Auditor,
			AuditTime:        timeutil.FormatToDateTime(l.AuditTime),
			OrderNo:          l.OrderNo,
			TrdOrderNo:       l.TrdOrderNo,
			CreatedAt:        timeutil.FormatToDateTime(&l.CreatedAt),
			ArrivalTime:      timeutil.FormatToDateTime(l.ArrivalTime),
			CancelTime:       timeutil.FormatToDateTime(l.CancelTime),
			VoidTime:         timeutil.FormatToDateTime(l.VoidTime),
			VoidOperator:     l.VoidOperator,
			ChannelId:        l.ChannelId,
			ChannelCode:      channel.Code,
			TpFrequency:      l.TpFrequency,
			RmFrequency:      l.RmFrequency,
			HundredFrequency: l.HundredFrequency,
			SlotsFrequency:   l.SlotsFrequency,
			GiveMode:         l.GiveMode,
			Ifsc:             l.ChannelId,
			PayStatus:        l.PayStatus,
			TaxRate:          l.TaxRate,
		})
	}

	res.List = list
	res.Total = total

	return
}

// 批量通过
func BatchPass(req request.CheckIds, name any) (err error) {
	payDb := global.Pay

	obj := new(pay.GiveMoney)
	bankObj := new(pay.BankInfo)

	var (
		giveMoneyList []pay.GiveMoney
		bankList      []pay.BankInfo
		channelMp     = make(map[int][]PayCfgAndRate)
	)

	channelMp, err = GetChannelPassageListMp()

	if err != nil {
		global.Logger["err"].Infof("BatchPass GetChannelPassageListMp failed, err:%s", err.Error())
		return
	}

	giveMoneyList, err = obj.GetListByIds(payDb, req.Ids)

	if err != nil {
		global.Logger["err"].Infof("BatchPass GetListByIds failed, err:%s", err.Error())
		return
	}

	var userIds = make([]int, 0, len(giveMoneyList))

	for _, gl := range giveMoneyList {
		//不是待审核的
		if gl.Status != pay.GiveMoneyStatusToBeReviewed {
			return ecode.DataStatusError
		}

		if gl.Recharge <= 0 {
			return ecode.UserNotRecharge
		}

		userIds = append(userIds, gl.Uid)
	}

	userIds = slice.UniqueSlice(userIds)

	bankList, err = bankObj.GetListByUserIds(payDb, userIds)

	if err != nil {
		global.Logger["err"].Infof("BatchPass GetListByUserIds failed, err:%s", err.Error())
		return
	}

	bankMp := make(map[int]pay.BankInfo, 0)

	for _, b := range bankList {
		bankMp[int(b.Uid)] = b
	}

	now := time.Now()

	mp := make(map[string]interface{})

	mp["auditor"] = name.(string)
	mp["audit_time"] = timeutil.FormatToDateTime(&now)

	//申请打款
	for _, g := range giveMoneyList {
		bank, bankOk := bankMp[g.Uid]

		if !bankOk {
			global.Logger["err"].Infof("BatchPass uid: %v,bank info query err:%s", g.Uid, err.Error())
			continue
		}

		channelPayList, channelMpOk := channelMp[g.ChannelId]

		if !channelMpOk {
			global.Logger["err"].Infof("BatchPass !channelMpOk ChannelId [%v]", g.ChannelId)
			continue
		}

		var payRes response.PaymentResponse

		payRes, err = Payment(g, bank, channelPayList)

		mp["pay_cfg_id"] = payRes.CfgId

		if err != nil {
			global.Logger["err"].Infof("BatchPass Payment failed, id: %v , cfgId:%v err:[%s]", g.Id, payRes.CfgId, err.Error())
			mp["status"] = pay.GiveMoneyStatusUpstreamAbnormal //更新为异常，需人工审核
		} else {
			mp["status"] = pay.GiveMoneyStatusInPayment
		}

		err = obj.UpdateById(payDb, g.Id, mp)

		if err != nil {
			global.Logger["err"].Infof("BatchPass obj.UpdateById failed, err:[%v], id: %v ,status: %v", err.Error(), g.Id, mp["status"])
			return err
		}

	}

	return
}

// 批量拒绝
func BatchRepulse(req request.CheckIds, name any) (err error) {
	payDb := global.Pay
	userDb := global.User

	obj := new(pay.GiveMoney)

	var giveMoneyList []pay.GiveMoney

	giveMoneyList, err = obj.GetListByIds(payDb, req.Ids)

	if err != nil {
		return
	}

	//附件mp
	attachMp := make(map[int]response.Refund, 0)
	userIds := make([]int, 0, len(giveMoneyList))
	//新手嘉年华单号
	orderNoList := make([]string, 0)

	type NumsAmount struct {
		Num    int
		Amount int
	}

	amountMp := make(map[int]NumsAmount)
	gUserList := make([]pay.GaveMoneyUser, 0, len(req.Ids))

	for _, gl := range giveMoneyList {
		//不是待审核的
		if gl.Status != pay.GiveMoneyStatusToBeReviewed {
			return ecode.DataStatusError
		}

		if gl.Type == pay.GiveMoneyTypeNormal {
			attachMp[gl.Id] = response.Refund{
				Uid:       gl.Uid,
				HasAttach: 1,
				Amount:    gl.Amount,
			}

			numsAmount, amountMpOk := amountMp[gl.Uid]

			if !amountMpOk {
				numsAmount = NumsAmount{}
			}

			numsAmount.Num++
			numsAmount.Amount += gl.Amount

			amountMp[gl.Uid] = numsAmount

			userIds = append(userIds, gl.Uid)

			gUserList = append(gUserList, pay.GaveMoneyUser{
				Uid:    gl.Uid,
				GaveId: gl.GiveId,
			})

		} else if gl.Type == pay.GiveMoneyTypeNoviceCarnival {
			//只发邮件，不退钱，嘉年华那边自己处理退钱
			attachMp[gl.Id] = response.Refund{
				Uid:       gl.Uid,
				HasAttach: 0,
				Amount:    0,
			}

			orderNoList = append(orderNoList, gl.OrderNo)
		}

	}

	err = obj.BatchHandle(payDb, req, map[string]interface{}{"status": pay.GiveMoneyStatusRepulse, "auditor": name, "audit_time": time.Now()})

	if err != nil {
		return
	}

	//发邮件  附件退钱
	err = SendRefundEmail("en", attachMp)
	if err != nil {
		return
	}

	userIds = slice.UniqueSlice(userIds)

	if len(userIds) > 0 {
		//更新提现次数、提现额度
		vipObj := new(user.VipUserInfo)
		vipList := make([]user.VipUserInfo, 0, len(userIds))

		vipList, err = vipObj.GetListByUserIds(userDb, userIds)
		if err != nil {
			return
		}

		vipKeys := make([]string, 0, len(userIds))

		for i, v := range vipList {
			val, ok := amountMp[v.Uid]

			if !ok {
				continue
			}

			vipList[i].WithdrawMoney -= val.Amount
			vipList[i].WithdrawNum -= val.Num

			vipKeys = append(vipKeys, consts.RedisVipUser+fmt.Sprintf("%v", v.Uid))
		}

		err = vipObj.CreateInBatches(
			userDb,
			vipList,
		)

		if err != nil {
			global.Logger["err"].Infof("BatchRepulse vipObj.CreateInBatches 失败:[%s]", err.Error())
			return
		}

		err = rds.DelRedisCacheByKey(vipKeys...)
		if err != nil {
			//缓存删除失败，只记录日志，不处理
			global.Logger["err"].Infof("BatchRepulse DelRedisCacheByKey 失败:[%s]，keys:[%v]", err.Error(), vipKeys)
		}
	}

	if len(orderNoList) > 0 {
		for _, orderNo := range orderNoList {
			NoviceCarnivalWithdrawRefuse(orderNo)
		}
	}

	//提现审核中值的扣减在邮件中实现的

	gUserObj := new(pay.GaveMoneyUser)

	for _, mu := range gUserList {
		err = gUserObj.UpdateNumByUidGiveId(payDb, mu.Uid, mu.GaveId, -1)

		if err != nil {
			global.Logger["err"].Errorf("BatchRepulse gUserObj.UpdateNumByUidGiveId failed, err:[%s]", err.Error())
			return
		}
	}

	return
}

// 批量主动取消
func BatchCancel(req request.CheckIds) (err error) {
	db := global.Pay

	obj := new(pay.GiveMoney)

	var giveMoneyList []pay.GiveMoney

	giveMoneyList, err = obj.GetListByIds(db, req.Ids)

	if err != nil {
		return
	}

	for _, gl := range giveMoneyList {
		//不是待审核的
		if gl.Status != pay.GiveMoneyStatusToBeReviewed {
			return ecode.DataStatusError
		}
	}

	err = obj.UpdateStatusByIds(db, req, pay.GiveMoneyStatusUserCancel)
	return
}

// 批量作废
func BatchInvalid(req request.CheckIds, name any) (err error) {
	db := global.Pay

	obj := new(pay.GiveMoney)

	var giveMoneyList []pay.GiveMoney

	giveMoneyList, err = obj.GetListByIds(db, req.Ids)

	if err != nil {
		return
	}

	for _, gl := range giveMoneyList {
		//不是待审核的
		if gl.Status != pay.GiveMoneyStatusToBeReviewed {
			return ecode.DataStatusError
		}
	}

	err = obj.BatchHandle(db, req, map[string]interface{}{"status": pay.GiveMoneyStatusInvalid, "auditor": name, "audit_time": time.Now()})

	return
}

//设置自动审批规则

//设置自动审批开关

// 用户信息
func GameUserInfo(req request.UserId, channelIds []int) (mp map[string]interface{}, err error) {
	userDb := global.User
	gameDb := global.Game
	payDb := global.Pay
	gmDb := global.DB
	logDb := global.Log

	obj := new(user.User)

	var users user.User

	users, err = obj.GetOneByUid(userDb, req.Uid)

	if err != nil {
		global.Logger["err"].Errorf("GameUserInfo obj.GetOneByUid failed,err:[%v]", err.Error())
		return
	}

	channelIdMp := slice.SliceToMap(channelIds)

	_, ok := channelIdMp[users.ChannelId]

	if !ok {
		err = errors.New("channel err")
		global.Logger["err"].Errorf("当前后台管理员无权查看的渠道用户,admin urser role channels:[%v],user channel:[%v]", channelIds, users.ChannelId)
		return
	}

	var userInfo user.UserInfo

	infoObj := new(user.UserInfo)

	userInfo, err = infoObj.GetOneByUid(userDb, req.Uid)

	if err != nil {
		global.Logger["err"].Errorf("GameUserInfo infoObj.GetOneByUid failed,err:[%v]", err.Error())
		return
	}

	chanObj := new(user.UserChannel)
	chanInfo := user.UserChannel{}

	chanInfo, err = chanObj.GetFirstById(userDb, users.ChannelId)
	if err != nil {
		global.Logger["err"].Errorf("GameUserInfo chanObj.GetFirstById failed,err:[%v]", err.Error())
		return
	}

	bankObj := new(pay.BankInfo)
	bankInfo := pay.BankInfo{}

	bankInfo, err = bankObj.GetFirstByUid(payDb, int64(req.Uid))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.Logger["err"].Errorf("GameUserInfo bankObj.GetFirstByUid failed,err:[%v]", err.Error())
		return
	}

	loginLogObj := new(log.UserLoginLog)
	loginLogInfo := log.UserLoginLog{}

	loginLogInfo, err = loginLogObj.GetLastLoginInfoByUid(logDb, req.Uid)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.Logger["err"].Errorf("GameUserInfo loginLogObj.GetLastLoginInfoByUid failed,err:[%v]", err.Error())
		return
	}

	mp = gin.H{
		"id":          users.Id,
		"uid":         users.Uid,
		"is_guest":    users.IsGuest,
		"is_send":     users.IsSend,
		"device":      users.Device,
		"user_name":   users.UserName,
		"icon":        users.Icon,
		"phone":       users.Phone,
		"email":       bankInfo.Email,
		"created_at":  timeutil.FormatToDateTime(&users.CreatedAt),
		"updated_at":  timeutil.FormatToDateTime(&users.UpdatedAt),
		"channel_id":  users.ChannelId,
		"reg_ip":      users.RegIp,
		"reg_version": users.RegVersion,
		"win_cash":    userInfo.WinCash,
		"cash":        userInfo.Cash,
		"bonus":       userInfo.Bonus,
		"recharge":    userInfo.Recharge,
		"withdraw":    userInfo.WithdrawedMoney,
		"channel":     chanInfo.ChannelName,
		"gaid":        users.Gpcadid,
		"login_ip":    loginLogInfo.Ip,
		"name":        bankInfo.Name,
		"mobile":      bankInfo.Phone,
		"id_card":     bankInfo.Vpa,
		"account":     bankInfo.AccountNo,
	}

	var (
		giveObj  pay.GiveMoney
		orderObj pay.Order
		giveMp   = make(map[string]int)
		payMp    = make(map[string]int)
	)

	//充值
	payMp, err = orderObj.GetYesterdayTodayListByUid(payDb, req.Uid)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.Logger["err"].Errorf("GameUserInfo orderObj.GetYesterdayTodayListByUid failed,err:[%v]", err.Error())
		return
	}

	mp["yesterday_recharge"] = payMp["yesterday"]
	mp["today_recharge"] = payMp["today"]
	mp["today_recharge_max"] = payMp["today_recharge_max"]

	//提现
	giveMp, err = giveObj.GetYesterdayTodayListByUid(payDb, req.Uid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.Logger["err"].Errorf("GameUserInfo giveObj.GetYesterdayTodayListByUid failed,err:[%v]", err.Error())
		return
	}

	mp["yesterday_withdraw"] = giveMp["yesterday"]
	mp["today_withdraw"] = giveMp["today"]
	mp["today_withdraw_max"] = giveMp["today_withdraw_max"]

	situObj := new(game.GameSituation)

	var gameSituationList []game.GameSituation

	gameSituationList, err = situObj.GetList(gameDb, req.Uid)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.Logger["err"].Errorf("GameUserInfo situObj.GetList failed,err:[%v]", err.Error())
		return
	}

	list := make([]response.GameSituation, 0, len(gameSituationList))

	for _, s := range gameSituationList {
		list = append(list, response.GameSituation{
			Id:        s.Id,
			GameName:  s.GameName,
			RoomType:  s.RoomType,
			Total:     s.Total,
			WinNum:    s.WinNum,
			WinMoney:  s.WinMoney,
			LossMoney: s.LossMoney,
		})
	}

	var (
		gmObj                                      = gm.Banned{}
		gmBannedList                               = make([]gm.Banned, 0)
		regIpStatus, deviceStatus, lastLoginStatus int
	)

	gmBannedList, err = gmObj.GetListByUid(gmDb, req.Uid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.Logger["err"].Errorf("GameUserInfo gmObj.GetListByUid failed,err:[%v]", err.Error())
		return
	}

	for _, banned := range gmBannedList {
		switch banned.Cate {
		case 1:
			regIpStatus = 1

		case 2:
			deviceStatus = 1

		case 3:
			if banned.Info == loginLogInfo.Ip {
				lastLoginStatus = 1
			}
		}
	}

	mp["lastLoginStatus"] = lastLoginStatus
	mp["deviceStatus"] = deviceStatus
	mp["regIpStatus"] = regIpStatus

	mp["list"] = list

	return
}

func ChangeRecharge(req request.ChangeRecharge) (err error) {
	userDb := global.User

	userInfoObj := new(user.UserInfo)

	err = userInfoObj.UpdateByUid(userDb, req.Uid, map[string]interface{}{"recharge": *req.EditRecharge})
	return
}

func WithdrawInfoRecord(req request.WithdrawInfoRecord) (res response.WithdrawInfoRsp, err error) {
	payDb := global.Pay

	bankObj := new(pay.BankInfo)

	var (
		bankInfoList []pay.BankInfo
		total        int64
	)

	total, bankInfoList, err = bankObj.GetPageList(payDb, &pay.BankInfo{Uid: req.Uid}, req.Page, req.Size, req.ChannelIds)

	if err != nil {
		return
	}

	list := make([]response.WithdrawInfo, 0, len(bankInfoList))
	for _, b := range bankInfoList {
		list = append(list, response.WithdrawInfo{
			Id:        b.ID,
			Uid:       b.Uid,
			BankCode:  b.BankCode,
			BankName:  b.BankName,
			AccountNo: b.AccountNo,
			Ifsc:      b.Ifsc,
			Name:      b.Name,
			Email:     b.Email,
			Phone:     b.Phone,
			Address:   b.Address,
			Vpa:       b.Vpa,
			Remark:    b.Remark,
			UpdatedAt: b.UpdatedAt.Format(timeutil.TimeFormat),
		})
	}

	res.List = list
	res.Total = total

	return
}

// 修改用户金币
func EditUserCoin(req request.EditUserCoin) (err error) {

	emailUser := account.SendEmailParams{
		Title: req.Title,
		Msg:   req.Msg,
		Type:  1,
	}

	if req.OpType != 1 {
		// 扣钱
		moneyInfo := account.UpdateUserMoneyInfo{
			Uid: uint64(req.Uid),
		}

		num := int64(req.Num * -1)

		switch req.CoinType {
		case 1: //winCash
			moneyInfo.WinCash = num
		case 2: //cash
			moneyInfo.Cash = num
		case 3: //bonus
			moneyInfo.Bonus = num
		}

		_, err = account.AddUserMoney(
			dbconn.NDB,
			account.LogType_admin,
			moneyInfo,
			dbconn.LogDB,
			"remark",
			"管理后台扣款",
		)

		if err != nil {
			global.Logger["err"].Infof("EditUserCoin 管理后台扣款 失败,err:%s", err.Error())
			return err
		}

	} else {

		attaches := make([]account.Attachment, 0, 1)

		tmp := account.Attachment{
			Nums: req.Num,
		}

		switch req.CoinType {
		case 1: //winCash
			tmp.ItemId = consts.ItemIdWinCash
		case 2: //cash
			tmp.ItemId = consts.ItemIdCash
		case 3: //bonus
			tmp.ItemId = consts.ItemIdBonus
		}

		attaches = append(attaches, tmp)

		err = SendEmail(emailUser, req.Uid, attaches)

		if err != nil {
			global.Logger["err"].Infof("editUserCoin account sendEmail failed:" + err.Error())
			return
		}
	}

	return
}

// 封禁
func Banned(req request.Banned) (err error) {

	userDb := global.User
	gmDb := global.DB

	obj := user.Banned{}

	err = obj.Upsert(userDb, user.BannedCateType(req.Cate), req.Info, req.Admin)

	if err != nil {
		return
	}

	gmObj := gm.Banned{}

	err = gmObj.Upsert(gmDb, gm.BannedCateType(req.Cate), req.Info, req.Uid, req.Admin)
	if err != nil {
		return
	}

	return
}

// 解封
func Unseal(req request.Banned) (err error) {

	userDb := global.User
	gmDb := global.DB

	obj := user.Banned{}

	err = obj.DeleteByUserCateInfo(userDb, user.BannedCateType(req.Cate), req.Info)

	if err != nil {
		return
	}

	gmObj := gm.Banned{}

	err = gmObj.DeleteByUserCateInfo(gmDb, gm.BannedCateType(req.Cate), req.Info, req.Uid)
	if err != nil {
		return
	}

	return
}

// 赠送数据处理
func GiveMoneyHandle(req request.GiveMoneyHandle) (err error) {
	payDb := global.Pay

	giveObj := new(pay.GiveMoney)

	var (
		give   pay.GiveMoney
		dictMp map[string]string
	)

	give, err = giveObj.GetOneByOrderNo(payDb, req.OrderNo)

	if err != nil {
		global.Logger["err"].Errorf("GiveMoneyHandle giveObj.GetOneByOrderNo failed,err:[%v]", err.Error())
		return
	}

	//自动审批
	gameDb := global.Game
	dictObj := new(game.Dict)

	dictMp, err = dictObj.GetDictMpByTypeCode(gameDb, "autoApprovalRules")
	if err != nil {
		global.Logger["err"].Errorf("GiveMoneyHandle dictObj.GetDictMpByTypeCode failed,err:[%v]", err.Error())
		return
	}

	//自动审批 开关
	onOff, sOk := dictMp["switch"]
	if !sOk {
		onOff = "0"
	}

	//单笔最大金额
	maxAmount, maxOk := dictMp["singleMaxAmount"]
	if !maxOk {
		maxAmount = "0"
	}

	//单日总金额
	dTotal, dtOk := dictMp["dailyTotal"]
	if !dtOk {
		dTotal = "0"
	}

	//自动审批充提比 = (单用户) 总赠送金额 / 总充值金额
	rwRate, rwOk := dictMp["rechargeWithdrawRate"]
	if !rwOk {
		rwRate = "0"
	}

	now := time.Now()

	ymd := timeutil.FormatToDateNumber(now)

	var dayGiveTotal string

	redis := global.Redis

	dayUserTotalKey := constant.UserTodayGiveTotal + ymd + strconv.Itoa(give.Uid)

	dayGiveTotal, err = redis.Get(dayUserTotalKey)

	if err != nil {
		//未查到，去数据库中查询
		var sumAmount int
		sumAmount, err = giveObj.SumAmountByUidAndCreatedAt(payDb, give.Uid, time.Now())
		if err != nil {
			global.Logger["err"].Infof("GiveMoneyHandle giveObj.SumAmountByUidAndCreatedAt failed, err:%s", err.Error())
			return
		}
		dayGiveTotal = strconv.Itoa(sumAmount)
	} else if dayGiveTotal == "" {
		dayGiveTotal = "0"
	}

	if onOff == "1" {
		if give.Recharge <= 0 {
			return ecode.UserNotRecharge
		}

		amountStr := strconv.Itoa(give.Amount)
		if amountStr <= maxAmount && dayGiveTotal <= dTotal && fmt.Sprintf("%.2f", give.CommitGiveRate) <= rwRate {
			give.Auditor = "automatic approval"
			give.AuditTime = &now
			give.Status = pay.GiveMoneyStatusInPayment
		}

		var (
			bank    pay.BankInfo
			bankObj pay.BankInfo
		)

		bank, err = bankObj.GetFirstByUid(payDb, int64(give.Uid))

		if err == nil {
			channelMp := make(map[int][]PayCfgAndRate)
			channelMp, err = GetChannelPassageListMp()

			if err == nil {

				channelList, channelMpOk := channelMp[give.ChannelId]

				if channelMpOk {
					var payRes response.PaymentResponse

					//自动打款
					payRes, err = Payment(give, bank, channelList)

					if err != nil {
						global.Logger["err"].Infof("GiveMoneyHandle 调用 Payment err [%v]", err.Error())
						give.Status = pay.GiveMoneyStatusUpstreamAbnormal //更新为异常，需人工审核
					} else {
						give.PayCfgId = payRes.CfgId
						give.Status = pay.GiveMoneyStatusInPayment
					}
				} else {
					global.Logger["err"].Infof("GiveMoneyHandle !channelMpOk ChannelId [%v]", give.ChannelId)
				}

			} else {
				global.Logger["err"].Infof("GiveMoneyHandle 调用 GetChannelPassageListMp: 出错 err:%s", err.Error())
			}

		} else {
			global.Logger["err"].Infof("GiveMoneyHandle uid: %v,bank info query err:%s", give.Uid, err.Error())
		}

	}
	var total int

	total, err = strconv.Atoi(dayGiveTotal)

	if err != nil {
		global.Logger["err"].Infof("GiveMoneyHandle strconv.Atoi(%v) failed,err:[%v]", dayGiveTotal, err.Error())
	}

	total += give.Amount

	_, err = redis.Set(dayUserTotalKey, strconv.Itoa(total), 86400)
	if err != nil {
		global.Logger["err"].Infof("GiveMoneyHandle redis.Set failed,key:%v,value:%v err:%s", dayUserTotalKey, strconv.Itoa(total), err.Error())
		return
	}

	//更新数据
	err = giveObj.Save(payDb, &give)
	if err != nil {
		global.Logger["err"].Errorf("GiveMoneyHandle giveObj.Save failed,err:[%v]", err.Error())
		return
	}

	return
}

func Payment(give pay.GiveMoney, bank pay.BankInfo, cfgList []PayCfgAndRate) (payRes response.PaymentResponse, err error) {

	mp := make(map[int]struct{}) //支付配置根据id去重

	payRes = response.PaymentResponse{}

	mode := global.ServerConfig.Mode

	//正式&&预发布&&测试才会发起接口调用
	modeStatus := mode != "release" && mode != "pre" && mode != "test"

	//默认失败，接口后替换
	err = errors.New("all pay channel failed")

	for _, cfg := range cfgList {

		if _, ok := mp[cfg.PayCfg.ID]; ok {
			continue
		}

		amount := float64(give.Amount)

		//最低赠送额度 100卢比
		if amount < 10000 {
			global.Logger["err"].Errorf("xPay failed,err:[Minimum limit of 100]")
			err = errors.New("Minimum limit of 100")
			return
		} else if amount > 10000 {
			//超过100的才扣税，100不扣税
			amount = calculatePaymentAmount(give.Amount, cfg.Rate)
		} else {
			amount = amount / 100
		}

		//一定是平台明确的失败才会切换，其他的就卡在那里，人工介入。避免重复到三方平台下代付单
		switch cfg.PayCfg.Markers {
		case "xPay":

			if modeStatus {
				break
			}

			err = TransferOrder(bank, cfg.PayCfg, give.OrderNo, amount)

			if err == nil {
				payRes.CfgId = cfg.PayCfg.ID
				return payRes, nil
			} else if errors.Is(err, &errUtil.UnmarshalErr{}) {
				//errUtil.UnmarshalErr 自定义错误类型，json验证错误
				//如果是这个错误，则直接返回，不继续后续循环了
				return payRes, err
			}

			global.Logger["err"].Errorf("xPay failed,err:[%s]", err.Error())

		case "inPay":

			if modeStatus {
				break
			}

			err = InPayPayment(cfg.PayCfg, give.OrderNo, amount, bank)

			if err == nil {
				payRes.CfgId = cfg.PayCfg.ID
				return payRes, nil
			} else if errors.Is(err, &errUtil.UnmarshalErr{}) {
				//errUtil.UnmarshalErr 自定义错误类型，json验证错误
				//如果是这个错误，则直接返回，不继续后续循环了
				return payRes, err
			}

			global.Logger["err"].Errorf("inPay failed,err:[%s]", err.Error())

		case "luckyPay":

			if modeStatus {
				break
			}

			err = LuckyPayPayment(cfg.PayCfg, give.OrderNo, amount, bank)

			if err == nil {
				payRes.CfgId = cfg.PayCfg.ID
				return payRes, nil
			} else if errors.Is(err, &errUtil.UnmarshalErr{}) {
				//errUtil.UnmarshalErr 自定义错误类型，json验证错误
				//如果是这个错误，则直接返回，不继续后续循环了
				return payRes, err
			}
			global.Logger["err"].Errorf("luckyPay failed,err:[%s]", err.Error())

		case "simulated":
			err = SimulatedPay(cfg.PayCfg, give.OrderNo, amount, bank)

			if err == nil {
				payRes.CfgId = cfg.PayCfg.ID
				return payRes, nil
			} else if errors.Is(err, &errUtil.UnmarshalErr{}) {
				//errUtil.UnmarshalErr 自定义错误类型，json解析错误 如果是这个错误，则直接返回，不继续后续循环了
				return payRes, err
			}
			global.Logger["err"].Errorf("simulated failed,err:[%s]", err.Error())
		}

		mp[cfg.PayCfg.ID] = struct{}{}
	}

	return
}

// xPay代付回调
func XPaymentCallbackNew(params url.Values, req request.XPaymentCallback) (err error) {

	var payCfg pay.PayConfig

	payCfg, err = GetSecretByMarkers("xPay")
	if err != nil {
		return err
	}

	//验签
	if !CheckFunPaySign(params, payCfg.Secret) {
		return ecode.SignCheckError
	}

	switch req.State {
	//订单状态 1-代付中，2-代付成功，3-代付失败，4-代付撤销
	case 1: //1-代付中
		return
	case 2: //2-代付成功
		return SuccessCallback(req.MchOrderNo, req.TransferId, req.SuccessTime)

	case 3, 4: //3-代付失败，4-代付撤销
		//3-代付失败，4-代付上游撤销
		status := 0

		if req.State == 3 {
			status = pay.GiveMoneyStatusFailed
		} else if req.State == 4 {
			status = pay.GiveMoneyStatusUpstreamRevocation
		}

		return FailedCallback(req.MchOrderNo, req.TransferId, status)

	}

	return
}

// xPay代付回调
func XPaymentCallback(req request.XPaymentCallback) (err error) {

	var payCfg pay.PayConfig

	payCfg, err = GetSecretByMarkers("xPay")
	if err != nil {
		return err
	}

	//验签
	if !CheckPaymentCallbackSign(req, payCfg.Secret) {
		return ecode.SignCheckError
	}

	switch req.State {
	//订单状态 1-代付中，2-代付成功，3-代付失败，4-代付撤销
	case 1: //1-代付中
		return
	case 2: //2-代付成功
		return SuccessCallback(req.MchOrderNo, req.TransferId, req.SuccessTime)

	case 3, 4: //3-代付失败，4-代付撤销
		//3-代付失败，4-代付上游撤销
		status := 0

		if req.State == 3 {
			status = pay.GiveMoneyStatusFailed
		} else if req.State == 4 {
			status = pay.GiveMoneyStatusUpstreamRevocation
		}

		return FailedCallback(req.MchOrderNo, req.TransferId, status)

	}

	return
}

// xPay代付回调
func XPaymentCallbackSuccess(req request.XPaymentCallbackSuccess) (err error) {
	if global.ServerConfig.Mode == "release" {
		return ecode.IllegalRequest
	}

	successTime := time.Now().UnixMilli()

	return SuccessCallback(req.OrderNo, strconv.FormatInt(successTime, 10), int(successTime))

}

// xPay代付回调
func XPaymentCallbackFailed(req request.PaymentCallbackFailed) (err error) {

	if global.ServerConfig.Mode == "release" {
		return ecode.IllegalRequest
	}

	switch req.State {
	//订单状态 1-代付中，2-代付成功，3-代付失败，4-代付撤销
	case 1:
		return
	case 3, 4:
		//3-代付失败，4-代付上游撤销

		transferId := strconv.Itoa(time.Now().Nanosecond())

		return FailedCallback(req.MchOrderNo, transferId, pay.GiveMoneyStatusFailed)

	}

	return
}

// inPay代付回调
func InPaymentCallback(req request.InPaymentCallback) (err error) {
	//验签
	err = pubDecrypt(req.Sign)

	if err != nil {
		global.Logger["err"].Infof("InPaymentCallback 验签失败 err:", err.Error())
		err = ecode.SignCheckError
		return
	}

	switch req.Status {
	//状态:1=处理中,2=拒绝,3=失败,4=成功,5=撤销
	case float64(1): //1=处理中
		return
	case float64(2): //2=拒绝

		return
	case float64(3): //3=失败

		return FailedCallback(req.OrderNumber, req.PlatNumber, pay.GiveMoneyStatusFailed)

	case float64(4): //4=成功
		return SuccessCallback(req.OrderNumber, req.PlatNumber, int(time.Now().Unix()))

	}
	return
}

func LuckyPaymentCallback(req request.LuckyPayCallback) (err error) {
	err = LuckyPayCallbackCheckSign(req)

	if err != nil {
		global.Logger["err"].Infof("LuckyPaymentCallback 验签失败 err:", err.Error())
		return
	}

	//代付支付状态（1-代付中、2-代付成功 、3-代付失败）
	switch req.PayState {
	case 1:
		return

	case 2: //2-代付成功
		var successTime int
		successTime, err = strconv.Atoi(req.PayFinishTime)

		if err != nil {
			successTime = int(time.Now().Unix())
		}

		return SuccessCallback(req.MchOrderNo, req.PayOrderNo, successTime)
	case 3: //3-代付失败
		return FailedCallback(req.MchOrderNo, req.PayOrderNo, pay.GiveMoneyStatusFailed)

	}
	return
}

func calculatePaymentAmount(amount int, rate float64) float64 {
	return float64(amount/100) * (100 - rate) / 100
}
