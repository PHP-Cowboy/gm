package game

import (
	"gm/model"
	"gm/request"
	"gorm.io/gorm"
	"time"
)

// 救济金礼包配置
type BenefitGiftPackConfig struct {
	ID              int       `gorm:"primarykey;column:id;type:int(11);" json:"id"`
	UserType        string    `gorm:"column:user_type;type:varchar(64);not null;default:'';comment:用户类型;" json:"user_type"`                  //用户类型
	Minimum         int       `gorm:"column:minimum;type:int(11);not null;default:0;comment:累计充值最小额度;" json:"minimum"`                       //累计充值最小额度
	Maximum         int       `gorm:"column:maximum;type:int(11);not null;default:0;comment:累计充值最大额度;" json:"maximum"`                       //累计充值最大额度
	MiniTimes       int       `gorm:"column:mini_times;type:int(11);not null;default:0;comment:连续触发破产礼包但未购买次数下限;" json:"mini_times"`         //连续触发破产礼包但未购买次数
	MaxiTimes       int       `gorm:"column:maxi_times;type:int(11);not null;default:0;comment:连续触发破产礼包但未购买次数上限;" json:"maxi_times"`         //连续触发破产礼包但未购买次数
	Quota           int       `gorm:"column:quota;type:int(11);not null;default:0;comment:礼包额度;" json:"quota"`                               //礼包额度
	Value           int       `gorm:"column:value;type:int(11);not null;default:0;comment:价值;" json:"value"`                                 //价值
	BasicRewards    int       `gorm:"column:basic_rewards;type:int(11);not null;default:0;comment:基础奖励;" json:"basic_rewards"`               //基础奖励
	BasicType       uint8     `gorm:"column:basic_type;type:tinyint(4);not null;default:1;comment:金额类别1=可提现可下注,2=不可提现可下注" json:"basic_type"` //金额类别(1=可提现可下注,2=不可提现可下注)
	RewardGiveaways int       `gorm:"column:reward_giveaways;type:int(11);not null;default:0;comment:赠送奖励;" json:"reward_giveaways"`         //赠送奖励
	GiftType        uint8     `gorm:"column:gift_type;type:tinyint(4);not null;default:1;comment:金额类别1=可提现可下注,2=不可提现可下注" json:"gift_type"`   //金额类别(1=可提现可下注,2=不可提现可下注)
	Bonus           int       `gorm:"column:bonus;type:int(11);not null;default:0;comment:bonus;" json:"bonus"`                              //bonus
	BonusType       uint8     `gorm:"column:bonus_type;type:tinyint(4);not null;default:1;comment:金额类别1=可提现可下注,2=不可提现可下注" json:"bonus_type"` //金额类别(1=可提现可下注,2=不可提现可下注)
	Ratio           int       `gorm:"column:ratio;type:int(11);not null;default:0;comment:优惠比例(百分比,填整数值);" json:"ratio"`                     //优惠比例
	IsClose         int       `gorm:"column:is_close;type:tinyint(4);not null;default:0;comment:0:未关闭,1:已关闭;" json:"is_close"`               //是否关闭
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime;" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime;" json:"updated_at"`
}

func (t *BenefitGiftPackConfig) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *BenefitGiftPackConfig) UpdateById(db *gorm.DB, id int, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id = ?", id).Updates(mp).Error
	return
}

func (t *BenefitGiftPackConfig) UpdateByIds(db *gorm.DB, ids []int, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("id in (?)", ids).Updates(mp).Error
	return
}

func (t *BenefitGiftPackConfig) DeleteById(db *gorm.DB, id int) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *BenefitGiftPackConfig) GetFirstById(db *gorm.DB, id int64) (BenefitGiftPackConfig BenefitGiftPackConfig, err error) {
	err = db.Model(t).Where("id = ?", id).First(&BenefitGiftPackConfig).Error
	return
}

func (t *BenefitGiftPackConfig) GetFirstByIsClose(db *gorm.DB, isClose int) (BenefitGiftPackConfig BenefitGiftPackConfig, err error) {
	err = db.Model(t).Where("is_close = ?", isClose).First(&BenefitGiftPackConfig).Error
	return
}

func (t *BenefitGiftPackConfig) GetList(db *gorm.DB) (list []BenefitGiftPackConfig, err error) {
	err = db.Model(t).Find(&list).Error
	return
}

func (t *BenefitGiftPackConfig) GetPageList(db *gorm.DB, req request.GetBenefitList) (total int64, list []BenefitGiftPackConfig, err error) {
	localDb := db.Model(t)

	err = localDb.Count(&total).Error

	if err != nil {
		return
	}

	err = db.Model(t).Scopes(model.Paginate(req.Page, req.Size)).Find(&list).Error
	return
}
