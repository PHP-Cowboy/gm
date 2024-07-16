package daos

import (
	"gm/daos/rds"
	"gm/global"
	"gm/model/game"
	"gm/model/user"
	"gm/request"
	"gm/response"
	"gm/utils/slice"
	"time"
	"za.game/lib/account"
	accUser "za.game/lib/account/user"
	"za.game/lib/consts"
)

func ChannelPageList(req request.ChannelList) (rsp response.ChannelList, err error) {
	obj := &user.UserChannel{
		ChannelName: req.ChannelName,
		Remark:      req.Remark,
	}

	channelList := []user.UserChannel{}

	db := global.User

	channelList, err = obj.GetPageList(db, req)

	if err != nil {
		global.Logger["err"].Errorf("ChannelPageList obj.GetPageList failed,err:[%v]", err.Error())
		return
	}

	var total int64

	total, err = obj.Count(db)

	if err != nil {
		global.Logger["err"].Errorf("ChannelPageList obj.Count failed,err:[%v]", err.Error())
		return
	}

	rsp.Total = total

	list := make([]response.Channel, 0, len(channelList))

	for _, cl := range channelList {
		list = append(list, response.Channel{
			Id:          cl.ID,
			ChannelName: cl.ChannelName,
			Code:        cl.Code,
			Remark:      cl.Remark,
		})
	}

	rsp.List = list

	return
}

// 全部渠道列表
func AllChannelList() (list []response.Channel, err error) {

	channelList := []accUser.UserChannel{}

	channelList, err = rds.GetAllChannelList()

	if err != nil {
		global.Logger["err"].Errorf("AllChannelList cache.GetAllChannelList failed,err:[%v]", err.Error())
		return
	}

	list = make([]response.Channel, 0, len(channelList))

	for _, cl := range channelList {

		list = append(list, response.Channel{
			Id:          cl.Id,
			ChannelName: cl.ChannelName,
			Code:        cl.Code,
			Remark:      cl.Remark,
		})
	}

	return
}

func RoleChannelList(channelIds []int) (list []response.Channel, err error) {

	userRoleChannelMp := slice.SliceToMap(channelIds)

	channelList := []accUser.UserChannel{}

	channelList, err = rds.GetAllChannelList()

	if err != nil {
		return
	}

	list = make([]response.Channel, 0, len(channelList))

	for _, cl := range channelList {
		//不在用户角色的渠道id map 中的 直接跳过
		_, ok := userRoleChannelMp[cl.Id]

		if !ok {
			continue
		}

		list = append(list, response.Channel{
			Id:          cl.Id,
			ChannelName: cl.ChannelName,
			Code:        cl.Code,
			Remark:      cl.Remark,
		})
	}

	return
}

func SaveChannel(req request.SaveChannel) (err error) {
	userDb := global.User

	obj := &user.UserChannel{
		ChannelName: req.ChannelName,
		Code:        req.Code,
		Remark:      req.Remark,
	}

	now := time.Now()

	if req.Id > 0 {
		obj.ID = req.Id
	} else {
		obj.CreatedAt = &now
	}

	obj.UpdatedAt = &now

	err = obj.Save(userDb)

	if err != nil {
		global.Logger["err"].Errorf("SaveChannel obj.Save failed,err:[%v]", err.Error())
		return
	}

	rds.DelRedisCacheByKey(consts.UserChannelList)

	return
}

