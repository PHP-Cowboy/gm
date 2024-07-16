package response

type GetEmailList struct {
	Id          uint64   `json:"id"`
	Type        uint8    `json:"type"`       //邮件类型(1=系统邮件,2=奖励邮件)
	Title       string   `json:"title"`      //邮件标题
	Msg         string   `json:"msg"`        //邮件内容
	Status      uint8    `json:"status"`     //状态(0=未发送,1=已发送,2=发送中)
	IsAnnex     uint8    `json:"is_annex"`   //是否有附件(1=有,0=否)
	AnnexIds    []string `json:"annex_ids"`  //附件IDS
	Annex       string   `json:"annex"`      //附件
	SendType    uint8    `json:"send_type"`  //发送类型(1=全局发送,2=定向发送,3=在线玩家发送)
	UserIds     string   `json:"user_ids"`   //接受玩家IDS
	StartTime   string   `json:"start_time"` //开始时间
	IsPermanent bool     `json:"is_permanent"`
	EndTime     string   `json:"end_time"`  //结束时间
	EventId     uint64   `json:"event_id"`  //关联事件id
	Event       string   `json:"event"`     //关联事件
	Condition   uint64   `json:"condition"` //发放条件
}

type EmailAnnexList struct {
	Id         uint64 `json:"id"`
	Name       string `json:"name"`        //附件名称
	EnName     string `json:"en_name"`     //附件英文名
	Type       uint8  `json:"type"`        //类型(1=筹码)
	Amount     int    `json:"amount"`      //金额
	AmountType uint8  `json:"amount_type"` //筹码类型(1=不可提现可下注，2=可提现下注)
	Unit       string `json:"unit"`        //单位
	Remark     string `json:"remark"`      //备注
}

type EmailEventRsp struct {
	Total int64        `json:"total"`
	List  []EmailEvent `json:"list"`
}

type EmailEvent struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"` //事件名称
}
