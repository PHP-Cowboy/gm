package request

type ChannelList struct {
	Paging
	ChannelName string `json:"channel_name" form:"channel_name"`
	Remark      string `json:"remark" form:"remark"`
}

type SaveChannel struct {
	Id          int    `json:"id"`
	ChannelName string `json:"channel_name" binding:"required"` //渠道名称
	Code        string `json:"code" binding:"required"`         //code
	Remark      string `json:"remark"`                          //备注
}

type PayConfigList struct {
	Paging
	Name string `json:"name" form:"name"`
}
