package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gproc"
	"github.com/mojocn/base64Captcha"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
	"runtime"
	"time"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		rsp.ErrorJSON(c, ecode.CreateCaptchaError)
		return
	}
	rsp.SucJson(c, gin.H{
		"captchaId": id,
		"path":      b64s,
	})
}

func UpdateSystemDate(ctx *gin.Context) {
	param := request.UpdateTime{}

	if err := ctx.ShouldBind(&param); err != nil {
		rsp.ErrorJSON(ctx, ecode.ParamInvalid)
		return
	}
	ok := UpdateSystemDateTime(param.DateTime)

	if !ok {
		rsp.ErrorJSON(ctx, errors.New("err"))
		return
	}

	rsp.SucJson(ctx, time.Now().Format("2006-01-02"))
	return
}

func UpdateSystemDateTime(dateTime string) bool {
	system := runtime.GOOS
	switch system {
	case "windows":
		{
			_, err1 := gproc.ShellExec(`date  ` + dateTime)
			if err1 != nil {
				glog.Info("更新系统时间错误:请用管理员身份启动程序!")
				return false
			}
			return true
		}
	case "linux":
		{
			_, err1 := gproc.ShellExec(`date -s  "` + dateTime + `"`)
			if err1 != nil {
				glog.Info("更新系统时间错误:", err1.Error())
				return false
			}
			return true
		}
	case "darwin":
		{
			_, err1 := gproc.ShellExec(`date -s  "` + dateTime + `"`)
			if err1 != nil {
				glog.Info("更新系统时间错误:", err1.Error())
				return false
			}
			return true
		}
	}
	return false
}

func PayTest(c *gin.Context) {

	type ReturnOrderData struct {
		MchOrderNo  string `json:"mchOrderNo"`
		OrderState  int    `json:"orderState"`
		PayData     string `json:"payData"`
		PayDataType string `json:"payDataType"`
		PayOrderId  string `json:"payOrderId"`
		ErrCode     int    `json:"errCode"`
		ErrMsg      string `json:"errMsg"`
	}

	type OrderReturn struct {
		Code uint16          `json:"code"`
		Data ReturnOrderData `json:"data"`
		Msg  string          `json:"msg"`
		Sign string          `json:"sign"`
	}

	nano := time.Now().UnixNano()

	data := ReturnOrderData{
		MchOrderNo:  fmt.Sprintf("P1784%v", nano),
		OrderState:  0,
		PayData:     "https://www.baidu.com/",
		PayDataType: "",
		PayOrderId:  fmt.Sprintf("P1784%v", nano),
		ErrCode:     0,
		ErrMsg:      "",
	}

	res := OrderReturn{
		Code: 0,
		Data: data,
		Msg:  "SUCCESS",
		Sign: "SUCCESS",
	}

	c.JSON(0, res)
}
