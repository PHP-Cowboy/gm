package request

type AddMenu struct {
	Name      string `json:"name" binding:"required"`
	Path      string `json:"path" binding:"required"`
	Label     string `json:"label" binding:"required"`
	Icon      string `json:"icon"`
	Url       string `json:"url" binding:"required"`
	ParentId  *int   `json:"parent_id"`
	CreatorId int    `json:"creator_id"`
	Level     int    `json:"level"`
	Creator   string `json:"creator"`
}

type MenuList struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
}

type LevelList struct {
	Level int `json:"level" form:"level"`
}

type ChangeMenu struct {
	Id       int    `json:"id" binding:"required"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Label    string `json:"label"`
	Url      string `json:"url"`
	Desc     string `json:"desc"`
	Status   int    `json:"status"`
	Level    int    `json:"level"`
	ParentId *int   `json:"parent_id"`
}

type RoleMenuList struct {
	RoleId int `json:"role_id"`
	MenuId int `json:"menu_id"`
}

type AddRoleMenu struct {
	RoleId int `json:"role_id"`
	MenuId int `json:"menu_id"`
}
