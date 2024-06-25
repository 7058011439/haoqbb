package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqApi struct {
	Type   string `form:"type" json:"type"`
	Title  string `form:"title" json:"title"`
	Path   string `form:"path" json:"path"`
	Method string `form:"method" json:"method"`
	Log    string `form:"log" json:"log"`
	commonHttp.QueryParam
}

func (q *QueryReqApi) Condition() interface{} {
	model := &admin.Api{}
	commonHttp.ReqToDB(q, model)
	return model
}

type UpdateReqApi struct {
	commonHttp.UpdateId
	InsertReqApi
}

func (u *UpdateReqApi) Model() common.IUpdateData {
	model := &admin.Api{}
	commonHttp.ReqToDB(u, model)
	return model
}

type InsertReqApi struct {
	Title  string `json:"title" form:"title"`
	Path   string `json:"path" form:"path"`
	Type   string `json:"type" form:"type"`
	Method string `json:"method" form:"method"`
	Log    string `json:"log" form:"log"`
}

func (i *InsertReqApi) Model() common.IUpdateData {
	model := &admin.Api{}
	commonHttp.ReqToDB(i, model)
	return model
}

type DeleteReqApi struct {
	commonHttp.DeleteId
}

func (d *DeleteReqApi) Model() common.IDataDB {
	return &admin.Api{}
}
