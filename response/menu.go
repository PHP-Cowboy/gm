package response

type Menu struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	ComponentName string `json:"component_name"`
	Path          string `json:"path" `
	Label         string `json:"label" `
	Icon          string `json:"icon" `
	Url           string `json:"url" `
	Status        int    `json:"status"`
	ParentId      int    `json:"parent_id"`
	Level         int    `json:"level"`
}

type MenuList struct {
	Total int64  `json:"total"`
	List  []Menu `json:"list"`
}

type ParentList struct {
	Id    int    `json:"id"`
	Label string `json:"label"`
}
