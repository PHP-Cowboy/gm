package daos

import (
	"fmt"
	"gm/common/constant"
	"gm/daos/rds"
	"gm/global"
	"gm/model/gm"
	"gm/model/log"
	"gm/model/pay"
	"gm/model/user"
	"gm/msgcenter"
	"gm/response"
	"gm/utils/timeutil"
	"strconv"
	"time"
)

// 报表数据任务
func ReportDataTask(now time.Time) (err error) {
	userDb := global.User
	gmDb := global.DB
	logDb := global.Log
	payDb := global.Pay

	//传入时间的前一天数据
	prevDate := now.AddDate(0, 0, -1)

	start := timeutil.GetZeroTime(prevDate).Format(timeutil.TimeFormat)
	end := timeutil.GetLastTime(prevDate).Format(timeutil.TimeFormat)
	ym := timeutil.GetLastTime(prevDate).Format(timeutil.MonthNumberFormat)
	ymd := prevDate.Format(timeutil.DateNumberFormat)
	channelMp := make(map[int]string)

	channelMp, err = rds.GetChannelIdNameMp()

	if err != nil {
		TaskLog("GetChannelMp err :" + err.Error())
		return
	}

	page := 1
	size := 100

	userObj := user.User{}
	roomFundsLogObj := log.RoomFundsFlowLog{}
	loginLogObj := log.UserLoginLog{}
	loginMp := make(map[int]int)       //渠道 用户数量
	userChannelMp := make(map[int]int) //渠道 用户数量
	orderObj := pay.Order{}
	orderMp := make(map[int]pay.AmountAndCash)
	yesterdayOrderMp := make(map[int]int, 0) //uid channel
	giveObj := pay.GiveMoney{}
	giveMp := make(map[int]pay.PeoplesAndAmount)

	userInfosMp := make(map[int]response.ReportData) //渠道数据

	userList := make([]user.User, 0, size)
	playMp := make(map[int]struct{})

	//DAU 今日登录用户 俩msp都已经根据uid去重
	loginMp, userChannelMp, err = loginLogObj.GetListByCreatedAt(logDb, start, end, ym)
	if err != nil {
		global.Logger["err"].Errorf("ReportDataTask GetListByCreatedAt failed,err:[%v]", err.Error())
		return
	}

	//提现 giveMp[channel] = sum(Amount)
	giveMp, err = giveObj.MapChannelAmountByArrivalTime(payDb, prevDate)
	if err != nil {
		global.Logger["err"].Errorf("ReportDataTask MapChannelAmountByArrivalTime failed,err:[%v]", err.Error())
		return
	}

	//统计任务时的前一天充值成功的用户 uid channel 根据uid去重了
	yesterdayOrderMp, err = orderObj.GetUidChannelListByYmd(payDb, prevDate.AddDate(0, 0, -1))
	if err != nil {
		global.Logger["err"].Errorf("ReportDataTask GetUidChannelListByYmd failed,err:[%v]", err.Error())
		return
	}

	for id, name := range channelMp {
		userInfo, userInfosMpOk := userInfosMp[id]

		if !userInfosMpOk {
			userInfo = response.ReportData{}
		}

		num, loginMpOk := loginMp[id]

		if !loginMpOk {
			num = 0
		}

		userInfo.DailyActiveUser = num
		userInfo.Channel = id
		userInfo.ChannelName = name

		userInfosMp[id] = userInfo
	}

	type PayUserNextDayLogin struct {
		YesterdayPayNum           int
		YesterdayPayTodayLoginNum int
	}

	paidUserRetentionMp := make(map[int]PayUserNextDayLogin) // key -> channel

	//处理 付费用户留存率 = 昨日付费的用户 在 今日还活跃的用户数   /   昨日付费的用户
	//yesterdayOrderMp uid channel
	for uid, channelId := range yesterdayOrderMp {

		payUserRetention, paidOk := paidUserRetentionMp[channelId]

		if !paidOk {
			payUserRetention = PayUserNextDayLogin{}
		}

		payUserRetention.YesterdayPayNum++

		//userChannelMp uid channel
		_, userOk := userChannelMp[uid]

		if userOk {
			payUserRetention.YesterdayPayTodayLoginNum++
		}

		paidUserRetentionMp[channelId] = payUserRetention
	}

	for {
		userList, err = userObj.GetPageListByCreateTime(userDb, start, end, page, size)

		if err != nil {
			global.Logger["err"].Errorf("ReportDataTask GetListByUserIds failed,err:[%v]", err.Error())
			return
		}

		userIds := make([]int, 0, len(userList))

		for _, ul := range userList {
			userIds = append(userIds, ul.Uid)
		}

		//新注册用户中，有任意对局行为的用户 playMp[uid] = struct{}{}
		playMp, err = roomFundsLogObj.GetPlayMpByUserIds(logDb, userIds, prevDate)

		if err != nil {
			global.Logger["err"].Errorf("ReportDataTask GetPlayMpByUserIds failed,err:[%v]", err.Error())
			return
		}

		//新增付费 orderMp[uid] = {Amount,Cash}
		orderMp, err = orderObj.GetAmountAndCoinByUserIdsAndYmd(payDb, userIds, prevDate)
		if err != nil {
			global.Logger["err"].Errorf("ReportDataTask GetAmountAndCoinByUserIdsAndYmd failed,err:[%v]", err.Error())
			return
		}

		for _, info := range userList {

			userStatistics, ok := userInfosMp[info.ChannelId]

			if !ok {
				userStatistics = response.ReportData{}
			}

			_, playOk := playMp[info.Uid]

			if playOk {
				userStatistics.EffectiveAddNums++ //新注册用户中，有任意对局行为的用户 log中记录处理
			}

			orderVal, orderOk := orderMp[info.Uid]

			if orderOk {
				userStatistics.NewAddUserRechargeCoinNum += orderVal.Cash
				userStatistics.NewAddUserRechargeAmount += float64(orderVal.Amount)
				userStatistics.NewAddUserRechargePeople++
			}

			userStatistics.AddNums++ //新增用户数

			userInfosMp[info.ChannelId] = userStatistics
		}

		page += 1

		//查到的数据条数小于size，则是最后一页的数据,跳出for循环
		if len(userList) < size {
			break
		}
	}

	dayOrderMp := make(map[int]int)
	//今日付费数据 dayOrderMp[uid] = Amount
	dayOrderMp, err = orderObj.GetListByYmd(payDb, prevDate) //dayOrderMp 已根据 uid 去重

	if err != nil {
		global.Logger["err"].Errorf("ReportDataTask GetListByYmd failed,err:[%v]", err.Error())
		return
	}

	var (
		orderUserIds  []int
		orderUserList []user.User
	)

	for k, _ := range dayOrderMp {
		orderUserIds = append(orderUserIds, k)
	}

	orderUserList, err = userObj.GetListByUserIds(userDb, orderUserIds)

	if err != nil {
		global.Logger["err"].Errorf("ReportDataTask GetListByUserIds failed,err:[%v]", err.Error())
		return err
	}

	type ChannelRecharge struct {
		Num    int
		Amount int
	}

	channelRechargeMp := make(map[int]ChannelRecharge)

	for _, ou := range orderUserList { //todo
		val, ok := channelRechargeMp[ou.ChannelId]

		if !ok {
			val = ChannelRecharge{}
		}

		amount, dayOrderMpOk := dayOrderMp[ou.Uid]

		if dayOrderMpOk {
			val.Num++
			val.Amount += amount
		}

		channelRechargeMp[ou.ChannelId] = val
	}

	users := make([]gm.ReportData, 0, len(userInfosMp))

	for _, u := range userInfosMp {

		channelRecharge, channelRechargeMpOk := channelRechargeMp[u.Channel]

		if !channelRechargeMpOk {
			channelRecharge = ChannelRecharge{}
		}

		give, _ := giveMp[u.Channel]

		tmp := gm.ReportData{
			Ymd:                           ymd,
			Channel:                       u.Channel,
			ChannelName:                   u.ChannelName,
			DailyActiveUser:               u.DailyActiveUser,
			EffectiveAddNums:              u.EffectiveAddNums,
			AddNums:                       u.AddNums,
			NewAddUserRechargePeople:      u.NewAddUserRechargePeople,
			NewAddUserRechargeCoinNum:     u.NewAddUserRechargeCoinNum / 100,
			NewAddUserRechargeAmount:      u.NewAddUserRechargeAmount / 100,
			DailyActiveUserRechargeAmount: channelRecharge.Amount / 100,                     //DAU付费总额
			OldUserRechargePeopleNum:      channelRecharge.Num - u.NewAddUserRechargePeople, //老用户充值金币人数 = 总充值人数 - 新用户充值人数
			GiveMoneyPeople:               give.PeopleNum,                                   //赠送金币人数
			GiveMoneyAmount:               give.SumAmount / 100,                             //赠送金币额
			NextDayRetention:              0,
			ThreeDayRetention:             0,
			FourDayRetention:              0,
			FiveDayRetention:              0,
			SixDayRetention:               0,
			SevenDayRetention:             0,
			FourteenDayRetention:          0,
		}

		if tmp.AddNums > 0 {
			tmp.NewAddUserAverageRevenuePerUser = tmp.NewAddUserRechargeAmount / float64(tmp.AddNums) // 新用户Arpu
			tmp.NewAddUserRechargeRate = float64(tmp.NewAddUserRechargePeople) / float64(tmp.AddNums) //新用户充值金币率
			tmp.PlayRate = float64(tmp.EffectiveAddNums) / float64(tmp.AddNums)                       //新用户玩牌率 对局用户 / 今日新增
		}

		//ARPU = 每用户平均收入 ARPU 是 渠道今日总充值钱数 / 渠道今日DAU数
		//ARPPU = 每付费用户平均收益 渠道今日总充值钱数 / 渠道今日充值人数

		if tmp.DailyActiveUser > 0 {
			//DAU ARPU = 渠道当日总付费 / 渠道当日总活跃人数 ；
			tmp.DailyActiveUserAverageRevenuePerUser = float64(tmp.DailyActiveUserRechargeAmount) / float64(tmp.DailyActiveUser) // DAU Arpu
			//tmp.PlayRate = float64(tmp.EffectiveAddNums) / float64(tmp.DailyActiveUser) todo 后期改成活跃玩牌率 EffectiveAddNums 这个要换成活跃用户的玩牌人数
		}

		if tmp.DailyActiveUserRechargeAmount > 0 {
			tmp.GiveMoneyRate = float64(tmp.GiveMoneyAmount) / float64(tmp.DailyActiveUserRechargeAmount) // 赠送金币率 = 赠送金币额度 / DAU充值总额
		}

		if tmp.NewAddUserRechargePeople+tmp.OldUserRechargePeopleNum > 0 {
			tmp.NewAddUserAverageRevenuePerPayingUser = tmp.NewAddUserRechargeAmount / float64(tmp.NewAddUserRechargePeople+tmp.OldUserRechargePeopleNum) //新用户Arppu
			//DAU ARPPU = 渠道当日总付费 / 渠道当日总付费人数=
			tmp.DailyActiveUserAverageRevenuePerPayingUser = float64(tmp.DailyActiveUserRechargeAmount) / float64((tmp.NewAddUserRechargePeople + tmp.OldUserRechargePeopleNum)) // DAU Arppu
		}

		//日活 - 新增用户 = 老用户数
		if tmp.DailyActiveUser-tmp.AddNums > 0 {
			tmp.OldUserRechargeRate = float64(tmp.OldUserRechargePeopleNum) / float64(tmp.DailyActiveUser-tmp.AddNums) // 老用户充值金币率 = 老用户充值人数 / 老用户数
		}

		paidUserRetention, paidUserRetentionOk := paidUserRetentionMp[u.Channel]

		// 付费用户留存率 = 昨日付费的用户 在 今日还活跃的用户数   /   昨日付费的用户
		if paidUserRetentionOk && paidUserRetention.YesterdayPayNum > 0 {
			tmp.YesterdayPayNum = paidUserRetention.YesterdayPayNum
			tmp.YesterdayPayTodayLoginNum = paidUserRetention.YesterdayPayTodayLoginNum
			tmp.PaidUserRetentionRate = float64(paidUserRetention.YesterdayPayTodayLoginNum) / float64(paidUserRetention.YesterdayPayNum)
		} else {
			tmp.YesterdayPayNum = 0           //昨日支付数为0
			tmp.YesterdayPayTodayLoginNum = 0 //昨日支付今日登录数为0
			tmp.PaidUserRetentionRate = 0
		}

		users = append(users, tmp)
	}

	uStatObj := new(gm.ReportData)

	err = uStatObj.CreateInBatches(gmDb, users)
	if err != nil {
		TaskLog("data CreateInBatches err:" + err.Error())
		return
	}
	return
}

