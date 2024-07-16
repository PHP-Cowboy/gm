package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gm/daos"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
)

// 新增用户
func AddUser(c *gin.Context) {
	var req request.AddUser
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

	err := daos.CreateUser(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

// 修改 名称 密码 状态
func ChangeUser(c *gin.Context) {
	var req request.ChangeUser

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.ChangeUser(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
	return
}

// 获取用户列表
func UserList(c *gin.Context) {
	var req request.UserList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result, err := daos.UserList(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, result)
	return
}

// 获取用户权限菜单列表
func UserRoleMenuTree(c *gin.Context) {
	uid, ok := c.Get("uid")

	if !ok {
		rsp.ErrorJSON(c, ecode.UserNotLogin)
		return
	}

	//获取用户权限菜单列表
	tree, err := daos.UserRoleMenuTree(uid.(int))

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, tree)
}

// 更新谷歌验证码
func UpdateUserGoogleCaptcha(c *gin.Context) {
	var req request.Uid

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.Uid == 0 {
		uid, ok := c.Get("uid")

		if !ok {
			rsp.ErrorJSON(c, ecode.UserNotFound)
			return
		}

		req.Uid = uid.(int)
	}

	img, err := daos.UpdateUserGoogleCaptcha(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, gin.H{"img": img})
}

func GetCaptchaQrBySecret(c *gin.Context) {
	var req request.GetCaptchaQr

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.Uid == 0 {
		uid, ok := c.Get("uid")

		if !ok {
			rsp.ErrorJSON(c, ecode.UserNotFound)
			return
		}

		req.Uid = uid.(int)
	}

	img, err := daos.GetCaptchaQrBySecret(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, gin.H{"img": img})
}

func GetCaptchaQr(c *gin.Context) {
	var req request.GetCaptchaQr

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	img, err := daos.GetCaptchaQrBySecret(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, gin.H{"img": img})
}

func CheckGoogleCaptcha(c *gin.Context) {
	var req request.CheckGoogleCaptcha

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	ok := daos.CheckGoogleCaptcha(req.Uid, req.Code)

	if !ok {
		rsp.ErrorJSON(c, errors.New("failed"))
		return
	}

	rsp.Success(c)
}

func BindGoogleCaptcha(c *gin.Context) {
	var req request.CheckGoogleCaptcha

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if req.Uid == 0 {
		uid, ok := c.Get("uid")

		if !ok {
			rsp.ErrorJSON(c, ecode.UserNotFound)
			return
		}

		req.Uid = uid.(int)
	}

	err := daos.BindGoogleCaptcha(req)

	if err != nil {
		rsp.ErrorJSON(c, errors.New("bind failed"))
		return
	}

	rsp.Success(c)
}
