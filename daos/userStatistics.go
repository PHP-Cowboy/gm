package daos

import (
	"gm/global"
	"gm/model/gm"
	"gm/request"
	"gm/response"
	"gm/utils/slice"
	"gm/utils/timeutil"
	"sort"
	"time"
)

func UserStatistics(req request.UserStatistics) (res response.ReportRsp, err error) {
	obj := gm.ReportData{
		Ymd: req.Ymd,
	}

	if len(req.BetweenDate) > 0 {
		req.Start, req.End, err = timeutil.DateRangeToZeroAndLastTimeFormat(req.BetweenDate[0], req.BetweenDate[1], timeutil.DateNumberFormat)

		if err != nil {
			global.Logger["err"].Errorf("UserStatistics timeutil.DateRangeToZeroAndLastTime failed,err:[%v]", err.Error())
			return
		}
	}

	db := global.DB

	var (
		dataList    []gm.ReportData
		total       int64
		mergeDataMp = make(map[string]gm.ReportData, len(dataList))
	)

	req.Has = slice.Contains(req.Channel, 0)

	total, dataList, err = obj.GetPageList(db, req)

	if err != nil {
		global.Logger["err"].Errorf("UserStatistics obj.GetPageList failed,err:[%v]", err.Error())
		return
	}

	//数据按天合并
	if req.Has || len(req.Channel) > 1 {

		for _, dl := range dataList {
			mergeData, ok := mergeDataMp[dl.Ymd]

			if !ok {
				mergeData = gm.ReportData{}
			}

			mergeData.Ymd = dl.Ymd
			if len(req.Channel) > 1 {
				mergeData.ChannelName = "merge channel"
			} else {
				mergeData.ChannelName = "all channel"
			}

			mergeData.AddNums += dl.AddNums
			mergeData.EffectiveAddNums += dl.EffectiveAddNums
			mergeData.DailyActiveUser += dl.DailyActiveUser
			mergeData.Channel = 0 //合并数据的channel设为0
			mergeData.NewAddUserRechargeCoinNum += dl.NewAddUserRechargeCoinNum
			//NewAddUserRechargeRate
			mergeData.NewAddUserRechargeAmount += dl.NewAddUserRechargeAmount
			mergeData.NewAddUserRechargePeople += dl.NewAddUserRechargePeople
			//mergeData.NewAddUserAverageRevenuePerUser
			//mergeData.NewAddUserAverageRevenuePerPayingUser
			mergeData.DailyActiveUserRechargeAmount += dl.DailyActiveUserRechargeAmount
			//DailyActiveUserAverageRevenuePerUser
			//DailyActiveUserAverageRevenuePerPayingUser
			mergeData.OldUserRechargePeopleNum += dl.OldUserRechargePeopleNum
			//OldUserRechargeRate
			//PaidUserRetentionRate
			mergeData.GiveMoneyPeople += dl.GiveMoneyPeople
			mergeData.GiveMoneyAmount += dl.GiveMoneyAmount
			//GiveMoneyRate
			//PlayRate
			//NextDayRetention

			mergeData.NextDayPeople += dl.NextDayPeople
			mergeData.ThreeDayPeople += dl.ThreeDayPeople
			mergeData.FourDayPeople += dl.FourDayPeople
			mergeData.FiveDayPeople += dl.FiveDayPeople
			mergeData.SixDayPeople += dl.SixDayPeople
			mergeData.SevenDayPeople += dl.SevenDayPeople
			mergeData.FourteenDayPeople += dl.FourteenDayPeople

			mergeDataMp[dl.Ymd] = mergeData
		}
	}

	list := make(response.ReportDataDescByYmd, 0, len(dataList))

	if req.Has || len(req.Channel) > 1 {
		for _, md := range mergeDataMp {
			if md.AddNums > 0 {
				md.NewAddUserAverageRevenuePerUser = md.NewAddUserRechargeAmount / float64(md.AddNums) // 新用户Arpu
				md.NewAddUserRechargeRate = float64(md.NewAddUserRechargePeople) / float64(md.AddNums) //新用户充值金币率
				md.PlayRate = float64(md.EffectiveAddNums) / float64(md.AddNums)                       //新用户玩牌率 对局用户 / 今日新增

				md.NextDayRetention = float64(md.NextDayPeople) / float64(md.AddNums)
				md.ThreeDayRetention = float64(md.ThreeDayPeople) / float64(md.AddNums)
				md.FourDayRetention = float64(md.FourDayPeople) / float64(md.AddNums)
				md.FiveDayRetention = float64(md.FiveDayPeople) / float64(md.AddNums)
				md.SixDayRetention = float64(md.SixDayPeople) / float64(md.AddNums)
				md.SevenDayRetention = float64(md.SevenDayPeople) / float64(md.AddNums)
				md.FourteenDayRetention = float64(md.FourteenDayPeople) / float64(md.AddNums)
			}

			//ARPU = 每用户平均收入 ARPU 是 渠道今日总充值钱数 / 渠道今日DAU数
			//ARPPU = 每付费用户平均收益 渠道今日总充值钱数 / 渠道今日充值人数

			if md.DailyActiveUser > 0 {
				//DAU ARPU = 渠道当日总付费 / 渠道当日总活跃人数 ；
				md.DailyActiveUserAverageRevenuePerUser = float64(md.DailyActiveUserRechargeAmount) / float64(md.DailyActiveUser) // DAU Arpu
				//tmp.PlayRate = float64(tmp.EffectiveAddNums) / float64(tmp.DailyActiveUser) todo 后期改成活跃玩牌率 EffectiveAddNums 这个要换成活跃用户的玩牌人数
			}

			if md.DailyActiveUserRechargeAmount > 0 {
				md.GiveMoneyRate = float64(md.GiveMoneyAmount) / float64(md.DailyActiveUserRechargeAmount) // 赠送金币率 = 赠送金币额度 / DAU充值总额
			}

			if md.NewAddUserRechargePeople+md.OldUserRechargePeopleNum > 0 {
				md.NewAddUserAverageRevenuePerPayingUser = md.NewAddUserRechargeAmount / float64(md.NewAddUserRechargePeople+md.OldUserRechargePeopleNum) //新用户Arppu
				//DAU ARPPU = 渠道当日总付费 / 渠道当日总付费人数=
				md.DailyActiveUserAverageRevenuePerPayingUser = float64(md.DailyActiveUserRechargeAmount) / float64((md.NewAddUserRechargePeople + md.OldUserRechargePeopleNum)) // DAU Arppu
			}

			//日活 - 新增用户 = 老用户数
			if md.DailyActiveUser-md.AddNums > 0 {
				md.OldUserRechargeRate = float64(md.OldUserRechargePeopleNum) / float64(md.DailyActiveUser-md.AddNums) // 老用户充值金币率 = 老用户充值人数 / 老用户数
			}

			if md.YesterdayPayNum > 0 {
				md.PaidUserRetentionRate = float64(md.YesterdayPayTodayLoginNum) / float64(md.YesterdayPayNum)
			}

			list = append(list, response.ReportData{
				Ymd:                                        md.Ymd,
				ChannelName:                                md.ChannelName,
				AddNums:                                    md.AddNums,
				EffectiveAddNums:                           md.EffectiveAddNums,
				DailyActiveUser:                            md.DailyActiveUser,
				Channel:                                    md.Channel,
				NewAddUserRechargeCoinNum:                  md.NewAddUserRechargeCoinNum,
				NewAddUserRechargeRate:                     md.NewAddUserRechargeRate,
				NewAddUserRechargeAmount:                   md.NewAddUserRechargeAmount,
				NewAddUserRechargePeople:                   md.NewAddUserRechargePeople,
				NewAddUserAverageRevenuePerUser:            md.NewAddUserAverageRevenuePerUser,
				NewAddUserAverageRevenuePerPayingUser:      md.NewAddUserAverageRevenuePerPayingUser,
				DailyActiveUserRechargeAmount:              md.DailyActiveUserRechargeAmount,
				DailyActiveUserAverageRevenuePerUser:       md.DailyActiveUserAverageRevenuePerUser,
				DailyActiveUserAverageRevenuePerPayingUser: md.DailyActiveUserAverageRevenuePerPayingUser,
				OldUserRechargePeopleNum:                   md.OldUserRechargePeopleNum,
				OldUserRechargeRate:                        md.OldUserRechargeRate,
				PaidUserRetentionRate:                      md.PaidUserRetentionRate,
				GiveMoneyPeople:                            md.GiveMoneyPeople,
				GiveMoneyAmount:                            md.GiveMoneyAmount,
				GiveMoneyRate:                              md.GiveMoneyRate,
				PlayRate:                                   md.PlayRate,
				NextDayRetention:                           md.NextDayRetention,
				ThreeDayRetention:                          md.ThreeDayRetention,
				FourDayRetention:                           md.FourDayRetention,
				FiveDayRetention:                           md.FiveDayRetention,
				SixDayRetention:                            md.SixDayRetention,
				SevenDayRetention:                          md.SevenDayRetention,
				FourteenDayRetention:                       md.FourteenDayRetention,
				NextDayPeople:                              md.NextDayPeople,
				ThreeDayPeople:                             md.ThreeDayPeople,
				FourDayPeople:                              md.FourDayPeople,
				FiveDayPeople:                              md.FiveDayPeople,
				SixDayPeople:                               md.SixDayPeople,
				SevenDayPeople:                             md.SevenDayPeople,
				FourteenDayPeople:                          md.FourteenDayPeople,
			})
		}

		sort.Sort(list)
	} else {
		for _, rl := range dataList {

			tmp := response.ReportData{
				Ymd:                                        rl.Ymd,
				ChannelName:                                rl.ChannelName,
				AddNums:                                    rl.AddNums,
				EffectiveAddNums:                           rl.EffectiveAddNums,
				DailyActiveUser:                            rl.DailyActiveUser,
				Channel:                                    rl.Channel,
				NewAddUserRechargeCoinNum:                  rl.NewAddUserRechargeCoinNum,
				NewAddUserRechargeRate:                     rl.NewAddUserRechargeRate,
				NewAddUserRechargeAmount:                   rl.NewAddUserRechargeAmount,
				NewAddUserRechargePeople:                   rl.NewAddUserRechargePeople,
				NewAddUserAverageRevenuePerUser:            rl.NewAddUserAverageRevenuePerUser,
				NewAddUserAverageRevenuePerPayingUser:      rl.NewAddUserAverageRevenuePerPayingUser,
				DailyActiveUserRechargeAmount:              rl.DailyActiveUserRechargeAmount,
				DailyActiveUserAverageRevenuePerUser:       rl.DailyActiveUserAverageRevenuePerUser,
				DailyActiveUserAverageRevenuePerPayingUser: rl.DailyActiveUserAverageRevenuePerPayingUser,
				OldUserRechargePeopleNum:                   rl.OldUserRechargePeopleNum,
				OldUserRechargeRate:                        rl.OldUserRechargeRate,
				PaidUserRetentionRate:                      rl.PaidUserRetentionRate,
				GiveMoneyPeople:                            rl.GiveMoneyPeople,
				GiveMoneyAmount:                            rl.GiveMoneyAmount,
				GiveMoneyRate:                              rl.GiveMoneyRate,
				PlayRate:                                   rl.PlayRate,
				NextDayRetention:                           rl.NextDayRetention,
				ThreeDayRetention:                          rl.ThreeDayRetention,
				FourDayRetention:                           rl.FourDayRetention,
				FiveDayRetention:                           rl.FiveDayRetention,
				SixDayRetention:                            rl.SixDayRetention,
				SevenDayRetention:                          rl.SevenDayRetention,
				FourteenDayRetention:                       rl.FourteenDayRetention,
				NextDayPeople:                              rl.NextDayPeople,
				ThreeDayPeople:                             rl.ThreeDayPeople,
				FourDayPeople:                              rl.FourDayPeople,
				FiveDayPeople:                              rl.FiveDayPeople,
				SixDayPeople:                               rl.SixDayPeople,
				SevenDayPeople:                             rl.SevenDayPeople,
				FourteenDayPeople:                          rl.FourteenDayPeople,
			}

			list = append(list, tmp)

		}
	}

	res.Total = total
	res.List = list

	return
}

