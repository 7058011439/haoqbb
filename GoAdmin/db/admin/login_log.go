package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
)

type LoginLog struct {
	common.Model
	UserName         string `json:"userName" gorm:"Index:idx_userName"`                // 登录用户名(账号)
	Status           string `json:"status" gorm:"size:4;default:'2';Index:idx_status"` // 状态
	Remark           string `json:"remark""`                                           // 备注
	Ip               string `json:"ip" gorm:"Index:idx_ip"`                            // 登录ip
	common.ControlBy `json:"-" gorm:"-"`
}

func (l *LoginLog) TableName() string {
	return "login_log"
}

func (l *LoginLog) IsValid() bool {
	return l.Model.IsValid() && l.UserName != ""
}