func GameList() (list []response.GameBase, err error) {
	obj := game.GameBase{}

	dataList := []game.GameBase{}

	db := global.Game

	dataList, err = obj.GetList(db)

	if err != nil {
		return
	}

	list = make([]response.GameBase, 0, len(dataList))

	for _, cl := range dataList {
		list = append(list, response.GameBase{
			Id:     cl.ID,
			EnName: cl.EnName,
		})
	}

	list = append(
		list,
		response.GameBase{
			Id:     account.LogType_day_task + 10000,
			EnName: "每日任务",
		},
		response.GameBase{
			Id:     account.LogType_vip + 10000,
			EnName: "vip",
		},
		response.GameBase{
			Id:     account.LogType_email + 10000,
			EnName: "邮件",
		},
		response.GameBase{
			Id:     account.LogType_sign + 10000,
			EnName: "签到",
		},
		response.GameBase{
			Id:     account.LogType_recharge + 10000,
			EnName: "充值",
		},
		response.GameBase{
			Id:     account.LogType_frozen + 10000,
			EnName: "冻结",
		},
		response.GameBase{
			Id:     account.LogType_bank + 10000,
			EnName: "储钱罐",
		},
		response.GameBase{
			Id:     account.LogType_benefit + 10000,
			EnName: "救济金",
		},
		response.GameBase{
			Id:     account.LogType_giftbag + 10000,
			EnName: "礼包",
		},
		response.GameBase{
			Id:     account.LogType_register_send + 10000,
			EnName: "注册赠送",
		},
		response.GameBase{
			Id:     account.LogType_luck_spin + 10000,
			EnName: "luck_spin",
		},
		response.GameBase{
			Id:     account.LogType_withdraw + 10000,
			EnName: "赠送",
		},
		response.GameBase{
			Id:     account.LogType_admin + 10000,
			EnName: "管理后台操作",
		},
		response.GameBase{
			Id:     account.LogType_RechargeGift + 10000,
			EnName: "充值赠送礼包",
		},
		response.GameBase{
			Id:     account.LogType_EventGift + 10000,
			EnName: "二选一礼包",
		},
		response.GameBase{
			Id:     account.LogType_OnlyOneGift + 10000,
			EnName: "OnlyOne",
		},
		response.GameBase{
			Id:     account.LogType_BenefitGift + 10000,
			EnName: "救济金",
		},
		response.GameBase{
			Id:     account.LogType_RechargeRoomGift + 10000,
			EnName: "房间特惠",
		},
		response.GameBase{
			Id:     account.LogType_Activity_NewPlayer + 10000,
			EnName: "新手嘉年华",
		},
		response.GameBase{
			Id:     account.LogType_Bonus_Draw + 10000,
			EnName: "bonus 提取",
		},
		response.GameBase{
			Id:     account.LogType_Tyro_Cash + 10000,
			EnName: "新用户首次",
		},
		response.GameBase{
			Id:     account.LogType_luck_spin_recharge + 10000,
			EnName: "luckspin充值",
		},
		response.GameBase{
			Id:     account.LogType_WithdrawSend + 10000,
			EnName: "withdraw 添加首次添加账户",
		},
		response.GameBase{
			Id:     account.LogType_SurpriseGift + 10000,
			EnName: "Surprise礼包充值",
		},
		response.GameBase{
			Id:     account.LogType_Activity_NewPlayerActCoin + 10000,
			EnName: "新手嘉年华活动币提现",
		},
	)

	return
}

func CommonRoomList() (roomList []response.Room, err error) {
	db := global.Game
	rl := new(game.Roomview)

	var (
		pageList []game.Roomview
	)

	pageList, err = rl.GetList(db)
	if err != nil {
		return
	}

	roomList = make([]response.Room, 0)

	for _, l := range pageList {
		roomList = append(roomList, response.Room{
			Id:            l.Id,
			RoomId:        l.RoomId,
			SvrId:         l.SvrId,
			GameId:        l.GameId,
			RoomIndex:     l.RoomIndex,
			Base:          l.Base,
			MinEntry:      l.MinEntry,
			MaxEntry:      l.MaxEntry,
			RoomName:      l.RoomName,
			RoomType:      l.RoomType,
			RoomSwitch:    l.RoomSwitch,
			RoomWelfare:   l.RoomWelfare,
			Desc:          l.Desc,
			Tax:           l.Tax,
			BonusDiscount: l.BonusDiscount,
			AiSwitch:      l.AiSwitch,
			AiLimit:       l.AiLimit,
			RechargeLimit: l.RechargeLimit,
			ExtData:       l.ExtData,
		})
	}

	return
}
