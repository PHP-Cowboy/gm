package daos

import (
	"encoding/json"
	"errors"
	"gm/common/constant"
	"gm/global"
	"gm/model/game"
	"gm/model/gm"
	"gm/model/user"
	"gm/msgcenter"
	"gm/request"
	"gm/response"
	"gm/utils/slice"
	"gm/utils/timeutil"
	"time"
)

// 在线用户列表
func OnlineList(req request.OnlineList) (res response.OnlineListRsp, err error) {
	key := constant.OnlineUserIds + req.Admin

	var (
		jsonTxt       []byte
		onlineUserIds []uint64
		offset        int
	)

	list := make([]response.Online, 0, req.Size)

	if req.RoomId > 0 {
		inRoomUsers := msgcenter.GetInRoomUsers()

		for uid, rMp := range inRoomUsers {
			for roomId, _ := range rMp {
				if uint64(roomId) == req.RoomId {
					onlineUserIds = append(onlineUserIds, uint64(uid))
				}
			}
		}

		jsonTxt, err = json.Marshal(onlineUserIds)

		global.GoCache.Set(key, jsonTxt, time.Hour)
	} else {
		onlineUserIds = msgcenter.GetOnlineUsers()

		jsonTxt, err = json.Marshal(onlineUserIds)

		global.GoCache.Set(key, jsonTxt, time.Hour)
	}

	if len(onlineUserIds) == 0 {
		res.List = list
		return
	}

	//渠道，充值金币范围 过滤

	if req.Page > 1 {
		val, ok := global.GoCache.Get(key)

		if !ok {
			err = errors.New("get cache failed")
		}

		jsonTxt = val.([]byte)

		err = json.Unmarshal(jsonTxt, &onlineUserIds)
		if err != nil {
			return
		}
	}

	offset = (req.Page - 1) * req.Size

	total := len(onlineUserIds)

	if total > 0 {
		limit := req.Page * req.Size

		if limit > total {
			limit = total
		}

		var userIds []uint64

		if req.Page > (total/req.Size + 1) {
			res.List = list
			return
		}

		if req.Uid > 0 {
			var find bool
			for _, id := range onlineUserIds {
				if req.Uid == id {
					find = true
					break
				}
			}

			if find {
				userIds = []uint64{req.Uid}
			} else {
				res.List = list
				return
			}
		} else {
			userIds = onlineUserIds[offset:limit]
		}

		var (
			ids        []int
			userMp     = make(map[int]user.User)
			uidMp      = make(map[int][]int)
			userInfos  []user.UserInfo
			infoMp     = make(map[int]user.UserInfo)
			userDataMp = make(map[int]gm.UserData)
			roomIds    []uint64
			userRoomMp = make(map[int]uint64)
		)

		for _, id := range userIds {
			intId := int(id)

			ids = append(ids, intId)

			uidMpKey := intId % 5

			uids, uidMpOk := uidMp[uidMpKey]

			if !uidMpOk {
				uids = make([]int, 0)
			}

			uids = append(uids, intId)

			uidMp[uidMpKey] = uids

			inRoomMp := msgcenter.GetInRoomUsersByUid(id)

			for roomId, _ := range inRoomMp {
				roomIds = append(roomIds, uint64(roomId))
				userRoomMp[intId] = uint64(roomId)
				//如果有多个，只拿第一个，应该不会有多个，一个用户不能同时在多个房间
				break
			}

		}

		userDb := global.User
		gmDb := global.DB
		gameDb := global.Game

		roomIds = slice.UniqueSlice(roomIds)

		//房间名称、游戏名称
		viewObj := new(game.Roomview)
		viewList := make([]game.Roomview, 0, len(roomIds))

		viewList, err = viewObj.GetListByRoomIds(gameDb, roomIds)
		if err != nil {
			return
		}

		type RoomInfo struct {
			RoomName string
			GameId   int
		}

		roomMp := make(map[int]RoomInfo)
		gameMp := make(map[int]string)
		gameIds := make([]int, 0, len(viewList))

		for _, rv := range viewList {
			roomMp[int(rv.RoomId)] = RoomInfo{
				RoomName: rv.RoomName,
				GameId:   int(rv.GameId),
			}

			gameIds = append(gameIds, int(rv.GameId))
		}

		gameIds = slice.UniqueSlice(gameIds)

		gameBaseObj := new(game.GameBase)
		gameBaseList := make([]game.GameBase, 0, len(gameIds))

		gameBaseList, err = gameBaseObj.GetListByIds(gameDb, gameIds)
		if err != nil {
			return
		}

		for _, base := range gameBaseList {
			gameMp[base.ID] = base.GameName
		}

		userObj := new(user.User)

		//查用户基本数据
		userMp, err = userObj.GetMpByUserIds(userDb, ids)
		if err != nil {
			return
		}

		infoObj := new(user.UserInfo)

		userInfos, err = infoObj.GetListByUserIds(userDb, uidMp)

		if err != nil {
			return
		}

		for _, info := range userInfos {
			infoMp[info.Uid] = info
		}

		userDataObj := new(gm.UserData)

		userDataMp, err = userDataObj.GetListByUidMp(gmDb, uidMp)

		if err != nil {
			return
		}

		for _, u := range userMp {

			if req.ChannelId > 0 && u.ChannelId != req.ChannelId {
				continue
			} else if req.ChannelId == 0 {
				channelIdsMp := slice.SliceToMap(req.ChannelIds)

				_, ok := channelIdsMp[u.ChannelId]
				//用户不属于当前登录后台用户的角色渠道，跳过
				if !ok {
					continue
				}
			}

			infoVal, infoMpOk := infoMp[u.Uid]

			if !infoMpOk {
				infoVal = user.UserInfo{}
			}

			if req.RechargeMin > 0 && infoVal.Recharge < req.RechargeMin {
				continue
			}

			if req.RechargeMax > 0 && infoVal.Recharge > req.RechargeMax {
				continue
			}

			userDataVal, userDataMpOk := userDataMp[u.Uid]

			if !userDataMpOk {
				userDataVal = gm.UserData{}
			}

			var (
				gameName string
				roomName string
			)

			userRoomId, userRoomMpOk := userRoomMp[u.Uid]

			if userRoomMpOk {
				roomInfo, roomMpOk := roomMp[int(userRoomId)]

				if roomMpOk {
					roomName = roomInfo.RoomName

					gameName, _ = gameMp[roomInfo.GameId]
				}
			}

			list = append(list, response.Online{
				Uid:           u.Uid,
				Nickname:      u.UserName,
				GameName:      gameName,
				RoomName:      roomName,
				Coin:          infoVal.Cash,
				WinCash:       infoVal.WinCash,
				RechargeTotal: infoVal.Recharge,
				GiveTotal:     userDataVal.WithdrawTotal,
				RegTime:       u.CreatedAt.Format(timeutil.TimeFormat),
				ChannelId:     u.ChannelId,
				ChannelName:   "",
				ChannelCode:   "",
			})
		}
	}

	res.List = list
	res.Total = total

	return
}
