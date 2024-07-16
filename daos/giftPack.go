package daos

import (
	"errors"
	"fmt"
	"gm/common/constant"
	"gm/daos/rds"
	"gm/global"
	"gm/model/game"
	"gm/request"
	"gm/response"
	"gorm.io/gorm"
	"strings"
	"time"
	"za.game/lib/consts"
)

// 获取vip配置列表
func GetVipConfigList() (list []game.VipConfig, err error) {
	//数据库查询
	db := global.Game
	vip := new(game.VipConfig)
	//读取配置数据
	list, err = vip.GetList(db)
	return
}

// 保存vip配置
func SaveVipConfig(req request.SaveVipConfig) (err error) {
	//保存vip配置
	db := global.Game

	vip := game.VipConfig{
		Level:          req.Level,
		NeedExp:        req.NeedExp,
		WithdrawNums:   req.WithdrawNums,
		WithdrawMoney:  req.WithdrawMoney,
		DayPrizeType:   req.DayPrizeType,
		DayPrizeNums:   req.DayPrizeNums,
		WeekPrizeType:  req.WeekPrizeType,
		WeekPrizeNums:  req.WeekPrizeNums,
		MonthPrizeType: req.MonthPrizeType,
		MonthPrizeNums: req.MonthPrizeNums,
		CreatedAt:      time.Now(),
	}

	if req.Id > 0 {
		vip.ID = req.Id
	}

	err = vip.Save(db)
	if err != nil {
		return
	}

	//删除缓存
	err = rds.DelRedisCacheByKey(constant.VipConfig)

	return
}

// 删除vip配置
func DelVipConfig(req request.DelVipConfig) (err error) {
	db := global.Game

	vip := new(game.VipConfig)

	err = vip.DeleteById(db, req.Id)
	if err != nil {
		return
	}
	//删除缓存
	err = rds.DelRedisCacheByKey(constant.VipConfig)

	return
}

// 二选一 配置 列表
func GetEventGiftList() (list []game.EventGiftConfig, err error) {
	db := global.Game

	vip := new(game.EventGiftConfig)

	return vip.GetList(db)
}

// 保存二选一配置
func SaveEventConfig(req request.SaveEventGiftConfig) (err error) {
	event := &game.EventGiftConfig{
		Grade:     req.Grade,
		Coin:      req.Coin,
		CoinType:  req.CoinType,
		GiftLimit: req.GiftLimit,
		GiftType:  req.GiftType,
		Bonus:     req.Bonus,
		BonusType: 3,
		Ratio:     req.Ratio,
		Type:      req.Type,
		CreatedAt: time.Now(),
	}

	if req.Id > 0 {
		event.ID = req.Id
	}

	db := global.Game

	err = event.Save(db)

	if err != nil {
		if err != nil {
			if strings.Contains(err.Error(), "1062") {
				err = errors.New("会员等级不能重复")
				return
			}
			return
		}
	}

	err = rds.DelRedisCacheByKey(constant.EventGiftConfig)

	return
}

// 删除二选一 配置
func DelEventConfig(req request.DelEventConfig) (err error) {
	db := global.Game
	event := new(game.EventGiftConfig)
	err = event.DeleteById(db, req.Id)
	if err != nil {
		return
	}
	err = rds.DelRedisCacheByKey(constant.EventGiftConfig)
	return
}

// 二选一开关
func OnOffEvent(req request.OnOffBenefit) (err error) {
	db := global.Game

	benefit := new(game.EventGiftConfig)

	var (
		dataList []game.EventGiftConfig
		ids      []int
		isClose  int
	)

	dataList, err = benefit.GetList(db)

	for _, d := range dataList {
		ids = append(ids, d.ID)
	}

	if req.Status > 0 {
		isClose = 1
	}

	mp := make(map[string]interface{})

	mp["is_close"] = isClose

	err = benefit.UpdateByIds(db, ids, mp)

	if err != nil {
		return
	}

	//删除缓存
	err = rds.DelRedisCacheByKey(consts.EventGiftConfig)

	return
}

// true 开启 falase 关闭
func EventStatus() (status bool, err error) {
	db := global.Game

	benefit := new(game.EventGiftConfig)

	var cfg game.EventGiftConfig

	cfg, err = benefit.GetFirstByIsClose(db, 0)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return
	}

	if cfg.ID > 0 {
		status = true
	}

	return

}

// 充200送200 配置列表
func GetRechargeGiftList() (list []game.RechargeGiftConfig, err error) {
	db := global.Game
	recharge := new(game.RechargeGiftConfig)
	return recharge.GetList(db)
}

