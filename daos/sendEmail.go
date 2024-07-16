package daos

import (
	"git.dev666.cc/external/dreamgo/xy"
	"gm/global"
	"gm/msgcenter"
	"gm/response"
	"za.game/lib/account"
	"za.game/lib/consts"
)

// 退款邮件
func SendRefundEmail(language string, attachMp map[int]response.Refund) (err error) {

	sendEmail := response.SendEmail{}
	if language == "en" {
		sendEmail = response.FailedMp["en"]
	} else if language == "hi" {
		sendEmail = response.FailedMp["hi"]
	}

	emailUser := account.SendEmailParams{
		Title: sendEmail.Title,
		Msg:   sendEmail.Msg,
		Type:  2,
	}

	for _, r := range attachMp {
		emailUser.Uid = uint64(r.Uid)

		attaches := make([]account.Attachment, 0, 1)

		//有附件，构造附件
		if r.HasAttach == 1 {
			tmp := account.Attachment{
				ItemId: consts.ItemIdWinCash,
				Nums:   r.Amount,
			}

			attaches = append(attaches, tmp)

			emailUser.Attachments = attaches
		}

		err = account.SendEmail(&emailUser)
		if err != nil {
			global.Logger["err"].Errorf("SendRefundEmail account sendEmail failed:" + err.Error())
			return
		}
	}

	return
}

// 发送邮件
func SendPassEmail(ids []int, language string) (err error) {
	sendEmail := response.SendEmail{}

	if language == "en" {
		sendEmail = response.SuccessMp["en"]
	} else if language == "hi" {
		sendEmail = response.SuccessMp["hi"]
	}

	emailUser := account.SendEmailParams{
		Title: sendEmail.Title,
		Msg:   sendEmail.Msg,
		Type:  1,
	}

	for _, id := range ids {
		emailUser.Uid = uint64(id)

		err = account.SendEmail(&emailUser)
		if err != nil {
			global.Logger["err"].Errorf("SendPassEmail account sendEmail failed:" + err.Error())
			return
		}
	}

	return
}

func SendEmail(params account.SendEmailParams, uid int, attachments []account.Attachment) (err error) {
	params.Uid = uint64(uid)
	params.Attachments = attachments

	err = account.SendEmail(&params)
	if err != nil {
		global.Logger["err"].Infof("editUserCoin account sendEmail failed:" + err.Error())
		return
	}

	return
}

func SendNotice() {
	//发送通知
	dicJson := make(map[string]interface{}, 0)
	dicJson["rdid"] = consts.RedDot_Email
	dicJson["count"] = 1
	msgcenter.SendNoticeEvent("reddot", dicJson, xy.UserId_(0))
}
