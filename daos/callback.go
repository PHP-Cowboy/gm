package daos

import (
	"encoding/json"
	"errors"
	"fmt"
	dream "git.dev666.cc/external/dreamgo"
	"git.dev666.cc/external/dreamgo/xy"
	"github.com/gin-gonic/gin"
	"gm/common/constant"
	"gm/daos/rds"
	"gm/global"
	"gm/model/pay"
	"gm/model/user"
	"gm/msgcenter"
	"gm/response"
	"gm/utils/ecode"
	req2 "gm/utils/request"
	"gm/utils/timeutil"
	"gorm.io/gorm"
	"time"
	"za.game/lib/account"
	"za.game/lib/consts"
	"za.game/lib/dbconn"
)

func NotifyLamp(userid, amount int, userName string) {
	dicJson := dream.NewJSON()
	dicJson["txt"] = fmt.Sprintf("Congratulations!<color=#4ab3f5> %s</color> withdraws <color=#ffe000>₹%.2f</color> successfully!",
		userName, float64(amount/100))
	dicJson["userid"] = 0    //发给谁的，0就是发给全用户 基本上是服务器在用，客户端不用管
	dicJson["type"] = 2      //类型 1游戏 2活动 3系统
	dicJson["times"] = 1     //播放次数
	dicJson["from"] = userid //触发这条信息的用户id
	dicJson["priority"] = 1  //优先级 数值越大的优先级越低
	v := &xy.NoticeEvent{
		Evt: "NoticeMsg",
		Dic: dicJson,
	}
	v.Uid = append(v.Uid, xy.UserId_(0))
	msgcenter.SendBackPackageJson(xy.XyNoticeEvent, v)
}

// 根据 markers 获取 支付配置
func GetSecretByMarkers(markers string) (payCfg pay.PayConfig, err error) {
	var (
		jsonTxt string
		payList = make([]pay.PayConfig, 0)
	)

	jsonTxt, err = global.Redis.Get(consts.PayConfig)
	if err != nil || jsonTxt == "" {
		payDb := global.Pay
		payCfgObj := new(pay.PayConfig)

		payList, err = payCfgObj.GetList(payDb)
		if err != nil {
			global.Logger["err"].Infof("select pay_config failed, err:", err.Error())
			return
		}

		b := make([]byte, 0)

		b, err = json.Marshal(payList)
		if err != nil {
			global.Logger["err"].Infof("GetSecretByMarkers json.Marshal failed, err:", err.Error())
			return
		}

		jsonTxt = string(b)

		_, _ = global.Redis.Set(consts.PayConfig, jsonTxt, 86400)

	} else {
		err = json.Unmarshal([]byte(jsonTxt), &payList)
		if err != nil {
			global.Logger["err"].Infof("GetSecretByMarkers json.Unmarshal failed, err:", err.Error())
			return
		}
	}

	for _, cfg := range payList {
		if cfg.Markers == markers {
			payCfg = cfg
			break
		}
	}

	if payCfg.ID == 0 {
		global.Logger["err"].Infof("GetSecretByMarkers markers 相关数据未找到")
		err = errors.New("markers not found")
		return
	}

	return
}

func SuccessCallback(orderNo, mOrderNo string, successTime int) (err error) {

	payDb := global.Pay
	var (
		giveObj pay.GiveMoney
		give    pay.GiveMoney
	)

	give, err = giveObj.GetOneByOrderNo(payDb, orderNo)
	if err != nil {
		global.Logger["err"].Infof("SuccessCallback giveObj.GetOneByOrderNo failed,err:[%s]", err.Error())
		return
	}

	if give.Status == pay.GiveMoneyStatusComplete {
		err = ecode.OrderStatusAbnormal
		global.Logger["err"].Infof("SuccessCallback give.Status == pay.GiveMoneyStatusComplete")
		return
	}
	var arrivalTime time.Time

	if successTime > 0 {
		arrivalTime = time.Unix(int64(successTime/1000), 0)
	} else {
		arrivalTime = time.Now()
	}

	err = giveObj.UpdateById(payDb, give.Id, map[string]interface{}{
		"trd_order_no": mOrderNo,
		"status":       pay.GiveMoneyStatusComplete,
		"arrival_time": arrivalTime,
	})

	if err != nil {
		global.Logger["err"].Infof("giveObj.UpdateById failed, err:%s", err.Error())
		return
	}

	userInfoObj := new(account.UserInfo)

	//修改用户赠送审核中，提现到账金额
	err = userInfoObj.UpdateWithdrawByUid(dbconn.NDB, -1*give.Amount, give.Amount, uint64(give.Uid))

	if err != nil {
		global.Logger["err"].Infof("account.UpdateWithdrawByUid failed, err:[%s]", err.Error())
		return
	}

	err = SendPassEmail([]int{give.Uid}, "en")

	if err != nil {
		global.Logger["err"].Infof("SuccessCallback SendPassEmail failed, orderNo:[%s], err:[%s]", orderNo, err.Error())
		return
	}

	userDb := global.User

	var u user.User

	uObj := new(user.User)

	u, err = uObj.GetOneByUid(userDb, give.Uid)
	if err != nil {
		global.Logger["err"].Infof("uObj.GetOneByUid failed, err:[%s]", err.Error())
		return
	}

	if give.Amount >= 100000 {
		NotifyLamp(give.Uid, give.Amount, u.UserName)
	}

	return
}

