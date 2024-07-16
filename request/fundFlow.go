package request

import "time"

type SlotFundsFlowLog struct {
	Paging
}

type TpFundsFlowLog struct {
	Paging
}

type RoomFundsFlowLog struct {
	Paging
	Uid            int      `json:"uid" form:"uid"`
	GameId         int      `json:"game_id" form:"game_id"`
	CreatedAt      []string `json:"created_at" form:"created_at[]"`
	StartCreatedAt string   `json:"start_created_at"`
	EndCreatedAt   string   `json:"end_created_at"`
}
type FundsFlowLogView struct {
	Paging
	Uid       int      `json:"uid" form:"uid"`
	GameId    int      `json:"game_id" form:"game_id"`
	CreatedAt []string `json:"created_at" form:"created_at[]"`
	Start     time.Time
	End       time.Time
}
