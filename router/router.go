package router

import (
	"github.com/gin-gonic/gin"
	"gm/api"
)

func Router(g *gin.RouterGroup) {
	// 基础
	baseGroup := g.Group("base/")
	{
		baseGroup.GET("captcha", api.GetCaptcha)            //验证码
		baseGroup.GET("time", api.UpdateSystemDate)         //更新服务器时间
		baseGroup.POST("api/pay/unifiedOrder", api.PayTest) //支付测试
	}

	t := g.Group("test/")
	{
		t.GET("online", api.Online)
		t.GET("gv", api.GetVersion)
	}

	cb := g.Group("callback/")
	{
		cb.POST("xPay", api.XPaymentCallback)
		cb.GET("xPaySuccess", api.XPaymentSuccessCallback)
		cb.GET("xPayFailed", api.XPaymentCallbackFailed)
		cb.POST("inPay", api.InPaymentCallback)
		cb.POST("luckyPay", api.LuckyPaymentCallback)
	}

	give := g.Group("give/")
	{
		give.POST("/giveMoneyHandle", api.GiveMoneyHandle) //赠送数据处理
	}

	// 后台用户
	user := g.Group("user")
	{
		user.GET("getCaptchaQr", api.GetCaptchaQr)
		user.POST("bindCaptcha", api.BindGoogleCaptcha)
	}
}
