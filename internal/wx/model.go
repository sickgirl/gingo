package wx

import (
	"github.com/songcser/gingo/config"
	"github.com/songcser/gingo/pkg/model"
)

//TWxUsers 微信用户表
type TWxUsers struct {
	model.BaseModel
	Openid  string `gorm:"column:openid;type:varchar(255);comment:微信用户的openid;NOT NULL" json:"openid"`
	Appid   int    `gorm:"column:appid;type:int(11);comment:对应的appid;NOT NULL" json:"appid"`
	Unionid string `gorm:"column:unionid;type:varchar(255);comment:对应微信用户的unionid;NOT NULL" json:"unionid"`
}

func (t *TWxUsers) TableName() string {
	return "t_wx_users"
}

func (t *TWxUsers) FindOneByOpenid(openid string) (TWxUsers, error) {
	res := TWxUsers{}
	query := config.GVA_DB.Model(&t).Where("openid = ?", openid).First(&res)
	return res, query.Error
}

func (t *TWxUsers) Create() (int64, error) {
	query := config.GVA_DB.Model(&t).Create(&t)
	return t.ID, query.Error
}

// 小程序 员工表
type TUsers struct {
	model.BaseModel
	Name  string `gorm:"column:name;type:varchar(255);comment:姓名;NOT NULL" json:"name"`
	Phone string `gorm:"column:phone;type:varchar(255);comment:电话" json:"phone"`
	WxId  int    `gorm:"column:wx_id;type:int(11);comment:关联微信用户表id" json:"wx_id"`
}

func (t *TUsers) TableName() string {
	return "t_users"
}

func (t *TUsers) FindOneByPhone(phone string) (TUsers, error) {
	res := TUsers{}
	query := config.GVA_DB.Model(&t).Where("phone = ?", phone).First(&res)
	return res, query.Error
}

func (t *TUsers) BindUser(phone string, wxId int) error {
	return config.GVA_DB.Model(&t).Where("phone = ?", phone).Update("wx_id", wxId).Error
}