// 报表数据全渠道合并任务
func MergeAllChannelReportDataTask(now time.Time) (err error) {
	//当天前一天数据
	prevDate := now.AddDate(0, 0, -1)

	ymd := prevDate.Format(timeutil.DateNumberFormat)

	db := global.DB

	var (
		reportDataObj = new(gm.ReportData)
		dataList      []gm.ReportData
	)

	dataList, err = reportDataObj.GetListByYmd(db, ymd)

	if err != nil {
		TaskLog(fmt.Sprintf("date: %s, reportData select err:%s", ymd, err.Error()))
	}

	mergeData := gm.ReportData{}

	for _, l := range dataList {
		if l.Channel == 0 {
			continue
		}
		mergeData.Ymd = l.Ymd
		mergeData.ChannelName = "all channel" // mergeData.ChannelName + l.ChannelName + ","
		mergeData.DailyActiveUser += l.DailyActiveUser
		mergeData.EffectiveAddNums += l.EffectiveAddNums
		mergeData.AddNums += l.AddNums
		mergeData.NewAddUserRechargeCoinNum += l.NewAddUserRechargeCoinNum
		mergeData.NewAddUserRechargeAmount += l.NewAddUserRechargeAmount
		mergeData.NewAddUserRechargePeople += l.NewAddUserRechargePeople
		mergeData.NewAddUserAverageRevenuePerPayingUser += l.NewAddUserAverageRevenuePerPayingUser
		mergeData.DailyActiveUserRechargeAmount += l.DailyActiveUserRechargeAmount
		mergeData.DailyActiveUserAverageRevenuePerPayingUser += l.DailyActiveUserAverageRevenuePerPayingUser
		mergeData.OldUserRechargePeopleNum += l.OldUserRechargePeopleNum
		mergeData.GiveMoneyPeople += l.GiveMoneyPeople
		mergeData.GiveMoneyAmount += l.GiveMoneyAmount
		mergeData.NextDayPeople += l.NextDayPeople
		mergeData.ThreeDayPeople += l.ThreeDayPeople
		mergeData.FourDayPeople += l.FourDayPeople
		mergeData.FiveDayPeople += l.FiveDayPeople
		mergeData.SixDayPeople += l.SixDayPeople
		mergeData.SevenDayPeople += l.SevenDayPeople
		mergeData.FourteenDayPeople += l.FourteenDayPeople
		mergeData.YesterdayPayNum += l.YesterdayPayNum
		mergeData.YesterdayPayTodayLoginNum += l.YesterdayPayTodayLoginNum
	}

	//mergeData.ChannelName = strings.TrimRight(mergeData.ChannelName, ",")

	if mergeData.AddNums > 0 {
		mergeData.NewAddUserRechargeRate = float64(mergeData.NewAddUserRechargePeople) / float64(mergeData.AddNums) //新用户充值金币率
	}

	if mergeData.DailyActiveUser-mergeData.AddNums > 0 {
		mergeData.OldUserRechargeRate = float64(mergeData.OldUserRechargePeopleNum) / float64(mergeData.DailyActiveUser-mergeData.AddNums)
	}

	if mergeData.DailyActiveUser > 0 {
		//mergeData.PlayRate = float64(mergeData.EffectiveAddNums) / float64(mergeData.DailyActiveUser) todo 后期改成活跃玩牌率 EffectiveAddNums 这个要换成活跃用户的玩牌人数

		//DAU ARPU = 渠道当日总付费 / 渠道当日总活跃人数
		mergeData.DailyActiveUserAverageRevenuePerUser = float64(mergeData.DailyActiveUserRechargeAmount) / float64(mergeData.DailyActiveUser) // DAU Arpu
	}

	if mergeData.DailyActiveUserRechargeAmount > 0 {
		mergeData.GiveMoneyRate = float64(mergeData.GiveMoneyAmount) / float64(mergeData.DailyActiveUserRechargeAmount) // 赠送金币率 = 赠送金币额度 / DAU充值总额
	}

	if mergeData.AddNums > 0 {
		mergeData.NewAddUserAverageRevenuePerUser = mergeData.NewAddUserRechargeAmount / float64(mergeData.AddNums) // 新用户Arpu
		mergeData.PlayRate = float64(mergeData.EffectiveAddNums) / float64(mergeData.AddNums)                       //新用户玩牌率 对局用户 / 今日新增
		mergeData.NextDayRetention = float64(mergeData.NextDayPeople) / float64(mergeData.AddNums)
		mergeData.ThreeDayRetention = float64(mergeData.ThreeDayPeople) / float64(mergeData.AddNums)
		mergeData.FourDayRetention = float64(mergeData.FourDayPeople) / float64(mergeData.AddNums)
		mergeData.FiveDayRetention = float64(mergeData.FiveDayPeople) / float64(mergeData.AddNums)
		mergeData.SixDayRetention = float64(mergeData.SixDayPeople) / float64(mergeData.AddNums)
		mergeData.SevenDayRetention = float64(mergeData.SevenDayPeople) / float64(mergeData.AddNums)
		mergeData.FourteenDayRetention = float64(mergeData.FourteenDayPeople) / float64(mergeData.AddNums)
	}

	if mergeData.YesterdayPayNum > 0 {
		// 付费用户留存率 = 昨日付费的用户 在 今日还活跃的用户数   /   昨日付费的用户
		mergeData.PaidUserRetentionRate = float64(mergeData.YesterdayPayTodayLoginNum) / float64(mergeData.YesterdayPayNum)
	}

	err = reportDataObj.CreateInBatches(db, []gm.ReportData{mergeData})
	if err != nil {
		TaskLog(fmt.Sprintf("date: %s, reportData save err:%s", ymd, err.Error()))
		return
	}

	return
}

