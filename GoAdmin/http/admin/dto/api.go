package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqApi struct {
	InsertReqApi
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
	Title  string `json:"title" form:"title"`   // 标题
	Path   string `json:"path" form:"path"`     // api路径
	Type   string `json:"type" form:"type"`     // 类型(BUS, SYS)
	Method string `json:"method" form:"method"` // 请求方式(GET, PUT, DELETE, POST)
	Log    string `json:"log" form:"log"`       // 是否记录日志(Y-记录日志;N-不记录日志)
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
