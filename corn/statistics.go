package corn

import (
	"fmt"
	"gm/daos"
	"gm/global"
	"gm/utils/timeutil"
	"time"
)

func Consumer() {
	for {
		select {
		case <-time.After(1 * time.Second):
			now := time.Now()

			hms := now.Format("15:04:05")

			if now.Minute()%5 == 0 && now.Second() == 0 {
				global.Logger["task"].Info("五分数据统计 start ,时间:%s", time.Now().Format(timeutil.TimeFormat))
				err := daos.FiveMinuteDataTask()

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("五分数据统计失败,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}
				global.Logger["task"].Info("五分数据统计 end ,时间:%s", time.Now().Format(timeutil.TimeFormat))
			}

			//每半小时统计一次今日报表数据
			if now.Minute()%30 == 0 && now.Second() == 0 {
				global.Logger["task"].Info(fmt.Sprintf("报表数据任务 start,时间:%s", time.Now().Format(timeutil.TimeFormat)))

				nextDay := now.AddDate(0, 0, 1)

				err := daos.ReportDataTask(nextDay)

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("报表数据任务,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}

				err = daos.MergeAllChannelReportDataTask(nextDay)

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("报表数据全渠道合并任务,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}

				global.Logger["task"].Info(fmt.Sprintf("报表数据任务 end,时间:%s", time.Now().Format(timeutil.TimeFormat)))

			}

			if hms == "00:00:05" {
				//todo
			}

			if hms == "00:00:10" {
				global.Logger["task"].Info(fmt.Sprintf("充值数据统计任务 start,时间:%s", time.Now().Format(timeutil.TimeFormat)))

				err := daos.RechargeStatisticsTask(now)

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("充值数据统计任务失败,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}

				global.Logger["task"].Info(fmt.Sprintf("充值数据统计任务 end,时间:%s", time.Now().Format(timeutil.TimeFormat)))
			}

			if hms == "00:00:15" {
				global.Logger["task"].Info(fmt.Sprintf("提现用户统计任务 start,时间:%s", time.Now().Format(timeutil.TimeFormat)))

				err := daos.WithdrawalStatisticsTask(now)

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("提现用户统计任务失败,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}
				global.Logger["task"].Info(fmt.Sprintf("提现用户统计任务 end,时间:%s", time.Now().Format(timeutil.TimeFormat)))
			}

			if hms == "00:00:20" {
				global.Logger["task"].Info(fmt.Sprintf("付费用户留存统计任务 start,时间:%s", time.Now().Format(timeutil.TimeFormat)))

				err := daos.PaidUserRetentionTask(now)

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("付费用户留存统计任务失败,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}
				global.Logger["task"].Info(fmt.Sprintf("付费用户留存统计任务 end,时间:%s", time.Now().Format(timeutil.TimeFormat)))
			}

			if hms == "00:00:25" {
				global.Logger["task"].Info(fmt.Sprintf("用户留存统计任务 start,时间:%s", time.Now().Format(timeutil.TimeFormat)))

				err := daos.UserRetentionTask(now)

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("用户留存统计任务失败,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}

				global.Logger["task"].Info(fmt.Sprintf("用户留存统计任务 end,时间:%s", time.Now().Format(timeutil.TimeFormat)))
			}

			if hms == "00:00:30" {
				global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-注册人数 start,时间:%s", time.Now().Format(timeutil.TimeFormat)))

				err := daos.HourUserDataTaskRegNum()

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-注册人数失败,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}

				global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-注册人数 end,时间:%s", time.Now().Format(timeutil.TimeFormat)))

			}

			if hms == "00:00:35" {
				global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-实时在线人数 start,时间:%s", time.Now().Format(timeutil.TimeFormat)))

				err := daos.HourUserDataTaskOnline()

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-实时在线人数失败,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}

				global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-实时在线人数 end,时间:%s", time.Now().Format(timeutil.TimeFormat)))

			}

			if hms == "00:00:40" {
				global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-活跃人数 start,时间:%s", time.Now().Format(timeutil.TimeFormat)))

				err := daos.HourUserDataTaskActive()

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-活跃人数失败,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}

				global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-活跃人数 end,时间:%s", time.Now().Format(timeutil.TimeFormat)))

			}

			if hms == "00:00:45" {
				global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-活跃人数 start,时间:%s", time.Now().Format(timeutil.TimeFormat)))

				err := daos.HourUserDataTaskPay()

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务失败-付费人数,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}

				global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-活跃人数 end,时间:%s", time.Now().Format(timeutil.TimeFormat)))

			}

			if hms == "00:00:50" {
				global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-付费额度 start,时间:%s", time.Now().Format(timeutil.TimeFormat)))

				err := daos.HourUserDataTaskPayAmount()

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-付费额度失败,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}

				global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-付费额度 end,时间:%s", time.Now().Format(timeutil.TimeFormat)))

			}

			if hms == "00:00:55" {
				global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-提现额度 start,时间:%s", time.Now().Format(timeutil.TimeFormat)))

				err := daos.HourUserDataTaskWithdraw()

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-提现额度失败,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
				}

				global.Logger["task"].Info(fmt.Sprintf("小时用户数据任务-提现额度 end,时间:%s", time.Now().Format(timeutil.TimeFormat)))
			}

			if hms == "00:01:00" {
				global.Logger["task"].Info(fmt.Sprintf("小时游戏用户数据任务 start,时间:%s", time.Now().Format(timeutil.TimeFormat)))

				err := daos.HourGameUserDataTask()

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("小时游戏用户数据任务失败,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
				}

				global.Logger["task"].Info(fmt.Sprintf("小时游戏用户数据任务 end,时间:%s", time.Now().Format(timeutil.TimeFormat)))
			}

			if hms == "00:01:05" {
				global.Logger["task"].Info("报表数据任务 start ")
				err := daos.ReportDataTask(now)

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("报表数据任务,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}

				global.Logger["task"].Info("报表数据任务 end next start MergeAllChannelReportDataTask")

				global.Logger["task"].Info("报表数据全渠道合并任务 start ")

				err = daos.MergeAllChannelReportDataTask(now)

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("报表数据全渠道合并任务,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}

				global.Logger["task"].Info("报表数据全渠道合并任务 end ")
			}

			if hms == "00:01:10" {
				global.Logger["task"].Info("次留任务 start ")
				err := daos.NextDayRemainedTask(now)

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("次留任务,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}
				global.Logger["task"].Info("次留任务 end ")
			}

			if hms == "00:01:15" {
				global.Logger["task"].Info("3留任务 start ")
				err := daos.ThreeDayRemainedTask(now)

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("3留任务,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}
				global.Logger["task"].Info("3留任务 end ")
			}

			if hms == "00:01:20" {
				global.Logger["task"].Info("4留任务 start ")
				err := daos.FourDayRemainedTask(now)

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("4留任务,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}
				global.Logger["task"].Info("4留任务 end ")
			}

			if hms == "00:01:25" {
				global.Logger["task"].Info("5留任务 start ")
				err := daos.FiveDayRemainedTask(now)

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("5留任务,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}
				global.Logger["task"].Info("5留任务 end ")
			}

			if hms == "00:01:30" {
				global.Logger["task"].Info("6留任务 start ")
				err := daos.SixDayRemainedTask(now)

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("6留任务,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}
				global.Logger["task"].Info("6留任务 end ")
			}

			if hms == "00:01:35" {
				global.Logger["task"].Info("7留任务 start ")
				err := daos.SevenDayRemainedTask(now)

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("7留任务,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}
				global.Logger["task"].Info("7留任务 end ")
			}

			if hms == "00:01:40" {
				global.Logger["task"].Info("14留任务 start ")
				err := daos.FourteenDayRemainedTask(now)

				if err != nil {
					global.Logger["task"].Info(fmt.Sprintf("14留任务,时间:%s,err:%v", time.Now().Format(timeutil.TimeFormat), err.Error()))
					return
				}
				global.Logger["task"].Info("14留任务 end ")
			}

		}
	}
}
