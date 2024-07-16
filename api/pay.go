package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gm/daos"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
)

// 充值礼包列表
func PayList(c *gin.Context) {
	var req request.PayList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	list, err := daos.PayList(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, list)
}

// 保存充值礼包
func SaveGift(c *gin.Context) {
	var req request.SaveGift
	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.SaveGift(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 删除充值礼包
func DelGift(c *gin.Context) {
	var req request.DeleteId
	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.DelGift(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 支付渠道配置列表
func ConfigList(c *gin.Context) {
	var req request.ConfigList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	res, err := daos.ConfigList(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}
	rsp.SucJson(c, res)
}

// 保存支付渠道配置
func SaveConfig(c *gin.Context) {
	var req request.SaveConfig
	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.SaveConfig(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 删除支付渠道配置
func DelConfig(c *gin.Context) {
	var req request.DeleteId
	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.DelConfig(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 获取用户银行卡列表
func BankList(c *gin.Context) {
	var req request.BankList

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

	list, err := daos.BankList(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, list)
}

// 获取用户银行卡列表
func OrderList(c *gin.Context) {
	var req request.OrderList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	list, err := daos.OrderList(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, list)
}

// 充值记录
func RechargeRecords(c *gin.Context) {
	var req request.RechargeRecords

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

	res, err := daos.RechargeRecords(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)
}

// 赠送配置列表
func GaveConfigPageList(c *gin.Context) {
	var req request.GaveConfigPageList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	res, err := daos.GaveConfigPageList(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)
}

// 赠送配置列表
func GaveConfigList(c *gin.Context) {

	var req request.GaveConfigList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	res, err := daos.GaveConfigList(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)
}

// 保存赠送配置
func SaveGiveConfig(c *gin.Context) {
	var req request.SaveGaveConfig
	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.SaveGaveConfig(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 删除赠送配置
func DelGiveConfig(c *gin.Context) {
	var req request.DeleteId
	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.DelGaveConfig(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}
