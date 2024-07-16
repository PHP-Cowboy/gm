package daos

import (
	"fmt"
	"gm/daos/rds"
	"gm/global"
	"gm/model/game"
	"gm/model/gm"
	"gm/request"
	"gm/response"
	"gm/utils/ecode"
	"time"
)

func DictTypeList(req request.DictTypeList) (data response.DictTypeList, err error) {
	db := global.DB

	var (
		total    int64
		dictType []gm.DictType
	)

	obj := new(gm.DictType)

	total, err = obj.Count(db, req)
	if err != nil {
		return
	}

	dictType, err = obj.GetPageList(db, req)
	if err != nil {
		return
	}

	list := make([]response.DictType, 0, len(dictType))

	for _, dt := range dictType {
		list = append(list, response.DictType{
			Id:   dt.Id,
			Code: dt.Code,
			Name: dt.Name,
		})
	}

	data.Total = total
	data.List = list

	return
}

func DictList(req request.DictList) (list []response.DictList, err error) {
	db := global.Game
	dict := new(game.Dict)
	data, err := dict.GetList(db, game.Dict{TypeCode: req.TypeCode})
	if err != nil {
		return
	}

	for _, d := range data {
		list = append(list, response.DictList{
			Id:       d.Id,
			TypeCode: d.TypeCode,
			Code:     d.Code,
			Name:     d.Name,
			Value:    d.Value,
			IsEdit:   d.IsEdit,
		})
	}
	return
}

func GetOneDict(req request.GetOneDict) (dict response.DictList, err error) {
	db := global.Game
	obj := new(game.Dict)

	var d game.Dict

	d, err = obj.GetOneByUnique(db, req.TypeCode, req.Code)

	if err != nil {
		return
	}

	dict = response.DictList{
		Id:       d.Id,
		TypeCode: d.TypeCode,
		Code:     d.Code,
		Name:     d.Name,
		Value:    d.Value,
		IsEdit:   d.IsEdit,
	}

	return
}

func SaveDict(req request.EditDict) (err error) {
	db := global.Game
	dict := game.Dict{
		TypeCode:   req.TypeCode,
		Code:       req.Code,
		Name:       req.Name,
		Value:      req.Value,
		IsEdit:     req.IsEdit,
		CreateTime: time.Now(),
	}
	err = dict.Save(db)

	if err != nil {
		return
	}

	//根据type code 删除缓存
	err = rds.DelRedisCacheByKey(fmt.Sprintf("hall:%s:config", dict.TypeCode))

	return
}

func ChangeValues(values map[string]string) (err error) {

	tCode, ok := values["type_code"]

	if !ok {
		err = ecode.ParamInvalid
		return
	}

	tx := global.Game.Begin()
	obj := new(game.Dict)

	for k, v := range values {
		//tCode 跳过
		if k == "type_code" {
			continue
		}

		err = obj.UpdateValueByTCode(tx, tCode, k, v)

		if err != nil {
			tx.Rollback()
			return
		}
	}

	tx.Commit()

	return
}
