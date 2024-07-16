package response

type InPayWithdrawRsp struct {
	Code *int              `json:"code"`
	Data InPayWithdrawData `json:"data"`
	Msg  string            `json:"msg"`
	Time int               `json:"time"`
}

type InPayWithdrawData struct {
	Status      int    `json:"status"`
	PlatNumber  string `json:"plat_number"`
	OrderNumber string `json:"order_number"`
	OrderAmount string `json:"order_amount"`
}

type InPayBalanceRsp struct {
	Code int          `json:"code"`
	Data InPayBalance `json:"data"`
	Msg  string       `json:"msg"`
	Time int          `json:"time"`
}

type InPayBalance struct {
	Balance string `json:"balance"`
}