func TaskLog(info string) {
	l, ok := global.Logger["task"]
	if !ok {
		panic("task日志加载失败")
	}
	l.Infof(info)
}

func GetDiffDays(ymd int) (days int, err error) {
	var date time.Time

	// 将 ymd 解析为 time.Time
	date, err = time.ParseInLocation(timeutil.DateNumberFormat, strconv.Itoa(ymd), time.Local)

	if err != nil {
		return
	}
	// 获取当前时间
	now := time.Now()

	// 计算两个时间之间的差值
	duration := now.Sub(date)

	// 注意：duration是time.Duration类型，它表示一个时间段，但我们需要的是天数
	// 由于time.Duration是以纳秒为单位的，我们可以通过除以24小时内的纳秒数来得到天数
	days = int(duration.Hours() / 24)

	return

}

func RechargeStatisticsTask(now time.Time) (err error) {
	userDb := global.User
	payDb := global.Pay
	gmDb := global.DB

	//当天前一天数据
	prevDate := now.AddDate(0, 0, -1)

	start := timeutil.GetZeroTime(prevDate).Format(timeutil.TimeFormat)
	end := timeutil.GetLastTime(prevDate).Format(timeutil.TimeFormat)
	ymd := prevDate.Format(timeutil.DateNumberFormat)
	channelMp := make(map[int]string)

	channelMp, err = rds.GetChannelIdNameMp()

	if err != nil {
		TaskLog("GetChannelMp err :" + err.Error())
		return
	}

	page := 1
	size := 100

	userObj := user.User{}
	orderObj := pay.Order{}

	userOrderMp := make(map[int]int)

	//充值数据
	userOrderMp, err = orderObj.GetListByYmd(payDb, prevDate)

	userList := make([]user.User, 0, size)

	userInfosMp := make(map[int]response.Recharge) //渠道数据

	for {
		//新注册用户
		userList, err = userObj.GetPageListByCreateTime(userDb, start, end, page, size)

		if err != nil {
			global.Logger["err"].Errorf("RechargeStatisticsTask GetPageListByCreateTime failed,err:[%v]", err.Error())
			return
		}

		for _, u := range userList {
			userInfos, userInfosMpOk := userInfosMp[u.ChannelId]

			if !userInfosMpOk {
				userInfos = response.Recharge{}
			}

			channelName := channelMp[u.ChannelId]

			userInfos.ChannelName = channelName

			amount, userOrderMpOk := userOrderMp[u.Uid]

			if userOrderMpOk {
				userInfos.NewUserRechargeNums++
				userInfos.NewUserRechargeTotal += amount
			}

			userInfos.AddUserNum++

			userInfosMp[u.ChannelId] = userInfos

			//删除今日充值的新用户，剩余即为老用户
			delete(userOrderMp, u.Uid)
		}

		page += 1

		//查到的数据条数小于size，则是最后一页的数据,跳出for循环
		if len(userList) < size {
			break
		}
	}

	userIds := make([]int, 0, len(userOrderMp))
	oldUserList := make([]user.User, 0, size)

	//新用户已被删除
	for id, _ := range userOrderMp {
		userIds = append(userIds, id)
	}

	//老用户数据
	oldUserList, err = userObj.GetListByUserIds(userDb, userIds)

	if err != nil {
		global.Logger["err"].Errorf("RechargeStatisticsTask GetListByUserIds failed,err:[%v]", err.Error())
		return
	}

	for _, ou := range oldUserList {
		userInfos, userInfosMpOk := userInfosMp[ou.ChannelId]

		if !userInfosMpOk {
			userInfos = response.Recharge{}
		}

		channelName := channelMp[ou.ChannelId]

		userInfos.ChannelName = channelName

		amount, userOrderMpOk := userOrderMp[ou.Uid]

		if userOrderMpOk {
			userInfos.OldUserRechargeNums++
			userInfos.OldUserRechargeTotal += amount
		}

		userInfosMp[ou.ChannelId] = userInfos
	}

	rechargeList := make([]gm.RechargeStatistics, 0, len(userInfosMp))

	for channelId, recharge := range userInfosMp {

		tmp := gm.RechargeStatistics{
			Ymd:                      ymd,
			Channel:                  channelId,
			ChannelName:              recharge.ChannelName,
			NewUserRechargeNums:      recharge.NewUserRechargeNums,
			NewUserRechargeTotal:     recharge.NewUserRechargeTotal,
			OldUserRechargeNums:      recharge.OldUserRechargeNums,
			OldUserFirstRechargeNums: 0, //老用户首次付费人数
			OldUserRechargeTotal:     recharge.OldUserRechargeTotal,
			OldUserRechargeRate:      0, //老用户充值率
			AverageRevenuePerOldUser: 0, //老用户每用户平均收入

		}

		if recharge.NewUserRechargeNums > 0 {
			tmp.AverageRevenuePerPayingNewUser = recharge.NewUserRechargeTotal / recharge.NewUserRechargeNums
		}

		if recharge.AddUserNum > 0 {
			tmp.NewUserRechargeRate = float64(recharge.NewUserRechargeNums / recharge.AddUserNum)
			tmp.AverageRevenuePerNewUser = recharge.NewUserRechargeTotal / recharge.AddUserNum
		}

		if recharge.OldUserRechargeNums > 0 {
			tmp.AverageRevenuePerPayingOldUser = recharge.OldUserRechargeTotal / recharge.OldUserRechargeNums
		}

		rechargeList = append(rechargeList, tmp)
	}

	uStatObj := new(gm.RechargeStatistics)

	err = uStatObj.CreateInBatches(gmDb, rechargeList)
	if err != nil {
		TaskLog("data CreateInBatches err:" + err.Error())
		return
	}

	return
}

// 提现用户统计任务
func WithdrawalStatisticsTask(now time.Time) (err error) {
	userDb := global.User
	gmDb := global.DB
	payDb := global.Pay

	//当天前一天数据
	prevDate := now.AddDate(0, 0, -1)

	start := timeutil.GetZeroTime(prevDate).Format(timeutil.TimeFormat)
	end := timeutil.GetLastTime(prevDate).Format(timeutil.TimeFormat)
	ymd := prevDate.Format(timeutil.DateNumberFormat)
	channelMp := make(map[int]string)

	channelMp, err = rds.GetChannelIdNameMp()

	if err != nil {
		TaskLog("GetChannelMp err :" + err.Error())
		return
	}

	page := 1
	size := 100

	userObj := user.User{}
	giveObj := pay.GiveMoney{}

	giveMp := make(map[int]int)

	//充值数据
	giveMp, err = giveObj.GetListByCreatedAt(payDb, prevDate)

	userList := make([]user.User, 0, size)

	userInfosMp := make(map[int]response.Withdrawal) //渠道数据

	for {
		//新注册用户
		userList, err = userObj.GetPageListByCreateTime(userDb, start, end, page, size)

		if err != nil {
			global.Logger["err"].Errorf("WithdrawalStatisticsTask GetPageListByCreateTime failed,err:[%v]", err.Error())
			return
		}

		for _, u := range userList {
			userInfos, userInfosMpOk := userInfosMp[u.ChannelId]

			if !userInfosMpOk {
				userInfos = response.Withdrawal{}
			}

			channelName := channelMp[u.ChannelId]

			userInfos.ChannelName = channelName

			amount, giveMpOk := giveMp[u.Uid]

			if giveMpOk {
				userInfos.NewUserWithdrawalNums++
				userInfos.NewUserWithdrawalTotal += amount
			}

			userInfosMp[u.ChannelId] = userInfos

			//删除今日充值的新用户，剩余即为老用户
			delete(giveMp, u.Uid)
		}

		page += 1

		//查到的数据条数小于size，则是最后一页的数据,跳出for循环
		if len(userList) < size {
			break
		}
	}

	userIds := make([]int, 0, len(giveMp))
	oldUserList := make([]user.User, 0, size)

	//新用户已被删除
	for id, _ := range giveMp {
		userIds = append(userIds, id)
	}

	//老用户数据
	oldUserList, err = userObj.GetListByUserIds(userDb, userIds)

	if err != nil {
		global.Logger["err"].Errorf("WithdrawalStatisticsTask GetListByUserIds failed,err:[%v]", err.Error())
		return
	}

	for _, ou := range oldUserList {
		userInfos, userInfosMpOk := userInfosMp[ou.ChannelId]

		if !userInfosMpOk {
			userInfos = response.Withdrawal{}
		}

		channelName := channelMp[ou.ChannelId]

		userInfos.ChannelName = channelName

		amount, userOrderMpOk := giveMp[ou.Uid]

		if userOrderMpOk {
			userInfos.OldUserWithdrawalNums++
			userInfos.OldUserWithdrawalTotal += amount
		}

		userInfosMp[ou.ChannelId] = userInfos
	}

	withdrawalList := make([]gm.WithdrawalStatistics, 0, len(userInfosMp))

	for channelId, recharge := range userInfosMp {

		tmp := gm.WithdrawalStatistics{
			Ymd:                    ymd,
			Channel:                channelId,
			ChannelName:            recharge.ChannelName,
			WithdrawalNums:         recharge.NewUserWithdrawalNums + recharge.OldUserWithdrawalNums,
			WithdrawalTotal:        recharge.NewUserWithdrawalTotal + recharge.OldUserWithdrawalTotal,
			NewUserWithdrawalNums:  recharge.NewUserWithdrawalNums,
			NewUserWithdrawalTotal: recharge.NewUserWithdrawalTotal,
			OldUserWithdrawalNums:  recharge.OldUserWithdrawalNums,
			OldUserWithdrawalTotal: recharge.OldUserWithdrawalTotal,
			OldUserWithdrawalRate:  0, //老用户提现率
		}

		if tmp.WithdrawalNums > 0 {
			tmp.NewUserWithdrawalRate = float64(tmp.NewUserWithdrawalNums / tmp.WithdrawalNums)
			tmp.OldUserWithdrawalRate = float64(tmp.OldUserWithdrawalNums / tmp.WithdrawalNums)
		}

		withdrawalList = append(withdrawalList, tmp)
	}

	uStatObj := new(gm.WithdrawalStatistics)

	err = uStatObj.CreateInBatches(gmDb, withdrawalList)
	if err != nil {
		TaskLog("data CreateInBatches err:" + err.Error())
		return
	}

	return
}

