package api

import (
	"github.com/gin-gonic/gin"
	"gm/daos"
	"gm/global"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
)

func ChannelPageList(c *gin.Context) {
	req := request.ChannelList{}

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	res, err := daos.ChannelPageList(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)
}

// 渠道列表
func AllChannelList(c *gin.Context) {

	res, err := daos.AllChannelList()

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)
}

// 角色渠道列表
func RoleChannelList(c *gin.Context) {
	channelIds, ok := c.Get("roleChannelIds")

	if !ok {
		rsp.ErrorJSON(c, ecode.UserRoleChannelNotFound)
		return
	}

	res, err := daos.RoleChannelList(channelIds.([]int))

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)
}

// 保存渠道数据
func SaveChannel(c *gin.Context) {
	var req request.SaveChannel

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.SaveChannel(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.Success(c)
}

func GameList(c *gin.Context) {

	res, err := daos.GameList()

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)
}

func CommonRoomList(c *gin.Context) {
	list, err := daos.CommonRoomList()
	if err != nil {
		global.Logger["err"].Errorf("CommonRoomList failed,err:[%v]", err.Error())
		return
	}

	rsp.SucJson(c, list)
}
