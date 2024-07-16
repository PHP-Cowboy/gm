package daos

import (
	"gm/global"
	"gm/model/log"
	"gm/model/user"
	"gm/utils/timeutil"
	"time"
)

// 次日留存 eg 0605 的次留 => 0605注册的用户 在 0606登录的人数
func NextDay(twoDaysAgo time.Time, size int) (mp map[int]int, err error) {
	return GetNumDayRemained(twoDaysAgo, 1, size)
}

// 三日留存
func ThreeDay(threeDaysAgo time.Time, size int) (mp map[int]int, err error) {
	return GetNumDayRemained(threeDaysAgo, 2, size)
}

// 四日留存
func FourDay(fourDaysAgo time.Time, size int) (mp map[int]int, err error) {
	return GetNumDayRemained(fourDaysAgo, 3, size)
}

// 五日留存
func FiveDay(fiveDaysAgo time.Time, size int) (mp map[int]int, err error) {
	return GetNumDayRemained(fiveDaysAgo, 4, size)

}

// 六日留存
func SixDay(sixDaysAgo time.Time, size int) (mp map[int]int, err error) {
	return GetNumDayRemained(sixDaysAgo, 5, size)
}

// 七日日留存
func SevenDay(sevenDaysAgo time.Time, size int) (mp map[int]int, err error) {
	return GetNumDayRemained(sevenDaysAgo, 6, size)
}

// 十四日留存
func FourteenDay(fourteenDaysAgo time.Time, size int) (mp map[int]int, err error) {
	return GetNumDayRemained(fourteenDaysAgo, 13, size)
}

// 二十一日留存
func TwentyOneDay(twentyOneDaysAgo time.Time, size int) (mp map[int]int, err error) {
	return GetNumDayRemained(twentyOneDaysAgo, 20, size)
}

// 三十日留存
func ThirtyDay(thirtyDaysAgo time.Time, size int) (mp map[int]int, err error) {
	return GetNumDayRemained(thirtyDaysAgo, 29, size)
}

// 六十日留存
func SixtyDay(sixtyDaysAgo time.Time, size int) (mp map[int]int, err error) {
	return GetNumDayRemained(sixtyDaysAgo, 59, size)
}

// 九十日留存
func NinetyDay(ninetyDaysAgo time.Time, size int) (mp map[int]int, err error) {
	return GetNumDayRemained(ninetyDaysAgo, 89, size)
}

// 统计 t 日期 注册的用户，在 num 日之后 登录的人数，按渠道分组
func GetNumDayRemained(t time.Time, num, size int) (mp map[int]int, err error) {

	//查找 num 天 之后 登录 的用户
	loginAfterNumDay := t.AddDate(0, 0, num)

	//注册用户开始时间&&结束时间
	start := timeutil.GetZeroTime(t).Format(timeutil.TimeFormat)
	end := timeutil.GetLastTime(t).Format(timeutil.TimeFormat)

	userDb := global.User
	logDb := global.Log

	userObj := new(user.User)

	loginLogObj := log.UserLoginLog{}

	mp = make(map[int]int)
	loginUserMp := make(map[int]struct{})

	page := 0

	userIds := make([]int, 0, size)

	for {
		userList := make([]user.User, 0, size)

		userList, err = userObj.GetPageListByCreateTime(userDb, start, end, page, size)

		if err != nil {
			global.Logger["err"].Errorf("GetNumDayRemained num:%v GetPageListByCreateTime failed,err:[%v]", num, err.Error())
			return
		}

		for _, u := range userList {
			userIds = append(userIds, u.Uid)
		}

		//登录日志数据处理
		loginUserMp, err = loginLogObj.GetListByUserIds(logDb, userIds, loginAfterNumDay)

		if err != nil {
			global.Logger["err"].Errorf("GetNumDayRemained GetListByUserIds failed,err:[%v]", err.Error())
			return
		}

		for _, ul := range userList {
			_, loginOk := loginUserMp[ul.Uid]

			if !loginOk {
				continue
			}

			mpCount, mpOk := mp[ul.ChannelId]

			if !mpOk {
				mpCount = 0
			}

			mpCount++

			mp[ul.ChannelId] = mpCount
		}

		page++

		//查到的数据条数小于size，则是最后一页的数据,跳出for循环
		if len(userList) < size {
			break
		}

	}

	return
}
