package request

type Paging struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type UpdateTime struct {
	DateTime string `json:"date_time" form:"date_time"`
}

type DeleteId struct {
	Id uint64 `json:"id"`
}

type CheckIds struct {
	Ids []int `json:"ids"`
}

type UserId struct {
	Uid int `json:"uid" form:"uid" binding:"required"`
}

type Uid struct {
	Uid int `json:"uid" form:"uid"`
}
