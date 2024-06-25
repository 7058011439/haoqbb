package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqDept struct {
	DeptName string `form:"deptName" json:"deptName"`
	Status   string `form:"status" json:"status"`
	commonHttp.QueryParam
}

func (q *QueryReqDept) Condition() interface{} {
	model := &admin.Dept{}
	commonHttp.ReqToDB(q, model)
	return model
}

func (q *QueryReqDept) Order() string {
	return "sort asc"
}

type UpdateReqDept struct {
	commonHttp.UpdateId
	InsertReqDept
}

func (u *UpdateReqDept) Model() common.IUpdateData {
	model := &admin.Dept{}
	commonHttp.ReqToDB(u, model)
	return model
}

type InsertReqDept struct {
	ParentId int64  `json:"parentId"` //上级部门
	DeptName string `json:"deptName"` //部门名称
	Sort     int    `json:"sort"`     //排序
	Leader   string `json:"leader"`   //负责人
	Phone    string `json:"phone"`    //手机
	Email    string `json:"email"`    //邮箱
	Status   string `json:"status"`   //状态
}

func (u *InsertReqDept) Model() common.IUpdateData {
	model := &admin.Dept{}
	commonHttp.ReqToDB(u, model)
	return model
}

type DeleteReqDept struct {
	commonHttp.DeleteId
}

func (d *DeleteReqDept) Model() common.IDataDB {
	return &admin.Dept{}
}
