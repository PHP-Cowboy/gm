package response

type ReportRsp struct {
	Total int64        `json:"total"`
	List  []ReportData `json:"list"`
}

type ReportData struct {
	Ymd                                        string  `json:"ymd"`                                               //日期
	Channel                                    int     `json:"channel"`                                           //渠道
	ChannelName                                string  `json:"channel_name"`                                      //渠道名称
	DailyActiveUser                            int     `json:"daily_active_user"`                                 //DAU
	EffectiveAddNums                           int     `json:"effective_add_nums"`                                //新注册用户中，有任意对局行为的用户,有效注册
	AddNums                                    int     `json:"add_nums"`                                          //新增用户数
	NewAddUserRechargeCoinNum                  int     `json:"new_add_user_recharge_coin_num"`                    //新用户充值金币数
	NewAddUserRechargeRate                     float64 `json:"new_add_user_recharge_rate"`                        //新用户充值金币率
	NewAddUserRechargeAmount                   float64 `json:"new_add_user_recharge_amount"`                      //新用户付费总额
	NewAddUserRechargePeople                   int     `json:"new_add_user_recharge_people"`                      //新用户付费人数
	NewAddUserAverageRevenuePerUser            float64 `json:"new_add_user_average_revenue_per_user"`             //新用户Arpu
	NewAddUserAverageRevenuePerPayingUser      float64 `json:"new_add_user_average_revenue_per_paying_user"`      //新用户Arppu
	DailyActiveUserRechargeAmount              int     `json:"daily_active_user_recharge_amount"`                 //DAU付费总额
	DailyActiveUserAverageRevenuePerUser       float64 `json:"daily_active_user_average_revenue_per_user"`        //DAU Arpu
	DailyActiveUserAverageRevenuePerPayingUser float64 `json:"daily_active_user_average_revenue_per_paying_user"` //DAU Arppu
	OldUserRechargePeopleNum                   int     `json:"old_user_recharge_people_num"`                      //老用户充值金币人数
	OldUserRechargeRate                        float64 `json:"old_user_recharge_rate"`                            //老用户充值金币率
	PaidUserRetentionRate                      float64 `json:"paid_user_retention_rate"`                          //付费用户留存率
	GiveMoneyPeople                            int     `json:"give_money_people"`                                 //赠送金币人数
	GiveMoneyAmount                            int     `json:"give_money_amount"`                                 //赠送金币额
	GiveMoneyRate                              float64 `json:"give_money_rate"`                                   //赠送金币率
	PlayRate                                   float64 `json:"play_rate"`                                         //玩牌率
	NextDayRetention                           float64 `json:"next_day_retention"`                                //次留
	ThreeDayRetention                          float64 `json:"three_day_retention"`                               //3留
	FourDayRetention                           float64 `json:"four_day_retention"`                                //4留
	FiveDayRetention                           float64 `json:"five_day_retention"`                                //5留
	SixDayRetention                            float64 `json:"six_day_retention"`                                 //6留
	SevenDayRetention                          float64 `json:"seven_day_retention"`                               //7留
	FourteenDayRetention                       float64 `json:"fourteen_day_retention"`                            //14留
	NextDayPeople                              int     `json:"next_day_people"`
	ThreeDayPeople                             int     `json:"three_day_people"`
	FourDayPeople                              int     `json:"four_day_people"`
	FiveDayPeople                              int     `json:"five_day_people"`
	SixDayPeople                               int     `json:"six_day_people"`
	SevenDayPeople                             int     `json:"seven_day_people"`
	FourteenDayPeople                          int     `json:"fourteen_day_people"`
}

type ReportDataDescByYmd []ReportData

func (a ReportDataDescByYmd) Len() int           { return len(a) }
func (a ReportDataDescByYmd) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ReportDataDescByYmd) Less(i, j int) bool { return a[i].Ymd > a[j].Ymd } // 降序

type UserStatistics struct {
	Total int64                `json:"total"`
	List  []UserStatisticsList `json:"list"`
}

