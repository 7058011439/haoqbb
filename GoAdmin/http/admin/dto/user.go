package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqUser struct {
	UserName string `form:"username" json:"username"` // 用户名
	Phone    string `form:"phone" json:"phone"`       // 电话
	Status   string `form:"status" json:"status"`     // 状态(1-禁止;2-启用)
	DeptId   int64  `form:"deptId" json:"deptId"`     // 部门id
	commonHttp.QueryParam
}

func (q *QueryReqUser) Condition() interface{} {
	return &admin.User{
		UserName: q.UserName,
		Phone:    q.Phone,
		Status:   q.Status,
		DeptId:   q.DeptId,
	}
}

type UpdateReqUser struct {
	commonHttp.UpdateId
	UserName string `json:"username" form:"username"`         // 用户名
	NickName string `json:"nickName" form:"nickName"`         // 用户昵称
	Phone    string `json:"phone" form:"phone"`               // 电话
	RoleId   int    `json:"roleId" form:"roleId"`             // 所属角色id
	Avatar   string `json:"avatar" form:"avatar"`             // 头像
	Sex      string `json:"sex" form:"sex"`                   // 性别
	Email    string `json:"email" form:"email"`               // 邮箱
	DeptId   int    `json:"deptId" form:"deptId"`             // 部门id
	PostId   int    `json:"postId" form:"postId"`             // 职位id
	Remark   string `json:"remark" form:"remark"`             // 备注
	Status   string `json:"status" form:"status" default:"1"` // 状态(1-禁止;2-启用)
}

func (u *UpdateReqUser) Model() common.IUpdateData {
	model := &admin.User{}
	commonHttp.ReqToDB(u, model)
	return model
}

type UpdateReqUserStatus struct {
	commonHttp.UpdateId
	Status string `json:"status" form:"status" default:"1"` // 状态(1-禁止;2-启用)
}

func (u *UpdateReqUserStatus) Model() common.IUpdateData {
	return &admin.User{
		Model:  common.Model{ID: u.ID},
		Status: u.Status,
	}
}

type UpdateReqUserPassword struct {
	commonHttp.UpdateId
	PassWord string `json:"password" form:"password"`
}

func (u *UpdateReqUserPassword) Model() common.IUpdateData {
	return &admin.User{
		Model:    common.Model{ID: u.ID},
		PassWord: u.PassWord,
	}
}

type InsertReqUser struct {
	UserName string `json:"username" form:"username"`         // 用户名
	Password string `json:"password" form:"password"`         // 密码
	NickName string `json:"nickName" form:"nickName"`         // 用户昵称
	Phone    string `json:"phone" form:"phone"`               // 电话
	RoleId   int    `json:"roleId" form:"roleId"`             // 所属角色id
	Avatar   string `json:"avatar" form:"avatar"`             // 头像
	Sex      string `json:"sex" form:"sex"`                   // 性别
	Email    string `json:"email" form:"email"`               // 邮箱
	DeptId   int    `json:"deptId" form:"deptId"`             // 部门id
	PostId   int    `json:"postId" form:"postId"`             // 职位id
	Remark   string `json:"remark" form:"remark"`             // 备注
	Status   string `json:"status" form:"status" default:"1"` // 状态(1-禁止;2-启用)
}

func (i *InsertReqUser) Model() common.IUpdateData {
	model := &admin.User{}
	commonHttp.ReqToDB(i, model)
	if model.Avatar == "" {
		model.Avatar = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	}
	model.PassWord = i.Password
	return model
}

type DeleteReqUser struct {
	commonHttp.DeleteId
}

func (d *DeleteReqUser) Model() common.IDataDB {
	return &admin.User{}
}
