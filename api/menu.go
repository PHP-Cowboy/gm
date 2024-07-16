package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gm/daos"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
)

// 新增菜单
func AddMenu(c *gin.Context) {
	var req request.AddMenu
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

	err := daos.CreateMenu(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 修改菜单
func ChangeMenu(c *gin.Context) {
	var req request.ChangeMenu

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.ChangeMenu(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
	return
}

// 菜单列表
func MenuList(c *gin.Context) {
	var req request.MenuList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result, err := daos.MenuList(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, result)
	return
}

// 菜单树列表
func MenuTree(c *gin.Context) {

	result, err := daos.MenuTree()

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, result)
	return
}

// 上级菜单列表
func LevelList(c *gin.Context) {
	var req request.LevelList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result, err := daos.LevelList(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, result)
	return
}

// 角色菜单列表
func RoleMenuList(c *gin.Context) {
	var req request.RoleMenuList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result, err := daos.RoleMenuList(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, result)
	return
}

// 新增角色菜单
func AddRoleMenu(c *gin.Context) {
	var req request.AddRoleMenu

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.AddRoleMenu(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
	return
}
