package msgcenter

import (
	"encoding/json"
	dream "git.dev666.cc/external/dreamgo"
	"git.dev666.cc/external/dreamgo/xy"
	"gm/config"
	"gm/global"
	"strconv"
	"sync"
	"time"
	"za.game/lib/account"
)

type mCenter struct {
	mct dream.IAutoConnect
}

type SvrMacAddr struct {
	Svrid xy.SvrId_ `gorm:"column:svrid;type:int(11);" json:"svrid"`
	Macid int       `gorm:"column:macid;type:int(11);" json:"macid"`
	Port  int       `gorm:"column:port;type:int(11);" json:"port"`
	InIp  string    `gorm:"column:inip;type:varchar(20);" json:"inip"`   //内网ip
	OutIp string    `gorm:"column:outip;type:varchar(20);" json:"outip"` //外围ip
}

var (
	_mct        *mCenter
	_wg         sync.WaitGroup
	_rwMutexMct sync.RWMutex

	onlineRooms  = make(map[xy.RoomId_]xy.RU2T_Room)
	inRoomUsers  = make(map[xy.UserId_]map[xy.RoomId_]xy.RU2T_User)
	connectUsers = make(map[xy.UserId_]map[xy.SvrId_]map[xy.ConnectId_]xy.RU2T_Connect)

	_rwMutexMac  sync.RWMutex
	cfgSvrMac    = make(map[xy.SvrId_]SvrMacAddr)
	lastLoadTime = time.Now()
)

func InitMct(c *config.Mct) {
	_mct = &mCenter{}
	_wg = sync.WaitGroup{}
	dream.LogStartup(c.LogCfg)
	_mct.mct = dream.AutoConnect("msgCenter", c.Addr, c.Auth, c.TypeIdx, c.SvrId, _mct)
	_mct.mct.Run(&_wg)

	account.RegisterSendNoticeEvent(SendNoticeEvent)
}

func Shutdown() {
	_mct.mct.Stop()
	_wg.Wait()
	dream.LogShutdown()
}

func (p *mCenter) OnPackage(data []byte, xyc dream.XY) {
	dream.Log("OnPackage xy " + strconv.Itoa(int(xyc)))
	if xyc >= xy.XyRoomUser2Tool_Room && xyc <= xy.XyRoomUser2Tool_Connect_Del {
		p.ru2tPackage(data, xyc)
	} else {
	}

}
func (p *mCenter) OnReady(id int) {
	dream.Log("OnReady id " + string(id))
	_rwMutexMct.Lock()
	defer _rwMutexMct.Unlock()
	clear(onlineRooms)
	clear(inRoomUsers)
	clear(connectUsers)
}
func (p *mCenter) OnClose() {
	dream.Log("OnClose ")
}

func (p *mCenter) ru2tPackage(data []byte, xyc dream.XY) {
	_rwMutexMct.Lock()
	defer _rwMutexMct.Unlock()
	switch xyc {
	case xy.XyRoomUser2Tool_Room:
		ro := xy.RU2T_Room{}
		if err := json.Unmarshal(data, &ro); err == nil {
			onlineRooms[ro.Rid] = ro
		}
	case xy.XyRoomUser2Tool_Room_Del:
		ro := xy.RU2T_Room_Del{}
		if err := json.Unmarshal(data, &ro); err == nil {
			delete(onlineRooms, ro.Rid)
		}
	case xy.XyRoomUser2Tool_User:
		ro := xy.RU2T_User{}
		if err := json.Unmarshal(data, &ro); err == nil {
			tempInfo, exists := inRoomUsers[ro.Uid]
			if !exists {
				tempInfo = make(map[xy.RoomId_]xy.RU2T_User)
			}
			tempInfo[ro.Rid] = ro
			inRoomUsers[ro.Uid] = tempInfo
		}
	case xy.XyRoomUser2Tool_User_Del:
		ro := xy.RU2T_User_Del{}
		if err := json.Unmarshal(data, &ro); err == nil {
			tempInfo, exists := inRoomUsers[ro.Uid]
			if exists {
				delete(tempInfo, ro.Rid)
				if len(inRoomUsers[ro.Uid]) == 0 {
					delete(inRoomUsers, ro.Uid)
				}
			}
		}
	case xy.XyRoomUser2Tool_Connect:
		ro := xy.RU2T_Connect{}
		if err := json.Unmarshal(data, &ro); err == nil {
			tempUser, exists := connectUsers[ro.Uid]
			if !exists {
				tempUser = make(map[xy.SvrId_]map[xy.ConnectId_]xy.RU2T_Connect)
				tempSvr := make(map[xy.ConnectId_]xy.RU2T_Connect)
				tempSvr[ro.ConnectId] = ro
				tempUser[ro.Sid] = tempSvr
			} else {
				tempSvr, exists := tempUser[ro.Sid]
				if !exists {
					tempSvr = make(map[xy.ConnectId_]xy.RU2T_Connect)
				}
				tempSvr[ro.ConnectId] = ro
				tempUser[ro.Sid] = tempSvr
			}
			connectUsers[ro.Uid] = tempUser
		}
	case xy.XyRoomUser2Tool_Connect_Del:
		ro := xy.RU2T_Connect_Del{}
		if err := json.Unmarshal(data, &ro); err == nil {
			itu, exists := connectUsers[ro.Uid]
			if exists {
				itsvr, existss := itu[ro.Sid]
				if existss {
					delete(itsvr, ro.ConnectId)
					if len(itu[ro.Sid]) == 0 {
						delete(itu, ro.Sid)
						if len(connectUsers[ro.Uid]) == 0 {
							delete(connectUsers, ro.Uid)
						}
					}
				}
			}
		}
	default:

	}
}
func GetInRoomUser(uid xy.UserId_) (map[xy.RoomId_]xy.RU2T_User, bool) {
	_rwMutexMct.RLock()
	defer _rwMutexMct.RUnlock()
	tempInfo, exists := inRoomUsers[uid]
	return tempInfo, exists
}

