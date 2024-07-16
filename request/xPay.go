package request

type XPayBalanceRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type XPayOrderStatusRequest struct {
	ID       int    `json:"id" validate:"required"`
	OrderNo  string `json:"orderNo" validate:"required"`
	MOrderNo string `json:"mOrderNo" validate:"required"`
}

type XPayPaymentOrderRequest struct {
	Uid         uint64 `json:"uid" validate:"required"`
	PayId       int    `json:"payId" validate:"required"`
	OrderNo     string `json:"orderNo" validate:"required"`
	OrderAmount int    `json:"OrderAmount" validate:"required"`
	BankCode    int8   `json:"bankCode" validate:"required"`
	BankName    string `json:"bankName"`
	AccountNo   string `json:"accountNo"`
	Ifsc        string `json:"ifsc"`
	Name        string `json:"accountName" validate:"required"`
	Email       string `json:"accountEmail" validate:"required"`
	Phone       string `json:"accountMobileNo" validate:"required"`
	Address     string `json:"accountAddress"`
	Vpa         string `json:"vpa"`
	Remark      string `json:"transferDesc"`
}

type XPayPaymentStatusRequest struct {
	ID       int    `json:"id" validate:"required"`
	OrderNo  string `json:"orderNo" validate:"required"`
	MOrderNo string `json:"mOrderNo" validate:"required"`
}

type PayOrder struct {
	PayId    int    `json:"id" validate:"required"`
	GiftId   int    `json:"giftId" validate:"required"`
	GiftType int8   `json:"giftType" validate:"required"`
	Remark   string `json:"remark"`
	RoomId   string `json:"roomid,omitempty"`
}
