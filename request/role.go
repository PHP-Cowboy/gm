package request

type CreateRole struct {
	Name        string `json:"name" binding:"required"`
	Desc        string `json:"desc"`
	CreatorId   int    `json:"creator_id"`
	Creator     string `json:"creator"`
	MenuList    []int  `json:"menuList"`
	ChannelList []int  `json:"channelList"`
}

type ChangeRole struct {
	Id          int    `json:"id" binding:"required"`
	Name        string `json:"name"`
	Desc        string `json:"desc"`
	Status      int    `json:"status"`
	MenuList    []int  `json:"menuList"`
	ChannelList []int  `json:"channelList"`
}

type RoleList struct {
	Paging
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Status int    `json:"status"`
}

type UserRoleList struct {
	UserId int `json:"user_id"`
}

type AddUserRole struct {
	UserId int `json:"user_id"`
	RoleId int `json:"role_id"`
}