// 用户数据统计
//func UserStatistics(req request.UserStatistics) (res response.UserStatistics, err error) {
//	obj := gm.UserStatistics{
//		Ymd:     req.Ymd,
//		Channel: req.Channel,
//	}
//
//	db := global.DB
//
//	var rechargeList []gm.UserStatistics
//
//	rechargeList, err = obj.GetPageList(db, req)
//
//	if err != nil {
//		return
//	}
//
//	var total int64
//
//	total, err = obj.Count(db)
//
//	if err != nil {
//		return
//	}
//
//	list := make([]response.UserStatisticsList, 0, len(rechargeList))
//
//	for _, rl := range rechargeList {
//		list = append(list, response.UserStatisticsList{
//			Ymd:                         rl.Ymd,
//			ChannelName:                 rl.ChannelName,
//			AddNums:                     rl.AddNums,
//			EffectiveAddNums:            rl.EffectiveAddNums,
//			EffectiveIncreaseRate:       rl.EffectiveIncreaseRate,
//			AddPayUserNextDayRetention:  rl.AddPayUserNextDayRetention,
//			DailyActiveUser:             rl.DailyActiveUser,
//			PayingSubscribers:           rl.PayingSubscribers,
//			PayoutRate:                  rl.PayoutRate,
//			TotalPayment:                rl.TotalPayment,
//			AverageRevenuePerUser:       rl.AverageRevenuePerUser,
//			AverageRevenuePerPayingUser: rl.AverageRevenuePerPayingUser,
//			WithdrawalRate:              rl.WithdrawalRate,
//			TotalWithdrawal:             rl.TotalWithdrawal,
//			NextDayRetention:            rl.NextDayRetention,
//			ThreeDayRetention:           rl.ThreeDayRetention,
//			SevenDayRetention:           rl.SevenDayRetention,
//			FourteenDayRetention:        rl.FourteenDayRetention,
//			ThirtyDayRetention:          rl.ThirtyDayRetention,
//			NinetyDayRetention:          rl.NinetyDayRetention,
//			NewPayingSubscribers:        rl.NewPayingSubscribers,
//			NewDeviceNum:                rl.NewDeviceNum,
//			NewPayingMoney:              rl.NewPayingMoney,
//			StrongActiveNum:             rl.StrongActiveNum,
//		})
//	}
//
//	res.Total = total
//	res.List = list
//
//	return
//}

