package request

import "za.game/lib/account"

type GetEmailList struct {
	Paging
	EventId uint64 `json:"event_id" form:"event_id"`
	Type    uint8  `json:"type" form:"type"`
}

type SaveEmail struct {
	Id          uint64               `json:"id"`
	Type        uint8                `json:"type"`         //邮件类型(1=即发邮件,2=预设邮件)
	Title       string               `json:"title"`        //邮件标题
	Msg         string               `json:"msg"`          //邮件内容
	AnnexIds    []string             `json:"annex_ids"`    //附件IDS
	SendType    uint8                `json:"send_type"`    //发送类型(1=全局发送,2=在线玩家发送,3=指定玩家)
	UserIds     string               `json:"user_ids"`     //接受玩家IDS
	IsPermanent bool                 `json:"is_permanent"` //永久有效
	StartTime   string               `json:"start_time"`   //开始时间
	EndTime     string               `json:"end_time"`     //结束时间
	EventId     uint64               `json:"event_id"`     //关联事件id
	Condition   uint64               `json:"condition"`    //发放条件
	Attachments []account.Attachment `json:"attachments"`
}

type SaveAnnex struct {
	Id         uint64 `json:"id"`
	Name       string `json:"name"`        //附件名称
	EnName     string `json:"en_name"`     //附件英文名
	Icon       string `json:"icon"`        //图标
	Type       uint8  `json:"type"`        //类型(1=筹码)
	Amount     int    `json:"amount"`      //金额
	AmountType uint8  `json:"amount_type"` //筹码类型(1=不可提现可下注，2=可提现下注)
	Unit       string `json:"unit"`        //单位
	Remark     string `json:"remark"`      //备注
}

type GetEmailEventList struct {
	Paging
}

type SaveEmailEvent struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"` //附件名称
}
