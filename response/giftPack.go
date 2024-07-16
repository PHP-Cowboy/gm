package response

type UserEventGiftConfig struct {
	Id        string `json:"id"`
	Grade     int    `json:"grade"`
	Coin      int    `json:"coin"`
	CoinType  uint8  `json:"coin_type"`
	GiftLimit int    `json:"gift_limit"`
	GiftType  uint8  `json:"gift_type"`
	Bonus     int    `json:"bonus"`
	BonusType uint8  `json:"bonus_type"`
	Ratio     int    `json:"ratio"`
	Type      int    `json:"type"`
}

type RechargeGiftConfig struct {
	Id           string `json:"id"`
	BasicRewards int    `json:"basic_rewards"` //基础奖励
	BasicType    uint8  `json:"basic_type"`    //金额类别(1=可提现可下注,2=不可提现可下注)
	GiftRewards  int    `json:"gift_rewards"`  // 赠送额度
	GiftType     uint8  `json:"gift_type"`     //金额类别(1=可提现可下注,2=不可提现可下注)
	Bonus        int    `json:"bonus"`         //bonus
	BonusType    uint8  `json:"bonus_type"`    //金额类别(1=可提现可下注,2=不可提现可下注)
	Total        int    `json:"total"`         //总额
	Price        int    `json:"price"`         //价格
	Times        int    `json:"times"`         //出现次数
	Interval     int    `json:"interval"`      //间隔时长(s)
	Ratio        int    `json:"ratio"`         //折扣比例
}

type BenefitRsp struct {
	Total int64     `json:"total"`
	List  []Benefit `json:"list"`
}

type Benefit struct {
	Id              int    `json:"id"`
	UserType        string `json:"user_type"`
	Minimum         int    `json:"minimum"`
	Maximum         int    `json:"maximum"`
	MiniTimes       int    `json:"mini_times"`
	MaxiTimes       int    `json:"maxi_times"`
	Quota           int    `json:"quota"`
	Value           int    `json:"value"`
	BasicRewards    int    `json:"basic_rewards"`
	BasicType       uint8  `json:"basic_type"`
	RewardGiveaways int    `json:"reward_giveaways"`
	GiftType        uint8  `json:"gift_type"`
	Bonus           int    `json:"bonus"`
	BonusType       uint8  `json:"bonus_type"`
	Ratio           int    `json:"ratio"`
}
