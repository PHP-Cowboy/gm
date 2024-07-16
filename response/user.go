package response

type Login struct {
	UserId   int                 `json:"user_id"`
	Username string              `json:"username"`
	Name     string              `json:"name"`
	IsBind   int                 `json:"is_bind"`
	Token    string              `json:"token"`
	Menu     []Group             `json:"menu"`
	MenuMp   map[string]struct{} `json:"menu_mp"`
}

type UserList struct {
	Total int64  `json:"total"`
	List  []User `json:"list"`
}

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Status    int    `json:"status"`
	CreatorId int    `json:"creator_id"`
	Creator   string `json:"creator"`
	RoleList  []int  `json:"role_list"`
	Secret    string `json:"secret"`
	IsBind    int    `json:"is_bind"`
}

type Group struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	Label    string  `json:"label"`
	Icon     string  `json:"icon"`
	Url      string  `json:"url"`
	ParentId int     `json:"parent_id"`
	Children []Group `json:"children"`
	Value    int     `json:"value"`
}
