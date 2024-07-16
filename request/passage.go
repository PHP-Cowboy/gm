package request

type PassageList struct {
	Paging
	CollectionStatus int `json:"collection_status" form:"collection_status"`
	PaymentStatus    int `json:"payment_status" form:"payment_status"`
	Entrance         int `json:"entrance" form:"entrance"`
}

type PassageSave struct {
	Id                   int      `json:"id"`
	CollectionStatus     int      `json:"collection_status" form:"collection_status"`
	PaymentStatus        int      `json:"payment_status" form:"payment_status"`
	Entrance             int      `json:"entrance" form:"entrance"`
	PassageId            int      `json:"passage_id"`
	CollectionRate       float64  `json:"collection_rate"`
	CollectionSort       int      `json:"collection_sort"`
	CollectionChannelIds []string `json:"collection_channel_ids"`
	PaymentRate          float64  `json:"payment_rate"`
	PaymentSort          int      `json:"payment_sort"`
	PaymentChannelIds    []string `json:"payment_channel_ids"`
}

type PassageChange struct {
	Id                   int      `json:"id"`
	Type                 string   `json:"type" form:"type"`
	CollectionStatus     int      `json:"collection_status" form:"collection_status"`
	PaymentStatus        int      `json:"payment_status" form:"payment_status"`
	CollectionChannelIds []string `json:"collection_channel_ids" form:"collection_channel_ids"`
	PaymentChannelIds    []string `json:"payment_channel_ids" form:"payment_channel_ids"`
}