type UserStatisticsList struct {
	Ymd                         string  `json:"ymd"`                             //日期
	Channel                     int     `json:"channel"`                         //渠道
	ChannelName                 string  `json:"channel_name"`                    //渠道名称
	AddNums                     int     `json:"add_nums"`                        //新增用户数
	EffectiveAddNums            int     `json:"effective_add_nums"`              //新注册用户中，有任意对局行为的用户
	EffectiveIncreaseRate       float64 `json:"effective_increase_rate"`         //新注册用户中，有任意对局行为的用户/新增用户
	AddPayUserNextDayRetention  float64 `json:"add_pay_user_next_day_retention"` //新增付费用户次日留存率
	DailyActiveUser             int     `json:"daily_active_user"`               //日活跃用户 当日登录
	PayingSubscribers           int     `json:"paying_subscribers"`              //付费用户数 [=渠道今日付费用户数]
	PayoutRate                  float64 `json:"payout_rate"`                     //付费率
	TotalPayment                int     `json:"total_payment"`                   //总付费
	AverageRevenuePerUser       int     `json:"average_revenue_per_user"`        //每用户平均收入
	AverageRevenuePerPayingUser int     `json:"average_revenue_per_paying_user"` //每付费用户平均收益
	WithdrawalRate              float64 `json:"withdrawal_rate"`                 //提现率 = 提现总额/总付费
	TotalWithdrawal             int     `json:"total_withdrawal"`                //提现总额
	WithdrawalPeople            int     `json:"withdrawal_people"`               //提现人数
	NextDayRetention            float64 `json:"next_day_retention"`              //次日留存
	ThreeDayRetention           float64 `json:"three_day_retention"`             //三日留存
	SevenDayRetention           float64 `json:"seven_day_retention"`             //七日留存
	FourteenDayRetention        float64 `json:"fourteen_day_retention"`          //十四日留存
	ThirtyDayRetention          float64 `json:"thirty_day_retention"`            //三十日留存
	NinetyDayRetention          float64 `json:"ninety_day_retention"`            //九十日留存
	NewPayingSubscribers        int     `json:"new_paying_subscribers"`          //新增付费人数
	NewDeviceNum                int     `json:"new_device_num"`                  //新增设备
	NewPayingMoney              int     `json:"new_paying_money"`                //新增付费额度
	StrongActiveNum             int     `json:"strong_active_num"`               //强活跃人数[登录并且游戏时间超过5分钟]
}

type RechargeStatistics struct {
	Total int64      `json:"total"`
	List  []Recharge `json:"list"`
}

type Recharge struct {
	Ymd                            string  `json:"ymd"`                                 //日期
	ChannelName                    string  `json:"channel_name"`                        //渠道
	NewUserRechargeNums            int     `json:"new_user_recharge_nums"`              //新用户充值人数
	NewUserRechargeTotal           int     `json:"new_user_recharge_total"`             //新用户充值总额
	NewUserRechargeRate            float64 `json:"new_user_recharge_rate"`              //新用户充值率
	AverageRevenuePerNewUser       int     `json:"average_revenue_per_new_user"`        //新用户每用户平均收入
	AverageRevenuePerPayingNewUser int     `json:"average_revenue_per_paying_new_user"` //新用户每付费用户平均收益
	OldUserRechargeNums            int     `json:"old_user_recharge_nums"`              //老用户充值人数
	OldUserFirstRechargeNums       int     `json:"old_user_first_recharge_nums"`        //老用户首次付费人数
	OldUserRechargeTotal           int     `json:"old_user_recharge_total"`             //老用户充值总额
	OldUserRechargeRate            float64 `json:"old_user_recharge_rate"`              //老用户充值率
	AverageRevenuePerOldUser       int     `json:"average_revenue_per_old_user"`        //老用户每用户平均收入
	AverageRevenuePerPayingOldUser int     `json:"average_revenue_per_paying_old_user"` //老用户每付费用户平均收益
	AddUserNum                     int
}

type WithdrawalStatistics struct {
	Total int64        `json:"total"`
	List  []Withdrawal `json:"list"`
}

type Withdrawal struct {
	Ymd                    string  `json:"ymd"`                       //日期
	ChannelName            string  `json:"channel_name"`              //渠道名称
	WithdrawalNums         int     `json:"withdrawal_nums"`           //提现用户数
	WithdrawalTotal        int     `json:"withdrawal_total"`          //提现总额
	NewUserWithdrawalNums  int     `json:"new_user_withdrawal_nums"`  //新用户提现用户数
	NewUserWithdrawalTotal int     `json:"new_user_withdrawal_total"` //新用户提现总额
	NewUserWithdrawalRate  float64 `json:"new_user_withdrawal_rate"`  //新用户提现率
	OldUserWithdrawalNums  int     `json:"old_user_withdrawal_nums"`  //老用户提现用户数
	OldUserWithdrawalTotal int     `json:"old_user_withdrawal_total"` //老用户提现总额
	OldUserWithdrawalRate  float64 `json:"old_user_withdrawal_rate"`  //老用户提现率
}

type PaidUserRetention struct {
	Total int64                   `json:"total"`
	List  []PaidUserRetentionList `json:"list"`
}

type PaidUserRetentionList struct {
	Ymd                   string `json:"ymd"`                      //日期
	ChannelName           string `json:"channel_name"`             //渠道名称
	UserNums              int    `json:"user_nums"`                //用户数
	NextDayRetention      int    `json:"next_day_retention"`       //次日留存
	TwoDayRetention       int    `json:"two_day_retention"`        //二日留存
	ThreeDayRetention     int    `json:"three_day_retention"`      //三日留存
	FourDayRetention      int    `json:"four_day_retention"`       //四日留存
	FiveDayRetention      int    `json:"five_day_retention"`       //五日留存
	SixDayRetention       int    `json:"six_day_retention"`        //六日留存
	SevenDayRetention     int    `json:"seven_day_retention"`      //七日留存
	FourteenDayRetention  int    `json:"fourteen_day_retention"`   //十四日留存
	TwentyOneDayRetention int    `json:"twenty_one_day_retention"` //二十一日留存
	ThirtyDayRetention    int    `json:"thirty_day_retention"`     //三十日留存
	SixtyDayRetention     int    `json:"sixty_day_retention"`      //六十日留存
	NinetyDayRetention    int    `json:"ninety_day_retention"`     //九十日留存
}