// 付费用户留存统计任务 todo
func PaidUserRetentionTaskNew(now time.Time) (err error) {
	userDb := global.User
	gmDb := global.DB
	payDb := global.Pay

	//当天前一天数据
	prevDate := now.AddDate(0, 0, -1)

	start := timeutil.GetZeroTime(prevDate).Format(timeutil.TimeFormat)
	end := timeutil.GetLastTime(prevDate).Format(timeutil.TimeFormat)
	ymd := prevDate.Format(timeutil.DateNumberFormat)
	channelMp := make(map[int]string)

	channelMp, err = rds.GetChannelIdNameMp()

	if err != nil {
		TaskLog("GetChannelMp err :" + err.Error())
		return
	}

	page := 1
	size := 100

	userObj := user.User{}

	giveObj := pay.GiveMoney{}

	giveMp := make(map[int]int)

	//充值数据
	giveMp, err = giveObj.GetListByCreatedAt(payDb, prevDate)

	userList := make([]user.User, 0, size)

	userInfosMp := make(map[int]response.Withdrawal) //渠道数据

	for {
		//新注册用户
		userList, err = userObj.GetPageListByCreateTime(userDb, start, end, page, size)

		if err != nil {
			global.Logger["err"].Errorf("PaidUserRetentionTaskNew GetPageListByCreateTime failed,err:[%v]", err.Error())
			return
		}

		for _, u := range userList {
			userInfos, userInfosMpOk := userInfosMp[u.ChannelId]

			if !userInfosMpOk {
				userInfos = response.Withdrawal{}
			}

			channelName := channelMp[u.ChannelId]

			userInfos.ChannelName = channelName

			amount, giveMpOk := giveMp[u.Uid]

			if giveMpOk {
				userInfos.NewUserWithdrawalNums++
				userInfos.NewUserWithdrawalTotal += amount
			}

			userInfosMp[u.ChannelId] = userInfos

			//删除今日充值的新用户，剩余即为老用户
			delete(giveMp, u.Uid)
		}

		page += 1

		//查到的数据条数小于size，则是最后一页的数据,跳出for循环
		if len(userList) < size {
			break
		}
	}

	userIds := make([]int, 0, len(giveMp))
	oldUserList := make([]user.User, 0, size)

	//新用户已被删除
	for id, _ := range giveMp {
		userIds = append(userIds, id)
	}

	//老用户数据
	oldUserList, err = userObj.GetListByUserIds(userDb, userIds)

	if err != nil {
		global.Logger["err"].Errorf("PaidUserRetentionTaskNew GetListByUserIds failed,err:[%v]", err.Error())
		return
	}

	for _, ou := range oldUserList {
		userInfos, userInfosMpOk := userInfosMp[ou.ChannelId]

		if !userInfosMpOk {
			userInfos = response.Withdrawal{}
		}

		channelName := channelMp[ou.ChannelId]

		userInfos.ChannelName = channelName

		amount, userOrderMpOk := giveMp[ou.Uid]

		if userOrderMpOk {
			userInfos.OldUserWithdrawalNums++
			userInfos.OldUserWithdrawalTotal += amount
		}

		userInfosMp[ou.ChannelId] = userInfos
	}

	withdrawalList := make([]gm.WithdrawalStatistics, 0, len(userInfosMp))

	for channelId, recharge := range userInfosMp {

		tmp := gm.WithdrawalStatistics{
			Ymd:                    ymd,
			Channel:                channelId,
			ChannelName:            recharge.ChannelName,
			WithdrawalNums:         recharge.NewUserWithdrawalNums + recharge.OldUserWithdrawalNums,
			WithdrawalTotal:        recharge.NewUserWithdrawalTotal + recharge.OldUserWithdrawalTotal,
			NewUserWithdrawalNums:  recharge.NewUserWithdrawalNums,
			NewUserWithdrawalTotal: recharge.NewUserWithdrawalTotal,
			OldUserWithdrawalNums:  recharge.OldUserWithdrawalNums,
			OldUserWithdrawalTotal: recharge.OldUserWithdrawalTotal,
			OldUserWithdrawalRate:  0, //老用户提现率
		}

		if tmp.WithdrawalNums > 0 {
			tmp.NewUserWithdrawalRate = float64(tmp.NewUserWithdrawalNums / tmp.WithdrawalNums)
			tmp.OldUserWithdrawalRate = float64(tmp.OldUserWithdrawalNums / tmp.WithdrawalNums)
		}

		withdrawalList = append(withdrawalList, tmp)
	}

	uStatObj := new(gm.WithdrawalStatistics)

	err = uStatObj.CreateInBatches(gmDb, withdrawalList)
	if err != nil {
		TaskLog("data CreateInBatches err:" + err.Error())
		return
	}
	return
}

