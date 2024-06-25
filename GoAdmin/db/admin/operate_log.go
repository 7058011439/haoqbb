package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
)

type OperateLog struct {
	common.Model
	UserName         string  `json:"userName" gorm:"Index:idx_userName"` // 操作管理员名
	Method           string  `json:"method" gorm:"Index:idx_method"`     // 请求方式
	Url              string  `json:"url"`                                // 请求api
	FullPath         string  `json:"fullPath" gorm:"Index:idx_fullPath"` // 全路径
	Ip               string  `json:"ip"`                                 // 操作ip
	Param            string  `json:"param" gorm:"type:longtext"`         // 操作参数
	LatencyTime      float64 `json:"latencyTime"`                        // 操作耗时
	Result           string  `json:"result" gorm:"type:longtext"`        // 返回参数
	Status           string  `json:"status" gorm:"size:4;default:'2'"`   // 状态
	common.ControlBy `json:"-" gorm:"-"`
}

func (l *OperateLog) TableName() string {
	return "operate_log"
}

func (l *OperateLog) IsValid() bool {
	return l.Model.IsValid() && l.UserName != ""
}
