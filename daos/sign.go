package daos

import (
	"errors"
	"gm/common/constant"
	"gm/daos/rds"
	"gm/global"
	"gm/model/game"
	"gm/request"
	"gm/response"
	"strconv"
	"strings"
	"time"
)

// 获取签到配置列表
func GetSignConfigList() (list []response.SignConfigList, err error) {
	//数据库查询
	db := global.Game
	sign := new(game.SignConfig)

	var (
		signList  []game.SignConfig
		prizeList []game.Prize
	)

	signList, err = sign.GetList(db)

	if err != nil {
		return
	}

	prize := new(game.Prize)
	prizeList, err = prize.GetList(db)

	if err != nil {
		return
	}

	prizeMp := make(map[uint64]string, len(prizeList))

	for _, pl := range prizeList {
		prizeMp[pl.ID] = pl.Name
	}

	for _, sl := range signList {
		var (
			cash  = ""
			bonus = ""
		)

		cash, _ = prizeMp[uint64(sl.PrizeCashId)]
		bonus, _ = prizeMp[uint64(sl.PrizeBonusId)]

		list = append(list, response.SignConfigList{
			Id:           sl.ID,
			Name:         sl.Name,
			SignNum:      sl.SignNum,
			PrizeIds:     sl.PrizeIds,
			PrizeCashId:  sl.PrizeCashId,
			PrizeBonusId: sl.PrizeBonusId,
			Unit:         sl.Unit,
			Remark:       sl.Remark,
			Cash:         cash,
			Bonus:        bonus,
		})
	}

	return
}

// 保存签到配置
func SaveSign(req request.SaveSign) (err error) {
	db := global.Game
	sign := game.SignConfig{
		Name:         req.Name,
		SignNum:      req.SignNum,
		PrizeCashId:  req.PrizeCashId,
		PrizeBonusId: req.PrizeBonusId,
		Unit:         req.Unit,
		Remark:       req.Remark,
		CreatedAt:    time.Now(),
	}

	if req.PrizeCashId > 0 {
		sign.PrizeIds = append(sign.PrizeIds, strconv.Itoa(req.PrizeCashId))
	}

	if req.PrizeBonusId > 0 {
		sign.PrizeIds = append(sign.PrizeIds, strconv.Itoa(req.PrizeBonusId))
	}

	if req.ID > 0 {
		sign.ID = req.ID
	}

	err = sign.Save(db)
	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			err = errors.New("签到次数和签到频次不能重复")
			return
		}
		return
	}

	//删除缓存数据，客户端接口查不到自动更新
	err = rds.DelRedisCacheByKey(constant.SignConfig)

	if err != nil {
		return
	}

	signPrizeDaysKey := make([]string, 0, 8)

	//删除七天签到
	for i := 1; i < 8; i++ {
		signPrizeDaysKey = append(signPrizeDaysKey, constant.SignPrizeDay+strconv.Itoa(i))
	}

	err = rds.DelRedisCacheByKey(signPrizeDaysKey...)

	return
}

// 删除签到配置
func DelSign(req request.DelSign) (err error) {
	db := global.Game
	sign := new(game.SignConfig)

	err = sign.DeleteById(db, req.Id)
	if err != nil {
		return
	}
	//删除缓存数据，客户端接口查不到自动更新
	err = rds.DelRedisCacheByKey(constant.SignConfig)

	if err != nil {
		return
	}

	signPrizeDaysKey := make([]string, 0, 8)

	//删除七天签到
	for i := 1; i < 8; i++ {
		signPrizeDaysKey = append(signPrizeDaysKey, constant.SignPrizeDay+strconv.Itoa(i))
	}

	err = rds.DelRedisCacheByKey(signPrizeDaysKey...)
	return
}

// 获取签到奖励配置列表
func GetSingPrizeList(req request.GetSingPrizeList) (res response.SingPrizeRsp, err error) {
	//数据库查询
	db := global.Game
	prize := new(game.Prize)

	var (
		total    int64
		dataList []game.Prize
	)

	total, dataList, err = prize.GetPageList(db, req)

	if err != nil {
		return
	}

	res.Total = total

	list := make([]response.SingPrize, 0, len(dataList))

	for _, d := range dataList {
		list = append(list, response.SingPrize{
			Id:        d.ID,
			Name:      d.Name,
			EnName:    d.EnName,
			Type:      d.Type,
			GoodsNum:  d.GoodsNum,
			GoodsType: d.GoodsType,
			Unit:      d.Unit,
			Remark:    d.Remark,
		})
	}

	res.List = list

	return
}

// 保存签到奖励配置
func SaveSingPrize(req request.SavePrize) (err error) {
	db := global.Game
	prize := game.Prize{
		Name:      req.Name,
		EnName:    req.EnName,
		Type:      req.Type,
		GoodsNum:  req.GoodsNum,
		GoodsType: req.GoodsType,
		Remark:    req.Remark,
		CreatedAt: time.Now(),
	}

	if req.ID > 0 {
		prize.ID = req.ID
	}

	err = prize.Save(db)
	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			err = errors.New("奖励名称或金额和金额类别不能重复")
			return
		}
		return
	}

	//删除缓存数据，客户端接口查不到自动更新
	err = rds.DelRedisCacheByKey(constant.SignConfig)
	if err != nil {
		return
	}

	//删除客户端接口层【获取签到奖励】缓存
	err = rds.DelRedisCacheByKey(constant.Prize)

	if err != nil {
		return
	}

	signPrizeDaysKey := make([]string, 0, 8)

	//删除七天签到
	for i := 1; i < 8; i++ {
		signPrizeDaysKey = append(signPrizeDaysKey, constant.SignPrizeDay+strconv.Itoa(i))
	}

	err = rds.DelRedisCacheByKey(signPrizeDaysKey...)
	return
}

// 删除签到配置
func DelSignPrize(req request.DelSingPrize) (err error) {
	db := global.Game
	prize := new(game.Prize)
	sign := new(game.SignConfig)
	list, err := sign.GetList(db)
	if err != nil {
		return
	}

	var intId = int(req.Id)

	for _, l := range list {
		if l.PrizeCashId == intId || l.PrizeBonusId == intId {
			err = errors.New("已经被配置到签到中的奖励不允许删除")
			return
		}
	}

	err = prize.DeleteById(db, req.Id)
	if err != nil {
		return
	}
	//删除缓存数据，客户端接口查不到自动更新
	err = rds.DelRedisCacheByKey(constant.SignConfig)
	if err != nil {
		return
	}

	//删除客户端接口层【获取签到奖励】缓存
	err = rds.DelRedisCacheByKey(constant.Prize)

	if err != nil {
		return
	}

	signPrizeDaysKey := make([]string, 0, 8)

	//删除七天签到
	for i := 1; i < 8; i++ {
		signPrizeDaysKey = append(signPrizeDaysKey, constant.SignPrizeDay+strconv.Itoa(i))
	}

	err = rds.DelRedisCacheByKey(signPrizeDaysKey...)
	return
}
