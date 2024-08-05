package admin

import "github.com/7058011439/haoqbb/GoAdmin/db/common"

type DictData struct {
	common.Model
	DictLabel string `json:"dictLabel"`                         // 数据标签
	DictValue string `json:"dictValue"`                         // 数据键值
	DictType  string `json:"dictType"`                          // 字典类型
	Status    string `json:"status" gorm:"size:4;default:'2';"` // 状态(1-禁止;2-启用)
	Sort      int    `json:"sort"`                              // 显示顺序
	Remark    string `json:"remark"`                            // 备注
	common.ControlBy
}

func (d *DictData) TableName() string {
	return "dict_data"
}
