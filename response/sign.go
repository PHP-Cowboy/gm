package response

type SignConfigList struct {
	Id           uint64   `json:"id"`
	Name         string   `json:"name"`           //签到名称
	SignNum      uint8    `json:"sign_num"`       //累计签到次数
	PrizeIds     []string `json:"prize_ids"`      //奖励ID
	PrizeCashId  int      `json:"prize_cash_id"`  //cash
	PrizeBonusId int      `json:"prize_bonus_id"` //bonus
	Unit         string   `json:"unit"`           //单位
	Remark       string   `json:"remark"`         //备注
	Cash         string   `json:"cash"`
	Bonus        string   `json:"bonus"`
}

type SingPrizeRsp struct {
	Total int64       `json:"total"`
	List  []SingPrize `json:"list"`
}

type SingPrize struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	EnName    string `json:"en_name"`
	Type      uint8  `json:"type"`
	GoodsNum  int    `json:"goods_num"`
	GoodsType uint8  `json:"goods_type"`
	Unit      string `json:"unit"`
	Remark    string `json:"remark"`
}
