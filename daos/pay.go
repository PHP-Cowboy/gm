package daos

import (
	"errors"
	"gm/daos/rds"
	"gm/global"
	"gm/model/pay"
	"gm/request"
	"gm/response"
	"gm/utils/timeutil"
	"gorm.io/gorm"
	"time"
	accUser "za.game/lib/account/user"
	"za.game/lib/consts"
)

// 充值礼包列表
func PayList(req request.PayList) (list []response.PayGift, err error) {
	db := global.Pay
	gift := &pay.PayGift{
		Name:    req.Name,
		Cash:    req.Cash,
		Account: req.Account,
		Status:  req.Status,
		Type:    req.Type,
	}

	var payGiftList []pay.PayGift

	payGiftList, err = gift.GetList(db, gift)

	list = make([]response.PayGift, 0, len(payGiftList))

	for _, l := range payGiftList {
		list = append(list, response.PayGift{
			ID:           l.ID,
			Name:         l.Name,
			Cash:         l.Cash,
			Account:      l.Account,
			Status:       l.Status,
			AddMoney:     l.AddMoney,
			AddMoneyType: l.AddMoneyType,
			AddCash:      l.AddCash,
			AddCashType:  l.AddCashType,
			Bonus:        l.Bonus,
			BonusType:    l.BonusType,
			Ratio:        l.Ratio,
			Remark:       l.Remark,
			Type:         l.Type,
			ReplaceId:    l.ReplaceId,
		})
	}

	return
}

// 保存充值礼包
func SaveGift(req request.SaveGift) (err error) {
	db := global.Pay

	gift := &pay.PayGift{
		Name:         req.Name,
		Cash:         req.Cash,
		Account:      req.Account,
		Status:       req.Status,
		AddMoney:     req.AddMoney,
		AddMoneyType: req.AddMoneyType,
		AddCash:      req.AddCash,
		AddCashType:  req.AddCashType,
		Bonus:        req.Bonus,
		BonusType:    req.BonusType,
		Remark:       req.Remark,
		Ratio:        req.Ratio,
		CreatedAt:    time.Now(),
		Type:         req.Type,
		ReplaceId:    req.ReplaceId,
	}

	obj := new(pay.PayGift)

	//正常 时 替换 id 为0
	if req.Type == pay.GaveMoneyConfigNormal {
		gift.ReplaceId = 0
	} else if req.Type == pay.GaveMoneyConfigFirst {

		if req.ID > 0 {
			if req.ID == uint64(req.ReplaceId) {
				err = errors.New("The target of replacement cannot be oneself")
				global.Logger["err"].Errorf("SaveGift failed,err:[%s]", err.Error())
				return err
			}

			var cfg pay.PayGift

			cfg, err = obj.GetOneByReplaceId(db, req.ID)

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				global.Logger["err"].Errorf("SaveGift obj.GetOneByReplaceId failed,err:[%s]", err.Error())
				return err
			}

			if cfg.ID > 0 {
				err = errors.New("The current data is already the target being replaced")
				global.Logger["err"].Errorf("SaveGift check failed,err:[%s]", err.Error())
				return err
			}
		}

		if req.ReplaceId > 0 {
			var data pay.PayGift

			data, err = obj.GetFirstById(db, uint64(req.ReplaceId))

			if err != nil {
				global.Logger["err"].Errorf("SaveGift failed,err:[%s]", err.Error())
				return err
			}

			if data.Type == pay.GaveMoneyConfigFirst {
				err = errors.New("The target being replaced cannot be a replacement type")
				global.Logger["err"].Errorf("SaveGift failed,err:[%s]", err.Error())
				return err
			}

			//被替换的数据的替换的目标需要为0
			if data.ReplaceId != 0 {
				err = errors.New("The replacement target for the replaced data needs to be 0")
				global.Logger["err"].Errorf("SaveGift failed,err:[%s]", err.Error())
				return err
			}
		}

	}

	if req.ID > 0 {
		gift.ID = req.ID
	}

	err = gift.Save(db)

	return
}

