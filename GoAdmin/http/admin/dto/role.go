package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqRole struct {
	RoleName string `form:"roleName" json:"roleName"`
	RoleKey  string `form:"roleKey" json:"roleKey"`
	Status   string `form:"status" json:"status"`
	commonHttp.QueryParam
}

func (q *QueryReqRole) Condition() interface{} {
	model := &admin.Role{}
	commonHttp.ReqToDB(q, model)
	return model
}

func (q *QueryReqRole) Order() string {
	return "sort asc"
}

type UpdateReqRole struct {
	commonHttp.UpdateId
	InsertReqRole
}

func (u *UpdateReqRole) Model() common.IUpdateData {
	model := &admin.Role{}
	commonHttp.ReqToDB(u, model)
	return model
}

type UpdateReqRoleStatus struct {
	commonHttp.UpdateId
	Status string `json:"status" form:"status" default:"1"`
}

func (u *UpdateReqRoleStatus) Model() common.IUpdateData {
	return &admin.Role{
		Model:  common.Model{ID: u.ID},
		Status: u.Status,
	}
}

type InsertReqRole struct {
	RoleName string  `form:"roleName" json:"roleName"`
	RoleKey  string  `form:"roleKey" json:"roleKey"`
	Status   string  `form:"status" json:"status"`
	Sort     int     `form:"sort" json:"sort"`
	Remark   string  `form:"remark" json:"remark"`
	MenuIds  []int64 `form:"menuIds" json:"menuIds"`
}

func (i *InsertReqRole) Model() common.IUpdateData {
	model := &admin.Role{}
	commonHttp.ReqToDB(i, model)
	admin.MysqlDB().Find(&model.SysMenu, i.MenuIds)
	return model
}

type DeleteReqRole struct {
	commonHttp.DeleteId
}

func (d *DeleteReqRole) Model() common.IDataDB {
	return &admin.Role{}
}
