package request

type SaveVipConfig struct {
	Id             int    `json:"id"`               //id
	Level          uint8  `json:"level"`            //等级
	NeedExp        uint64 `json:"need_exp"`         // 升级所需经验
	WithdrawNums   uint8  `json:"withdraw_nums"`    // 每日提现次数
	WithdrawMoney  uint64 `json:"withdraw_money"`   // 每日提现金额
	DayPrizeType   uint8  `json:"day_prize_type"`   // 每日奖品类型
	DayPrizeNums   uint64 `json:"day_prize_nums"`   // 每日奖品数量
	WeekPrizeType  uint8  `json:"week_prize_type"`  // 每周奖品类型
	WeekPrizeNums  uint64 `json:"week_prize_nums"`  // 每周奖品数量
	MonthPrizeType uint8  `json:"month_prize_type"` // 每月奖品类型
	MonthPrizeNums uint64 `json:"month_prize_nums"` // 每月奖品数量
}

type DelVipConfig struct {
	Id int `json:"id"`
}

type GetBenefitList struct {
	Paging
}

type SaveEventGiftConfig struct {
	Id        int   `json:"id"` //id
	Grade     int   `json:"grade" binding:"required"`
	Type      int   `json:"type" binding:"required"`
	Coin      int   `json:"coin" binding:"required"`
	CoinType  uint8 `json:"coin_type" binding:"required"`
	GiftLimit int   `json:"gift_limit" binding:"required"`
	GiftType  uint8 `json:"gift_type" binding:"required"`
	Bonus     int   `json:"bonus"`
	BonusType uint8 `json:"bonus_type"`
	Ratio     int   `json:"ratio" binding:"required"`
}

type DelEventConfig struct {
	Id int `json:"id"` //id
}

type SaveRechargeGiftConfig struct {
	Id           int   `json:"id"`
	BasicRewards int   `json:"basic_rewards"` //基础奖励
	BasicType    uint8 `json:"basic_type"`    //金额类别(1=可提现可下注,2=不可提现可下注)
	GiftRewards  int   `json:"gift_rewards"`  // 赠送额度
	GiftType     uint8 `json:"gift_type"`     //金额类别(1=可提现可下注,2=不可提现可下注)
	Gift2Rewards int   `json:"gift2_rewards"` //赠送额度2
	Gift2Type    uint8 `json:"gift2_type"`    //金额类别(1=可提现可下注,2=不可提现可下注)
	Total        int   `json:"total"`         //总额度
	Price        int   `json:"price"`         //价格
	Times        int   `json:"times"`         //出现次数
	Interval     int   `json:"interval"`      //间隔时长(s)
	Ratio        int   `json:"ratio"`         //折扣比例
}

type DelRechargeGift struct {
	Id int `json:"id"` //id
}

type SaveRechargePack struct {
	Id           int   `json:"id"`
	GameId       int   `json:"game_id"`       //游戏id
	RoomId       int   `json:"room_id"`       //房间id
	BasicRewards int   `json:"basic_rewards"` //基础奖励
	BasicType    uint8 `json:"basic_type"`    //金额类别(1=可提现可下注,2=不可提现可下注)
	GiftRewards  int   `json:"gift_rewards"`  // 赠送额度
	GiftType     uint8 `json:"gift_type"`     //金额类别(1=可提现可下注,2=不可提现可下注)
	Bonus        int   `json:"bonus"`         //bonus
	BonusType    uint8 `json:"bonus_type"`    //金额类别(1=可提现可下注,2=不可提现可下注)
	Total        int   `json:"total"`         //总额度
	Price        int   `json:"price"`         //价格
	Times        int   `json:"times"`         //出现次数
	Interval     int   `json:"interval"`      //间隔时长(s)
	Ratio        int   `json:"ratio"`         //折扣比例
}

type DelRechargePack struct {
	Id uint64 `json:"id"`
}

type SaveBenefit struct {
	Id              int    `json:"id"`
	UserType        string `json:"user_type"`        //用户类型
	Minimum         int    `json:"minimum"`          //累计充值最小额度
	Maximum         int    `json:"maximum"`          //累计充值最大额度
	MiniTimes       int    `json:"mini_times"`       //连续触发破产礼包但未购买次数
	MaxiTimes       int    `json:"maxi_times"`       //连续触发破产礼包但未购买次数
	Quota           int    `json:"quota"`            //礼包额度
	Value           int    `json:"value"`            //价值
	BasicRewards    int    `json:"basic_rewards"`    //基础奖励
	BasicType       uint8  `json:"basic_type"`       //金额类别(1=可提现可下注,2=不可提现可下注)
	RewardGiveaways int    `json:"reward_giveaways"` //赠送奖励
	GiftType        uint8  `json:"gift_type"`        //金额类别(1=可提现可下注,2=不可提现可下注)
	Bonus           int    `json:"bonus"`            //bonus
	BonusType       uint8  `json:"bonus_type"`       //金额类别(1=可提现可下注,2=不可提现可下注)
	Ratio           int    `json:"ratio"`            //优惠比例
}

type OnOffBenefit struct {
	Status int `json:"status"`
}

type DelBenefit struct {
	Id int `json:"id"`
}

type SaveOnly struct {
	Id        int   `json:"id"`
	Grade     int   `json:"grade"`     // 档位
	FirstDay  int   `json:"first_day"` // 第一天领取
	FirstType uint8 `json:"first_type"`
	NextDay   int   `json:"next_day"` // 第二天领取
	NextType  uint8 `json:"next_type"`
	ThirdDay  int   `json:"third_day"` // 第三天领取
	ThirdType uint8 `json:"third_type"`
	Ratio     int   `json:"ratio"` //折扣比例
}

type DelOnly struct {
	Id int `json:"id"`
}
