package response

type XPayBalanceReturn struct {
	Balance     float64 `json:"balance"`
	MoneyFrozen float64 `json:"moneyFrozen"`
}

type XPayOrderStatusReturn struct {
	Uid          uint64 `json:"uid"`
	OrderNo      string `json:"orderNo"`
	MOrderNo     string `json:"mOrderNo"`
	Account      int    `json:"account"`
	RequestTime  string `json:"requestTime"`
	RedirectTime string `json:"redirectTime"`
	CompleteTime string `json:"completeTime"`
	Status       int8   `json:"status"`
}

type XPayPaymentOrderReturn struct {
	Uid         uint64 `json:"uid"`
	OrderNo     string `json:"orderNo"`
	MOrderNo    string `json:"mOrderNo"`
	OrderAmount int    `json:"orderAmount"`
	RequestTime string `json:"request_time"` //下单时间
	OrderTime   string `json:"order_time"`
}

type XPayPaymentStatusReturn struct {
	Uid          uint64 `json:"uid"`
	OrderNo      string `json:"orderNo"`
	MOrderNo     string `json:"mOrderNo"`
	OrderAmount  int    `json:"orderAmount"`
	RequestTime  string `json:"request_time"` //下单时间
	OrderTime    string `json:"order_time"`
	CompleteTime string `json:"completeTime"`
	Status       int8   `json:"status"`
}

//-------------------------------------------------

// 查询余额请求结构体
type BalanceRequest struct {
	MchNo   string `json:"mchNo"`
	AppId   string `json:"appId"`
	ReqTime string `json:"reqTime"`
	Sign    string `json:"sign"`
}

// 查询余额返回
type BalanceReturn struct {
	Code uint16      `json:"code"`
	Data BalanceData `json:"data"`
	Msg  string      `json:"msg"`
	Sign string      `json:"sign"`
}

// 查询余额返回里面数据
type BalanceData struct {
	MchBalance     float64 `json:"mchBalance"`
	MchMoneyFrozen float64 `json:"mchMoneyFrozen"`
	MchName        string  `json:"mchName"`
	MchNo          string  `json:"mchNo"`
}

// 查询代收订单请
type OrderStatusRequest struct {
	MchNo      string `json:"mchNo"`
	AppId      string `json:"appId"`
	PayOrderId string `json:"payOrderId"`
	ReqTime    string `json:"reqTime"`
	Sign       string `json:"sign"`
}

// 查询代收订单返回
type OrderStatusReturn struct {
	Code uint16          `json:"code"`
	Data OrderStatusData `json:"data"`
	Msg  string          `json:"msg"`
	Sign string          `json:"sign"`
}

// 查询代收订单返回数据
type OrderStatusData struct {
	AppId       string `json:"appId"`
	CreatedAt   int64  `json:"createdAt"`
	Currency    string `json:"currency"`
	MchNo       string `json:"mchNo"`
	MchOrderNo  string `json:"mchOrderNo"`
	OrderAmount int    `json:"orderAmount"`
	PayOrderId  string `json:"payOrderId"`
	State       int    `json:"state"`
	SuccessTime int64  `json:"successTime"`
	ErrCode     int    `json:"errCode"`
	ErrMsg      string `json:"errMsg"`
}

// 代付订单请求
type PaymentOrderRequest struct {
	MchNo           string `json:"mchNo"`
	AppId           string `json:"appId"`
	MchOrderNo      string `json:"mchOrderNo"`
	Currency        string `json:"currency"`
	OrderAmount     string `json:"orderAmount"`
	BankCode        string `json:"bankCode"`
	BankName        string `json:"bankName,omitempty"`
	AccountNo       string `json:"accountNo,omitempty"`
	Ifsc            string `json:"ifsc,omitempty"`
	AccountName     string `json:"accountName"`
	AccountEmail    string `json:"accountEmail"`
	AccountMobileNo string `json:"accountMobileNo"`
	AccountAddress  string `json:"accountAddress,omitempty"`
	Vpa             string `json:"vpa,omitempty"`
	TransferDesc    string `json:"transferDesc,omitempty"`
	NotifyUrl       string `json:"notifyUrl"`
	ReqTime         string `json:"reqTime"`
	Sign            string `json:"sign"`
}

// 代付下单请求返回
type PaymentOrderReturn struct {
	Code *uint16          `json:"code"`
	Data PaymentOrderData `json:"data"`
	Msg  string           `json:"msg"`
	Sign string           `json:"sign"`
}

// 代付下单返回数据
type PaymentOrderData struct {
	AccountName string `json:"accountName"`
	AccountNo   string `json:"accountNo"`
	ErrCode     int    `json:"errCode"`
	ErrMsg      string `json:"errMsg"`
	MchOrderNo  string `json:"mchOrderNo"`
	OrderAmount int    `json:"orderAmount"`
	State       int    `json:"state"`
	TransferId  string `json:"transferId"`
}

// 查询代付订单请求
type PaymentStatusRequest struct {
	MchNo      string `json:"mchNo"`
	AppId      string `json:"appId"`
	TransferId string `json:"transferId"`
	ReqTime    string `json:"reqTime"`
	Sign       string `json:"sign"`
}

// 查询代付订单返回
type PaymentStatusReturn struct {
	Code uint16            `json:"code"`
	Data PaymentStatusData `json:"data"`
	Msg  string            `json:"msg"`
	Sign string            `json:"sign"`
}

// 代付订单返回数据
type PaymentStatusData struct {
	AppId       string `json:"appId"`
	CreatedAt   int64  `json:"createdAt"`
	Currency    string `json:"currency"`
	MchNo       string `json:"mchNo"`
	MchOrderNo  string `json:"mchOrderNo"`
	OrderAmount int    `json:"orderAmount"`
	PayOrderId  string `json:"payOrderId"`
	State       int    `json:"state"`
	SuccessTime int64  `json:"successTime"`
	ErrCode     int    `json:"errCode"`
	ErrMsg      string `json:"errMsg"`
}
