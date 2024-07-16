package response

type LuckyPayPaymentRsp struct {
	Code           *string             `json:"code"`
	Message        string              `json:"message"`
	Data           LuckyPayPaymentData `json:"data"`
	MonitorTrackId string              `json:"monitorTrackId"`
	Sign           string              `json:"sign"`
	Success        bool                `json:"success"`
	Timestamp      string              `json:"timestamp"`
}

type LuckyPayPaymentData struct {
	Currency        string `json:"currency"`
	MchNo           string `json:"mchNo"`
	MchOrderNo      string `json:"mchOrderNo"`
	PayAmount       string `json:"payAmount"`
	PayInitiateTime string `json:"payInitiateTime"`
	PayOrderNo      string `json:"payOrderNo"`
	PayState        int    `json:"payState"`
}
