package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gm/daos"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
)

// 获取二选一 配置列表
func GetEventList(c *gin.Context) {
	//获取二选一 配置列表
	list, err := daos.GetEventGiftList()

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}
	rsp.SucJson(c, list)
}

// 保存二选一配置
func SaveEventConfig(c *gin.Context) {
	req := request.SaveEventGiftConfig{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//保存二选一 配置
	err := daos.SaveEventConfig(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 删除二选一 配置
func DelEventConfig(c *gin.Context) {
	req := request.DelEventConfig{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//删除二选一 配置
	err := daos.DelEventConfig(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 二选一开关
func OnOffEvent(c *gin.Context) {
	var req request.OnOffBenefit

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//二选一开关
	err := daos.OnOffEvent(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// true 开启 falase 关闭
func EventStatus(c *gin.Context) {
	//更新救济金礼包开关状态
	status, err := daos.EventStatus()

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, status)
}

// 充200送200 配置列表
func GetRechargeGiftList(c *gin.Context) {
	list, err := daos.GetRechargeGiftList()

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}
	rsp.SucJson(c, list)
}

// 保存充200送200配置
func SaveRechargeGiftConfig(c *gin.Context) {
	req := request.SaveRechargeGiftConfig{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//保存充200送200 配置
	err := daos.SaveRechargeGiftConfig(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 删除充200送200 配置
func DelRechargeGift(c *gin.Context) {
	req := request.DelRechargeGift{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//删除充200送200 配置
	err := daos.DelRechargeGift(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 充200送200开关
func OnOffRechargeGift(c *gin.Context) {
	var req request.OnOffBenefit

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//充200送200开关
	err := daos.OnOffRechargeGift(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// true 开启 falase 关闭
func RechargeGiftStatus(c *gin.Context) {
	//更新救济金礼包开关状态
	status, err := daos.RechargeGiftStatus()

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, status)
}

// 充值礼包 配置列表
func GetRechargePackList(c *gin.Context) {
	list, err := daos.GetRechargePackList()

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}
	rsp.SucJson(c, list)
}

// 保存充值礼包 配置
func SaveRechargePack(c *gin.Context) {
	req := request.SaveRechargePack{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//保存充200送200 配置
	err := daos.SaveRechargePack(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}
	rsp.Success(c)
}

// 删除充值礼包 配置
func DelRechargePack(c *gin.Context) {
	req := request.DelRechargePack{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//删除充值礼包 配置
	err := daos.DelRechargePack(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}
	rsp.Success(c)
}

// 获取vip配置列表
func GetVipConfigList(c *gin.Context) {
	list, err := daos.GetVipConfigList()
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, list)
}

// 保存vip配置
func SaveVipConfig(c *gin.Context) {
	req := request.SaveVipConfig{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//保存vip配置
	err := daos.SaveVipConfig(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 删除vip 配置
func DelVipConfig(c *gin.Context) {
	req := request.DelVipConfig{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//删除vip 配置
	err := daos.DelVipConfig(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}
	rsp.Success(c)
}

// 救济金 配置列表
func GetBenefitList(c *gin.Context) {
	var req request.GetBenefitList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	res, err := daos.GetBenefitList(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)
}

// 救济金 配置
func SaveBenefit(c *gin.Context) {
	req := request.SaveBenefit{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//保存救济金配置
	err := daos.SaveBenefit(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}
	rsp.Success(c)
}

func OnOffBenefit(c *gin.Context) {
	var req request.OnOffBenefit

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//更新救济金礼包开关
	err := daos.OnOffBenefit(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// true 开启 falase 关闭
func BenefitStatus(c *gin.Context) {
	//更新救济金礼包开关状态
	status, err := daos.BenefitStatus()

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, status)
}

// 删除救济金 配置
func DelBenefit(c *gin.Context) {
	req := request.DelBenefit{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//删除救济金 配置
	err := daos.DelBenefit(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}
	rsp.Success(c)
}

// 三选一礼包 配置列表
func GetOnlyList(c *gin.Context) {
	list, err := daos.GetOnlyList()
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, list)
}

// 保存三选一礼包 配置
func SaveOnly(c *gin.Context) {
	req := request.SaveOnly{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//保存三选一 配置
	err := daos.SaveOnly(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}
	rsp.Success(c)
}

// 删除三选一礼包 配置
func DelOnly(c *gin.Context) {
	req := request.DelOnly{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//删除三选一  配置
	err := daos.DelOnly(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}
	rsp.Success(c)
}

// 三选一开关
func OnOffOnly(c *gin.Context) {
	var req request.OnOffBenefit

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//三选一开关
	err := daos.OnOffOnly(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// true 开启 falase 关闭
func OnlyStatus(c *gin.Context) {
	//更新救济金礼包开关状态
	status, err := daos.OnlyStatus()

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, status)
}
