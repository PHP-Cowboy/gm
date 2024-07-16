package response

type PayGift struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`           // 支付包名称
	Cash         int    `json:"cash"`           // 到账金额
	Account      int    `json:"account"`        // 支付金额
	Status       int8   `json:"status"`         // 状态（1=可用，0=不可用）
	AddMoney     int    `json:"add_money"`      //额外赠送无限制
	AddMoneyType int8   `json:"add_money_type"` //无限制赠送类别（1=金额，2=比例）
	AddCash      int    `json:"add_cash"`       //赠送cash
	AddCashType  int8   `json:"add_cash_type"`  //赠送cash金额（1=金额，2=比例）
	Bonus        int    `json:"bonus"`          //储钱罐
	BonusType    int8   `json:"bonus_type"`     //储钱罐类型（1=金额，2=比例）
	Ratio        int    `json:"ratio"`          //优惠比例
	Remark       string `json:"remark"`         //备注
	Type         int    `json:"type"`
	ReplaceId    int    `json:"replace_id"`
}

type BankRsp struct {
	Total int64  `json:"total"`
	List  []Bank `json:"list"`
}

type Bank struct {
	ID        uint64 `json:"id"`
	Uid       uint64 `json:"uid"`        // 用户ID
	BankCode  string `json:"bank_code"`  // 银行类型
	BankName  string `json:"bank_name"`  // 银行名称
	AccountNo string `json:"account_no"` // 银行账号
	Ifsc      string `json:"ifsc"`       // ifsc号
	Name      string `json:"name"`       //客户姓名
	Email     string `json:"email"`      //客户邮箱
	Phone     string `json:"phone"`      //客户手机
	Address   string `json:"address"`    //客户地址
	Vpa       string `json:"vpa"`        //vpa
	Remark    string `json:"remark"`     //备注
	CreatedAt string `json:"created_at"`
}

type OrderList struct {
	ID           int    `json:"id"`
	Uid          int    `json:"uid"`           // 用户ID
	Ymd          int    `json:"ymd"`           // 年月日
	OrderNo      string `json:"order_no"`      // 订单ID
	MOrderNo     string `json:"m_order_no"`    // 商户订单ID
	Account      int    `json:"account"`       // 支付金额
	Cash         int    `json:"cash"`          //cash金额
	GiftCash     int    `json:"gift_cash"`     //额外赠送cash
	Bonus        int    `json:"bonus"`         //储钱罐
	RequestTime  string `json:"request_time"`  //下单时间
	Email        string `json:"email"`         // 邮箱地址
	Name         string `json:"name"`          // 用户名
	Phone        string `json:"phone"`         //手机
	RedirectTime string `json:"redirect_time"` //下单拿h5地址时间
	Status       int8   `json:"status"`        // 状态 0=等待支付，1=支付完成，2=下单失败
	CompleteTime int    `json:"complete_time"` //订单完成时间
	H5Url        string `json:"h_5_url"`       // h5地址
	Type         int    `json:"type"`          // 类型 1:充值;2:礼包
	Remark       string `json:"remark"`        //备注
}

type PayConfigRsp struct {
	Total int64       `json:"total"`
	List  []PayConfig `json:"list"`
}

type PayConfig struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`             // 支付渠道名称
	Url            string `json:"url"`              // 支付host
	BackUrl        string `json:"back_url"`         // 代收回调地址
	PaymentBackUrl string `json:"payment_back_url"` // 代付回调地址
	AppId          string `json:"app_id"`           //appID
	Secret         string `json:"secret"`           //secret 秘钥
	Merchant       string `json:"merchant"`         //商户号
	Status         int8   `json:"status"`           //状态 (1=可用，0=不可用)
	Remark         string `json:"remark"`           //备注
	Markers        string `json:"markers"`          //调用接口标记
}

type RechargeRecords struct {
	Total int64                 `json:"total"`
	List  []RechargeRecordsList `json:"list"`
}

type RechargeRecordsList struct {
	Id           int    `json:"id"`
	Uid          int    `json:"uid"`
	Ymd          int    `json:"ymd"`
	Channel      string `json:"channel"`
	ChannelNo    string `json:"channel_no"`
	OrderNo      string `json:"order_no"`
	MOrderNo     string `json:"m_order_no"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Account      int    `json:"account"`
	Cash         int    `json:"cash"`      //cash金额
	GiftCash     int    `json:"gift_cash"` //额外赠送cash
	Bonus        int    `json:"bonus"`     //bonus
	Name         string `json:"name"`
	Status       int8   `json:"status"` //状态 0=等待支付，1=支付完成，2=下单失败
	CompleteTime string `json:"complete_time"`
	Type         int    `json:"type"`
	Remark       string `json:"remark"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type GaveConfigRsp struct {
	Total int64        `json:"total"`
	List  []GaveConfig `json:"list"`
}

type GaveConfig struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Account   int    `json:"account"`
	Status    int    `json:"status"`
	Type      int    `json:"type"`
	ReplaceId int    `json:"replace_id"`
}
