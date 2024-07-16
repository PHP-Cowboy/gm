package api

import (
	"github.com/gin-gonic/gin"
	"gm/daos"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
)

// 用户数据统计
func UserStatistics(c *gin.Context) {
	req := request.UserStatistics{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.Size == 0 {
		req.Size = 20
	}

	channelIds, ok := c.Get("roleChannelIds")

	if !ok {
		rsp.ErrorJSON(c, ecode.UserRoleChannelNotFound)
		return
	}

	req.ChannelIds = channelIds.([]int)

	res, err := daos.UserStatistics(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)
}

// 充值统计
func RechargeStatistics(c *gin.Context) {
	req := request.RechargeStatistics{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//充值统计
	list, err := daos.RechargeStatistics(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, list)
}

// 提现用户统计
func WithdrawalStatistics(c *gin.Context) {
	req := request.WithdrawalStatistics{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	list, err := daos.WithdrawalStatistics(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, list)
}

// 付费用户留存
func PaidUserRetention(c *gin.Context) {
	req := request.PaidUserRetention{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	list, err := daos.PaidUserRetention(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, list)
}

// 用户留存
func UserRetention(c *gin.Context) {
	req := request.UserRetention{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	list, err := daos.UserRetention(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, list)
}

// 每小时数据统计
func PerHourDataNum(c *gin.Context) {
	req := request.PerHourDataNum{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	list, err := daos.PerHourDataNum(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, list)
}

// 每小时数据统计
func PerHourGameNum(c *gin.Context) {
	req := request.PerHourGameNum{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	list, err := daos.PerHourGameNum(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, list)
}

// 五分钟数据统计
func FiveMinuteData(c *gin.Context) {

	req := request.FiveMinuteData{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	channelIds, ok := c.Get("roleChannelIds")

	if !ok {
		rsp.ErrorJSON(c, ecode.UserRoleChannelNotFound)
		return
	}

	req.ChannelIds = channelIds.([]int)

	list, err := daos.FiveMinuteData(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, list)
}