// 付费用户留存统计任务
func PaidUserRetentionTask(now time.Time) (err error) {
	userDb := global.User
	gmDb := global.DB
	logDb := global.Log
	payDb := global.Pay

	//当天前一天数据
	prevDate := now.AddDate(0, 0, -1)

	start := timeutil.GetZeroTime(prevDate).Format(timeutil.TimeFormat)
	end := timeutil.GetLastTime(prevDate).Format(timeutil.TimeFormat)
	ymd := prevDate.Format(timeutil.DateNumberFormat)
	channelMp := make(map[int]string)

	channelMp, err = rds.GetChannelIdNameMp()

	if err != nil {
		TaskLog("GetChannelMp err :" + err.Error())
		return
	}

	page := 1
	size := 100

	userObj := user.User{}
	roomFundsLogObj := log.RoomFundsFlowLog{}
	loginLogObj := log.UserLoginLog{}
	loginMp := make(map[int]struct{})
	orderObj := pay.Order{}
	orderMp := make(map[int]int)
	giveObj := pay.GiveMoney{}
	giveMp := make(map[int]int)
	//渠道设备数量
	deviceNumMp := make(map[int]int)
	//设备
	deviceMp := make(map[string]struct{})

	userList := make([]user.User, 0, size)
	playMp := make(map[int]struct{})
	userInfosMp := make(map[int]response.UserStatisticsList) //渠道数据

	for {
		userList, err = userObj.GetPageListByCreateTime(userDb, start, end, page, size)

		if err != nil {
			global.Logger["err"].Errorf("PaidUserRetentionTask GetPageListByCreateTime failed,err:[%v]", err.Error())
			return
		}

		userIds := make([]int, 0, len(userList))

		for _, ul := range userList {
			userIds = append(userIds, ul.Uid)
		}

		//新注册用户中，有任意对局行为的用户
		playMp, err = roomFundsLogObj.GetPlayMpByUserIds(logDb, userIds, prevDate)

		if err != nil {
			global.Logger["err"].Errorf("PaidUserRetentionTask GetPlayMpByUserIds failed,err:[%v]", err.Error())
			return
		}

		//日活 即为当日登录用户
		loginMp, err = loginLogObj.GetListByUserIdsAndCreatedAt(logDb, userIds, prevDate)
		if err != nil {
			global.Logger["err"].Errorf("PaidUserRetentionTask GetListByUserIdsAndCreatedAt failed,err:[%v]", err.Error())
			return
		}

		//新增付费
		orderMp, err = orderObj.GetListByUserIdsAndYmd(payDb, userIds, prevDate)
		if err != nil {
			global.Logger["err"].Errorf("PaidUserRetentionTask GetListByUserIdsAndYmd failed,err:[%v]", err.Error())
			return
		}

		//提现
		giveMp, err = giveObj.GetListByUserIdsAndArrivalTime(payDb, userIds, prevDate)
		if err != nil {
			global.Logger["err"].Errorf("PaidUserRetentionTask GetListByUserIdsAndCreatedAt failed,err:[%v]", err.Error())
			return
		}

		for _, info := range userList {

			_, deviceOk := deviceMp[info.Device]
			if !deviceOk {
				deviceNumMp[info.ChannelId]++
			}

			userStatistics, ok := userInfosMp[info.ChannelId]

			if !ok {
				userStatistics.AddNums = 0
			}

			cName, cOk := channelMp[info.ChannelId]

			if !cOk {
				cName = ""
			}

			_, playOk := playMp[info.Uid]

			if playOk {
				userStatistics.EffectiveAddNums++ //新注册用户中，有任意对局行为的用户 log中记录处理
			}

			_, loginOk := loginMp[info.Uid]

			if loginOk {
				userStatistics.DailyActiveUser++ //日活跃用户 当日登录
			}

			orderAmount, orderOk := orderMp[info.Uid]

			if orderOk {
				userStatistics.NewPayingSubscribers++        //付费用户数
				userStatistics.NewPayingMoney += orderAmount //总付费
			}

			giveAmount, giveOk := giveMp[info.Uid]

			if !giveOk {
				userStatistics.TotalWithdrawal += giveAmount //提现总额
			}

			userStatistics.Channel = info.ChannelId //渠道id
			userStatistics.ChannelName = cName      //渠道名称
			userStatistics.AddNums++                //新增用户数
			userStatistics.NextDayRetention = 0     //次日留存
			userStatistics.ThreeDayRetention = 0    //三日留存
			userStatistics.SevenDayRetention = 0    //七日留存
			userStatistics.FourteenDayRetention = 0 //十四日留存
			userStatistics.ThirtyDayRetention = 0   //三十日留存

			userInfosMp[info.ChannelId] = userStatistics
		}

		page += 1

		//查到的数据条数小于size，则是最后一页的数据,跳出for循环
		if len(userList) < size {
			break
		}
	}

	dayOrderMp := make(map[int]int)
	//今日付费数据,uid 已去重
	dayOrderMp, err = orderObj.GetListByYmd(payDb, prevDate)

	if err != nil {
		global.Logger["err"].Errorf("PaidUserRetentionTask GetListByYmd failed,err:[%v]", err.Error())
		return
	}

	var (
		orderUserIds  []int
		orderUserList []user.User
	)

	for k, _ := range dayOrderMp {
		orderUserIds = append(orderUserIds, k)
	}

	orderUserList, err = userObj.GetListByUserIds(userDb, orderUserIds)

	if err != nil {
		global.Logger["err"].Errorf("PaidUserRetentionTask GetListByUserIds failed,err:[%v]", err.Error())
		return
	}

	type ChannelRecharge struct {
		Num    int
		Amount int
	}

	channelRechargeMp := make(map[int]ChannelRecharge)

	for _, ou := range orderUserList {
		val, ok := channelRechargeMp[ou.ChannelId]

		if !ok {
			val = ChannelRecharge{}
		}

		amount := dayOrderMp[ou.Uid]

		val.Num++
		val.Amount += amount
	}

	//users := make([]gm.UserStatistics, 0, len(userInfosMp))
	users := make([]gm.PaidUserRetention, 0, len(userInfosMp))

	for _, u := range userInfosMp {

		tmp := gm.PaidUserRetention{
			Ymd:                   ymd,
			Channel:               u.Channel,     //渠道id
			ChannelName:           u.ChannelName, //渠道名称
			UserNums:              u.AddNums,
			NextDayRetention:      0,
			TwoDayRetention:       0,
			ThreeDayRetention:     0,
			FourDayRetention:      0,
			FiveDayRetention:      0,
			SixDayRetention:       0,
			SevenDayRetention:     0,
			FourteenDayRetention:  0,
			TwentyOneDayRetention: 0,
			ThirtyDayRetention:    0,
			SixtyDayRetention:     0,
			NinetyDayRetention:    0,
		}

		users = append(users, tmp)
	}

	uStatObj := new(gm.PaidUserRetention)

	err = uStatObj.CreateInBatches(gmDb, users)
	if err != nil {
		TaskLog("data CreateInBatches err:" + err.Error())
		return
	}

	return
}

// 用户留存统计任务
func UserRetentionTask(now time.Time) (err error) {
	userDb := global.User
	gmDb := global.DB

	//当天前一天数据
	prevDate := now.AddDate(0, 0, -1)

	start := timeutil.GetZeroTime(prevDate).Format(timeutil.TimeFormat)
	end := timeutil.GetLastTime(prevDate).Format(timeutil.TimeFormat)
	ymd := prevDate.Format(timeutil.DateNumberFormat)
	channelMp := make(map[int]string)

	channelMp, err = rds.GetChannelIdNameMp()

	if err != nil {
		TaskLog("GetChannelMp err :" + err.Error())
		return
	}

	page := 1
	size := 100

	userObj := user.User{}

	userList := make([]user.User, 0, size)

	userInfosMp := make(map[int]response.UserRetentionList) //渠道数据

	for {
		//新注册用户
		userList, err = userObj.GetPageListByCreateTime(userDb, start, end, page, size)

		if err != nil {
			global.Logger["err"].Errorf("UserRetentionTask GetPageListByCreateTime failed,err:[%v]", err.Error())
			return
		}

		for _, u := range userList {
			userInfos, userInfosMpOk := userInfosMp[u.ChannelId]

			if !userInfosMpOk {
				userInfos = response.UserRetentionList{}
			}

			channelName := channelMp[u.ChannelId]

			userInfos.ChannelName = channelName

			userInfos.UserNums++

			userInfosMp[u.ChannelId] = userInfos

		}

		page += 1

		//查到的数据条数小于size，则是最后一页的数据,跳出for循环
		if len(userList) < size {
			break
		}
	}

	userRetentionList := make([]gm.UserRetention, 0)

	for channelId, u := range userInfosMp {
		userRetentionList = append(userRetentionList, gm.UserRetention{
			Ymd:                   ymd,
			Channel:               channelId,
			ChannelName:           u.ChannelName,
			UserNums:              u.UserNums,
			NextDayRetention:      0,
			TwoDayRetention:       0,
			ThreeDayRetention:     0,
			FourDayRetention:      0,
			FiveDayRetention:      0,
			SixDayRetention:       0,
			SevenDayRetention:     0,
			FourteenDayRetention:  0,
			TwentyOneDayRetention: 0,
			ThirtyDayRetention:    0,
			SixtyDayRetention:     0,
			NinetyDayRetention:    0,
		})
	}

	uStatObj := new(gm.UserRetention)

	err = uStatObj.CreateInBatches(gmDb, userRetentionList)
	if err != nil {
		TaskLog("data CreateInBatches err:" + err.Error())
		return
	}

	return
}

