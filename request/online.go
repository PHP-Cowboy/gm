package request

type OnlineList struct {
	Paging
	Uid         uint64 `json:"uid" form:"uid"`
	RoomId      uint64 `json:"room_id" form:"room_id"`
	ChannelId   int    `json:"channel_id" form:"channel_id"`
	RechargeMin int    `json:"recharge_min" form:"recharge_min"`
	RechargeMax int    `json:"recharge_max" form:"recharge_max"`
	Admin       string `json:"admin"`
	ChannelIds  []int
}
