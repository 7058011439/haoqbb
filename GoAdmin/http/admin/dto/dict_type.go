package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqDictType struct {
	DictName string `form:"dictName" json:"dictName"`
	DictType string `form:"dictType" json:"dictType"`
	Status   string `form:"status" json:"status"`
	commonHttp.QueryParam
}

func (q *QueryReqDictType) Condition() interface{} {
	model := &admin.DictType{}
	commonHttp.ReqToDB(q, model)
	return model
}

type UpdateReqDictType struct {
	commonHttp.UpdateId
	InsertReqDictType
}

func (u *UpdateReqDictType) Model() common.IUpdateData {
	model := &admin.DictType{}
	commonHttp.ReqToDB(u, model)
	return model
}

type InsertReqDictType struct {
	DictName string `json:"dictName" form:"dictName"`
	DictType string `json:"dictType" form:"dictType"`
	Status   string `json:"status" form:"status"`
	Remark   string `json:"remark" form:"remark"`
}

func (i *InsertReqDictType) Model() common.IUpdateData {
	model := &admin.DictType{}
	commonHttp.ReqToDB(i, model)
	return model
}

type DeleteReqDictType struct {
	commonHttp.DeleteId
}

func (d *DeleteReqDictType) Model() common.IDataDB {
	return &admin.DictType{}
}
