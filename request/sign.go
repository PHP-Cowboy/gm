package request

type SaveSign struct {
	ID           uint64   `json:"id" form:"id"`
	Name         string   `json:"name" form:"name"`           //签到名称
	SignNum      uint8    `json:"sign_num" form:"sign_num"`   //累计签到次数
	PrizeIds     []string `json:"prize_ids" form:"prize_ids"` //奖励ID
	PrizeCashId  int      `json:"prize_cash_id"`
	PrizeBonusId int      `json:"prize_bonus_id"`
	Unit         string   `json:"unit" form:"unit"`     //单位
	Remark       string   `json:"remark" form:"remark"` //备注
}

type DelSign struct {
	Id uint64 `json:"id" binding:"required"`
}

type GetSingPrizeList struct {
	Paging
}

type SavePrize struct {
	ID        uint64 `json:"id" form:"id"`
	Name      string `json:"name" form:"name"`                                //奖励名称
	EnName    string `json:"en_name" form:"en_name"`                          //奖励英文名称
	Type      uint8  `json:"type" form:"type"`                                //类型(1=筹码)
	GoodsNum  int    `json:"goods_num" form:"goods_num" binding:"required"`   //金额
	GoodsType uint8  `json:"goods_type" form:"goods_type" binding:"required"` //金额类别(1=可提现可下注,2=不可提现可下注)
	Remark    string `json:"remark" form:"remark"`                            //备注
}

type DelSingPrize struct {
	Id uint64 `json:"id" binding:"required"`
}
