package request

type RechargeStatistics struct {
	Paging
	Ymd     string `json:"ymd" form:"ymd"`
	Channel int    `json:"channel" form:"channel"`
}

type PaidUserRetention struct {
	Paging
	Ymd     string `json:"ymd" form:"ymd"`
	Channel int    `json:"channel" form:"channel"`
}

type UserRetention struct {
	Paging
	Ymd     string `json:"ymd" form:"ymd"`
	Channel int    `json:"channel" form:"channel"`
}

type UserStatistics struct {
	Paging
	Ymd         string   `json:"ymd" form:"ymd"`
	BetweenDate []string `json:"betweenDate" form:"betweenDate[]"`
	Channel     []int    `json:"channel[]" form:"channel[]"`
	Start       string
	End         string
	ChannelIds  []int
	Has         bool
}

type WithdrawalStatistics struct {
	Paging
	Ymd     string `json:"ymd" form:"ymd"`
	Channel int    `json:"channel" form:"channel"`
}

type OnlineStatistics struct {
	Paging
	Ymd     int `json:"ymd" form:"ymd"`
	Channel int `json:"channel" form:"channel"`
}

type PerHourDataNum struct {
	Paging
	Type int `json:"type" form:"type"`
}

type PerHourGameNum struct {
	Paging
	GameId int `json:"game_id" form:"game_id"`
	RoomId int `json:"room_id" form:"room_id"`
	Chip   int `json:"chip" form:"chip"`
}

type FiveMinuteData struct {
	Channel    []int `json:"channel[]" form:"channel[]"`
	Has        bool
	ChannelIds []int
}