// 充值统计
func RechargeStatistics(req request.RechargeStatistics) (res response.RechargeStatistics, err error) {
	obj := gm.RechargeStatistics{
		Ymd:     req.Ymd,
		Channel: req.Channel,
	}

	db := global.DB

	var rechargeList []gm.RechargeStatistics

	rechargeList, err = obj.GetPageList(db, req)

	if err != nil {
		return
	}

	var total int64

	total, err = obj.Count(db)

	if err != nil {
		return
	}

	list := make([]response.Recharge, 0, len(rechargeList))

	for _, rl := range rechargeList {
		list = append(list, response.Recharge{
			Ymd:                            rl.Ymd,
			ChannelName:                    rl.ChannelName,
			NewUserRechargeNums:            rl.NewUserRechargeNums,
			NewUserRechargeTotal:           rl.NewUserRechargeTotal,
			NewUserRechargeRate:            rl.NewUserRechargeRate,
			AverageRevenuePerNewUser:       rl.AverageRevenuePerNewUser,
			AverageRevenuePerPayingNewUser: rl.AverageRevenuePerPayingNewUser,
			OldUserRechargeNums:            rl.OldUserRechargeNums,
			OldUserFirstRechargeNums:       rl.OldUserFirstRechargeNums,
			OldUserRechargeTotal:           rl.OldUserRechargeTotal,
			OldUserRechargeRate:            rl.OldUserRechargeRate,
			AverageRevenuePerOldUser:       rl.AverageRevenuePerOldUser,
			AverageRevenuePerPayingOldUser: rl.AverageRevenuePerPayingOldUser,
		})
	}

	res.Total = total
	res.List = list

	return
}

