package request

type RoomList struct {
	Paging
	RoomName string `json:"room_name" form:"room_name"`
}

type SaveRoom struct {
	ID            uint64 `json:"id"`
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
	RechargeLimit *int   `json:"recharge_limit"` //充值准入值
	PoolID        *int   `json:"pool_id"`        //充值准入值
	ExtData       string `json:"ext_data"`       //特殊配置
}

type UpdateExtDataByExcel struct {
	RoomId            uint64    `json:"RoomId"`            //房间
	BetRange          []float64 `json:"BetRange"`          //下注范围
	GoldPDownLimit    int       `json:"GoldPDownLimit"`    //金池下限
	GoldPUpLimit      int       `json:"GoldPUpLimit"`      //金池上限
	MinWinContral     int       `json:"MinWinContral"`     //最低正控参数/千分数
	MaxWinContral     int       `json:"MaxWinContral"`     //最高正控参数/千分数
	MinLoseContral    int       `json:"MinLoseContral"`    //最低负控参数/千分数
	MaxLoseContral    int       `json:"MaxLoseContral"`    //最高负控参数/千分数
	RaiseRange        []float64 `json:"RaiseRange"`        //增长值范围
	RaiseTime         []float64 `json:"RaiseTime"`         //增长间隔时间/秒
	AutoJackpotValue  int       `json:"AutoJackpotValue"`  //自动爆奖池值
	AutoJackpotRate   int       `json:"AutoJackpotRate"`   //自动爆奖池几率/千分比
	AutoJackpotRange  []float64 `json:"AutoJackpotRange"`  //自动爆出比例/千分比
	MinJackpot        int       `json:"MinJackpot"`        //奖池保底值
	UnlockTimes       int       `json:"UnlockTimes"`       //解锁次数
	AITake            []float64 `json:"AITake"`            //AI带入范围
	PoolTax           int       `json:"PoolTax"`           //暗税比例/千分比
	AINumber          []float64 `json:"AINumber"`          //AI人数区间
	AIPlayTime        []float64 `json:"AIPlayTime"`        //AI时间(秒)
	JackPotFloorTime  int       `json:"JackPotFloorTime"`  //奖池体验的判断轮数
	FreeSpinFloorTime int       `json:"FreeSpinFloorTime"` //免费游戏体验的判断轮数
	JackPotRate       int       `json:"JackPotRate"`       //奖池体验的几率/千分比
	FreeSpinRate      int       `json:"FreeSpinRate"`      //免费游戏体验的几率/千分比
	JackPotOdds       int       `json:"JackPotOdds"`       //奖池体验触发倍数
	FreeSpinOdds      int       `json:"FreeSpinOdds"`      //免费游戏体验触发倍数
}
