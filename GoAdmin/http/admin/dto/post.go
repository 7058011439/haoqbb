package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqPost struct {
	PostCode string `form:"postCode" json:"postCode"`
	PostName string `form:"postName" json:"postName"`
	Status   string `form:"status" json:"status"`
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
	PostCode string `json:"postCode" form:"postCode"`
	PostName string `json:"postName" form:"postName"`
	Sort     int    `json:"sort" form:"sort"`
	Status   string `json:"status" form:"status"`
	Remark   string `json:"remark" form:"remark"`
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