func GetOnlineUsers() (userIds []uint64) {
	userIds = make([]uint64, 0, len(connectUsers))
	for userId, _ := range connectUsers {
		userIds = append(userIds, uint64(userId))
	}

	return
}

func GetInRoomUsers() map[xy.UserId_]map[xy.RoomId_]xy.RU2T_User {
	return inRoomUsers
}

func GetInRoomUsersByUid(userId uint64) map[xy.RoomId_]xy.RU2T_User {
	return inRoomUsers[xy.UserId_(userId)]
}

func GetGameUsers() (userIds []uint64) {

	userIds = make([]uint64, 0, len(inRoomUsers))
	for userId, _ := range connectUsers {
		userIds = append(userIds, uint64(userId))
	}

	return
}

func GetOnlineRooms(ruid xy.RoomId_) (xy.RU2T_Room, bool) {
	_rwMutexMct.RLock()
	defer _rwMutexMct.RUnlock()
	if ruid == 0 { //如果不传id就随机取一个
		for _, v := range onlineRooms {
			return v, true
		}
		return xy.RU2T_Room{}, false
	}

	tempInfo, exists := onlineRooms[ruid]
	return tempInfo, exists
}

func GetOnlineRoomsByGameid(gameid xy.GameId_) (rooms []xy.RU2T_Room) {
	_rwMutexMct.RLock()
	defer _rwMutexMct.RUnlock()
	for _, val := range onlineRooms {
		if val.GameId == gameid {
			rooms = append(rooms, val)
		}
	}
	return rooms
}

func GetIpPortBySvrid(id xy.SvrId_) (string, int, bool) {
	if len(cfgSvrMac) <= 0 || lastLoadTime.Day() != time.Now().Day() { //重新加载
		var infos []SvrMacAddr

		db := global.Game

		err := db.Raw("select svrid, a.macid, port, b.inip, b.outip from cfg_servers a left join cfg_macs b  on A.MACID=B.MACID where enable <> 0").
			Scan(&infos).
			Error

		if err != nil {
			global.Logger["err"].Errorf("GetIpPortBySvrid  db.Raw failed,err:[%v]", err.Error())

			return "", 0, false
		}
		_rwMutexMac.Lock()
		for _, info := range infos {
			cfgSvrMac[info.Svrid] = info
		}
		_rwMutexMac.Unlock()

	}
	_rwMutexMac.RLock()
	defer _rwMutexMac.RUnlock()
	val, exist := cfgSvrMac[id]
	return val.OutIp, val.Port, exist
}
func SendBackPackage(data []byte, xy dream.XY) {
	go func() {
		if _mct != nil {
			_mct.mct.SendXY(xy, data)
		}
	}()
}

func SendBackPackageJson(xy dream.XY, v any) {
	go func() {
		if _mct != nil {
			if b, err := json.Marshal(v); err == nil {
				_mct.mct.SendXY(xy, b)
			}
		}
	}()
}

func SendUserEvent(userid, fromUserid uint64, ev string, val int32, val64 int64, dic any) {
	v := &xy.UserEvent{
		Id:       xy.UserIdItem{xy.UserId_(userid)},
		FromUser: xy.UserIdItem{xy.UserId_(fromUserid)},
		Val32:    val,
		Val64:    val64,
		Evt:      ev,
		Dic:      dic,
	}
	SendBackPackageJson(xy.XyUserEvent, v)
}

func SendNoticeEvent(ev string, dic any, userid ...xy.UserId_) {
	v := &xy.NoticeEvent{
		Evt: ev,
		Dic: dic,
	}
	v.Uid = append(v.Uid, userid...)
	SendBackPackageJson(xy.XyNoticeEvent, v)
}
