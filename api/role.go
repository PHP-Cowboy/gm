package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gm/daos"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
)

// 新增角色
func AddRole(c *gin.Context) {
	var req request.CreateRole
	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	uid, ok := c.Get("uid")

	if !ok {
		rsp.ErrorJSON(c, ecode.GetContextUserInfoFailed)
		return
	}

	req.CreatorId = uid.(int)

	name, ok := c.Get("name")

	if !ok {
		rsp.ErrorJSON(c, ecode.GetContextUserInfoFailed)
		return
	}

	req.Creator = name.(string)

	err := daos.CreateRole(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 修改角色
func ChangeRole(c *gin.Context) {
	var req request.ChangeRole

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.ChangeRole(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
	return
}

// 角色列表
func RoleList(c *gin.Context) {
	var req request.RoleList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result, err := daos.RoleList(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, result)
	return
}

// 用户角色列表
func UserRoleList(c *gin.Context) {
	var req request.UserRoleList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result, err := daos.UserRoleList(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, result)
	return
}

// 用户新增角色
func AddUserRole(c *gin.Context) {
	var req request.AddUserRole

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.AddUserRole(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
	return
}
