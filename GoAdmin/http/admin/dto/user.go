package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqUser struct {
	UserName string `form:"username" json:"username"`
	Phone    string `form:"phone" json:"phone"`
	Status   string `form:"status" json:"status"`
	DeptId   int64  `form:"deptId" json:"deptId"`
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
	UserName string `json:"username" form:"username"`
	NickName string `json:"nickName" form:"nickName"`
	Phone    string `json:"phone" form:"phone"`
	RoleId   int    `json:"roleId" form:"roleId"`
	Avatar   string `json:"avatar" form:"avatar"`
	Sex      string `json:"sex" form:"sex"`
	Email    string `json:"email" form:"email"`
	DeptId   int    `json:"deptId" form:"deptId"`
	PostId   int    `json:"postId" form:"postId"`
	Remark   string `json:"remark" form:"remark"`
	Status   string `json:"status" form:"status" default:"1"`
}

func (u *UpdateReqUser) Model() common.IUpdateData {
	model := &admin.User{}
	commonHttp.ReqToDB(u, model)
	return model
}

type UpdateReqUserStatus struct {
	commonHttp.UpdateId
	Status string `json:"status" form:"status" default:"1"`
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
	UserName string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	NickName string `json:"nickName" form:"nickName"`
	Phone    string `json:"phone" form:"phone"`
	RoleId   int    `json:"roleId" form:"roleId"`
	Avatar   string `json:"avatar" form:"avatar"`
	Sex      string `json:"sex" form:"sex"`
	Email    string `json:"email" form:"email"`
	DeptId   int    `json:"deptId" form:"deptId"`
	PostId   int    `json:"postId" form:"postId"`
	Remark   string `json:"remark" form:"remark"`
	Status   string `json:"status" form:"status" default:"1"`
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
