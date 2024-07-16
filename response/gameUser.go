package response

type GameUserRsp struct {
	Total int64      `json:"total"`
	List  []GameUser `json:"list"`
}

type GameUser struct {
	Id         int    `json:"id"`
	Uid        int    `json:"uid"`        //用户ID
	IsGuest    uint8  `json:"is_guest"`   //是否游客
	IsSend     uint8  `json:"is_send"`    //是否赠送
	Device     string `json:"device"`     //设备码
	UserName   string `json:"user_name"`  //用户名
	Icon       int8   `json:"icon"`       //头像
	Phone      string `json:"phone"`      //电话
	Email      string `json:"email"`      //邮箱
	ChannelId  int    `json:"channel_id"` //渠道ID
	TpNew      int    `json:"tp_new"`
	RegIp      string `json:"reg_ip"`      //注册ip
	RegVersion string `json:"reg_version"` //注册版本号
	PayStatus  int    `json:"pay_status"`
	Asset      int    `json:"asset"`
	Recharge   int    `json:"recharge"`
	GiftCash   int    `json:"gift_cash"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type LoginLogRsp struct {
	Total int64      `json:"total"`
	List  []LoginLog `json:"list"`
}

type LoginLog struct {
	Id                 int    `json:"id"`
	Uid                int    `json:"uid"`                 //用户id
	Nickname           string `json:"nickname"`            //用户昵称
	ChannelId          int    `json:"channel_id"`          //渠道id
	Channel            string `json:"channel"`             //渠道名称
	Assets             int    `json:"assets"`              //总资产cash+winCash
	ReferralCommission int    `json:"referral_commission"` //推荐佣金
	Ip                 string `json:"ip"`                  //登录ip
	Device             string `json:"device"`              //登录设备号
	Version            string `json:"version"`             //登录版本号
	LoginMode          int    `json:"login_mode"`          //登录方式(1:游客,2:手机号)
	LoginTime          string `json:"login_time"`          //登录时间
	RegTime            string `json:"reg_time"`            //注册时间
}

type GameGiveRsp struct {
	Total int64       `json:"total"`
	List  []GiveMoney `json:"list"`
}

type GiveMoney struct {
	Id               int     `json:"id"`
	Uid              int     `json:"uid"`               //用户id
	UserName         string  `json:"user_name"`         //昵称
	PayChannelId     int     `json:"pay_channel_id"`    //赠送金币通道id
	PayChannelName   string  `json:"pay_channel_name"`  //赠送金币通道名称
	Amount           int     `json:"amount"`            //赠送金币
	AmountTotal      int     `json:"amount_total"`      //赠送总币
	Recharge         int     `json:"recharge"`          //充值总币
	GiveRate         float64 `json:"give_rate"`         //实际赠送金币比
	CommitGiveRate   float64 `json:"commit_give_rate"`  //提交赠送金币比
	Status           int     `json:"status"`            //订单状态
	Auditor          string  `json:"auditor"`           //审核人
	AuditTime        string  `json:"audit_time"`        //审核时间
	OrderNo          string  `json:"order_no"`          //订单号
	TrdOrderNo       string  `json:"trd_order_no"`      //三方单号
	CreatedAt        string  `json:"created_at"`        //申请时间
	ArrivalTime      string  `json:"arrival_time"`      //到账时间
	CancelTime       string  `json:"cancel_time"`       //取消时间
	VoidTime         string  `json:"void_time"`         //作废时间
	VoidOperator     string  `json:"void_operator"`     //作废操作人
	ChannelId        int     `json:"channel_id"`        //渠道id
	ChannelCode      string  `json:"channel_code"`      //渠道
	TpFrequency      int     `json:"tp_frequency"`      //TP总局数
	RmFrequency      int     `json:"rm_frequency"`      //RM总局数
	HundredFrequency int     `json:"hundred_frequency"` //百人总局数
	SlotsFrequency   int     `json:"slots_frequency"`   //Slots总局数
	GiveMode         int     `json:"give_mode"`         //赠送金币方式
	Ifsc             int     `json:"ifsc"`              //IFSC
	PayStatus        int     `json:"pay_status"`        //付费状态
	TaxRate          float64 `json:"tax_rate"`          //税收(%)
}

type GameSituation struct {
	Id        int    `json:"id"`
	GameName  string `json:"game_name"`
	RoomType  string `json:"room_type"`
	Total     int    `json:"total"`
	WinNum    int    `json:"win_num"`
	WinMoney  int    `json:"win_money"`
	LossMoney int    `json:"loss_money"`
}

type WithdrawInfoRsp struct {
	Total int64          `json:"total"`
	List  []WithdrawInfo `json:"list"`
}

type WithdrawInfo struct {
	Id        uint64 `json:"id"`
	Uid       uint64 `json:"uid"`        // 用户ID
	BankCode  string `json:"bank_code"`  // 银行类型
	BankName  string `json:"bank_name"`  // 银行名称
	AccountNo string `json:"account_no"` // 银行账号
	Ifsc      string `json:"ifsc"`       // ifsc号
	Name      string `json:"name"`       //客户姓名
	Email     string `json:"email"`      //客户邮箱
	Phone     string `json:"phone"`      //客户手机
	Address   string `json:"address"`    //客户地址
	Vpa       string `json:"vpa"`        //vpa
	Remark    string `json:"remark"`     //备注
	UpdatedAt string `json:"updated_at"`
}

type Refund struct {
	Uid       int
	HasAttach int
	Amount    int
}

type SendEmail struct {
	Title string
	Msg   string
}

var FailedMp = map[string]SendEmail{
	"en": {
		Title: "Withdrawal failure notification",
		Msg: `Dear:
We are sorry to inform you that your withdrawal failed, the amount has been refunded.
We judge the following two possibilities:
1. The withdrawal information is incorrect, please check the completed withdrawal information;
2. The payment channel is congested, causing the order to fail.
If you have confirmed that you have filled in the information correctly, you can try to initiate the request again.`,
	},
	"hi": {
		Title: "निकासी विफलता अधिसूचना",
		Msg: `प्रिय:
हमें आपको यह बताते हुए दुख हो रहा है कि आपकी निकासी विफल हो गई है, राशि वापस कर दी गई है।
हम निम्नलिखित दो संभावनाओं का आकलन करते हैं:
1. निकासी की जानकारी गलत है, कृपया पूर्ण निकासी जानकारी की जांच करें;
2. भुगतान चैनल भीड़भाड़ वाला है, जिसके कारण ऑर्डर विफल हो गया है।
यदि आपने पुष्टि कर दी है कि आपने जानकारी सही ढंग से भरी है, तो आप फिर से अनुरोध शुरू करने का प्रयास कर सकते हैं।`,
	},
}

var SuccessMp = map[string]SendEmail{
	"en": {
		Title: "withdrawal success!",
		Msg: `Dear:
Your withdrawal has been successfully deposited into your account, please check.
We wish you happiness and big win.`,
	},
	"hi": {
		Title: "वापसी सफलता!",
		Msg: `प्रिय:
आपकी निकासी सफलतापूर्वक आपके खाते में जमा हो गई है, कृपया जांच लें।
हम आपकी ख़ुशी और बड़ी जीत की कामना करते हैं।`,
	},
}
