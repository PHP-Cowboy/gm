package response

type SlotFundsFlowLog struct {
	Total int64                  `json:"total"`
	List  []SlotFundsFlowLogList `json:"list"`
}

type SlotFundsFlowLogList struct {
	Id            int    `json:"id"`
	UserName      string `json:"user_name"`       //用户昵称
	Uid           int    `json:"uid"`             //用户ID
	RoomId        int    `json:"room_id"`         //房间id
	DeskId        int    `json:"desk_id"`         //桌子id
	BeforeCash    int    `json:"before_cash"`     //变化前cash
	BeforeWinCash int    `json:"before_win_cash"` //变化前winCash
	Cash          int    `json:"cash"`            //赠送金币账户
	WinCash       int    `json:"win_cash"`        // 充值金币账户
	Nums          int    `json:"nums"`            // 变动数额
	Tax           int    `json:"tax"`             //税收
	Balance       int    `json:"balance"`         //账户余额
	Remark        string `json:"remark"`          //备注
}

type TpFundsFlowLog struct {
	Total int64                `json:"total"`
	List  []TpFundsFlowLogList `json:"list"`
}

type TpFundsFlowLogList struct {
	Id            int    `json:"id"`
	UserName      string `json:"user_name"`       //用户昵称
	Uid           int    `json:"uid"`             //用户ID
	RoomId        int    `json:"room_id"`         //房间id
	DeskId        int    `json:"desk_id"`         //桌子id
	BeforeCash    int    `json:"before_cash"`     //变化前cash
	BeforeWinCash int    `json:"before_win_cash"` //变化前winCash
	Cash          int    `json:"cash"`            //赠送金币账户
	WinCash       int    `json:"win_cash"`        // 充值金币账户
	Nums          int    `json:"nums"`            // 变动数额
	Tax           int    `json:"tax"`             //税收
	Balance       int    `json:"balance"`         //账户余额
	Remark        string `json:"remark"`          //备注
}

type RoomFundsFlowLog struct {
	Total int64                  `json:"total"`
	List  []RoomFundsFlowLogList `json:"list"`
}

type RoomFundsFlowLogList struct {
	Id            int    `json:"id"`
	UserName      string `json:"user_name"`       //用户昵称
	Uid           int    `json:"uid"`             //用户ID
	CreatedAt     string `json:"created_at"`      //变动时间
	GameId        int    `json:"game_id"`         //房间id
	RoomId        int    `json:"room_id"`         //房间id
	GameIdX       int64  `json:"game_id_x"`       //对局标识
	DeskId        int    `json:"desk_id"`         //桌子id
	Pos           int    `json:"pos"`             //座位号
	Nums          int64  `json:"nums"`            //变动数额
	BeforeCash    int64  `json:"before_cash"`     //变化前cash
	Cash          int64  `json:"cash"`            //当前cash余额
	CashNums      int64  `json:"cash_nums"`       // cash变动数额
	BeforeWinCash int64  `json:"before_win_cash"` //变化前winCash
	WinCash       int64  `json:"win_cash"`        // 当前winCash余额
	WinCashNums   int64  `json:"win_cash_nums"`   // winCash变动数额
	Withdraw      int64  `json:"withdraw"`        //赠送储钱罐提现额度
	WithdrawNums  int64  `json:"withdraw_nums"`   //变动的储钱罐提现额度
	ExtraGift     int64  `json:"extra_gift"`      //额外赠送
	BeforeBonus   int64  `json:"before_bonus"`    //变化前 bonus
	Bonus         int64  `json:"bonus"`           //bonus
	BonusNums     int64  `json:"bonus_nums"`      //bonus变动数额
	Tax           int64  `json:"tax"`             //税收
	Balance       int64  `json:"balance"`         //账户余额
	Remark        string `json:"remark"`          //备注
	Details       string `json:"details"`
}