// 五分钟数据
func FiveMinuteDataTask() (err error) {
	gmDb := global.DB
	userDb := global.User
	payDb := global.Pay
	logDb := global.Log

	now := time.Now().Truncate(time.Minute) // 获取当前时间，截断到分钟级别

	remainder := now.Minute() % 5

	t := now.Add(time.Duration(remainder) * time.Minute * -1)

	//start := t.Add(-5 * time.Minute).Format(timeutil.TimeFormat)
	//end := t.Format(timeutil.TimeFormat)

	var (
		channelRegNumMp     = make(map[int]int)
		channelOnlineNumMp  = make(map[int]int)
		channelDauNumMp     = make(map[int]int)
		channelPayMp        = make(map[int]pay.ChannelPeopleNumAndTotalAmount)
		channelGiveAmountMp = make(map[int]int)
		channelMp           = make(map[int]string)
	)

	channelMp, err = rds.GetChannelIdNameMp()

	if err != nil {
		global.Logger["err"].Errorf("FiveMinuteDataTask rds.GetChannelIdNameMp failed,err:[%v]", err.Error())
		return
	}

	tNow := time.Now()

	todayStart := timeutil.GetZeroTime(tNow).Format(timeutil.TimeFormat)
	todayEnd := timeutil.GetLastTime(tNow).Format(timeutil.TimeFormat)

	//渠道新增注册
	userObj := new(user.User)
	channelRegNumMp, err = userObj.CountChannelRegNumByCreateTime(userDb, todayStart, todayEnd)

	if err != nil {
		global.Logger["err"].Errorf("FiveMinuteDataTask userObj.CountByCreateTime failed,err:[%v]", err.Error())
		return
	}
	//实时在线
	onlineIds := msgcenter.GetOnlineUsers()

	channelOnlineNumMp, err = userObj.CountChannelNumByUserIds(userDb, onlineIds)

	if err != nil {
		global.Logger["err"].Errorf("FiveMinuteDataTask userObj.CountChannelNumByUserIds failed,err:[%v]", err.Error())
		return
	}

	//活跃人数
	loginLogObj := new(log.UserLoginLog)

	channelDauNumMp, err = loginLogObj.CountChannelLoginNumByCreateAt(logDb, todayStart, todayEnd, timeutil.FormatToMonthNumber(tNow))

	if err != nil {
		global.Logger["err"].Errorf("FiveMinuteDataTask loginLogObj.CountChannelLoginNumByCreateAt failed,err:[%v]", err.Error())
		return
	}

	orderObj := new(pay.Order)
	//付费人数 付费额度
	channelPayMp, err = orderObj.CountChannelRechargeByTime(payDb, todayStart, todayEnd)

	if err != nil {
		global.Logger["err"].Errorf("FiveMinuteDataTask orderObj.CountChannelRechargeByTime failed,err:[%v]", err.Error())
		return
	}

	//赠送金币额度
	giveObj := new(pay.GiveMoney)

	channelGiveAmountMp, err = giveObj.SumChannelAmountByCreatedAt(payDb, todayStart, todayEnd)

	if err != nil {
		global.Logger["err"].Errorf("FiveMinuteDataTask giveObj.SumChannelAmountByCreatedAt failed,err:[%v]", err.Error())
		return
	}

	fiveObj := new(gm.FiveMinuteData)

	dataList := make([]gm.FiveMinuteData, 0, len(channelMp))

	for id, _ := range channelMp {

		regNum, _ := channelRegNumMp[id]
		onlineNum, _ := channelOnlineNumMp[id]
		activeNum, _ := channelDauNumMp[id]
		channelPay, payOk := channelPayMp[id]

		if !payOk {
			channelPay = pay.ChannelPeopleNumAndTotalAmount{}
		}

		giveAmount, _ := channelGiveAmountMp[id]

		tmp := gm.FiveMinuteData{
			TimePoint:  t.Format(timeutil.TimePointMinute),
			Channel:    id,
			RegNum:     regNum,
			OnlineNum:  onlineNum,
			ActiveNum:  activeNum,
			PayNum:     channelPay.PeopleNum,
			PayAmount:  channelPay.TotalAmount,
			GiveAmount: giveAmount,
		}

		dataList = append(dataList, tmp)
	}

	err = fiveObj.CreateInBatches(gmDb, dataList)

	if err != nil {
		global.Logger["err"].Errorf("FiveMinuteDataTask fiveObj.CreateInBatches failed,err:[%v]", err.Error())
		return err
	}

	return
}

// 小时用户数据任务-注册人数
func HourUserDataTaskRegNum() (err error) {

	db := global.DB
	userDb := global.User

	now := time.Now()

	hour := now.Hour()

	obj := new(gm.PerHourDataNum)

	upData := make(map[string]interface{})

	userObj := new(user.User)

	zeroMinSec := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())

	// 格式化时间字符串
	start := zeroMinSec.Format(timeutil.TimeFormat)

	end := zeroMinSec.Add(1 * time.Hour).Format(timeutil.TimeFormat)

	var num int64

	num, err = userObj.CountByCreateTime(userDb, start, end)

	if err != nil {
		TaskLog(fmt.Sprintf("统计注册人数失败,开始时间:%v,err:%s", start, err.Error()))
		return
	}

	switch hour {
	case 0:
		upData["zero_num"] = num
	case 1:
		upData["one_num"] = num
	case 2:
		upData["two_num"] = num
	case 3:
		upData["three_num"] = num
	case 4:
		upData["four_num"] = num
	case 5:
		upData["five_num"] = num
	case 6:
		upData["six_num"] = num
	case 7:
		upData["seven_num"] = num
	case 8:
		upData["eight_num"] = num
	case 9:
		upData["nine_num"] = num
	case 10:
		upData["ten_num"] = num
	case 11:
		upData["eleven_num"] = num
	case 12:
		upData["twelve_num"] = num
	case 13:
		upData["thirteen_num"] = num
	case 14:
		upData["fourteen_num"] = num
	case 15:
		upData["fifteen_num"] = num
	case 16:
		upData["sixteen_num"] = num
	case 17:
		upData["seventeen_num"] = num
	case 18:
		upData["eighteen_num"] = num
	case 19:
		upData["nineteen_num"] = num
	case 20:
		upData["twenty_num"] = num
	case 21:
		upData["twenty_one_num"] = num
	case 22:
		upData["twenty_two_num"] = num
	case 23:
		upData["twenty_three_num"] = num
	}

	dateNumber := now.Format(timeutil.DateNumberFormat)

	err = obj.UpdateByYmdAndType(db, dateNumber, gm.PerHourDataNumTypeReg, upData)

	return

}

// 小时用户数据任务-实时在线人数
func HourUserDataTaskOnline() (err error) {

	db := global.DB

	userIds := msgcenter.GetOnlineUsers()

	now := time.Now()

	hour := now.Hour()

	obj := new(gm.PerHourDataNum)

	upData := make(map[string]interface{})

	num := len(userIds)

	switch hour {
	case 0:
		upData["zero_num"] = num
	case 1:
		upData["one_num"] = num
	case 2:
		upData["two_num"] = num
	case 3:
		upData["three_num"] = num
	case 4:
		upData["four_num"] = num
	case 5:
		upData["five_num"] = num
	case 6:
		upData["six_num"] = num
	case 7:
		upData["seven_num"] = num
	case 8:
		upData["eight_num"] = num
	case 9:
		upData["nine_num"] = num
	case 10:
		upData["ten_num"] = num
	case 11:
		upData["eleven_num"] = num
	case 12:
		upData["twelve_num"] = num
	case 13:
		upData["thirteen_num"] = num
	case 14:
		upData["fourteen_num"] = num
	case 15:
		upData["fifteen_num"] = num
	case 16:
		upData["sixteen_num"] = num
	case 17:
		upData["seventeen_num"] = num
	case 18:
		upData["eighteen_num"] = num
	case 19:
		upData["nineteen_num"] = num
	case 20:
		upData["twenty_num"] = num
	case 21:
		upData["twenty_one_num"] = num
	case 22:
		upData["twenty_two_num"] = num
	case 23:
		upData["twenty_three_num"] = num
	}

	dateNumber := now.Format(timeutil.DateNumberFormat)

	err = obj.UpdateByYmdAndType(db, dateNumber, gm.PerHourDataNumTypeOnline, upData)

	return

}

// 小时用户数据任务-活跃人数
func HourUserDataTaskActive() (err error) {

	db := global.DB
	logDb := global.Log

	now := time.Now()

	hour := now.Hour()

	obj := new(gm.PerHourDataNum)

	upData := make(map[string]interface{})

	zeroMinSec := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())

	// 格式化时间字符串
	start := zeroMinSec.Format(timeutil.TimeFormat)

	end := zeroMinSec.Add(1 * time.Hour).Format(timeutil.TimeFormat)

	//活跃人数
	loginLogObj := new(log.UserLoginLog)

	var num int64

	num, err = loginLogObj.CountByCreatedAt(logDb, start, end)

	if err != nil {
		TaskLog(fmt.Sprintf("统计活跃人数失败,开始时间:%v,err:%s", start, err.Error()))
		return
	}

	switch hour {
	case 0:
		upData["zero_num"] = num
	case 1:
		upData["one_num"] = num
	case 2:
		upData["two_num"] = num
	case 3:
		upData["three_num"] = num
	case 4:
		upData["four_num"] = num
	case 5:
		upData["five_num"] = num
	case 6:
		upData["six_num"] = num
	case 7:
		upData["seven_num"] = num
	case 8:
		upData["eight_num"] = num
	case 9:
		upData["nine_num"] = num
	case 10:
		upData["ten_num"] = num
	case 11:
		upData["eleven_num"] = num
	case 12:
		upData["twelve_num"] = num
	case 13:
		upData["thirteen_num"] = num
	case 14:
		upData["fourteen_num"] = num
	case 15:
		upData["fifteen_num"] = num
	case 16:
		upData["sixteen_num"] = num
	case 17:
		upData["seventeen_num"] = num
	case 18:
		upData["eighteen_num"] = num
	case 19:
		upData["nineteen_num"] = num
	case 20:
		upData["twenty_num"] = num
	case 21:
		upData["twenty_one_num"] = num
	case 22:
		upData["twenty_two_num"] = num
	case 23:
		upData["twenty_three_num"] = num
	}

	dateNumber := now.Format(timeutil.DateNumberFormat)

	err = obj.UpdateByYmdAndType(db, dateNumber, gm.PerHourDataNumTypeActive, upData)

	return

}

