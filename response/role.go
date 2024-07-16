package response

type Role struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Desc        string `json:"desc"`
	Status      int    `json:"status"`
	MenuList    []int  `json:"menuList"`
	ChannelList []int  `json:"channelList"`
}

type RoleList struct {
	Total int64  `json:"total"`
	List  []Role `json:"list"`
}
