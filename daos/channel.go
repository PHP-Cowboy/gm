package daos

import (
	"gm/global"
	"gm/model/pay"
	"strconv"
)

type PayCfgAndRate struct {
	PayCfg pay.PayConfig
	Rate   float64
}

func GetChannelPassageListMp() (channelMp map[int][]PayCfgAndRate, err error) {
	payDb := global.Pay

	channelMp = make(map[int][]PayCfgAndRate)

	var (
		passageList []pay.Passage
		payCfgList  []pay.PayConfig
		payCfgMp    = make(map[int]pay.PayConfig)
	)

	passageObj := new(pay.Passage)
	passageList, err = passageObj.GetListPaymentSort(payDb)

	if err != nil {
		global.Logger["err"].Infof("GetChannelPassageListMp query passage list err:[%v]", err.Error())
		return
	}

	payCfgObj := new(pay.PayConfig)

	payCfgList, err = payCfgObj.GetList(payDb)

	if err != nil {
		global.Logger["err"].Infof("GetChannelPassageListMp query payCfg list err:[%v]", err.Error())
		return
	}

	for _, pc := range payCfgList {
		payCfgMp[pc.ID] = pc
	}

	for _, pl := range passageList {
		for _, chanId := range pl.PaymentChannelIds {
			var cId int
			cId, err = strconv.Atoi(chanId)
			if err != nil {
				global.Logger["err"].Infof("GetChannelPassageListMp strconv.Atoi failed, err:[%v]", err.Error())

				return nil, err
			}

			channelPayList, ok := channelMp[cId]

			if !ok {
				channelPayList = make([]PayCfgAndRate, 0)
			}

			payCfg, payCfgOk := payCfgMp[pl.PassageId]

			if !payCfgOk {
				continue
			}

			tmp := PayCfgAndRate{PayCfg: payCfg, Rate: pl.PaymentRate}

			channelPayList = append(channelPayList, tmp)

			channelMp[cId] = channelPayList
		}
	}

	return

}
