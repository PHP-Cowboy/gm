package cache

import (
	"gm/global"
	"gm/model/pay"
	"time"
	"za.game/lib/consts"
)

func GetAllPayCfgMap() (mp map[int]pay.PayConfig, err error) {
	key := consts.MapPayCfgList

	mp = make(map[int]pay.PayConfig)

	cacheChannelMp, ok := global.GoCache.Get(key)

	if !ok {
		var (
			dataList []pay.PayConfig
		)

		dataList, err = GetAllPayCfgList()

		if err != nil {
			return
		}

		for _, p := range dataList {
			mp[p.ID] = p
		}

		global.GoCache.Set(key, mp, 1*time.Minute)

		return
	}

	mp = cacheChannelMp.(map[int]pay.PayConfig)

	return
}

// 根数据缓存，时间较久。后台更新user_channel时需删除这个缓存
func GetAllPayCfgList() (dataList []pay.PayConfig, err error) {
	key := consts.PayCfgList

	cachePassageList, ok := global.GoCache.Get(key)

	if !ok {
		var obj pay.PayConfig

		payDb := global.Pay

		dataList, err = obj.GetList(payDb)

		if err != nil {
			global.Logger["err"].Errorf("select passage failed, err:[%v]", err.Error())
			return
		}

		global.GoCache.Set(key, dataList, 5*time.Minute)

		return
	}

	dataList = cachePassageList.([]pay.PayConfig)

	return
}