// 提现用户统计
func WithdrawalStatistics(req request.WithdrawalStatistics) (res response.WithdrawalStatistics, err error) {
	obj := gm.WithdrawalStatistics{
		Ymd:     req.Ymd,
		Channel: req.Channel,
	}

	db := global.DB

	var withdrawalList []gm.WithdrawalStatistics

	withdrawalList, err = obj.GetPageList(db, req)

	if err != nil {
		return
	}

	var total int64

	total, err = obj.Count(db)

	if err != nil {
		return
	}

	list := make([]response.Withdrawal, 0, len(withdrawalList))

	for _, rl := range withdrawalList {
		list = append(list, response.Withdrawal{
			Ymd:                    rl.Ymd,
			ChannelName:            rl.ChannelName,
			WithdrawalNums:         rl.WithdrawalNums,
			WithdrawalTotal:        rl.WithdrawalTotal,
			NewUserWithdrawalNums:  rl.NewUserWithdrawalNums,
			NewUserWithdrawalTotal: rl.NewUserWithdrawalTotal,
			NewUserWithdrawalRate:  rl.NewUserWithdrawalRate,
			OldUserWithdrawalNums:  rl.OldUserWithdrawalNums,
			OldUserWithdrawalTotal: rl.OldUserWithdrawalTotal,
			OldUserWithdrawalRate:  rl.OldUserWithdrawalRate,
		})
	}

	res.Total = total
	res.List = list

	return
}