type UserRetention struct {
	Total int64               `json:"total"`
	List  []UserRetentionList `json:"list"`
}

type UserRetentionList struct {
	Ymd                   string `json:"ymd"`                      //日期
	ChannelName           string `json:"channel_name"`             //渠道名称
	UserNums              int    `json:"user_nums"`                //用户数
	NextDayRetention      int    `json:"next_day_retention"`       //次日留存
	TwoDayRetention       int    `json:"two_day_retention"`        //二日留存
	ThreeDayRetention     int    `json:"three_day_retention"`      //三日留存
	FourDayRetention      int    `json:"four_day_retention"`       //四日留存
	FiveDayRetention      int    `json:"five_day_retention"`       //五日留存
	SixDayRetention       int    `json:"six_day_retention"`        //六日留存
	SevenDayRetention     int    `json:"seven_day_retention"`      //七日留存
	FourteenDayRetention  int    `json:"fourteen_day_retention"`   //十四日留存
	TwentyOneDayRetention int    `json:"twenty_one_day_retention"` //二十一日留存
	ThirtyDayRetention    int    `json:"thirty_day_retention"`     //三十日留存
	SixtyDayRetention     int    `json:"sixty_day_retention"`      //六十日留存
	NinetyDayRetention    int    `json:"ninety_day_retention"`     //九十日留存
}

type PerHourDataNum struct {
	Total int64                `json:"total"`
	List  []PerHourDataNumList `json:"list"`
}

type PerHourDataNumList struct {
	Ymd            int `json:"ymd"`
	NumType        int `json:"num_type"`
	ZeroNum        int `json:"zero_num"`
	OneNum         int `json:"one_num"`
	TwoNum         int `json:"two_num"`
	ThreeNum       int `json:"three_num"`
	FourNum        int `json:"four_num"`
	FiveNum        int `json:"five_num"`
	SixNum         int `json:"six_num"`
	SevenNum       int `json:"seven_num"`
	EightNum       int `json:"eight_num"`
	NineNum        int `json:"nine_num"`
	TenNum         int `json:"ten_num"`
	ElevenNum      int `json:"eleven_num"`
	TwelveNum      int `json:"twelve_num"`
	ThirteenNum    int `json:"thirteen_num"`
	FourteenNum    int `json:"fourteen_num"`
	FifteenNum     int `json:"fifteen_num"`
	SixteenNum     int `json:"sixteen_num"`
	SeventeenNum   int `json:"seventeen_num"`
	EighteenNum    int `json:"eighteen_num"`
	NineteenNum    int `json:"nineteen_num"`
	TwentyNum      int `json:"twenty_num"`
	TwentyOneNum   int `json:"twenty_one_num"`
	TwentyTwoNum   int `json:"twenty_two_num"`
	TwentyThreeNum int `json:"twenty_three_num"`
}

type PerHourGameNum struct {
	Total int64                `json:"total"`
	List  []PerHourGameNumList `json:"list"`
}

type PerHourGameNumList struct {
	Ymd            int `json:"ymd"`
	GameId         int `json:"game_id"`
	RoomId         int `json:"room_id"`
	Chip           int `json:"chip"`
	NumType        int `json:"num_type"`
	ZeroNum        int `json:"zero_num"`
	OneNum         int `json:"one_num"`
	TwoNum         int `json:"two_num"`
	ThreeNum       int `json:"three_num"`
	FourNum        int `json:"four_num"`
	FiveNum        int `json:"five_num"`
	SixNum         int `json:"six_num"`
	SevenNum       int `json:"seven_num"`
	EightNum       int `json:"eight_num"`
	NineNum        int `json:"nine_num"`
	TenNum         int `json:"ten_num"`
	ElevenNum      int `json:"eleven_num"`
	TwelveNum      int `json:"twelve_num"`
	ThirteenNum    int `json:"thirteen_num"`
	FourteenNum    int `json:"fourteen_num"`
	FifteenNum     int `json:"fifteen_num"`
	SixteenNum     int `json:"sixteen_num"`
	SeventeenNum   int `json:"seventeen_num"`
	EighteenNum    int `json:"eighteen_num"`
	NineteenNum    int `json:"nineteen_num"`
	TwentyNum      int `json:"twenty_num"`
	TwentyOneNum   int `json:"twenty_one_num"`
	TwentyTwoNum   int `json:"twenty_two_num"`
	TwentyThreeNum int `json:"twenty_three_num"`
}

type FiveMinuteData struct {
	RegNum     int `json:"reg_num"`     //新增注册
	OnlineNum  int `json:"online_num"`  //实时在线
	ActiveNum  int `json:"active_num"`  //活跃人数
	PayNum     int `json:"pay_num"`     //付费人数
	PayAmount  int `json:"pay_amount"`  //付费额度
	GiveAmount int `json:"give_amount"` //赠送额度
}
