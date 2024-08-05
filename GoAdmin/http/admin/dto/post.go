package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqPost struct {
	PostCode string `form:"postCode" json:"postCode"` // 职位编号
	PostName string `form:"postName" json:"postName"` // 职位名称
	Status   string `form:"status" json:"status"`     // 状态(1-禁止;2-启用)
	commonHttp.QueryParam
}

func (q *QueryReqPost) Condition() interface{} {
	model := &admin.Post{}
	commonHttp.ReqToDB(q, model)
	return model
}

func (q *QueryReqPost) Order() string {
	return "sort asc"
}

type UpdateReqPost struct {
	commonHttp.UpdateId
	InsertReqPost
}

func (u *UpdateReqPost) Model() common.IUpdateData {
	model := &admin.Post{}
	commonHttp.ReqToDB(u, model)
	return model
}

type InsertReqPost struct {
	PostCode string `json:"postCode" form:"postCode"` // 职位编号
	PostName string `json:"postName" form:"postName"` // 职位名称
	Sort     int    `json:"sort" form:"sort"`         // 显示序号
	Status   string `json:"status" form:"status"`     // 状态(1-禁止;2-启用)
	Remark   string `json:"remark" form:"remark"`     // 备注
}

func (i *InsertReqPost) Model() common.IUpdateData {
	model := &admin.Post{}
	commonHttp.ReqToDB(i, model)
	return model
}

type DeleteReqPost struct {
	commonHttp.DeleteId
}

func (d *DeleteReqPost) Model() common.IDataDB {
	return &admin.Post{}
}
