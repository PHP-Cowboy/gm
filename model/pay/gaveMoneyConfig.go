package pay

import (
	"gm/model"
	"gm/request"
	"gorm.io/gorm"
)

// 赠送配置
type GaveMoneyConfig struct {
	Id        int    `gorm:"primaryKey;type:int(11) AUTO_INCREMENT;comment:id"`
	Name      string `gorm:"type:varchar(32);not null;default:'';comment:提现名称"`
	Account   int    `gorm:"type:int(11);not null;default:0;comment:提现金额"`
	Status    int    `gorm:"type:tinyint(4);not null;default:1;comment:状态(1=在用,2=停用)"` //状态(1=在用,2=停用)
	Type      int    `gorm:"column:_type;type:tinyint(4);not null;default:1;comment:赠送类型(1:正常，2:首提)"`
	ReplaceId int    `gorm:"type:int;not null;default:0;comment:替换的目标id"`
}

const (
	GaveMoneyConfigZero   = iota
	GaveMoneyConfigNormal //正常
	GaveMoneyConfigFirst  //首提
)

func (t *GaveMoneyConfig) Save(db *gorm.DB) (err error) {
	err = db.Model(t).Save(t).Error
	return
}

func (t *GaveMoneyConfig) DeleteById(db *gorm.DB, id uint64) (err error) {
	err = db.Model(t).Delete("id = ?", id).Error
	return
}

func (t *GaveMoneyConfig) GetOneById(db *gorm.DB, id int) (data GaveMoneyConfig, err error) {
	err = db.Model(t).First(&data, id).Error
	return
}

func (t *GaveMoneyConfig) GetOneByReplaceId(db *gorm.DB, replace_id int) (data GaveMoneyConfig, err error) {
	err = db.Model(t).Where("replace_id = ?", replace_id).First(&data).Error
	return
}

func (t *GaveMoneyConfig) GetPageList(db *gorm.DB, req request.GaveConfigPageList) (total int64, dataList []GaveMoneyConfig, err error) {
	localDb := db.Model(t).Where(GaveMoneyConfig{Name: req.Name, Status: req.Status, Type: req.Type})

	err = localDb.Count(&total).Error

	if err != nil {
		return
	}

	err = localDb.Scopes(model.Paginate(req.Page, req.Size)).Find(&dataList).Error
	return
}

func (t *GaveMoneyConfig) GetList(db *gorm.DB, req request.GaveConfigList) (list []GaveMoneyConfig, err error) {
	err = db.Model(t).Where(GaveMoneyConfig{Type: req.Type, Status: req.Status}).Find(&list).Error
	return
}
