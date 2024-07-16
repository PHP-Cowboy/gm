package response

type RoomList struct {
	Total int64  `json:"total"`
	List  []Room `json:"list"`
}

type Room struct {
	Id            uint64 `json:"id"`
	RoomId        uint64 `json:"room_id"`
	SvrId         uint32 `json:"svr_id"`         //服务器id
	GameId        uint32 `json:"game_id"`        //游戏ID
	RoomIndex     uint32 `json:"room_index"`     //房间index
	Base          uint32 `json:"base"`           //底注
	MinEntry      int    `json:"min_entry"`      //进入限制(下) 0代表无限制
	MaxEntry      int    `json:"max_entry"`      //进入限制(上) 0代表无限制
	RoomName      string `json:"room_name"`      //房间名称
	RoomType      uint8  `json:"room_type"`      //类型 1体验大厅 2正常大厅
	RoomSwitch    int    `json:"room_switch"`    //房间开关
	RoomWelfare   int    `json:"room_welfare"`   //房间赠送
	Desc          string `json:"desc"`           //房间描述
	Tax           int    `json:"tax"`            //税千分比
	BonusDiscount int    `json:"bonus_discount"` //比例千分比
	AiSwitch      int    `json:"ai_switch"`      //ai开关
	AiLimit       int    `json:"ai_limit"`       //ai人数限制
	RechargeLimit int    `json:"recharge_limit"` //充值准入值
	PoolID        int    `json:"pool_id"`        //奖金池ID
	ExtData       string `json:"ext_data"`       //特殊配置
}
