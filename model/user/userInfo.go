package user

import (
	"fmt"
	"gm/request"
	"gorm.io/gorm"
	"time"
)

type UserInfo struct {
	ID              int       `gorm:"primaryKey;column:id;type:int(11);"`
	Uid             int       `gorm:"column:uid;type:bigint(20);not null;default:0;"`       //用户ID
	Channel         int8      `gorm:"column:channel;type:tinyint(4);not null;default:0;"`   //渠道ID
	Sex             int8      `gorm:"column:sex;type:tinyint(4);not null;default:1;"`       //性别（1=男，2=女）
	Age             int8      `gorm:"column:age;type:tinyint(4);not null;default:0;"`       //年龄
	Vip             int       `gorm:"column:vip;type:tinyint(4);not null;default:0;"`       //vip等级
	WinCash         int       `gorm:"column:win_cash;type:int(11);not null;default:0;"`     //可用金额
	Cash            int       `gorm:"column:cash;type:int(11);not null;default:0;"`         //半冻结解金额
	Bonus           int       `gorm:"column:bonus;type:int(11);not null;default:0;"`        //冻结金额
	Recharge        int       `gorm:"column:recharge;type:int(11);not null;default:0;"`     //总充值
	Withdraw        int       `gorm:"column:withdraw;type:int(11);not null;default:0;"`     //总提现
	Bet             int       `gorm:"column:bet;type:bigint(20);not null;default:0;"`       //总下注
	Transport       int       `gorm:"column:transport;type:bigint(20);not null;default:0;"` //总输
	Win             int       `gorm:"column:win;type:bigint(20);not null;default:0;"`       //总赢
	Remark          string    `gorm:"column:remark;type:varchar(64);not null;default:'';"`  //备注
	CreatedAt       time.Time `gorm:"column:created_at;"`
	UpdatedAt       time.Time `gorm:"column:updated_at;"`
	WithdrawMoney   int       `gorm:"column:withdraw_money;type:int(11);default:0;comment:赠送审核中;"`
	WithdrawedMoney int       `gorm:"column:withdrawed_money;type:int(11);default:0;comment:赠送到账;"`
}

func (t *UserInfo) TableName(uid int) string {
	return fmt.Sprintf("user_info_0%v", uid%5)
}

func (t *UserInfo) GetOneByUid(db *gorm.DB, uid int) (info UserInfo, err error) {
	err = db.Table(t.TableName(uid)).Where("uid = ?", uid).First(&info).Error
	return
}

func (t *UserInfo) GetListByUserIds(db *gorm.DB, uMp map[int][]int) (userInfos []UserInfo, err error) {
	userInfos = make([]UserInfo, 0, 16)

	for k, userIds := range uMp {
		infos := make([]UserInfo, 0)
		err = db.Table(t.TableName(k)).
			Select("uid,recharge,win_cash,cash,withdrawed_money").
			Where("uid in (?)", userIds).
			Find(&infos).
			Error
		if err != nil {
			return
		}
		userInfos = append(userInfos, infos...)
	}

	return
}

func (t *UserInfo) GetListUnionAll(db *gorm.DB, req request.GetGameUserList) (userInfos []UserInfo, err error) {

	//没有查询条件，直接返回
	if req.PayStatus == nil && req.AssetMin == 0 && req.AssetMax == 0 && req.RechargeMin == 0 && req.RechargeMax == 0 {
		return
	}

	whereStr := ""

	if req.PayStatus != nil {
		if *req.PayStatus == 0 {
			whereStr += " and recharge = 0"
		} else {
			whereStr += " and recharge > 0"
		}
	}

	if req.AssetMin > 0 {
		whereStr += fmt.Sprintf(" and cash + win_cash >= %v", req.AssetMin)
	}

	if req.AssetMax > 0 {
		whereStr += fmt.Sprintf(" and cash + win_cash <= %v", req.AssetMax)
	}

	if req.RechargeMin > 0 {
		whereStr += fmt.Sprintf(" and recharge >= %v", req.RechargeMin)
	}

	if req.RechargeMax > 0 {
		whereStr += fmt.Sprintf(" and recharge <= %v", req.RechargeMax)
	}

	sql := fmt.Sprintf("select uid,recharge,win_cash,cash from user_info_00 where 1=1 %s "+
		"union all select uid,recharge,win_cash,cash from user_info_01 where 1=1 %s "+
		"union all select uid,recharge,win_cash,cash from user_info_02 where 1=1 %s "+
		"union all select uid,recharge,win_cash,cash from user_info_03 where 1=1 %s "+
		"union all select uid,recharge,win_cash,cash from user_info_04 where 1=1 %s ",
		whereStr,
		whereStr,
		whereStr,
		whereStr,
		whereStr,
	)

	err = db.Raw(sql).Scan(&userInfos).Error

	return
}

func (t *UserInfo) UpdateByUid(db *gorm.DB, uid int, mp map[string]interface{}) (err error) {
	err = db.Table(t.TableName(uid)).Where("uid = ?", uid).Updates(mp).Error
	return
}
