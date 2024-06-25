package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Role struct {
	common.Model
	RoleName         string  `json:"roleName" gorm:"unique;"`                 // 角色名称
	RoleKey          string  `json:"roleKey" gorm:"unique;"`                  // 角色代码
	Status           string  `json:"status" gorm:"default:'2';size:4;"`       // 角色状态
	Sort             int     `json:"sort"`                                    // 角色排序
	Remark           string  `json:"remark"`                                  // 备注
	SysMenu          []*Menu `json:"sysMenu" gorm:"many2many:role_menu_rule"` // 角色对应菜单权限
	common.ControlBy `json:"-"`
	MenuIds          []int64 `json:"menuIds" gorm:"-"`
}

func (r *Role) TableName() string {
	return "role"
}

func (r *Role) BeforeDelete(tx *gorm.DB) (err error) {
	if r.ID == common.AdminRoleId {
		return fmt.Errorf("不能删除[超级管理员]角色")
	}
	userCount := int64(0)
	tx.Model(&User{}).Where(&User{RoleId: r.ID}).Count(&userCount)
	if userCount > 0 {
		return fmt.Errorf("当前角色还有用户，请先将用户移到其他角色")
	}
	return tx.Model(r).Association("SysMenu").Clear().Error
}
