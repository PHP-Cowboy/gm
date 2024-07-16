package cache

import (
	"gm/global"
	"gm/model/pay"
	"time"
	"za.game/lib/consts"
)

func GetAllPassageMap() (mp map[int]pay.Passage, err error) {
	key := consts.MapPassageList

	mp = make(map[int]pay.Passage)

	cacheChannelMp, ok := global.GoCache.Get(key)

	if !ok {
		var (
			dataList []pay.Passage
		)

		dataList, err = GetAllPassageList()

		if err != nil {
			return
		}

		for _, p := range dataList {
			mp[p.Id] = p
		}

		global.GoCache.Set(key, mp, 1*time.Minute)

		return
	}

	mp = cacheChannelMp.(map[int]pay.Passage)

	return
}

// 根数据缓存，时间较久。后台更新user_channel时需删除这个缓存
func GetAllPassageList() (dataList []pay.Passage, err error) {
	key := consts.PassageList

	cachePassageList, ok := global.GoCache.Get(key)

	if !ok {
		var obj pay.Passage

		payDb := global.Pay

		dataList, err = obj.GetList(payDb)

		if err != nil {
			global.Logger["err"].Errorf("select passage failed, err:[%v]", err.Error())
			return
		}

		global.GoCache.Set(key, dataList, 5*time.Minute)

		return
	}

	dataList = cachePassageList.([]pay.Passage)

	return
}
