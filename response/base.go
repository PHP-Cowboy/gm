package response

type DictTypeList struct {
	Total int64      `json:"total"`
	List  []DictType `json:"list"`
}

type DictType struct {
	Id   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type DictList struct {
	Id       uint64 `json:"id"`
	TypeCode string `json:"type_code"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	IsEdit   int    `json:"is_edit"`
}
type NoviceCarnivalWithdraw struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