// 删除充值礼包
func DelGift(req request.DeleteId) (err error) {
	db := global.Pay

	gift := new(pay.PayGift)

	var cfg pay.PayGift

	cfg, err = gift.GetOneByReplaceId(db, req.Id)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.Logger["err"].Errorf("DelGift obj.GetOneByReplaceId failed,err:[%s]", err.Error())
		return err
	}

	if cfg.ID > 0 {
		err = errors.New("The current data is potentially replaceable and cannot be deleted")
		global.Logger["err"].Errorf("DelGift check failed,err:[%s]", err.Error())
		return err
	}

	err = gift.DeleteById(db, req.Id)

	return
}

func ConfigList(req request.ConfigList) (res response.PayConfigRsp, err error) {
	db := global.Pay

	payCfg := new(pay.PayConfig)

	total, getList, err := payCfg.GetPageList(db, req)
	if err != nil {
		return
	}

	list := make([]response.PayConfig, 0, len(getList))

	for _, l := range getList {
		list = append(list, response.PayConfig{
			ID:             l.ID,
			Name:           l.Name,
			Url:            l.Url,
			BackUrl:        l.BackUrl,
			PaymentBackUrl: l.PaymentBackUrl,
			AppId:          l.AppId,
			Secret:         l.Secret,
			Merchant:       l.Merchant,
			Status:         l.Status,
			Remark:         l.Remark,
			Markers:        l.Markers,
		})
	}

	res.Total = total
	res.List = list

	return
}

// 保存支付渠道配置
func SaveConfig(req request.SaveConfig) error {
	db := global.Pay

	payCfg := &pay.PayConfig{
		Name:           req.Name,
		Icon:           req.Icon,
		Url:            req.Url,
		BackUrl:        req.BackUrl,
		PaymentBackUrl: req.PaymentBackUrl,
		AppId:          req.AppId,
		Secret:         req.Secret,
		Merchant:       req.Merchant,
		Status:         req.Status,
		Remark:         req.Remark,
		Markers:        req.Markers,
		CreatedAt:      time.Now(),
	}

	if req.ID > 0 {
		payCfg.ID = req.ID
	}

	return payCfg.Save(db)
}

func DelConfig(req request.DeleteId) (err error) {
	db := global.Pay

	payCfg := &pay.PayConfig{}

	err = payCfg.DeleteById(db, req.Id)
	if err != nil {
		global.Logger["err"].Errorf("DelConfig payCfg.DeleteById failed,err:[%v]", err.Error())
		return
	}

	return
}

// 获取用户银行卡列表
func BankList(req request.BankList) (res response.BankRsp, err error) {
	db := global.Pay

	bankInfo := &pay.BankInfo{
		Uid:       req.Uid,
		BankCode:  req.BankCode,
		BankName:  req.BankName,
		AccountNo: req.AccountNo,
		Ifsc:      req.Ifsc,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
	}

	total, pageList, err := bankInfo.GetPageList(db, bankInfo, req.Page, req.Size, req.ChannelIds)
	if err != nil {
		return
	}

	list := make([]response.Bank, 0, len(pageList))

	for _, l := range pageList {
		list = append(list, response.Bank{
			ID:        l.ID,
			Uid:       l.Uid,
			BankCode:  l.BankCode,
			BankName:  l.BankName,
			AccountNo: l.AccountNo,
			Ifsc:      l.Ifsc,
			Name:      l.Name,
			Email:     l.Email,
			Phone:     l.Phone,
			Address:   l.Address,
			Vpa:       l.Vpa,
			Remark:    l.Remark,
			CreatedAt: l.CreatedAt.Format(timeutil.MinuteFormat),
		})
	}
	res.Total = total
	res.List = list
	return
}

// 获取用户银行卡列表
func OrderList(req request.OrderList) (list []response.OrderList, err error) {
	db := global.Pay

	order := &pay.Order{
		Uid:     req.Uid,
		Ymd:     req.Ymd,
		OrderNo: req.OrderNo,
		Name:    req.Name,
		Phone:   req.Phone,
		Status:  req.Status,
	}

	pageList, err := order.GetPageList(db, order, req.Page, req.Size)
	if err != nil {
		return
	}

	list = make([]response.OrderList, 0, len(pageList))

	for _, l := range pageList {
		list = append(list, response.OrderList{
			ID:           l.ID,
			Uid:          l.Uid,
			Ymd:          l.Ymd,
			OrderNo:      l.OrderNo,
			MOrderNo:     l.MOrderNo,
			Account:      l.Account,
			Cash:         l.Cash,
			GiftCash:     l.GiftCash,
			Bonus:        l.Bonus,
			RequestTime:  l.RequestTime.Format(timeutil.MinuteFormat),
			Email:        l.Email,
			Name:         l.Name,
			Phone:        l.Phone,
			RedirectTime: l.RedirectTime.Format(timeutil.MinuteFormat),
			Status:       l.Status,
			CompleteTime: l.CompleteTime,
			Type:         l.Type,
			Remark:       l.Remark,
		})
	}
	return
}

