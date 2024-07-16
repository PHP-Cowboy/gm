package api

import (
	"github.com/gin-gonic/gin"
	"gm/daos"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
)

// tp 游戏流水
func TP(c *gin.Context) {
	req := request.FundsFlowLogView{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	list, err := daos.TpFundFlowLog(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, list)
}

// slot 游戏流水
func SLOT(c *gin.Context) {
	req := request.SlotFundsFlowLog{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	list, err := daos.SlotFundFlowLog(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, list)
}
