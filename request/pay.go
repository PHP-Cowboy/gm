package request

type PayList struct {
	Name    string `json:"name" form:"name"`
	Cash    int    `json:"cash" form:"cash"`
	Account int    `json:"account" form:"account"`
	Status  int8   `json:"status" form:"status"`
	Type    int    `json:"type" form:"type"`
}

type SaveGift struct {
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

type SaveConfig struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`             // 支付渠道名称
	Icon           string `json:"icon"`             // 支付图标
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

type BankList struct {
	Paging
	Uid        uint64 `json:"uid" form:"uid"`               // 用户ID
	BankCode   string `json:"bank_code" form:"bank_code"`   // 银行类型
	BankName   string `json:"bank_name" form:"bank_name"`   // 银行名称
	AccountNo  string `json:"account_no" form:"account_no"` // 银行账号
	Ifsc       string `json:"ifsc" form:"ifsc"`             // ifsc号
	Name       string `json:"name" form:"name"`             //客户姓名
	Email      string `json:"email" form:"email"`           //客户邮箱
	Phone      string `json:"phone" form:"phone"`           //客户手机
	ChannelIds []int
}

type OrderList struct {
	Paging
	Uid     int    `json:"uid" form:"uid"` // 用户ID
	Ymd     int    `json:"ymd"`            // 年月日
	OrderNo string `json:"order_no"`       // 订单ID
	Name    string `json:"name"`           // 用户名
	Phone   string `json:"phone"`          //手机
	Status  int8   `json:"status"`         // 状态 0=等待支付，1=支付完成，2=下单失败
}

type ConfigList struct {
	Paging
	Name string `json:"name" form:"name"`
}

type RechargeRecords struct {
	Paging
	Uid        int      `json:"uid" form:"uid"`
	Ymd        []string `json:"ymd[]" form:"ymd[]"`
	Start      string   `json:"start"`
	End        string   `json:"end"`
	OrderNo    string   `json:"order_no" form:"order_no"`
	Status     *int8    `json:"status" form:"status"`
	Channel    int      `json:"channel" form:"channel"`
	ChannelIds []int
}

type GaveConfigList struct {
	Status int `json:"status" form:"status"`
	Type   int `json:"type" form:"type"`
}

type GaveConfigPageList struct {
	Paging
	Name   string `json:"name" form:"name"`
	Status int    `json:"status" form:"status"`
	Type   int    `json:"type" form:"type"`
}

type SaveGaveConfig struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Account   int    `json:"account"`
	Status    int    `json:"status"`
	Type      int    `json:"type"`
	ReplaceId int    `json:"replace_id"`
}
