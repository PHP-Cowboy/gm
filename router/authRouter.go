package router

import (
	"github.com/gin-gonic/gin"
	"gm/api"
	"gm/middlewares"
)

func AuthRouter(g *gin.RouterGroup) {
	g.Use(middlewares.JWTAuth())

	// 角色
	role := g.Group("role")
	{
		role.POST("/add", api.AddRole)                //新增角色
		role.POST("/change", api.ChangeRole)          //修改角色
		role.GET("/list", api.RoleList)               //角色列表
		role.GET("/user_role_list", api.UserRoleList) //用户角色列表
		role.POST("/add_user_role", api.AddUserRole)  //用户新增角色
	}

	// 菜单
	menu := g.Group("menu")
	{
		menu.POST("/add", api.AddMenu)                //新增菜单
		menu.POST("/change", api.ChangeMenu)          //修改菜单
		menu.GET("/list", api.MenuList)               //菜单列表
		menu.GET("/tree", api.MenuTree)               //菜单树列表
		menu.GET("/getLevelList", api.LevelList)      //菜单列表
		menu.GET("/role_menu_list", api.RoleMenuList) //角色菜单列表
		menu.POST("/add_role_menu", api.AddRoleMenu)  //新增角色菜单
	}

	// 公共模块
	common := g.Group("common")
	{
		common.GET("/channelPageList", api.ChannelPageList) //渠道列表
		common.GET("/allChannel", api.AllChannelList)       //角色渠道列表
		common.GET("/roleChannel", api.RoleChannelList)     //角色渠道列表
		common.POST("/saveChannel", api.SaveChannel)        //角色渠道列表
		common.GET("/gameList", api.GameList)               //游戏列表
		common.GET("/roomList", api.CommonRoomList)         //房间列表
	}

	// 字典
	dict := g.Group("dict")
	{
		dict.GET("/typeList", api.DictTypeList)      //字典类型列表
		dict.GET("/list", api.DictList)              //字典列表
		dict.GET("/one", api.GetOneDict)             //单条字典数据
		dict.POST("/save", api.SaveDict)             //字典列表
		dict.POST("/changeValues", api.ChangeValues) //变更字典数据
	}

	// 邮件
	email := g.Group("email")
	{
		email.GET("/list", api.GetEmailList)           //邮件列表
		email.POST("/save", api.SaveEmail)             //保存邮件
		email.POST("/del", api.DelEmail)               //删除邮件
		email.GET("/annexList", api.GetAnnexList)      //附件列表
		email.POST("/saveAnnex", api.SaveAnnex)        //保存附件
		email.POST("/delAnnex", api.DelAnnex)          //删除附件
		email.GET("/eventList", api.GetEmailEventList) //邮件关联事件列表
		email.POST("/saveEvent", api.SaveEmailEvent)   //保存邮件关联事件
		email.POST("/delEvent", api.DelEmailEvent)     //删除邮件关联事件
	}

	// 资金流水
	fundsFlow := g.Group("fundsFlow")
	{
		fundsFlow.GET("/tp", api.TP)     //tp 游戏流水
		fundsFlow.GET("/slot", api.SLOT) //slot 游戏流水
	}

	// 游戏用户数据
	gu := g.Group("gameUser")
	{
		gu.GET("/list", api.GetGameUserList)                  //游戏用户列表
		gu.GET("/info", api.GameUserInfo)                     //用户信息
		gu.POST("/changeRecharge", api.ChangeRecharge)        //修改用户累计充值总额
		gu.POST("/editUserCoin", api.EditUserCoin)            //修改用户金币
		gu.GET("/withdrawInfoRecord", api.WithdrawInfoRecord) //修改用户累计充值总额
		gu.POST("/banned", api.Banned)                        //封禁
		gu.POST("/unseal", api.Unseal)                        //解封
		gu.GET("/loginLogList", api.GetLoginLogList)          // 登录日志列表
		gu.GET("/giveList", api.GiveList)                     // 登录日志列表
	}

	//游戏用户数据 批量操作
	b := g.Group("batch")
	{
		b.POST("/pass", api.BatchPass)       //批量通过
		b.POST("/repulse", api.BatchRepulse) //批量拒绝
		b.POST("/cancel", api.BatchCancel)   //批量主动取消
		b.POST("/invalid", api.BatchInvalid) //批量作废
	}

	// 礼包
	giftPack := g.Group("giftPack")
	{
		//二选一
		giftPack.GET("/getEventList", api.GetEventList)  //获取二选一 配置列表
		giftPack.POST("/saveEvent", api.SaveEventConfig) //保存二选一 配置
		giftPack.POST("/delEvent", api.DelEventConfig)   //删除二选一 配置
		giftPack.POST("/onOffEvent", api.OnOffEvent)     //二选一开关
		giftPack.POST("/eventStatus", api.EventStatus)   //开关状态
		//充200送200
		giftPack.GET("/getRechargeGiftList", api.GetRechargeGiftList)  //充200送200 配置列表
		giftPack.POST("/saveRechargeGift", api.SaveRechargeGiftConfig) //保存充200送200 配置
		giftPack.POST("/delRechargeGift", api.DelRechargeGift)         //删除充200送200 配置
		giftPack.POST("/onOffRechargeGift", api.OnOffRechargeGift)     //充200送200开关
		giftPack.POST("/rechargeGiftStatus", api.RechargeGiftStatus)   //开关状态
		//充值礼包(游戏内充值)
		giftPack.GET("/getRechargePackList", api.GetRechargePackList) //充值礼包 配置列表
		giftPack.POST("/saveRechargePack", api.SaveRechargePack)      //保存充值礼包 配置
		giftPack.POST("/delRechargePack", api.DelRechargePack)        //删除充值礼包 配置
		//vip
		giftPack.GET("/getVipList", api.GetVipConfigList) //获取vip配置列表
		giftPack.POST("/saveVip", api.SaveVipConfig)      //保存vip配置
		giftPack.POST("/delVip", api.DelVipConfig)        //删除vip 配置
		//救济金
		giftPack.GET("/getBenefitList", api.GetBenefitList) //救济金 配置列表
		giftPack.POST("/saveBenefit", api.SaveBenefit)      //保存救济金 配置
		giftPack.POST("/onOffBenefit", api.OnOffBenefit)    //更新救济金礼包开关
		giftPack.POST("/benefitStatus", api.BenefitStatus)  //开关状态
		giftPack.POST("/delBenefit", api.DelBenefit)        //删除救济金 配置
		//三选一礼包
		giftPack.GET("/getOnlyList", api.GetOnlyList) //三选一 配置列表
		giftPack.POST("/saveOnly", api.SaveOnly)      //保存三选一 配置
		giftPack.POST("/delOnly", api.DelOnly)        //删除三选一 配置
		giftPack.POST("/onOffOnly", api.OnOffOnly)    //三选一开关
		giftPack.POST("/onlyStatus", api.OnlyStatus)  //开关状态
	}

	// 充值
	pay := g.Group("pay")
	{
		pay.GET("/list", api.PayList)                         //充值礼包列表
		pay.POST("/save", api.SaveGift)                       // 保存充值礼包
		pay.POST("/del", api.DelGift)                         // 删除充值礼包
		pay.GET("/configList", api.ConfigList)                //支付渠道配置列表
		pay.POST("/saveConfig", api.SaveConfig)               // 保存支付渠道配置
		pay.POST("/delConfig", api.DelConfig)                 // 删除支付渠道配置
		pay.GET("/bankList", api.BankList)                    //获取用户银行卡列表
		pay.GET("/orderList", api.OrderList)                  //获取用户充值记录列表
		pay.GET("/rechargeRecords", api.RechargeRecords)      //充值记录
		pay.GET("gaveConfigPageList", api.GaveConfigPageList) //赠送配置
		pay.GET("gaveConfigList", api.GaveConfigList)         //赠送配置
		pay.POST("saveGiveConfig", api.SaveGiveConfig)        //保存赠送配置
		pay.POST("delGiveConfig", api.DelGiveConfig)          //删除赠送配置
	}

	// 房间
	room := g.Group("room")
	{
		room.GET("/list", api.RoomList)                              //房间列表
		room.GET("/getColumnComment", api.GetColumnComment)          //字段 && comment
		room.POST("/save", api.SaveRoom)                             //房间保存
		room.POST("/updateExtDataByExcel", api.UpdateExtDataByExcel) //根据excel更新ExtData的值
	}

	// 签到
	sign := g.Group("sign")
	{
		sign.GET("/getList", api.GetSignConfigList)  //获取签到配置列表
		sign.POST("/saveSign", api.SaveSign)         //保存签到配置
		sign.POST("/delSign", api.DelSign)           //删除签到配置
		sign.GET("/prizeList", api.GetSingPrizeList) //获取签到奖励配置列表
		sign.POST("/savePrize", api.SaveSingPrize)   //保存签到奖励配置
		sign.POST("/delPrize", api.DelSingPrize)     //删除签到奖励配置
	}

	// 后台用户
	user := g.Group("user")
	{
		user.POST("/loginOut", api.LoginOut)                              //登出
		user.POST("/add", api.AddUser)                                    //新增用户
		user.POST("/change", api.ChangeUser)                              //修改 名称 密码 状态
		user.GET("/list", api.UserList)                                   //获取用户列表
		user.GET("/tree", api.UserRoleMenuTree)                           //获取用户权限菜单列表
		user.POST("updateUserGoogleCaptcha", api.UpdateUserGoogleCaptcha) //更新用户谷歌验证码
		user.GET("getCaptchaQrBySecret", api.GetCaptchaQrBySecret)
		user.POST("checkGoogleCaptcha", api.CheckGoogleCaptcha)
		user.POST("bindGoogleCaptcha", api.BindGoogleCaptcha)
	}

	task := g.Group("/task")
	{
		task.GET("/reportDataTask", api.ReportDataTask)                               //报表数据任务
		task.GET("/mergeAllChannelReportDataTask", api.MergeAllChannelReportDataTask) //报表数据全渠道合并任务
		task.GET("/rechargeStatisticsTask", api.RechargeStatisticsTask)               //充值数据统计任务
		task.GET("/withdrawalStatisticsTask", api.WithdrawalStatisticsTask)           //提现用户统计任务
		task.GET("/paidUserRetentionTask", api.PaidUserRetentionTask)                 //付费用户留存统计任务
		task.GET("/userRetentionTask", api.UserRetentionTask)                         //用户留存统计任务
		task.GET("/hourUserDataTaskOnline", api.HourUserDataTaskOnline)               //小时用户数据任务
		task.GET("/hourGameUserDataTask", api.HourGameUserDataTask)                   //小时游戏用户数据任务
		task.GET("/fiveMinuteDataTask", api.FiveMinuteDataTask)                       //五分钟数据
		task.GET("/nextDayRemainedTask", api.NextDayRemainedTask)                     //次留任务
		task.GET("/threeDayRemainedTask", api.ThreeDayRemainedTask)                   //3留任务
		task.GET("/fourDayRemainedTask", api.FourDayRemainedTask)                     //4留任务
		task.GET("/fiveDayRemainedTask", api.FiveDayRemainedTask)                     //5留任务
		task.GET("/sixDayRemainedTask", api.SixDayRemainedTask)                       //6留任务
		task.GET("/sevenDayRemainedTask", api.SevenDayRemainedTask)                   //7留任务
		task.GET("/fourteenDayRemainedTask", api.FourteenDayRemainedTask)             //14留任务
	}

	// 游戏用户统计
	statistics := g.Group("statistics")
	{
		statistics.GET("/user", api.UserStatistics)                 // 用户数据统计[大盘数据]
		statistics.GET("/recharge", api.RechargeStatistics)         //充值统计
		statistics.GET("/withdrawal", api.WithdrawalStatistics)     // 提现用户统计
		statistics.GET("/paidUserRetention", api.PaidUserRetention) // 付费用户留存
		statistics.GET("/userRetention", api.UserRetention)         // 用户留存
		statistics.GET("/perHourDataNum", api.PerHourDataNum)       // 每小时数据统计
		statistics.GET("/perHourGameNum", api.PerHourGameNum)       // 每小时游戏数据统计
		statistics.GET("/fiveMinuteData", api.FiveMinuteData)       // 五分钟数据统计
	}

	//代收代付
	passage := g.Group("passage/")
	{
		passage.GET("list", api.PassageList)
		passage.POST("save", api.PassageSave)
		passage.POST("del", api.PassageDel)
		passage.POST("change", api.PassageChange)
	}

	//在线用户
	ol := g.Group("online/")
	{
		ol.GET("list", api.OnlineList) //在线用户列表
	}
}
