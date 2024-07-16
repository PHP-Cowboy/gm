package daos

import (
	"gm/global"
	"gm/model/pay"
	req2 "gm/utils/request"
	"strconv"
)

type SimulatedOrder struct {
	Millisecond int `json:"millisecond"`
	Success     int `json:"success"`
}

type SimulatedPayRes struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SimulatedPay(cfg pay.PayConfig, orderNo string, amount float64, bank pay.BankInfo) (err error) {
	path := "http://127.0.0.1:8086/v1/payment/simulated"

	var (
		millisecond int
		success     int
	)

	millisecond, err = strconv.Atoi(cfg.AppId)

	if err != nil {
		global.Logger["err"].Errorf("SimulatedPay strconv.Atoi to millisecond failed, err:[%v]", err.Error())
		return
	}

	success, err = strconv.Atoi(cfg.Merchant)

	if err != nil {
		global.Logger["err"].Errorf("SimulatedPay strconv.Atoi to success failed, err:[%v]", err.Error())
		return
	}

	params := SimulatedOrder{
		Millisecond: millisecond,
		Success:     success,
	}

	var (
		jsonTxt string
		res     SimulatedPayRes
	)

	jsonTxt, err = JsonString(params)

	if err != nil {
		global.Logger["err"].Errorf("SimulatedPay params:[%v], JsonString failed! error:[%v]", params, err)
		return
	} else {
		global.Logger["err"].Errorf("SimulatedPay jsonTxt:[%v]", jsonTxt)
	}

	err = req2.Call(path, jsonTxt, &res)

	if err != nil {
		global.Logger["err"].Errorf("SimulatedPay req2.Call failed! error:[%v]", err)
		return
	}

	return
}
