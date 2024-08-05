package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/jinzhu/gorm"
)

type DictType struct {
	common.Model
	DictName string `json:"dictName" gorm:"unique;"`          // 字典名称
	DictType string `json:"dictType" gorm:"unique;"`          // 字典类型
	Status   string `json:"status" gorm:"size:4;default:'2'"` // 状态(1-禁止;2-启用)
	Remark   string `json:"remark"`                           // 备注
	common.ControlBy
}

func (d *DictType) TableName() string {
	return "dict_type"
}

func (d *DictType) BeforeDelete(tx *gorm.DB) error {
	data := LoadForumData(d).(*DictType)
	return tx.Where(&DictData{DictType: data.DictType}).Delete(&DictData{}).Error
}
