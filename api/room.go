package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gm/daos"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
)

// 房间列表
func RoomList(c *gin.Context) {
	req := request.RoomList{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	res, err := daos.RoomList(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)
}

// 字段 && comment
func GetColumnComment(c *gin.Context) {
	req := request.RoomList{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	res, err := daos.GetColumnComment(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)
}

func SaveRoom(c *gin.Context) {
	var req request.SaveRoom
	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.SaveRoom(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 根据excel更新ExtData的值
func UpdateExtDataByExcel(c *gin.Context) {
	var req []request.UpdateExtDataByExcel
	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.UpdateExtDataByExcel(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}
