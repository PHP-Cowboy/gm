package daos

import (
	"encoding/json"
	"errors"
	"gm/global"
	"gm/model/game"
	"gm/request"
	"gm/response"
	"gm/utils/slice"
	"za.game/lib/dbconn"
	sqlInfoGame "za.game/lib/sqlInfo/game"
)

// 房间列表
func RoomList(req request.RoomList) (res response.RoomList, err error) {
	db := global.Game
	rl := new(game.Roomview)

	var (
		total    int64
		pageList []game.Roomview
	)

	total, pageList, err = rl.GetPageList(db, req)
	if err != nil {
		return
	}

	roomList := make([]response.Room, 0)

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
			PoolID:        l.PoolID,
			ExtData:       l.ExtData,
		})
	}

	res.Total = total
	res.List = roomList
	return
}

// 字段 && comment
func GetColumnComment(req request.RoomList) (res []sqlInfoGame.ColumnComment, err error) {
	var dataList []sqlInfoGame.ColumnComment

	dataList, err = sqlInfoGame.GetRoomViewColumnComment(dbconn.GameDB)

	if err != nil {
		return
	}

	fields := make([]string, 0, len(dataList))

	for _, comment := range dataList {
		fields = append(fields, comment.ColumnName)
	}

	res = make([]sqlInfoGame.ColumnComment, 0)

	diffMp := slice.MapOfElementsInANotInB(fields, sqlInfoGame.RoomviewSlice)

	for _, d := range dataList {
		_, ok := diffMp[d.ColumnName]

		if ok {
			res = append(res, d)
		}
	}

	return
}

// 保存房间信息
func SaveRoom(req request.SaveRoom) (err error) {
	db := global.Game

	if req.ExtData != "" {
		var j interface{}

		err = json.Unmarshal([]byte(req.ExtData), &j)

		if err != nil {
			err = errors.New("特殊配置格式不正确:" + err.Error())
			return
		}
	}

	room := &game.Roomlist{
		SvrId:         req.SvrId,
		GameId:        req.GameId,
		RoomIndex:     req.RoomIndex,
		Base:          req.Base,
		MinEntry:      req.MinEntry,
		MaxEntry:      req.MaxEntry,
		RoomName:      req.RoomName,
		RoomType:      req.RoomType,
		RoomSwitch:    req.RoomSwitch,
		RoomWelfare:   req.RoomWelfare,
		Desc:          req.Desc,
		Tax:           req.Tax,
		BonusDiscount: req.BonusDiscount,
		AiSwitch:      req.AiSwitch,
		AiLimit:       req.AiLimit,
		RechargeLimit: *req.RechargeLimit,
		PoolID:        *req.PoolID,
		ExtData:       req.ExtData,
	}

	if req.ID > 0 {
		room.ID = req.ID
	}

	err = room.Save(db)

	return
}

// 生成ExtData
func UpdateExtDataByExcel(req []request.UpdateExtDataByExcel) (err error) {
	db := global.Game
	roomView := new(game.Roomview)

	mp := make(map[uint64][]request.UpdateExtDataByExcel)

	for _, data := range req {
		//多行数据相同的roomId 合并
		mp[data.RoomId] = append(mp[data.RoomId], data)
	}

	var jsonStr []byte

	for i, d := range mp {
		jsonStr, err = json.Marshal(d)
		if err != nil {
			return
		}

		err = roomView.UpdateByRoomId(db, i, map[string]interface{}{"ExtData": jsonStr})
		if err != nil {
			return
		}

	}

	return
}