// 小时用户数据任务-付费人数
func HourUserDataTaskPay() (err error) {

	db := global.DB
	payDb := global.Pay

	now := time.Now()

	hour := now.Hour()

	obj := new(gm.PerHourDataNum)

	upData := make(map[string]interface{})

	zeroMinSec := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())

	// 格式化时间字符串
	start := zeroMinSec.Format(timeutil.TimeFormat)

	end := zeroMinSec.Add(1 * time.Hour).Format(timeutil.TimeFormat)

	orderObj := new(pay.Order)

	var num int

	num, _, err = orderObj.CountRechargeByTime(payDb, start, end)

	if err != nil {
		TaskLog(fmt.Sprintf("统计付费人数失败,开始时间:%v,err:%s", start, err.Error()))
		return
	}

	switch hour {
	case 0:
		upData["zero_num"] = num
	case 1:
		upData["one_num"] = num
	case 2:
		upData["two_num"] = num
	case 3:
		upData["three_num"] = num
	case 4:
		upData["four_num"] = num
	case 5:
		upData["five_num"] = num
	case 6:
		upData["six_num"] = num
	case 7:
		upData["seven_num"] = num
	case 8:
		upData["eight_num"] = num
	case 9:
		upData["nine_num"] = num
	case 10:
		upData["ten_num"] = num
	case 11:
		upData["eleven_num"] = num
	case 12:
		upData["twelve_num"] = num
	case 13:
		upData["thirteen_num"] = num
	case 14:
		upData["fourteen_num"] = num
	case 15:
		upData["fifteen_num"] = num
	case 16:
		upData["sixteen_num"] = num
	case 17:
		upData["seventeen_num"] = num
	case 18:
		upData["eighteen_num"] = num
	case 19:
		upData["nineteen_num"] = num
	case 20:
		upData["twenty_num"] = num
	case 21:
		upData["twenty_one_num"] = num
	case 22:
		upData["twenty_two_num"] = num
	case 23:
		upData["twenty_three_num"] = num
	}

	dateNumber := now.Format(timeutil.DateNumberFormat)

	err = obj.UpdateByYmdAndType(db, dateNumber, gm.PerHourDataNumTypePay, upData)

	return

}

// 小时用户数据任务-付费额度
func HourUserDataTaskPayAmount() (err error) {

	db := global.DB
	payDb := global.Pay

	now := time.Now()

	hour := now.Hour()

	obj := new(gm.PerHourDataNum)

	upData := make(map[string]interface{})

	zeroMinSec := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())

	// 格式化时间字符串
	start := zeroMinSec.Format(timeutil.TimeFormat)

	end := zeroMinSec.Add(1 * time.Hour).Format(timeutil.TimeFormat)

	orderObj := new(pay.Order)

	var num int

	_, num, err = orderObj.CountRechargeByTime(payDb, start, end)

	if err != nil {
		TaskLog(fmt.Sprintf("统计付费额度失败,开始时间:%v,err:%s", start, err.Error()))
		return
	}

	switch hour {
	case 0:
		upData["zero_num"] = num
	case 1:
		upData["one_num"] = num
	case 2:
		upData["two_num"] = num
	case 3:
		upData["three_num"] = num
	case 4:
		upData["four_num"] = num
	case 5:
		upData["five_num"] = num
	case 6:
		upData["six_num"] = num
	case 7:
		upData["seven_num"] = num
	case 8:
		upData["eight_num"] = num
	case 9:
		upData["nine_num"] = num
	case 10:
		upData["ten_num"] = num
	case 11:
		upData["eleven_num"] = num
	case 12:
		upData["twelve_num"] = num
	case 13:
		upData["thirteen_num"] = num
	case 14:
		upData["fourteen_num"] = num
	case 15:
		upData["fifteen_num"] = num
	case 16:
		upData["sixteen_num"] = num
	case 17:
		upData["seventeen_num"] = num
	case 18:
		upData["eighteen_num"] = num
	case 19:
		upData["nineteen_num"] = num
	case 20:
		upData["twenty_num"] = num
	case 21:
		upData["twenty_one_num"] = num
	case 22:
		upData["twenty_two_num"] = num
	case 23:
		upData["twenty_three_num"] = num
	}

	dateNumber := now.Format(timeutil.DateNumberFormat)

	err = obj.UpdateByYmdAndType(db, dateNumber, gm.PerHourDataNumTypePayAmount, upData)

	return

}

// 小时用户数据任务-提现额度
func HourUserDataTaskWithdraw() (err error) {

	db := global.DB
	payDb := global.Pay

	now := time.Now()

	hour := now.Hour()

	obj := new(gm.PerHourDataNum)

	upData := make(map[string]interface{})

	zeroMinSec := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())

	// 格式化时间字符串
	start := zeroMinSec.Format(timeutil.TimeFormat)

	end := zeroMinSec.Add(1 * time.Hour).Format(timeutil.TimeFormat)

	giveObj := new(pay.GiveMoney)

	var num int

	num, err = giveObj.SumAmountByCreatedAt(payDb, start, end)

	if err != nil {
		TaskLog(fmt.Sprintf("统计提现额度失败,开始时间:%v,err:%s", start, err.Error()))
		return
	}

	switch hour {
	case 0:
		upData["zero_num"] = num
	case 1:
		upData["one_num"] = num
	case 2:
		upData["two_num"] = num
	case 3:
		upData["three_num"] = num
	case 4:
		upData["four_num"] = num
	case 5:
		upData["five_num"] = num
	case 6:
		upData["six_num"] = num
	case 7:
		upData["seven_num"] = num
	case 8:
		upData["eight_num"] = num
	case 9:
		upData["nine_num"] = num
	case 10:
		upData["ten_num"] = num
	case 11:
		upData["eleven_num"] = num
	case 12:
		upData["twelve_num"] = num
	case 13:
		upData["thirteen_num"] = num
	case 14:
		upData["fourteen_num"] = num
	case 15:
		upData["fifteen_num"] = num
	case 16:
		upData["sixteen_num"] = num
	case 17:
		upData["seventeen_num"] = num
	case 18:
		upData["eighteen_num"] = num
	case 19:
		upData["nineteen_num"] = num
	case 20:
		upData["twenty_num"] = num
	case 21:
		upData["twenty_one_num"] = num
	case 22:
		upData["twenty_two_num"] = num
	case 23:
		upData["twenty_three_num"] = num
	}

	dateNumber := now.Format(timeutil.DateNumberFormat)

	err = obj.UpdateByYmdAndType(db, dateNumber, gm.PerHourDataNumTypeWithdraw, upData)

	return

}

// 小时游戏用户数据任务
func HourGameUserDataTask() (err error) {
	db := global.DB
	now := time.Now()
	dateNumber := now.Format(timeutil.DateNumberFormat)
	hour := now.Hour()
	key := constant.HourGameUserData + dateNumber + strconv.Itoa(hour)
	val, ok := global.GoCache.Get(key)

	var mp map[uint64]struct{}

	if !ok {
		mp = make(map[uint64]struct{})
	} else {
		mp = val.(map[uint64]struct{})
	}

	userIds := msgcenter.GetGameUsers()

	for _, id := range userIds {
		mp[id] = struct{}{}
	}

	global.GoCache.Set(key, mp, 60*2*time.Minute)

	if now.Minute() >= 55 {
		//
		obj := new(gm.PerHourGameNum)

		upData := make(map[string]interface{})

		num := len(mp)

		switch hour {
		case 0:
			upData["zero_num"] = num
		case 1:
			upData["one_num"] = num
		case 2:
			upData["two_num"] = num
		case 3:
			upData["three_num"] = num
		case 4:
			upData["four_num"] = num
		case 5:
			upData["five_num"] = num
		case 6:
			upData["six_num"] = num
		case 7:
			upData["seven_num"] = num
		case 8:
			upData["eight_num"] = num
		case 9:
			upData["nine_num"] = num
		case 10:
			upData["ten_num"] = num
		case 11:
			upData["eleven_num"] = num
		case 12:
			upData["twelve_num"] = num
		case 13:
			upData["thirteen_num"] = num
		case 14:
			upData["fourteen_num"] = num
		case 15:
			upData["fifteen_num"] = num
		case 16:
			upData["sixteen_num"] = num
		case 17:
			upData["seventeen_num"] = num
		case 18:
			upData["eighteen_num"] = num
		case 19:
			upData["nineteen_num"] = num
		case 20:
			upData["twenty_num"] = num
		case 21:
			upData["twenty_one_num"] = num
		case 22:
			upData["twenty_two_num"] = num
		case 23:
			upData["twenty_three_num"] = num
		}

		err = obj.UpdateByYmdAndType(db, dateNumber, 2, upData)

		return
	}

	return
}

