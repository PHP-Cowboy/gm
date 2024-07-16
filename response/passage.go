package response

type PassageRsp struct {
	Total int64     `json:"total"`
	List  []Passage `json:"list"`
}

type Passage struct {
	Id                   int      `json:"id"`
	PassageId            int      `json:"passage_id"`
	Entrance             int      `json:"entrance"`
	CollectionRate       float64  `json:"collection_rate"`
	CollectionStatus     int      `json:"collection_status"`
	CollectionSort       int      `json:"collection_sort"`
	CollectionChannelIds []string `json:"collection_channel_ids"`
	PaymentRate          float64  `json:"payment_rate"`
	PaymentStatus        int      `json:"payment_status"`
	PaymentSort          int      `json:"payment_sort"`
	PaymentChannelIds    []string `json:"payment_channel_ids"`
	CreatedAt            string   `json:"created_at"`
	UpdatedAt            string   `json:"updated_at"`
}