// 保存充200送200配置
func SaveRechargeGiftConfig(req request.SaveRechargeGiftConfig) (err error) {

	recharge := new(game.RechargeGiftConfig)

	mp := make(map[string]interface{})

	mp["basic_rewards"] = req.BasicRewards
	mp["basic_type"] = req.BasicType
	mp["gift_rewards"] = req.GiftRewards
	mp["gift_type"] = req.GiftType
	mp["gift2_rewards"] = req.Gift2Rewards
	mp["gift2_type"] = req.Gift2Type
	mp["total"] = req.Total
	mp["price"] = req.Price
	mp["times"] = req.Times
	mp["interval"] = req.Interval
	mp["ratio"] = req.Ratio

	db := global.Game

	//只更新不添加
	err = recharge.UpdateById(db, req.Id, mp)

	if err != nil {
		return
	}

	err = rds.DelRedisCacheByKey(constant.RechargeGiftConfig)

	return
}

// 删除充200送200 配置
func DelRechargeGift(req request.DelRechargeGift) (err error) {
	db := global.Game
	recharge := new(game.RechargeGiftConfig)
	err = recharge.DeleteById(db, req.Id)
	if err != nil {
		return
	}

	err = rds.DelRedisCacheByKey(constant.RechargeGiftConfig)
	return
}

// 充200送200开关
func OnOffRechargeGift(req request.OnOffBenefit) (err error) {
	db := global.Game

	benefit := new(game.RechargeGiftConfig)

	var (
		dataList []game.RechargeGiftConfig
		ids      []int
		isClose  int
	)

	dataList, err = benefit.GetList(db)

	for _, d := range dataList {
		ids = append(ids, d.ID)
	}

	if req.Status > 0 {
		isClose = 1
	}

	mp := make(map[string]interface{})

	mp["is_close"] = isClose

	err = benefit.UpdateByIds(db, ids, mp)

	if err != nil {
		return
	}

	//删除缓存
	err = rds.DelRedisCacheByKey(consts.RechargeGiftConfig)

	return
}

// true 开启 falase 关闭
func RechargeGiftStatus() (status bool, err error) {
	db := global.Game

	benefit := new(game.RechargeGiftConfig)

	var cfg game.RechargeGiftConfig

	cfg, err = benefit.GetFirstByIsClose(db, 0)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return
	}

	if cfg.ID > 0 {
		status = true
	}

	return

}

// 充值礼包 配置列表
func GetRechargePackList() (list []game.RechargePackConfig, err error) {
	db := global.Game
	recharge := new(game.RechargePackConfig)
	return recharge.GetList(db)
}

// 保存充值礼包 配置
func SaveRechargePack(req request.SaveRechargePack) (err error) {
	db := global.Game
	recharge := game.RechargePackConfig{
		GameId:       req.GameId,
		RoomId:       req.RoomId,
		BasicRewards: req.BasicRewards,
		BasicType:    req.BasicType,
		GiftRewards:  req.GiftRewards,
		GiftType:     req.GiftType,
		Bonus:        req.Bonus,
		BonusType:    3,
		Total:        req.Total,
		Price:        req.Price,
		Times:        req.Times,
		Interval:     req.Interval,
		Ratio:        req.Ratio,
		CreatedAt:    time.Now(),
	}

	if req.Id > 0 {
		recharge.ID = req.Id
	}

	err = recharge.Save(db)

	if err != nil {
		if err != nil {
			if strings.Contains(err.Error(), "1062") {
				err = errors.New("游戏房间内礼包价格不能重复")
				return
			}
			return
		}
	}

	//删除缓存
	err = rds.DelRedisCacheByKey(constant.RechargePackConfigGameRoom + fmt.Sprintf("%v%v", req.GameId, req.RoomId))

	return
}

// 删除充值礼包
func DelRechargePack(req request.DelRechargePack) (err error) {
	db := global.Game
	recharge := new(game.RechargePackConfig)
	info, err := recharge.GetFirstById(db, req.Id)
	if err != nil {
		return
	}

	err = recharge.DeleteById(db, req.Id)
	if err != nil {
		return
	}

	//删除缓存
	err = rds.DelRedisCacheByKey(constant.RechargePackConfigGameRoom + fmt.Sprintf("%v%v", info.GameId, info.RoomId))
	return
}

// 救济金 配置列表
func GetBenefitList(req request.GetBenefitList) (res response.BenefitRsp, err error) {
	db := global.Game
	benefit := new(game.BenefitGiftPackConfig)

	total, dataList, err := benefit.GetPageList(db, req)
	if err != nil {
		return
	}

	res.Total = total

	list := make([]response.Benefit, 0, len(dataList))

	for _, d := range dataList {
		list = append(list, response.Benefit{
			Id:              d.ID,
			UserType:        d.UserType,
			Minimum:         d.Minimum,
			Maximum:         d.Maximum,
			MiniTimes:       d.MiniTimes,
			MaxiTimes:       d.MaxiTimes,
			Quota:           d.Quota,
			Value:           d.Value,
			BasicRewards:    d.BasicRewards,
			BasicType:       d.BasicType,
			RewardGiveaways: d.RewardGiveaways,
			GiftType:        d.GiftType,
			Bonus:           d.Bonus,
			BonusType:       d.BonusType,
			Ratio:           d.Ratio,
		})
	}

	res.List = list

	return
}