func FailedCallback(orderNo, mOrderNo string, status int) (err error) {
	payDb := global.Pay
	var (
		giveObj pay.GiveMoney
		give    pay.GiveMoney
	)

	give, err = giveObj.GetOneByOrderNo(payDb, orderNo)
	if err != nil {
		global.Logger["err"].Infof("orderNo not found orderNo: %v", orderNo)
		return
	}

	if give.Status != pay.GiveMoneyStatusInPayment {
		global.Logger["err"].Infof("give.Status != pay.GiveMoneyStatusInPayment give_money id: %v", give.Id)
		err = ecode.OrderStatusAbnormal
		return
	}

	err = giveObj.UpdateById(payDb, give.Id, map[string]interface{}{
		"trd_order_no": mOrderNo,
		"status":       status,
	})

	if err != nil {
		global.Logger["err"].Infof("giveObj.UpdateById failed, err:%s", err.Error())
		return
	}

	if give.Type == pay.GiveMoneyTypeNormal {
		//普通类型提现失败，退币 && 返还提现次数
		attachMp := make(map[int]response.Refund, 0)

		attachMp[give.Id] = response.Refund{
			Uid:       give.Uid,
			HasAttach: 1,
			Amount:    give.Amount,
		}

		//发邮件  附件退钱
		err = SendRefundEmail("en", attachMp)
		if err != nil {
			global.Logger["err"].Infof("FailedCallback refund SendRefundEmail failed, err:[%s]", err.Error())
			return
		}

		if time.Now().Format(timeutil.DateNumberFormat) == give.CreatedAt.Format(timeutil.DateNumberFormat) {
			userDb := global.User
			//更新提现次数、提现额度
			vipObj := new(user.VipUserInfo)

			err = vipObj.UpdateWithdrawByUid(
				userDb,
				give.Uid,
				map[string]interface{}{
					"withdraw_num":   gorm.Expr("withdraw_num = withdraw_num - ?", 1),
					"withdraw_money": gorm.Expr("withdraw_money = withdraw_money - ?", give.Amount),
				},
			)

			if err != nil {
				global.Logger["err"].Infof("FailedCallback vipObj.UpdateWithdrawByUid failed, err:[%s]", err.Error())
				//这个失败了不让三方重新调用，改成nil
				err = nil
			}

			err = rds.DelRedisCacheByKey(consts.RedisVipUser + fmt.Sprintf("%v", give.Uid))
			if err != nil {
				global.Logger["err"].Infof("FailedCallback rds.DelRedisCacheByKey failed, err:[%s]", err.Error())
				//删缓存失败不处理
				return nil
			}
		}

		//更新用户提现记录表
		gUserObj := new(pay.GaveMoneyUser)

		err = gUserObj.UpdateNumByUidGiveId(payDb, give.Uid, give.GiveId, -1)

		if err != nil {
			global.Logger["err"].Errorf("FailedCallback gUserObj.UpdateNumByUidGiveId failed, err:[%s]", err.Error())
			return
		}

	} else if give.Type == pay.GiveMoneyTypeNoviceCarnival {
		//新手嘉年华 类型提现失败，发提现失败邮件
		attachMp := make(map[int]response.Refund, 0)

		attachMp[give.Id] = response.Refund{
			Uid:       give.Uid,
			HasAttach: 0,
			Amount:    0,
		}

		//只发邮件  不退钱
		err = SendRefundEmail("en", attachMp)

		if err != nil {
			global.Logger["err"].Infof("FailedCallback No refund SendRefundEmail failed, err:[%s]", err.Error())
			return
		}
		//新手嘉年华提现失败，调用大厅方法
		NoviceCarnivalWithdrawRefuse(give.OrderNo)
	}

	return
}

// 新手嘉年华 赠送 失败 或 拒绝
func NoviceCarnivalWithdrawRefuse(orderNo string) {

	url := global.ServerConfig.Srv.Hall + constant.NoviceCarnivalWithdraw

	var (
		err  error
		res  response.NoviceCarnivalWithdraw
		body []byte
	)

	param := gin.H{"order_no": orderNo}

	body, err = req2.SrvPost(url, param)

	if err != nil {
		global.Logger["err"].Infof("新手嘉年华提现失败，调用大厅方法 失败:[%s]", err.Error())
		return
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		global.Logger["err"].Infof("SrvPost 返回值解析失败 body:[%s],res:[%v]", string(body), err.Error())
		return
	}

	if res.Code != 200 {
		global.Logger["err"].Infof("新手嘉年华提现失败，调用大厅方法 失败:[%s]", err.Error())
	}

	return
}
