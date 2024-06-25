package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Post struct {
	common.Model
	PostName string `json:"postName" gorm:"unique;"`          //岗位名称
	PostCode string `json:"postCode" gorm:"unique;"`          //岗位代码
	Sort     int    `json:"sort"`                             //岗位排序
	Status   string `json:"status" gorm:"size:4;default:'2'"` //状态
	Remark   string `json:"remark"`                           //描述
	common.ControlBy
}

func (p *Post) TableName() string {
	return "post"
}

func (p *Post) IsValid() bool {
	return p.Model.IsValid() && p.PostName != ""
}

func (p *Post) BeforeDelete(tx *gorm.DB) (err error) {
	userCount := int64(0)
	tx.Model(&User{}).Where(&User{PostId: p.ID}).Count(&userCount)
	if userCount > 0 {
		return fmt.Errorf("当前岗位还有用户，请先将用户移到其他岗位")
	}
	return nil
}
