package rds

import (
	"encoding/json"
	"errors"
	"gm/global"
	"za.game/lib/account/user"
	"za.game/lib/consts"
)

func GetChannelById(channelId int) (userCh user.UserChannel, err error) {
	var (
		mp   = make(map[int]user.UserChannel)
		mpOk bool
	)

	mp, err = GetAllChannelMap()

	if err != nil {
		global.Logger["err"].Errorf("GetChannelById GetAllChannelMap failed,err:[%v]", err.Error())
		return
	}

	userCh, mpOk = mp[channelId]

	if !mpOk {
		err = errors.New("channel info not found")
		global.Logger["err"].Errorf("GetChannelById channel info not found")
		return
	}

	return

}

// user_channel map 的 key 值为渠道id
func GetAllChannelMap() (mp map[int]user.UserChannel, err error) {

	mp = make(map[int]user.UserChannel)

	var (
		dataList []user.UserChannel
	)

	dataList, err = GetAllChannelList()

	if err != nil {
		global.Logger["err"].Errorf("GetAllChannelMap GetAllChannelList failed,err:[%v]", err.Error())
		return
	}

	for _, ch := range dataList {
		mp[ch.Id] = ch
	}

	return
}

// 根数据缓存，时间较久。后台更新user_channel时需删除这个缓存
func GetAllChannelList() (dataList []user.UserChannel, err error) {
	key := consts.UserChannelList

	r := global.Redis

	val, err := r.Get(key)
	if err != nil {
		global.Logger["err"].Errorf("GetAllChannelList r.Get failed,err:[%v]", err.Error())
		return
	}

	if val == "" {
		dataList, err = user.GetUserChannelList()

		if err != nil {
			global.Logger["err"].Errorf("select user_channel failed, err:" + err.Error())
			return
		}

		var b []byte

		b, err = json.Marshal(dataList)
		if err != nil {
			global.Logger["err"].Errorf("GetAllChannelList json.Marshal failed,err:[%v]", err.Error())
			return
		}

		_, err = r.Set(key, string(b), 86400)
		if err != nil {
			global.Logger["err"].Errorf("GetAllChannelList r.Set failed,err:[%v]", err.Error())
			return
		}

		return
	}

	err = json.Unmarshal([]byte(val), &dataList)

	if err != nil {
		global.Logger["err"].Errorf("GetAllChannelList json.Unmarshal failed,err:[%v]", err.Error())
		return
	}

	return
}

func GetChannelIdNameMp() (map[int]string, error) {

	channelList, err := GetAllChannelList()

	if err != nil {
		return nil, err
	}

	mp := make(map[int]string)

	for _, l := range channelList {
		mp[l.Id] = l.ChannelName
	}

	return mp, nil
}
