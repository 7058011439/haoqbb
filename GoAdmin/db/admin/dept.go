package admin

import (
	"fmt"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/jinzhu/gorm"
)

type Dept struct {
	common.Model
	ParentId int64  `json:"parentId"`                         // 上级部门
	DeptName string `json:"deptName"  gorm:"unique;"`         // 部门名称
	Sort     int    `json:"sort"`                             // 排序
	Leader   string `json:"leader"`                           // 负责人
	Phone    string `json:"phone" gorm:"size:11;"`            // 手机
	Email    string `json:"email"`                            // 邮箱
	Status   string `json:"status" gorm:"size:4;default:'2'"` // 状态(1-禁止;2-启用)
	common.ControlBy
	Children []*Dept `json:"children" gorm:"-"`
}

func (d *Dept) TableName() string {
	return "dept"
}

func (d *Dept) IsValid() bool {
	return d.Model.IsValid() && d.DeptName != ""
}

func (d *Dept) GetParentId() int64 {
	return d.ParentId
}

func (d *Dept) AddChild(child common.IChild) {
	d.Children = append(d.Children, child.(*Dept))
}

func (d *Dept) Summary() map[string]interface{} {
	ret := map[string]interface{}{}
	ret["id"] = d.ID
	ret["label"] = d.DeptName
	if len(d.Children) > 0 {
		var child []interface{}
		for _, c := range d.Children {
			child = append(child, c.Summary())
		}
		ret["children"] = child
	}
	return ret
}

func (d *Dept) BeforeDelete(tx *gorm.DB) (err error) {
	if d.ID == common.SuperDeptId {
		return fmt.Errorf("该部门无法删除")
	}
	childCount := int64(0)
	tx.Model(&Dept{}).Where(&Dept{ParentId: d.ID}).Count(&childCount)
	if childCount > 0 {
		return fmt.Errorf("当前部门还有子部门，请处理子部门")
	}

	userCount := int64(0)
	tx.Model(&User{}).Where(&User{DeptId: d.ID}).Count(&userCount)
	if userCount > 0 {
		return fmt.Errorf("当前部门还有用户，请先将用户移到其他部门")
	}
	return nil
}