// 次留任务
func NextDayRemainedTask(now time.Time) (err error) {
	gmDb := global.DB

	/**
		eg today 0607 now = 0607
		统计 0605 的次留 0606的次留暂时无法统计，07过完后才有07的全部登录信息
		twoDaysAgo = 0605
	**/

	twoDaysAgo := now.AddDate(0, 0, -2)

	mp := make(map[int]int)
	//统计
	mp, err = NextDay(twoDaysAgo, 100)

	if err != nil {
		TaskLog(fmt.Sprintf("次留任务 err:%s", err.Error()))
		return
	}

	reportObj := new(gm.ReportData)
	// 查出 prevYmd 的报表数据 计算次留  = mp count / reportData.AddNums

	dataList := make([]gm.ReportData, 0, len(mp))

	//查询今天前两天的reportData数据
	dataList, err = reportObj.GetListByYmd(gmDb, twoDaysAgo.Format(timeutil.DateNumberFormat))

	if err != nil {
		TaskLog(fmt.Sprintf("次留任务 err:%s", err.Error()))
		return
	}

	for i, d := range dataList {
		count, ok := mp[d.Channel]

		if !ok {
			continue
		}

		if d.AddNums > 0 {
			dataList[i].NextDayRetention = float64(count) / float64(d.AddNums)
		}

		dataList[i].NextDayPeople = count
	}

	err = reportObj.CreateInBatchesNextDay(gmDb, dataList)

	if err != nil {
		TaskLog(fmt.Sprintf("次留任务 err:%s", err.Error()))
		return
	}

	return
}

// 3留任务
func ThreeDayRemainedTask(now time.Time) (err error) {
	gmDb := global.DB

	//今天的前2天，3留
	prevYmd := now.AddDate(0, 0, -3)

	mp := make(map[int]int)

	mp, err = ThreeDay(prevYmd, 100)

	if err != nil {
		TaskLog(fmt.Sprintf("3留任务 err:%s", err.Error()))
		return
	}

	reportObj := new(gm.ReportData)
	// 查出 prevYmd 的报表数据 计算3留  = mp count / reportData.AddNums

	dataList := make([]gm.ReportData, 0, len(mp))

	dataList, err = reportObj.GetListByYmd(gmDb, prevYmd.Format(timeutil.DateNumberFormat))

	if err != nil {
		TaskLog(fmt.Sprintf("3留任务 err:%s", err.Error()))
		return
	}

	for i, d := range dataList {
		count, ok := mp[d.Channel]

		if !ok {
			continue
		}

		if d.AddNums > 0 {
			dataList[i].ThreeDayRetention = float64(count) / float64(d.AddNums)
		}

		dataList[i].ThreeDayPeople = count
	}

	err = reportObj.CreateInBatchesThreeDay(gmDb, dataList)

	if err != nil {
		TaskLog(fmt.Sprintf("3留任务 err:%s", err.Error()))
		return
	}

	return
}

// 4留任务
func FourDayRemainedTask(now time.Time) (err error) {
	gmDb := global.DB

	//今天的前3天，4留
	prevYmd := now.AddDate(0, 0, -4)

	mp := make(map[int]int)

	mp, err = FourDay(prevYmd, 100)

	if err != nil {
		TaskLog(fmt.Sprintf("4留任务 err:%s", err.Error()))
		return
	}

	reportObj := new(gm.ReportData)
	// 查出 prevYmd 的报表数据 计算4留  = mp count / reportData.AddNums

	dataList := make([]gm.ReportData, 0, len(mp))

	dataList, err = reportObj.GetListByYmd(gmDb, prevYmd.Format(timeutil.DateNumberFormat))

	if err != nil {
		TaskLog(fmt.Sprintf("4留任务 err:%s", err.Error()))
		return
	}

	for i, d := range dataList {
		count, ok := mp[d.Channel]

		if !ok {
			continue
		}

		if d.AddNums > 0 {
			dataList[i].FourDayRetention = float64(count) / float64(d.AddNums)
		}

		dataList[i].FourDayPeople = count
	}

	err = reportObj.CreateInBatchesFourDay(gmDb, dataList)

	if err != nil {
		TaskLog(fmt.Sprintf("4留任务 err:%s", err.Error()))
		return
	}

	return
}

// 5留任务
func FiveDayRemainedTask(now time.Time) (err error) {
	gmDb := global.DB

	//今天的前4天，5留
	prevYmd := now.AddDate(0, 0, -5)

	mp := make(map[int]int)

	mp, err = FiveDay(prevYmd, 100)

	if err != nil {
		TaskLog(fmt.Sprintf("5留任务 err:%s", err.Error()))
		return
	}

	reportObj := new(gm.ReportData)
	// 查出 prevYmd 的报表数据 计算5留  = mp count / reportData.AddNums

	dataList := make([]gm.ReportData, 0, len(mp))

	dataList, err = reportObj.GetListByYmd(gmDb, prevYmd.Format(timeutil.DateNumberFormat))

	if err != nil {
		TaskLog(fmt.Sprintf("5留任务 err:%s", err.Error()))
		return
	}

	for i, d := range dataList {
		count, ok := mp[d.Channel]

		if !ok {
			continue
		}

		if d.AddNums > 0 {
			dataList[i].FiveDayRetention = float64(count) / float64(d.AddNums)
		}

		dataList[i].FiveDayPeople = count
	}

	err = reportObj.CreateInBatchesFiveDay(gmDb, dataList)

	if err != nil {
		TaskLog(fmt.Sprintf("5留任务 err:%s", err.Error()))
		return
	}

	return
}

// 6留任务
func SixDayRemainedTask(now time.Time) (err error) {
	gmDb := global.DB

	//今天的前5天，6留
	prevYmd := now.AddDate(0, 0, -6)

	mp := make(map[int]int)

	mp, err = SixDay(prevYmd, 100)

	if err != nil {
		TaskLog(fmt.Sprintf("6留任务 err:%s", err.Error()))
		return
	}

	reportObj := new(gm.ReportData)
	// 查出 prevYmd 的报表数据 计算6留  = mp count / reportData.AddNums

	dataList := make([]gm.ReportData, 0, len(mp))

	dataList, err = reportObj.GetListByYmd(gmDb, prevYmd.Format(timeutil.DateNumberFormat))

	if err != nil {
		TaskLog(fmt.Sprintf("6留任务 err:%s", err.Error()))
		return
	}

	for i, d := range dataList {
		count, ok := mp[d.Channel]

		if !ok {
			continue
		}

		if d.AddNums > 0 {
			dataList[i].SixDayRetention = float64(count) / float64(d.AddNums)
		}

		dataList[i].SixDayPeople = count
	}

	err = reportObj.CreateInBatchesSixDay(gmDb, dataList)

	if err != nil {
		TaskLog(fmt.Sprintf("6留任务 err:%s", err.Error()))
		return
	}

	return
}

// 7留任务
func SevenDayRemainedTask(now time.Time) (err error) {
	gmDb := global.DB

	//今天的前6天，7留
	prevYmd := now.AddDate(0, 0, -7)

	mp := make(map[int]int)

	mp, err = SevenDay(prevYmd, 100)

	if err != nil {
		TaskLog(fmt.Sprintf("7留任务 err:%s", err.Error()))
		return
	}

	reportObj := new(gm.ReportData)
	// 查出 prevYmd 的报表数据 计算7留  = mp count / reportData.AddNums

	dataList := make([]gm.ReportData, 0, len(mp))

	dataList, err = reportObj.GetListByYmd(gmDb, prevYmd.Format(timeutil.DateNumberFormat))

	if err != nil {
		TaskLog(fmt.Sprintf("7留任务 err:%s", err.Error()))
		return
	}

	for i, d := range dataList {
		count, ok := mp[d.Channel]

		if !ok {
			continue
		}
		if d.AddNums > 0 {
			dataList[i].SevenDayRetention = float64(count) / float64(d.AddNums)
		}
		dataList[i].SevenDayPeople = count
	}

	err = reportObj.CreateInBatchesSevenDay(gmDb, dataList)

	if err != nil {
		TaskLog(fmt.Sprintf("7留任务 err:%s", err.Error()))
		return
	}

	return
}

// 14留任务
func FourteenDayRemainedTask(now time.Time) (err error) {
	gmDb := global.DB

	//今天的前13天，14留
	prevYmd := now.AddDate(0, 0, -14)

	mp := make(map[int]int)

	mp, err = FourteenDay(prevYmd, 100)

	if err != nil {
		TaskLog(fmt.Sprintf("14留任务 err:%s", err.Error()))
		return
	}

	reportObj := new(gm.ReportData)
	// 查出 prevYmd 的报表数据 计算14留  = mp count / reportData.AddNums

	dataList := make([]gm.ReportData, 0, len(mp))

	dataList, err = reportObj.GetListByYmd(gmDb, prevYmd.Format(timeutil.DateNumberFormat))

	if err != nil {
		TaskLog(fmt.Sprintf("14留任务 err:%s", err.Error()))
		return
	}

	for i, d := range dataList {
		count, ok := mp[d.Channel]

		if !ok {
			continue
		}
		if d.AddNums > 0 {
			dataList[i].FourteenDayRetention = float64(count) / float64(d.AddNums)
		}

		dataList[i].FourteenDayPeople = count
	}

	err = reportObj.CreateInBatchesFourteenDay(gmDb, dataList)

	if err != nil {
		TaskLog(fmt.Sprintf("14留任务 err:%s", err.Error()))
		return
	}

	return
}
