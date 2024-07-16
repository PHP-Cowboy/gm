package api

import (
	"github.com/gin-gonic/gin"
	"gm/daos"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
	"strconv"
)

// 在线用户列表
func OnlineList(c *gin.Context) {
	var req request.OnlineList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	uid, ok := c.Get("uid")

	if !ok {
		rsp.ErrorJSON(c, ecode.UserNotFound)
		return
	}

	channelIds, ok := c.Get("roleChannelIds")

	if !ok {
		rsp.ErrorJSON(c, ecode.UserRoleChannelNotFound)
		return
	}

	req.ChannelIds = channelIds.([]int)

	req.Admin = strconv.Itoa(uid.(int))

	res, err := daos.OnlineList(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, res)
}
