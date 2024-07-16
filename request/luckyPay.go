package request

type LuckyPayment struct {
	AccountCode  string `json:"accountCode"`
	AccountEmail string `json:"accountEmail"`
	AccountName  string `json:"accountName"`
	AccountNo    string `json:"accountNo"`
	AccountPhone string `json:"accountPhone"`
	AccountType  string `json:"accountType"`
	Currency     string `json:"currency"`
	CustomerIp   string `json:"customerIp"`
	MchNo        string `json:"mchNo"`
	MchOrderNo   string `json:"mchOrderNo"`
	NotifyUrl    string `json:"notifyUrl"`
	PayAmount    string `json:"payAmount"`
	ReqTime      string `json:"reqTime"`
	Sign         string `json:"sign"`
	Summary      string `json:"summary"`
}

type LuckyPayCallback struct {
	Currency        string `json:"currency"`
	MchNo           string `json:"mchNo"`
	MchOrderNo      string `json:"mchOrderNo"`
	PayAmount       string `json:"payAmount"`
	PayInitiateTime string `json:"payInitiateTime"`
	PayFinishTime   string `json:"payFinishTime"`
	PayOrderNo      string `json:"payOrderNo"`
	PayState        int    `json:"payState"`
	ErrMsg          string `json:"errMsg"`
	Sign            string `json:"sign"`
}
