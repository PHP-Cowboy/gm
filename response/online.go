package response

type OnlineListRsp struct {
	Total int      `json:"total"`
	List  []Online `json:"list"`
}

type Online struct {
	Uid           int    `json:"uid"`
	Nickname      string `json:"nickname"`
	GameName      string `json:"game_name"`
	RoomName      string `json:"room_name"`
	Coin          int    `json:"coin"`
	WinCash       int    `json:"win_cash"`
	RechargeTotal int    `json:"recharge_total"`
	GiveTotal     int    `json:"give_total"`
	RegTime       string `json:"reg_time"`
	ChannelId     int    `json:"channel_id"`
	ChannelName   string `json:"channel_name"`
	ChannelCode   string `json:"channel_code"`
}
