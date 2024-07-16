package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gm/daos"
	"gm/global"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
	"net/url"
	"strings"
)

// 游戏用户列表
func GetGameUserList(c *gin.Context) {
	req := request.GetGameUserList{}

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

	res, err := daos.GetGameUserList(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)

}

// 登录日志列表
func GetLoginLogList(c *gin.Context) {
	var req request.GetLoginLogList

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

	res, err := daos.GetLoginLogList(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)
}

// 赠送记录
func GiveList(c *gin.Context) {
	req := request.GiveList{}

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

	res, err := daos.GiveList(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)
}

// 批量通过
func BatchPass(c *gin.Context) {
	req := request.CheckIds{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	name, nOk := c.Get("name")
	if !nOk {
		rsp.ErrorJSON(c, ecode.UserNotFound)
		return
	}

	err := daos.BatchPass(req, name)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 批量拒绝
func BatchRepulse(c *gin.Context) {
	req := request.CheckIds{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	name, nOk := c.Get("name")
	if !nOk {
		rsp.ErrorJSON(c, ecode.UserNotFound)
		return
	}

	err := daos.BatchRepulse(req, name)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 批量主动取消
func BatchCancel(c *gin.Context) {
	req := request.CheckIds{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.BatchCancel(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 批量作废
func BatchInvalid(c *gin.Context) {
	req := request.CheckIds{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	name, nOk := c.Get("name")
	if !nOk {
		rsp.ErrorJSON(c, ecode.UserNotFound)
		return
	}

	err := daos.BatchInvalid(req, name)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 赠送数据处理
func GiveMoneyHandle(c *gin.Context) {
	var req request.GiveMoneyHandle

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.GiveMoneyHandle(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// xPay代付回调
//func XPaymentCallback(c *gin.Context) {
//	var req request.XPaymentCallback
//
//	if err := c.ShouldBind(&req); err != nil {
//		rsp.ErrorJSON(c, ecode.ParamInvalid)
//		return
//	}
//
//	err := daos.XPaymentCallback(req)
//
//	if err != nil {
//		rsp.ErrorJSON(c, err)
//		return
//	}
//
//	rsp.String(c, "success")
//}

func XPaymentCallback(c *gin.Context) {
	var req request.XPaymentCallback

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}
	// 合并URL参数和POST参数
	// 合并URL查询参数和POST请求体中的参数
	params := make(url.Values)
	// 添加URL查询参数
	for key, values := range c.Request.URL.Query() {
		params.Add(key, strings.Join(values, "")) // 假设查询参数中没有多个同名参数
	}

	// 添加POST请求体中的参数
	if err := c.Request.ParseForm(); err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	for key, values := range c.Request.PostForm {
		params.Add(key, strings.Join(values, "")) // 假设POST参数中没有多个同名参数
	}

	err := daos.XPaymentCallbackNew(params, req)

	if err != nil {
		rsp.String(c, err.Error())
		return
	}

	rsp.String(c, "success")
}

// xPay代付回调
func XPaymentSuccessCallback(c *gin.Context) {
	var req request.XPaymentCallbackSuccess

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.XPaymentCallbackSuccess(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.String(c, "success")
}

// xPay代付回调
func XPaymentCallbackFailed(c *gin.Context) {
	var req request.PaymentCallbackFailed

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.XPaymentCallbackFailed(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.String(c, "success")
}

// inPay代付回调
func InPaymentCallback(c *gin.Context) {
	var req request.InPaymentCallback
	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	jTxt, err := json.Marshal(req)
	if err != nil {
		rsp.ErrorJSON(c, ecode.DataParsingFailure)
		return
	}

	global.Logger["info"].Info("inPay params :%s", string(jTxt))

	err = daos.InPaymentCallback(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// LuckyPay代付回调
func LuckyPaymentCallback(c *gin.Context) {
	var req request.LuckyPayCallback
	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	jTxt, err := json.Marshal(req)
	if err != nil {
		rsp.ErrorJSON(c, ecode.DataParsingFailure)
		return
	}

	global.Logger["info"].Info("LuckyPay params :%s", string(jTxt))

	err = daos.LuckyPaymentCallback(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.String(c, "success")
}

// 用户信息
func GameUserInfo(c *gin.Context) {
	req := request.UserId{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	channelIds, ok := c.Get("roleChannelIds")

	if !ok {
		rsp.ErrorJSON(c, ecode.UserRoleChannelNotFound)
		return
	}

	userInfo, err := daos.GameUserInfo(req, channelIds.([]int))

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, userInfo)
}

// 修改用户累计充值总额
func ChangeRecharge(c *gin.Context) {
	req := request.ChangeRecharge{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.ChangeRecharge(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 用户提现信息记录
func WithdrawInfoRecord(c *gin.Context) {
	req := request.WithdrawInfoRecord{}

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

	res, err := daos.WithdrawInfoRecord(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)
}

// 修改用户金币
func EditUserCoin(c *gin.Context) {
	req := request.EditUserCoin{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.EditUserCoin(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 封禁
func Banned(c *gin.Context) {
	req := request.Banned{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	admin, ok := c.Get("uid")

	if !ok {
		rsp.ErrorJSON(c, ecode.CommunalParamInvalid)
		return
	}

	req.Admin = admin.(int)

	err := daos.Banned(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 解封
func Unseal(c *gin.Context) {
	req := request.Banned{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.Unseal(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}
