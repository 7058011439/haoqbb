package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqDictData struct {
	DictType string `form:"dictType" json:"dictType"`
	commonHttp.QueryParam
}

func (q *QueryReqDictData) Condition() interface{} {
	model := &admin.DictData{}
	commonHttp.ReqToDB(q, model)
	return model
}

func (q *QueryReqDictData) Order() string {
	return "sort asc"
}

type UpdateReqDictData struct {
	commonHttp.UpdateId
	InsertReqDictData
}

func (u *UpdateReqDictData) Model() common.IUpdateData {
	model := &admin.DictData{}
	commonHttp.ReqToDB(u, model)
	return model
}

type InsertReqDictData struct {
	DictLabel string `json:"dictLabel" form:"dictLabel"`
	DictType  string `json:"dictType" form:"dictType"`
	DictValue string `json:"dictValue" form:"dictValue"`
	Sort      int    `json:"sort" form:"sort"`
	Status    string `json:"status" form:"status"`
	Remark    string `json:"remark" form:"remark"`
}

func (i *InsertReqDictData) Model() common.IUpdateData {
	model := &admin.DictData{}
	commonHttp.ReqToDB(i, model)
	return model
}

type DeleteReqDictData struct {
	commonHttp.DeleteId
}

func (d *DeleteReqDictData) Model() common.IDataDB {
	return &admin.DictData{}
}