func RechargeRecords(req request.RechargeRecords) (res response.RechargeRecords, err error) {
	db := global.Pay

	order := new(pay.Order)

	if len(req.Ymd) > 0 {
		req.Start, req.End, err = timeutil.DateRangeToZeroAndLastTimeFormat(req.Ymd[0], req.Ymd[1], timeutil.DateNumberFormat)

		if err != nil {
			global.Logger["err"].Errorf("TpFundFlowLog timeutil.DateRangeToZeroAndLastTimeFormat failed,err:[%v]", err.Error())
			return
		}
	}

	total, records, err := order.RechargeRecords(db, req)

	if err != nil {
		return
	}

	var userIds []int

	for _, r := range records {
		userIds = append(userIds, r.Uid)
	}

	list := make([]response.RechargeRecordsList, 0, len(records))

	for _, r := range records {
		completeTime := ""
		if r.CompleteTime > 0 {
			completeTime = time.Unix(int64(r.CompleteTime), 0).Format(timeutil.TimeFormat)
		}

		var userCh accUser.UserChannel

		userCh, err = rds.GetChannelById(r.Channel)

		if err != nil {
			userCh = accUser.UserChannel{}
		}

		tmp := response.RechargeRecordsList{
			Id:           r.ID,
			Uid:          r.Uid,
			Ymd:          r.Ymd,
			Channel:      userCh.ChannelName,
			ChannelNo:    userCh.Code,
			OrderNo:      r.OrderNo,
			MOrderNo:     r.MOrderNo,
			Email:        r.Email,
			Phone:        r.Phone,
			Account:      r.Account,
			GiftCash:     r.GiftCash,
			Cash:         r.Cash,
			Bonus:        r.Bonus,
			Name:         r.Name,
			Status:       r.Status,
			CompleteTime: completeTime,
			//Type:         r.Type,
			Remark:    r.Remark,
			CreatedAt: r.CreatedAt.Format(timeutil.TimeFormat),
			UpdatedAt: r.UpdatedAt.Format(timeutil.TimeFormat),
		}

		if r.GiftId > consts.ItemIdRechargeRoomGiftIDBegin && r.GiftId <= consts.ItemIdRechargeRoomGiftIDEnd {
			r.Type = 2 //房间特惠礼包
		} else if r.GiftId > consts.ItemIdRechargeGiftIDBegin && r.GiftId <= consts.ItemIdRechargeGiftIDEnd { //充值赠送礼包
			r.Type = 3 //充值赠送礼包
		} else if r.GiftId > consts.ItemIdEventGiftIDBegin && r.GiftId <= consts.ItemIdEventGiftIDEnd { //发货二选一礼包
			r.Type = 4 //二选一礼包
		} else if r.GiftId > consts.ItemIdOnlyOneGiftIDBegin && r.GiftId <= consts.ItemIdOnlyOneGiftIDEnd { //发货OnlyOne三选一礼包
			r.Type = 5 // OnlyOne三选一礼包
		} else if r.GiftId > consts.ItemIdBenefitGiftIDBegin && r.GiftId <= consts.ItemIdBenefitGiftIDEnd { //发货救济金礼包
			r.Type = 6 //救济金礼包
		} else if r.GiftId > consts.ItemIdLuckSpinGiftIDBegin && r.GiftId <= consts.ItemIdLuckSpinGiftIDEnd { //发货luckspin礼包
			r.Type = 7 //luckspin礼包
		} else { //正常的充值项发货
			r.Type = 1
		}

		tmp.Type = r.Type

		list = append(list, tmp)
	}

	res.Total = total
	res.List = list

	return
}

