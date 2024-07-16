package api

import (
	"github.com/gin-gonic/gin"
	"gm/daos"
	"gm/global"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
	"gm/utils/timeutil"
	"time"
)

// 报表数据任务
func ReportDataTask(c *gin.Context) {
	var (
		req request.DateTime
		err error
		now time.Time
	)

	if err = c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.DateTime != "" {
		now, err = time.ParseInLocation(timeutil.DateFormat, req.DateTime, time.Local)

		if err != nil {
			global.Logger["err"].Errorf("ReportDataTask time.ParseInLocation failed,err:[%v]", err.Error())
			return
		}
	} else {
		now = time.Now()
	}

	err = daos.ReportDataTask(now)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 报表数据全渠道合并任务
func MergeAllChannelReportDataTask(c *gin.Context) {
	var (
		req request.DateTime
		err error
		now time.Time
	)

	if err = c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.DateTime != "" {
		now, err = time.ParseInLocation(timeutil.DateFormat, req.DateTime, time.Local)

		if err != nil {
			global.Logger["err"].Errorf("MergeAllChannelReportDataTask time.ParseInLocation failed,err:[%v]", err.Error())
			return
		}
	} else {
		now = time.Now()
	}

	err = daos.MergeAllChannelReportDataTask(now)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 充值数据统计任务
func RechargeStatisticsTask(c *gin.Context) {
	var (
		req request.DateTime
		err error
		now time.Time
	)

	if err = c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.DateTime != "" {
		now, err = time.ParseInLocation(timeutil.DateFormat, req.DateTime, time.Local)

		if err != nil {
			global.Logger["err"].Errorf("RechargeStatisticsTask time.ParseInLocation failed,err:[%v]", err.Error())
			return
		}
	} else {
		now = time.Now()
	}

	err = daos.RechargeStatisticsTask(now)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 提现用户统计任务
func WithdrawalStatisticsTask(c *gin.Context) {
	var (
		req request.DateTime
		err error
		now time.Time
	)

	if err = c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.DateTime != "" {
		now, err = time.ParseInLocation(timeutil.DateFormat, req.DateTime, time.Local)

		if err != nil {
			global.Logger["err"].Errorf("WithdrawalStatisticsTask time.ParseInLocation failed,err:[%v]", err.Error())
			return
		}
	} else {
		now = time.Now()
	}

	err = daos.WithdrawalStatisticsTask(now)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 付费用户留存统计任务
func PaidUserRetentionTask(c *gin.Context) {
	var (
		req request.DateTime
		err error
		now time.Time
	)

	if err = c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.DateTime != "" {
		now, err = time.ParseInLocation(timeutil.DateFormat, req.DateTime, time.Local)

		if err != nil {
			global.Logger["err"].Errorf("PaidUserRetentionTask time.ParseInLocation failed,err:[%v]", err.Error())
			return
		}
	} else {
		now = time.Now()
	}

	err = daos.PaidUserRetentionTask(now)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 用户留存统计任务
func UserRetentionTask(c *gin.Context) {
	var (
		req request.DateTime
		err error
		now time.Time
	)

	if err = c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.DateTime != "" {
		now, err = time.ParseInLocation(timeutil.DateFormat, req.DateTime, time.Local)

		if err != nil {
			global.Logger["err"].Errorf("UserRetentionTask time.ParseInLocation failed,err:[%v]", err.Error())
			return
		}
	} else {
		now = time.Now()
	}

	err = daos.UserRetentionTask(now)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 次留任务
func NextDayRemainedTask(c *gin.Context) {
	var (
		req request.DateTime
		err error
		now time.Time
	)

	if err = c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.DateTime != "" {
		now, err = time.ParseInLocation(timeutil.DateFormat, req.DateTime, time.Local)

		if err != nil {
			global.Logger["err"].Errorf("NextDayRemainedTask time.ParseInLocation failed,err:[%v]", err.Error())

			return
		}
	} else {
		now = time.Now()
	}

	err = daos.NextDayRemainedTask(now)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 3留任务
func ThreeDayRemainedTask(c *gin.Context) {
	var (
		req request.DateTime
		err error
		now time.Time
	)

	if err = c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.DateTime != "" {
		now, err = time.ParseInLocation(timeutil.DateFormat, req.DateTime, time.Local)

		if err != nil {
			global.Logger["err"].Errorf("ThreeDayRemainedTask time.ParseInLocation failed,err:[%v]", err.Error())
			return
		}
	} else {
		now = time.Now()
	}

	err = daos.ThreeDayRemainedTask(now)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 4留任务
func FourDayRemainedTask(c *gin.Context) {
	var (
		req request.DateTime
		err error
		now time.Time
	)

	if err = c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.DateTime != "" {
		now, err = time.ParseInLocation(timeutil.DateFormat, req.DateTime, time.Local)

		if err != nil {
			return
		}
	} else {
		now = time.Now()
	}

	err = daos.FourDayRemainedTask(now)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 5留任务
func FiveDayRemainedTask(c *gin.Context) {
	var (
		req request.DateTime
		err error
		now time.Time
	)

	if err = c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.DateTime != "" {
		now, err = time.ParseInLocation(timeutil.DateFormat, req.DateTime, time.Local)

		if err != nil {
			return
		}
	} else {
		now = time.Now()
	}

	err = daos.FiveDayRemainedTask(now)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 6留任务
func SixDayRemainedTask(c *gin.Context) {
	var (
		req request.DateTime
		err error
		now time.Time
	)

	if err = c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.DateTime != "" {
		now, err = time.ParseInLocation(timeutil.DateFormat, req.DateTime, time.Local)

		if err != nil {
			return
		}
	} else {
		now = time.Now()
	}

	err = daos.SixDayRemainedTask(now)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 7留任务
func SevenDayRemainedTask(c *gin.Context) {
	var (
		req request.DateTime
		err error
		now time.Time
	)

	if err = c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.DateTime != "" {
		now, err = time.ParseInLocation(timeutil.DateFormat, req.DateTime, time.Local)

		if err != nil {
			return
		}
	} else {
		now = time.Now()
	}

	err = daos.SevenDayRemainedTask(now)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 14留任务
func FourteenDayRemainedTask(c *gin.Context) {
	var (
		req request.DateTime
		err error
		now time.Time
	)

	if err = c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.DateTime != "" {
		now, err = time.ParseInLocation(timeutil.DateFormat, req.DateTime, time.Local)

		if err != nil {
			return
		}
	} else {
		now = time.Now()
	}

	err = daos.FourteenDayRemainedTask(now)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 五分钟数据
func FiveMinuteDataTask(c *gin.Context) {

	err := daos.FiveMinuteDataTask()

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 小时用户数据任务-实时在线人数
func HourUserDataTaskOnline(c *gin.Context) {
	err := daos.HourUserDataTaskOnline()

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 小时游戏用户数据任务
func HourGameUserDataTask(c *gin.Context) {
	err := daos.HourGameUserDataTask()

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}