// 付费用户留存
func PaidUserRetention(req request.PaidUserRetention) (res response.PaidUserRetention, err error) {
	obj := gm.PaidUserRetention{
		Ymd:     req.Ymd,
		Channel: req.Channel,
	}

	db := global.DB

	var rechargeList []gm.PaidUserRetention

	rechargeList, err = obj.GetPageList(db, req)

	if err != nil {
		return
	}

	var total int64

	total, err = obj.Count(db)

	if err != nil {
		return
	}

	list := make([]response.PaidUserRetentionList, 0, len(rechargeList))

	for _, rl := range rechargeList {
		list = append(list, response.PaidUserRetentionList{
			Ymd:                   rl.Ymd,
			ChannelName:           rl.ChannelName,
			UserNums:              rl.UserNums,
			TwoDayRetention:       rl.TwoDayRetention,
			NextDayRetention:      rl.NextDayRetention,
			ThreeDayRetention:     rl.ThreeDayRetention,
			FourDayRetention:      rl.FourDayRetention,
			FiveDayRetention:      rl.FiveDayRetention,
			SixDayRetention:       rl.SixDayRetention,
			SevenDayRetention:     rl.SevenDayRetention,
			FourteenDayRetention:  rl.FourteenDayRetention,
			TwentyOneDayRetention: rl.TwentyOneDayRetention,
			ThirtyDayRetention:    rl.ThirtyDayRetention,
			SixtyDayRetention:     rl.SixtyDayRetention,
			NinetyDayRetention:    rl.NinetyDayRetention,
		})
	}

	res.Total = total
	res.List = list

	return
}

// 用户留存
func UserRetention(req request.UserRetention) (res response.UserRetention, err error) {
	obj := gm.UserRetention{
		Ymd:     req.Ymd,
		Channel: req.Channel,
	}

	db := global.DB

	var rechargeList []gm.UserRetention

	rechargeList, err = obj.GetPageList(db, req)

	if err != nil {
		return
	}

	var total int64

	total, err = obj.Count(db)

	if err != nil {
		return
	}

	list := make([]response.UserRetentionList, 0, len(rechargeList))

	for _, rl := range rechargeList {
		list = append(list, response.UserRetentionList{
			Ymd:                   rl.Ymd,
			ChannelName:           rl.ChannelName,
			UserNums:              rl.UserNums,
			TwoDayRetention:       rl.TwoDayRetention,
			NextDayRetention:      rl.NextDayRetention,
			ThreeDayRetention:     rl.ThreeDayRetention,
			FourDayRetention:      rl.FourDayRetention,
			FiveDayRetention:      rl.FiveDayRetention,
			SixDayRetention:       rl.SixDayRetention,
			SevenDayRetention:     rl.SevenDayRetention,
			FourteenDayRetention:  rl.FourteenDayRetention,
			TwentyOneDayRetention: rl.TwentyOneDayRetention,
			ThirtyDayRetention:    rl.ThirtyDayRetention,
			SixtyDayRetention:     rl.SixtyDayRetention,
			NinetyDayRetention:    rl.NinetyDayRetention,
		})
	}

	res.Total = total
	res.List = list

	return
}

// 每小时数据统计
func PerHourDataNum(req request.PerHourDataNum) (res response.PerHourDataNum, err error) {
	obj := gm.PerHourDataNum{
		NumType: req.Type,
	}

	db := global.DB

	var perHourDataNumList []gm.PerHourDataNum

	perHourDataNumList, err = obj.GetPageList(db, req)

	if err != nil {
		return
	}

	var total int64

	total, err = obj.Count(db)

	if err != nil {
		return
	}

	list := make([]response.PerHourDataNumList, 0, len(perHourDataNumList))

	for _, rl := range perHourDataNumList {
		list = append(list, response.PerHourDataNumList{
			Ymd:            rl.Ymd,
			NumType:        rl.NumType,
			ZeroNum:        rl.ZeroNum,
			OneNum:         rl.OneNum,
			TwoNum:         rl.TwoNum,
			ThreeNum:       rl.ThreeNum,
			FourNum:        rl.FourNum,
			FiveNum:        rl.FiveNum,
			SixNum:         rl.SixNum,
			SevenNum:       rl.SevenNum,
			EightNum:       rl.EightNum,
			NineNum:        rl.NineNum,
			TenNum:         rl.TenNum,
			ElevenNum:      rl.ElevenNum,
			TwelveNum:      rl.TwelveNum,
			ThirteenNum:    rl.ThirteenNum,
			FourteenNum:    rl.FourteenNum,
			FifteenNum:     rl.FifteenNum,
			SixteenNum:     rl.SixteenNum,
			SeventeenNum:   rl.SeventeenNum,
			EighteenNum:    rl.EighteenNum,
			NineteenNum:    rl.NineteenNum,
			TwentyNum:      rl.TwentyNum,
			TwentyOneNum:   rl.TwentyOneNum,
			TwentyTwoNum:   rl.TwentyTwoNum,
			TwentyThreeNum: rl.TwentyThreeNum,
		})
	}

	res.Total = total
	res.List = list

	return
}