func GaveConfigPageList(req request.GaveConfigPageList) (res response.GaveConfigRsp, err error) {

	db := global.Pay

	obj := new(pay.GaveMoneyConfig)

	var (
		total    int64
		dataList []pay.GaveMoneyConfig
	)

	total, dataList, err = obj.GetPageList(db, req)

	if err != nil {
		global.Logger["err"].Errorf("GaveConfigPageList obj.GetPageList failed,err:[%s]", err.Error())
		return
	}

	list := make([]response.GaveConfig, 0, len(dataList))

	for _, d := range dataList {
		list = append(list, response.GaveConfig{
			Id:        d.Id,
			Name:      d.Name,
			Account:   d.Account,
			Status:    d.Status,
			Type:      d.Type,
			ReplaceId: d.ReplaceId,
		})
	}

	res.Total = total
	res.List = list

	return
}

func GaveConfigList(req request.GaveConfigList) (list []response.GaveConfig, err error) {
	db := global.Pay

	obj := new(pay.GaveMoneyConfig)

	var dataList []pay.GaveMoneyConfig

	dataList, err = obj.GetList(db, req)
	if err != nil {
		global.Logger["err"].Errorf("GaveConfigList obj.GetList failed,err:[%v]", err.Error())
		return
	}

	list = make([]response.GaveConfig, 0, len(dataList))

	for _, d := range dataList {
		list = append(list, response.GaveConfig{
			Id:   d.Id,
			Name: d.Name,
		})
	}

	return
}

// 保存支付渠道配置
func SaveGaveConfig(req request.SaveGaveConfig) error {
	db := global.Pay

	obj := &pay.GaveMoneyConfig{
		Name:      req.Name,
		Account:   req.Account,
		Status:    req.Status,
		Type:      req.Type,
		ReplaceId: req.ReplaceId,
	}

	//正常 时 替换 id 为0
	if req.Type == pay.GaveMoneyConfigNormal {
		obj.ReplaceId = 0
	} else if req.Type == pay.GaveMoneyConfigFirst {

		if req.Id > 0 {
			if req.Id == req.ReplaceId {
				err := errors.New("The target of replacement cannot be oneself")
				global.Logger["err"].Errorf("SaveGaveConfig failed,err:[%s]", err.Error())
				return err
			}

			cfg, err := obj.GetOneByReplaceId(db, req.Id)

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				global.Logger["err"].Errorf("SaveGaveConfig obj.GetOneByReplaceId failed,err:[%s]", err.Error())
				return err
			}

			if cfg.Id > 0 {
				err = errors.New("The current data is already the target being replaced")
				global.Logger["err"].Errorf("SaveGaveConfig check failed,err:[%s]", err.Error())
				return err
			}
		}

		if req.ReplaceId > 0 {
			data, err := obj.GetOneById(db, req.ReplaceId)

			if err != nil {
				global.Logger["err"].Errorf("SaveGaveConfig failed,err:[%s]", err.Error())
				return err
			}

			if data.Type == pay.GaveMoneyConfigFirst {
				err = errors.New("The target being replaced cannot be a replacement type")
				global.Logger["err"].Errorf("SaveGaveConfig failed,err:[%s]", err.Error())
				return err
			}

			//被替换的数据的替换的目标需要为0
			if data.ReplaceId != 0 {
				err = errors.New("The replacement target for the replaced data needs to be 0")
				global.Logger["err"].Errorf("SaveGaveConfig failed,err:[%s]", err.Error())
				return err
			}
		}

	}

	if req.Id > 0 {
		obj.Id = req.Id
	}

	return obj.Save(db)
}

func DelGaveConfig(req request.DeleteId) (err error) {
	db := global.Pay

	obj := &pay.GaveMoneyConfig{}

	cfg, err := obj.GetOneByReplaceId(db, int(req.Id))

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.Logger["err"].Errorf("SaveGaveConfig obj.GetOneByReplaceId failed,err:[%s]", err.Error())
		return err
	}

	if cfg.Id > 0 {
		err = errors.New("The current data is already the target being replaced")
		global.Logger["err"].Errorf("SaveGaveConfig check failed,err:[%s]", err.Error())
		return err
	}

	err = obj.DeleteById(db, req.Id)
	if err != nil {
		global.Logger["err"].Errorf("DelConfig obj.DeleteById failed,err:[%v]", err.Error())
		return
	}

	return
}
