package request

// 查询余额
type InPayBalance struct {
	MerchantId string `json:"merchant_id"`
	Sign       string `json:"sign"`
}

type InPayWithdraw struct {
	MerchantId  string `json:"merchant_id"`
	OrderNumber string `json:"order_number"`
	OrderAmount string `json:"order_amount"`
	Type        string `json:"type"`
	Vpa         string `json:"vpa"`
	Email       string `json:"email"`
	Account     string `json:"account"`
	Name        string `json:"name"`
	Ifsc        string `json:"ifsc"`
	Phone       string `json:"phone"`
	NotifyUrl   string `json:"notify_url"`
	Sign        string `json:"sign,omitempty"`
}
