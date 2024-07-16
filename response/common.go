package response

type ChannelList struct {
	Total int64     `json:"total"`
	List  []Channel `json:"list"`
}

type Channel struct {
	Id          int    `json:"id"`
	ChannelName string `json:"channel_name"`
	Code        string `json:"code"`
	Remark      string `json:"remark"`
}

type GameBase struct {
	Id     int    `json:"id"`
	EnName string `json:"en_name"`
}
