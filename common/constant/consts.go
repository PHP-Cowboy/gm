package constant

// sign
const (
	SignConfig   = "hall:sign:config:list"
	Prize        = "hall:prize:list"
	SignPrizeDay = "hall:sign:prize:day:" //第x天签到奖励
)

// vip
const (
	VipConfig = "hall:vipConfig" //vip配置
)

// 活动礼包
const (
	UserEventGift   = "hall:eventGift:user:"
	EventGiftConfig = "hall:eventGiftConfig"
)

// 救济金&&破产礼包
const (
	Benefit               = "hall:benefit"
	UserClaimedTimes      = "hall:claimedTimes:user" //用户领取次数
	UserLastReceiveTime   = "hall:lastReceiveTime:user"
	BenefitGiftPackConfig = "hall:benefitGiftPackConfig" //破产礼包配置
	UserBankruptcyTimes   = "hall:bankruptcyTimes:user"  //用户破产未购买礼包次数
)

// onlyOne 三选一
const (
	OnlyOneConfig              = "hall:onlyOneConfig"
	UserOnlyOneRecordDayStatus = "hall:onlyOneRecord:day:status"
	UserOnlyOneRecordDay       = "hall:onlyOneRecord:user:day:"
	UserOnlyOneRecord          = "hall:onlyOneRecord:user:"
)

// 充值赠送礼包
const (
	RechargeGiftConfig              = "hall:rechargeGiftConfig"
	RechargeGiftInterval            = "hall:rechargeGiftInterval"
	DayUserRechargeGiftTimes        = "hall:rechargeGiftTimes:user:day:"
	UserLastReceiveRechargeGiftTime = "hall:lastReceiveRechargeGiftTime:user"
)

// 充值礼包
const (
	RechargePackConfigId                   = "hall:rechargePackConfig:id:"
	RechargePackConfigGameRoom             = "hall:rechargePackConfig:game:room:"
	UserRechargePackConfigGameRoomDayTimes = "hall:rechargePackConfig:game:room:day:userTimes" //游戏房间内用户购买礼包次数
)

// 签到
const (
	UserTodayHasSigned = "hall:hasSigned:day:user:" //用户今日是否已签到
)

// 邮件
const (
	CommonEmailList = "hall:email"
)

// 赠送
const (
	UserTodayGiveTotal = "backend:give:day:user:total:"
)

const (
	PayOrderUrl          = "/api/pay/unifiedOrder" //支付下单
	XPayBalanceUrl       = "/api/mch/balance"      //查询余额
	XPayOrderStatusUrl   = "/api/pay/query"        //查询代收订单状态
	XPayPaymentUrl       = "/api/transferOrder"    // 代付下单
	XPayPaymentStatusUrl = "/api/transfer/query"   //代付状态查询
)
