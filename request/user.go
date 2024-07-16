package request

type AddUser struct {
	Username  string `json:"username" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password" binding:"required"`
	CreatorId int    `json:"creator_id"`
	Status    int    `json:"status"`
	RoleList  []int  `json:"role_list"`
}

type ChangeUser struct {
	Id       int    `json:"id" binding:"required"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   int    `json:"status"`
	RoleList []int  `json:"role_list"`
}

type UserList struct {
	Paging
	Id       int    `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	Name     string `json:"name" form:"name"`
	Status   int    `json:"status" form:"status"`
}

type GetCaptchaQr struct {
	Uid    int    `json:"uid" form:"uid"`
	Secret string `json:"secret" form:"secret"`
}
