package request

type DictTypeList struct {
	Paging
	Code string `json:"code" form:"code"`
	Name string `json:"name" form:"name"`
}

type DictList struct {
	TypeCode string `json:"type_code" form:"type_code"`
}

type EditDict struct {
	TypeCode string `json:"type_code"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	IsEdit   int    `json:"is_edit"`
}

type GetOneDict struct {
	TypeCode string `json:"type_code" form:"type_code" binding:"required"`
	Code     string `json:"code" form:"type_code" binding:"required"`
}