// 保存救济金配置
func SaveBenefit(req request.SaveBenefit) (err error) {
	db := global.Game
	benefit := game.BenefitGiftPackConfig{
		UserType:        req.UserType,
		Minimum:         req.Minimum,
		Maximum:         req.Maximum,
		MiniTimes:       req.MiniTimes,
		MaxiTimes:       req.MaxiTimes,
		Quota:           req.Quota,
		Value:           req.Value,
		BasicRewards:    req.BasicRewards,
		BasicType:       req.BasicType,
		RewardGiveaways: req.RewardGiveaways,
		GiftType:        req.GiftType,
		Bonus:           req.Bonus,
		BonusType:       3,
		Ratio:           req.Ratio,
		CreatedAt:       time.Now(),
	}

	if req.Id > 0 {
		benefit.ID = req.Id
	}
	err = benefit.Save(db)

	if err != nil {
		return
	}

	err = rds.DelRedisCacheByKey(constant.BenefitGiftPackConfig)

	return
}

// 更新救济金礼包开关
func OnOffBenefit(req request.OnOffBenefit) (err error) {
	db := global.Game

	benefit := new(game.BenefitGiftPackConfig)

	var (
		dataList []game.BenefitGiftPackConfig
		ids      []int
		isClose  int
	)

	dataList, err = benefit.GetList(db)

	for _, d := range dataList {
		ids = append(ids, d.ID)
	}

	if req.Status > 0 {
		isClose = 1
	}

	mp := make(map[string]interface{})

	mp["is_close"] = isClose

	err = benefit.UpdateByIds(db, ids, mp)

	if err != nil {
		return
	}

	//删除缓存
	err = rds.DelRedisCacheByKey(consts.BenefitGiftPackConfig)

	return
}

// true 开启 falase 关闭
func BenefitStatus() (status bool, err error) {
	db := global.Game

	benefit := new(game.BenefitGiftPackConfig)

	var cfg game.BenefitGiftPackConfig

	cfg, err = benefit.GetFirstByIsClose(db, 0)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return
	}

	if cfg.ID > 0 {
		status = true
	}

	return

}

// 删除救济金
func DelBenefit(req request.DelBenefit) (err error) {
	db := global.Game
	benefit := new(game.BenefitGiftPackConfig)
	err = benefit.DeleteById(db, req.Id)
	if err != nil {
		return
	}

	err = rds.DelRedisCacheByKey(constant.BenefitGiftPackConfig)
	return
}

// 三选一礼包
func GetOnlyList() (list []game.OnlyOneConfig, err error) {
	db := global.Game
	only := new(game.OnlyOneConfig)
	return only.GetList(db)
}

// 保存三选一
func SaveOnly(req request.SaveOnly) (err error) {
	db := global.Game
	only := game.OnlyOneConfig{
		Grade:     req.Grade,
		FirstDay:  req.FirstDay,
		FirstType: req.FirstType,
		NextDay:   req.NextDay,
		NextType:  req.NextType,
		ThirdDay:  req.ThirdDay,
		ThirdType: req.ThirdType,
		Ratio:     req.Ratio,
		CreatedAt: time.Now(),
	}

	if req.Id > 0 {
		only.ID = req.Id
	}

	err = only.Save(db)

	if err != nil {
		return
	}

	err = rds.DelRedisCacheByKey(constant.OnlyOneConfig)
	return
}

func DelOnly(req request.DelOnly) (err error) {
	db := global.Game
	benefit := new(game.OnlyOneConfig)
	err = benefit.DeleteById(db, req.Id)

	if err != nil {
		return
	}

	err = rds.DelRedisCacheByKey(constant.OnlyOneConfig)

	return
}

// 更新三选一开关
func OnOffOnly(req request.OnOffBenefit) (err error) {
	db := global.Game

	benefit := new(game.OnlyOneConfig)

	var (
		dataList []game.OnlyOneConfig
		ids      []int
		isClose  int
	)

	dataList, err = benefit.GetList(db)

	for _, d := range dataList {
		ids = append(ids, d.ID)
	}

	if req.Status > 0 {
		isClose = 1
	}

	mp := make(map[string]interface{})

	mp["is_close"] = isClose

	err = benefit.UpdateByIds(db, ids, mp)

	if err != nil {
		return
	}

	//删除缓存
	err = rds.DelRedisCacheByKey(consts.OnlyOneConfig)

	return
}

// true 开启 falase 关闭
func OnlyStatus() (status bool, err error) {
	db := global.Game

	benefit := new(game.OnlyOneConfig)

	var cfg game.OnlyOneConfig

	cfg, err = benefit.GetFirstByIsClose(db, 0)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return
	}

	if cfg.ID > 0 {
		status = true
	}

	return

}
