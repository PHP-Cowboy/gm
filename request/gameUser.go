package request

type GetLoginLogList struct {
	Paging
	Uid        int      `json:"uid" form:"uid"`
	ChannelId  int      `json:"channel_id" form:"channel_id"`
	LoginMode  int      `json:"login_mode" form:"login_mode"`
	StartTime  string   `json:"start_time" form:"start_time"`
	EndTime    string   `json:"end_time" form:"end_time"`
	CreatedAt  []string `json:"created_at[]" form:"created_at[]"`
	ChannelIds []int
}

type GetGameUserList struct {
	Paging
	Uid            int      `json:"uid" form:"uid"`
	ChannelId      int      `json:"channel_id" form:"channel_id"`
	IsGuest        *uint8   `json:"is_guest" form:"is_guest"`
	PayStatus      *int     `json:"pay_status" form:"pay_status"`
	AssetMin       int      `json:"asset_min" form:"asset_min"`
	AssetMax       int      `json:"asset_max" form:"asset_max"`
	RechargeMin    int      `json:"recharge_min" form:"recharge_min"`
	RechargeMax    int      `json:"recharge_max" form:"recharge_max"`
	GiftCoinMin    int      `json:"gift_coin_min" form:"gift_coin_min"`
	GiftCoinMax    int      `json:"gift_coin_max" form:"gift_coin_max"`
	CreatedAt      []string `json:"created_at" form:"created_at[]"`
	StartCreatedAt string   `json:"start_created_at"`
	EndCreatedAt   string   `json:"end_created_at"`
	UpdatedAt      []string `json:"updated_at" form:"updated_at[]"`
	StartUpdatedAt string   `json:"start_updated_at"`
	EndUpdatedAt   string   `json:"end_updated_at"`
	UserIds        []int    `json:"user_ids"`
	ChannelIds     []int
}

type GiveList struct {
	Paging
	Uid               int      `json:"uid" form:"uid"`
	ChannelId         int      `json:"channel_id" form:"channel_id"`
	PayChannelId      int      `json:"pay_channel_id" form:"pay_channel_id"`
	Status            *int     `json:"status" form:"status"`
	PayStatus         *int     `json:"pay_status" form:"pay_status"`
	GiveRateMin       int      `json:"give_rate_min" form:"give_rate_min"`
	GiveRateMax       int      `json:"give_rate_max" form:"give_rate_max"`
	CommitGiveRateMin int      `json:"commit_give_rate_min" form:"commit_give_rate_min"`
	CommitGiveRateMax int      `json:"commit_give_rate_max" form:"commit_give_rate_max"`
	CreatedAt         []string `json:"created_at" form:"created_at[]"`
	StartCreatedAt    string   `json:"start_created_at"`
	EndCreatedAt      string   `json:"end_created_at"`
	ChannelIds        []int
}

type XPaymentCallback struct {
	TransferId  string  `json:"transferId" form:"transferId" binding:"required"`   //平台支付系统订单号
	MchNo       string  `json:"mchNo" form:"mchNo" binding:"required"`             //商户号
	AppId       string  `json:"appId" form:"appId" binding:"required"`             //商户应用Id
	MchOrderNo  string  `json:"mchOrderNo" form:"mchOrderNo" binding:"required"`   //返回商户传入的订单号
	OrderAmount float64 `json:"orderAmount" form:"orderAmount" binding:"required"` //金额
	Currency    string  `json:"currency" form:"currency" binding:"required"`       //货币代码
	State       int     `json:"state" form:"state" binding:"required"`             //订单状态 1-代付中，2-代付成功，3-代付失败，4-代付撤销
	CreatedAt   int     `json:"createdAt" form:"createdAt" binding:"required"`     //订单创建时间，13位时间戳
	AccountNo   string  `json:"accountNo" form:"accountNo"`                        //账号信息
	AccountName string  `json:"accountName" form:"accountName" binding:"required"` //客户姓名
	ErrMsg      string  `json:"errMsg" form:"errMsg"`                              //返回错误描述
	Utr         string  `json:"utr" form:"utr"`                                    //utr凭证
	Vpa         string  `json:"vpa" form:"vpa"`                                    //UPi账号
	SuccessTime int     `json:"successTime" form:"successTime"`                    //订单支付成功时间，13位时间戳
	ReqTime     int     `json:"reqTime" form:"reqTime" binding:"required"`         //13位时间戳
	Sign        string  `json:"sign" form:"sign" binding:"required"`               //签名值，详见签名算法，不参与签名
}

type XPaymentCallbackSuccess struct {
	OrderNo string `json:"order_no" form:"order_no"`
}

type PaymentCallbackFailed struct {
	MchOrderNo string `json:"mchOrderNo" form:"mchOrderNo" binding:"required"` //返回商户传入的订单号
	State      int    `json:"state" form:"state" binding:"required"`           //订单状态 1-代付中，2-代付成功，3-代付失败，4-代付撤销
}

type InPaymentCallback struct {
	Status      interface{} `json:"status"`
	Message     string      `json:"message"`
	Money       string      `json:"money"`
	PlatNumber  string      `json:"plat_number"`
	OrderNumber string      `json:"order_number"`
	Utr         interface{} `json:"utr"`
	Sign        string      `json:"sign"`
}

type ChangeRecharge struct {
	Uid          int  `json:"uid" binding:"required"`
	Recharge     *int `json:"recharge" binding:"required"`
	EditRecharge *int `json:"edit_recharge" binding:"required"`
}

type WithdrawInfoRecord struct {
	Paging
	Uid        uint64 `json:"uid" form:"uid" binding:"required"`
	ChannelIds []int
}

type EditUserCoin struct {
	Uid      int    `json:"uid" binding:"required"`
	CoinType int    `json:"coin_type"` //币类型
	OpType   int    `json:"op_type"`   //操作类型
	Num      int    `json:"num"`
	Title    string `json:"title"`
	Msg      string `json:"msg"`
}

type Banned struct {
	Admin int    `json:"admin"` //管理后台操作人id
	Uid   int    `json:"uid"`   //游戏用户
	Cate  int    `json:"cate"`
	Info  string `json:"info"`
}

type GiveMoneyHandle struct {
	OrderNo string `json:"order_no" form:"order_no" binding:"required"`
}
