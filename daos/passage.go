package daos

import (
	"gm/global"
	"gm/model"
	"gm/model/pay"
	"gm/request"
	"gm/response"
	"gm/utils/timeutil"
)

func PassageList(req request.PassageList) (rsp response.PassageRsp, err error) {
	db := global.Pay

	obj := new(pay.Passage)

	var (
		total    int64
		dataList []pay.Passage
	)

	total, dataList, err = obj.GetPageList(db, req)
	if err != nil {
		return
	}

	list := make([]response.Passage, 0, len(dataList))

	for _, d := range dataList {
		list = append(list, response.Passage{
			Id:                   d.Id,
			PassageId:            d.PassageId,
			Entrance:             d.Entrance,
			CollectionRate:       d.CollectionRate,
			CollectionStatus:     d.CollectionStatus,
			CollectionSort:       d.CollectionSort,
			CollectionChannelIds: d.CollectionChannelIds,
			PaymentRate:          d.PaymentRate,
			PaymentStatus:        d.PaymentStatus,
			PaymentSort:          d.PaymentSort,
			PaymentChannelIds:    d.PaymentChannelIds,
			CreatedAt:            timeutil.FormatToDateTime(&d.CreatedAt),
			UpdatedAt:            timeutil.FormatToDateTime(&d.UpdatedAt),
		})
	}

	rsp.Total = total
	rsp.List = list

	return
}

func PassageSave(req request.PassageSave) (err error) {
	db := global.Pay

	obj := new(pay.Passage)

	data := pay.Passage{
		PassageId:            req.PassageId,
		Entrance:             req.Entrance,
		CollectionRate:       req.CollectionRate,
		CollectionStatus:     req.CollectionStatus,
		CollectionSort:       req.CollectionSort,
		CollectionChannelIds: model.GormList(req.CollectionChannelIds),
		PaymentRate:          req.PaymentRate,
		PaymentStatus:        req.PaymentStatus,
		PaymentSort:          req.PaymentSort,
		PaymentChannelIds:    model.GormList(req.PaymentChannelIds),
	}

	if req.Id > 0 {
		data.Id = req.Id
	}

	err = obj.Create(db, data)

	return
}

func PassageDel(req request.DeleteId) (err error) {
	db := global.Pay

	obj := new(pay.Passage)

	err = obj.DeleteById(db, req.Id)

	return
}

func PassageChange(req request.PassageChange) (err error) {
	db := global.Pay

	obj := new(pay.Passage)

	mp := make(map[string]interface{})

	switch req.Type {
	case "paymentStatus":
		mp["payment_status"] = req.PaymentStatus
	case "paymentChannel":
		mp["payment_channel_ids"] = model.GormList(req.PaymentChannelIds)
	case "collectionStatus":
		mp["collection_status"] = req.CollectionStatus
	case "collectionChannel":
		mp["collection_channel_ids"] = model.GormList(req.CollectionChannelIds)
	}

	err = obj.UpdateById(db, req.Id, mp)

	return
}
