package daos

import (
	"encoding/json"
	"errors"
	"gm/common/constant"
	"gm/daos/rds"
	"gm/global"
	"gm/model/game"
	"gm/msgcenter"
	"gm/request"
	"gm/response"
	"gm/utils/timeutil"
	"strconv"
	"strings"
	"time"
	"za.game/lib/account"
)

// 邮件列表
func GetEmailList(req request.GetEmailList) (list []response.GetEmailList, err error) {
	db := global.Game
	eventTable := new(game.EmailEvent)

	eventList, err := eventTable.GetList(db)
	if err != nil {
		global.Logger["err"].Errorf("GetEmailList GetList failed,err:[%v]", err.Error())
		return
	}

	eventMp := make(map[uint64]string)

	for _, l := range eventList {
		eventMp[l.Id] = l.Name
	}

	annex := new(game.EmailAnnex)
	annexList, err := annex.GetList(db)

	if err != nil {
		global.Logger["err"].Errorf("GetEmailList GetList failed,err:[%v]", err.Error())
		return
	}
	annexMp := make(map[string]string)

	for _, al := range annexList {
		annexMp[strconv.FormatUint(al.ID, 10)] = al.Name
	}

	email := new(game.Email)

	if req.EventId > 0 {
		email.EventId = req.EventId
	}

	email.Type = 2 // 预设邮件

	emailList, err := email.GetPageList(db, req.Page, req.Size)

	if err != nil {
		global.Logger["err"].Errorf("GetEmailList GetPageList failed,err:[%v]", err.Error())
		return
	}

	list = make([]response.GetEmailList, 0, len(emailList))
	for _, e := range emailList {
		var (
			startTime = ""
			endTime   = ""
			annexs    []string
			event     string
		)

		if e.StartTime > 0 {
			startTime = timeutil.FormatToDate(time.Unix(e.StartTime, 0))
		}

		if e.EndTime > 0 {
			endTime = timeutil.FormatToDate(time.Unix(e.EndTime, 0))
		}

		annexIds := strings.Split(e.AnnexIds, ",")

		if len(e.AnnexIds) > 0 {
			for _, id := range annexIds {
				annexStr, ok := annexMp[id]
				if !ok {
					err = errors.New("附件不存在")
					return
				}
				annexs = append(annexs, annexStr)
			}
		}

		event, _ = eventMp[e.EventId]

		isPermanent := true

		if startTime != "" || endTime != "" {
			isPermanent = false
		}

		list = append(list, response.GetEmailList{
			Id:          e.ID,
			Type:        e.Type,
			Title:       e.Title,
			Msg:         e.Msg,
			Status:      e.Status,
			IsAnnex:     e.IsAnnex,
			AnnexIds:    annexIds,
			Annex:       strings.Join(annexs, ","),
			SendType:    e.SendType,
			UserIds:     e.UserIds,
			IsPermanent: isPermanent,
			StartTime:   startTime,
			EndTime:     endTime,
			Event:       event,
			EventId:     e.EventId,
			Condition:   e.Condition,
		})
	}

	return
}

// 保存邮件
func SaveEmail(req request.SaveEmail) (err error) {
	tx := global.Game.Begin()
	var (
		startTime time.Time
		endTime   time.Time
	)

	if req.SendType == 3 && req.UserIds == "" {
		err = errors.New("指定玩家id不能为空")
		return
	}

	if req.SendType != 3 {
		req.UserIds = ""
	}

	req.UserIds = strings.TrimRight(req.UserIds, "、")

	email := game.Email{
		Type:      req.Type,
		Title:     req.Title,
		Msg:       req.Msg,
		Status:    2,
		AnnexIds:  "",
		SendType:  req.SendType,
		UserIds:   req.UserIds,
		EventId:   req.EventId,
		Condition: req.Condition,
		CreatedAt: time.Now(),
	}

	attachment := make([]account.Attachment, 0)

	if len(req.Attachments) > 0 {
		for _, attach := range req.Attachments {
			if attach.Nums > 0 && attach.ItemId > 0 {
				attachment = append(attachment, attach)
			}
		}
	}

	if len(attachment) > 0 {
		var b []byte
		b, err = json.Marshal(attachment)
		if err != nil {
			global.Logger["err"].Errorf("SaveEmail json.Marshal failed,err:[%v]", err.Error())

			return err
		}
		email.Attachments = string(b)
	}

	if req.StartTime != "" && !req.IsPermanent {
		startTime, err = time.Parse(timeutil.DateFormat, req.StartTime)
		if err != nil {
			global.Logger["err"].Errorf("SaveEmail time.Parse StartTime failed,err:[%v]", err.Error())
			return
		}
		email.StartTime = startTime.Unix()
	} else {
		email.StartTime = 0
	}

	if req.EndTime != "" && !req.IsPermanent {
		endTime, err = time.Parse(timeutil.DateFormat, req.EndTime)
		if err != nil {
			global.Logger["err"].Errorf("SaveEmail time.Parse EndTime failed,err:[%v]", err.Error())
			return
		}
		email.EndTime = endTime.Unix()
	} else {
		email.EndTime = 0
	}

	if req.IsPermanent {
		email.StartTime = 0
		email.EndTime = 0
	}

	if req.Id > 0 {
		email.ID = req.Id
	}

	var userIds = make([]uint64, 0)

	if req.UserIds != "" {
		// 切割字符串，校验是否合规 、
		for _, uidStr := range strings.Split(req.UserIds, "、") {
			var uid int
			uid, err = strconv.Atoi(uidStr)

			if err != nil {
				err = errors.New("用户id格式错误:" + uidStr)
				global.Logger["err"].Errorf("SaveEmail strconv.Atoi failed,err:[%v]", err.Error())
				return
			}
			userIds = append(userIds, uint64(uid))
		}
	}

	//在线玩家
	if req.SendType == 2 {
		userIds = msgcenter.GetOnlineUsers()
	}

	err = email.Save(tx)

	if err != nil {
		global.Logger["err"].Errorf("SaveEmail Save failed,err:[%v]", err.Error())

		tx.Rollback()
		return
	}

	//分发到子email中去
	emailUser := account.SendEmailParams{
		Title: email.Title,
		Msg:   email.Msg,
		Type:  1,
	}

	for _, id := range userIds {

		err = SendEmail(emailUser, int(id), attachment)
		if err != nil {
			global.Logger["err"].Errorf("SaveEmail SendEmail failed,err:[%v]", err.Error())
			return err
		}
	}

	if req.SendType == 1 {
		//全局邮件发送消息事件
		SendNotice()
		//删除邮件缓存
		err = rds.DelRedisCacheByKey(constant.CommonEmailList)
	}

	tx.Commit()

	return
}

