package daos

import (
	"gm/msgcenter"
	req2 "gm/utils/request"
)

type OnlineRsp struct {
	UserIds         []uint64   `json:"user_ids"`
	LenConnectUsers int        `json:"len_connect_users"`
	LenInRoomUsers  int        `json:"len_InRoomUsers"`
	LenOnlineRooms  int        `json:"len_OnlineRooms"`
	OnlineInfo      OnlineInfo `json:"online_info"`
}

type OnlineInfo struct {
	ConnectUsers []uint64 `json:"ConnectUsers"`
	InRoomUsers  []uint64 `json:"InRoomUsers"`
	OnlineRooms  []uint64 `json:"OnlineRooms"`
}

func Online() (rsp OnlineRsp, err error) {
	userIds := msgcenter.GetOnlineUsers()

	//var b []byte

	//b, err = request.Get("http://192.168.0.254:5020/total?info=1")

	//if err != nil {
	//	return
	//}
	//
	//var info OnlineInfo
	//
	//err = json.Unmarshal(b, &info)
	//
	//if err != nil {
	//	return
	//}

	var info OnlineInfo

	err = req2.Call("http://192.168.0.254:5020/total?info=1", nil, &info)
	if err != nil {
		return OnlineRsp{}, err
	}

	rsp.LenConnectUsers = len(info.ConnectUsers)
	rsp.LenInRoomUsers = len(info.InRoomUsers)
	rsp.LenOnlineRooms = len(info.OnlineRooms)

	rsp.UserIds = userIds
	rsp.OnlineInfo = info

	return
}
