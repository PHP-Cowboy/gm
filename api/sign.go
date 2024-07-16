package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gm/daos"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
)

// 获取签到配置列表
func GetSignConfigList(c *gin.Context) {
	list, err := daos.GetSignConfigList()
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}
	rsp.SucJson(c, list)
	return
}

// 保存签到配置
func SaveSign(c *gin.Context) {
	req := request.SaveSign{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.SaveSign(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
	return
}

// 删除签到配置
func DelSign(c *gin.Context) {
	req := request.DelSign{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	//删除签到配置
	err := daos.DelSign(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
	return
}

// 获取签到奖励配置列表
func GetSingPrizeList(c *gin.Context) {
	var req request.GetSingPrizeList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	list, err := daos.GetSingPrizeList(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}
	rsp.SucJson(c, list)
	return
}

// 保存签到奖励配置
func SaveSingPrize(c *gin.Context) {
	req := request.SavePrize{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.SaveSingPrize(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
	return
}

// 删除签到奖励配置
func DelSingPrize(c *gin.Context) {
	req := request.DelSingPrize{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.DelSignPrize(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
	return
}