// 删除邮件
func DelEmail(req request.DeleteId) (err error) {
	db := global.Game
	email := new(game.Email)
	err = email.DeleteById(db, req.Id)
	if err != nil {
		global.Logger["err"].Errorf("DelEmail DeleteById failed,err:[%v]", err.Error())
		return
	}
	return
}

// 附件列表
func GetAnnexList() (list []response.EmailAnnexList, err error) {
	db := global.Game
	annex := new(game.EmailAnnex)
	annexList, err := annex.GetList(db)
	if err != nil {
		global.Logger["err"].Errorf("GetAnnexList GetList failed,err:[%v]", err.Error())

		return
	}

	list = make([]response.EmailAnnexList, 0, len(annexList))
	for _, al := range annexList {
		list = append(list, response.EmailAnnexList{
			Id:         al.ID,
			Name:       al.Name,
			EnName:     al.EnName,
			Type:       al.Type,
			Amount:     al.Amount,
			AmountType: al.AmountType,
			Unit:       al.Unit,
			Remark:     al.Remark,
		})
	}

	return
}

// 保存附件
func SaveAnnex(req request.SaveAnnex) (err error) {
	db := global.Game
	annex := game.EmailAnnex{
		Name:       req.Name,
		EnName:     req.EnName,
		Type:       req.Type,
		Amount:     req.Amount,
		AmountType: req.AmountType,
		Unit:       req.Unit,
		Remark:     req.Remark,
		CreatedAt:  time.Now(),
	}

	if req.Id > 0 {
		annex.ID = req.Id
	}

	err = annex.Save(db)

	return
}

// 删除附件
func DelAnnex(req request.DeleteId) (err error) {
	db := global.Game
	annex := new(game.EmailAnnex)
	err = annex.DeleteById(db, req.Id)
	if err != nil {
		global.Logger["err"].Errorf("SaveEmail SendEmail failed,err:[%v]", err.Error())
		return
	}
	return
}

// 附件列表
func GetEmailEventList(req request.GetEmailEventList) (res response.EmailEventRsp, err error) {
	db := global.Game
	event := new(game.EmailEvent)
	total, eventList, err := event.GetPageList(db, req)

	if err != nil {
		global.Logger["err"].Errorf("GetEmailEventList GetPageList failed,err:[%v]", err.Error())
		return
	}

	res.Total = total

	list := make([]response.EmailEvent, 0, len(eventList))
	for _, el := range eventList {
		list = append(list, response.EmailEvent{
			Id:   el.Id,
			Name: el.Name,
		})
	}

	res.List = list

	return
}

// 保存附件
func SaveEmailEvent(req request.SaveEmailEvent) (err error) {
	db := global.Game
	event := game.EmailEvent{
		Name: req.Name,
	}

	if req.Id > 0 {
		event.Id = req.Id
	}

	err = event.Save(db)

	return
}

// 删除附件
func DelEmailEvent(req request.DeleteId) (err error) {
	db := global.Game
	event := new(game.EmailEvent)
	err = event.DeleteById(db, req.Id)
	if err != nil {
		global.Logger["err"].Errorf("DelEmailEvent DeleteById failed,err:[%v]", err.Error())
		return
	}
	return
}