func PerHourGameNum(req request.PerHourGameNum) (res response.PerHourGameNum, err error) {
	obj := gm.PerHourGameNum{
		GameId: req.GameId,
		RoomId: req.RoomId,
		Chip:   req.Chip,
	}

	db := global.DB

	var perHourGameNumList []gm.PerHourGameNum

	perHourGameNumList, err = obj.GetPageList(db, req)

	if err != nil {
		return
	}

	var total int64

	total, err = obj.Count(db)

	if err != nil {
		return
	}

	list := make([]response.PerHourGameNumList, 0, len(perHourGameNumList))

	for _, rl := range perHourGameNumList {
		list = append(list, response.PerHourGameNumList{
			Ymd:            rl.Ymd,
			NumType:        rl.NumType,
			ZeroNum:        rl.ZeroNum,
			OneNum:         rl.OneNum,
			TwoNum:         rl.TwoNum,
			ThreeNum:       rl.ThreeNum,
			FourNum:        rl.FourNum,
			FiveNum:        rl.FiveNum,
			SixNum:         rl.SixNum,
			SevenNum:       rl.SevenNum,
			EightNum:       rl.EightNum,
			NineNum:        rl.NineNum,
			TenNum:         rl.TenNum,
			ElevenNum:      rl.ElevenNum,
			TwelveNum:      rl.TwelveNum,
			ThirteenNum:    rl.ThirteenNum,
			FourteenNum:    rl.FourteenNum,
			FifteenNum:     rl.FifteenNum,
			SixteenNum:     rl.SixteenNum,
			SeventeenNum:   rl.SeventeenNum,
			EighteenNum:    rl.EighteenNum,
			NineteenNum:    rl.NineteenNum,
			TwentyNum:      rl.TwentyNum,
			TwentyOneNum:   rl.TwentyOneNum,
			TwentyTwoNum:   rl.TwentyTwoNum,
			TwentyThreeNum: rl.TwentyThreeNum,
		})
	}

	res.Total = total
	res.List = list

	return
}

func FiveMinuteData(req request.FiveMinuteData) (res map[string]response.FiveMinuteData, err error) {
	gmDb := global.DB

	req.Has = slice.Contains(req.Channel, 0)

	now := time.Now().Truncate(time.Minute) // 获取当前时间，截断到分钟级别

	remainder := now.Minute() % 5

	t := now.Add(time.Duration(remainder) * time.Minute * -1)
	preDay := t.AddDate(0, 0, -1)

	res = make(map[string]response.FiveMinuteData)

	res["today"] = response.FiveMinuteData{}
	res["yesterday"] = response.FiveMinuteData{}

	fiveObj := new(gm.FiveMinuteData)

	var (
		dataList []gm.FiveMinuteData
	)

	today := t.Format(timeutil.TimePointMinute)

	yesterday := preDay.Format(timeutil.TimePointMinute)

	timePoints := []string{today, yesterday}

	dataList, err = fiveObj.GetList(gmDb, timePoints, req)

	if err != nil {
		return
	}

	for _, data := range dataList {
		if data.TimePoint == today {

			todayData, _ := res["today"]

			todayData.RegNum += data.RegNum
			todayData.OnlineNum += data.OnlineNum
			todayData.ActiveNum += data.ActiveNum
			todayData.PayNum += data.PayNum
			todayData.PayAmount += data.PayAmount
			todayData.GiveAmount += data.GiveAmount

			res["today"] = todayData
		} else if data.TimePoint == yesterday {

			yesterdayData, _ := res["yesterday"]

			yesterdayData.RegNum += data.RegNum
			yesterdayData.OnlineNum += data.OnlineNum
			yesterdayData.ActiveNum += data.ActiveNum
			yesterdayData.PayNum += data.PayNum
			yesterdayData.PayAmount += data.PayAmount
			yesterdayData.GiveAmount += data.GiveAmount

			res["yesterday"] = yesterdayData
		}
	}

	return
}
