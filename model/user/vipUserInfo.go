package user

import (
	"database/sql"
	"gm/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type VipUserInfo struct {
	Uid           int          `db:"uid" gorm:"primaryKey;column:uid;type:bigint(20);" json:"uid" `                    // 用户id
	DayGeted      time.Time    `db:"day_geted" gorm:"column:day_geted;type:datetime;" json:"day_geted" `               // 每天奖励，下次可领取时间
	WeekGeted     time.Time    `db:"week_geted" gorm:"column:week_geted;type:datetime;" json:"week_geted" `            // 每周奖励，下次可领取时间
	MonGeted      time.Time    `db:"mon_geted" gorm:"column:mon_geted;type:datetime;" json:"mon_geted" `               // 每月奖励，下次可领取时间
	WithdrawNum   int          `db:"withdraw_num" gorm:"column:withdraw_num;type:int(11);" json:"withdraw_num" `       // 每日已提现次数
	WithdrawMoney int          `db:"withdraw_money" gorm:"column:withdraw_money;type:int(11);" json:"withdraw_money" ` // 每日已提现额度
	WithdrawTime  sql.NullTime `db:"withdraw_time" gorm:"column:withdraw_time;type:datetime;" json:"withdraw_time" `   // 上次提现的时间
}

func (t VipUserInfo) UpdateWithdrawByUid(db *gorm.DB, uid int, mp map[string]interface{}) (err error) {
	err = db.Model(t).Where("uid = ?", uid).Updates(mp).Error
	return
}

func (t *VipUserInfo) CreateInBatches(db *gorm.DB, list []VipUserInfo) (err error) {
	err = db.Model(t).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "uid"}}, // 冲突时的列
				DoUpdates: clause.AssignmentColumns([]string{
					"withdraw_num",
					"withdraw_money",
				}), // 要更新的列
			},
		).
		CreateInBatches(&list, model.BatchSize).
		Error
	return
}

func (t VipUserInfo) GetListByUserIds(db *gorm.DB, userIds []int) (list []VipUserInfo, err error) {
	err = db.Model(t).Where("uid in (?)", userIds).Find(&list).Error
	return
}
