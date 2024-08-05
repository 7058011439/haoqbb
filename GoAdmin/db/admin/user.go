package admin

import (
	"fmt"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/jinzhu/gorm"
)

type User struct {
	common.Model
	UserName string  `json:"username" gorm:"unique;Index:idx_username"` // 管理员账号
	PassWord string  `json:"-"`                                         // 管理员密码
	NickName string  `json:"nickName"`                                  // 昵称
	Phone    string  `json:"phone" gorm:"size:11;comment:'手机号'"`        // 电话
	RoleId   int64   `json:"roleId" gorm:"comment:'角色ID'"`              // 角色id
	Avatar   string  `json:"avatar" gorm:"comment:'头像'"`                // 头像
	Sex      string  `json:"sex" gorm:"comment:'性别'"`                   // 性别
	Email    string  `json:"email" gorm:"size:128;comment:'邮箱'"`        // 邮箱
	DeptId   int64   `json:"deptId" gorm:"comment:'部门'"`                // 部门
	PostId   int64   `json:"postId" gorm:"comment:'岗位'"`                // 岗位
	Remark   string  `json:"remark" gorm:"comment:'备注'"`                // 备注
	Status   string  `json:"status" gorm:"size:4;comment:'状态'"`         // 状态(1-禁止;2-启用)
	DeptIds  []int64 `json:"deptIds" gorm:"-"`
	PostIds  []int64 `json:"postIds" gorm:"-"`
	RoleIds  []int64 `json:"roleIds" gorm:"-"`
	Dept     *Dept   `json:"dept" gorm:"-"`
	common.ControlBy
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) IsValid() bool {
	return u.Model.IsValid() && u.UserName != ""
}

func (u *User) AfterFind(_ *gorm.DB) error {
	u.DeptIds = []int64{u.DeptId}
	u.PostIds = []int64{u.PostId}
	u.RoleIds = []int64{u.RoleId}
	u.Dept = LoadForumData(&Dept{Model: common.Model{ID: u.DeptId}}).(*Dept)
	return nil
}

// 使用 AfterSave 钩子函数
func (u *User) AfterSave(tx *gorm.DB) (err error) {
	if u.Dept == nil {
		u.DeptIds = []int64{u.DeptId}
		u.PostIds = []int64{u.PostId}
		u.RoleIds = []int64{u.RoleId}
		u.Dept = LoadForumData(&Dept{Model: common.Model{ID: u.DeptId}}).(*Dept)
	}
	return
}

func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	if u.ID == common.AdminUserId {
		return fmt.Errorf("不能删除[超级管理员]")
	}
	return nil
}
