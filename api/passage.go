package api

import (
	"github.com/gin-gonic/gin"
	"gm/daos"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
)

// 代收代付列表
func PassageList(c *gin.Context) {
	var req request.PassageList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	res, err := daos.PassageList(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)
}

// 代收代付保存
func PassageSave(c *gin.Context) {
	var req request.PassageSave

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.PassageSave(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 代收代付删除
func PassageDel(c *gin.Context) {
	var req request.DeleteId

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.PassageDel(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

func PassageChange(c *gin.Context) {
	var req request.PassageChange

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.PassageChange(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}
