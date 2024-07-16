package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gm/daos"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
)

// 登录
func Login(c *gin.Context) {
	var req request.Login

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	dto, err := daos.Login(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, dto)
}

// 登出
func LoginOut(c *gin.Context) {
	var req request.Login

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}
	rsp.Success(c)
}
