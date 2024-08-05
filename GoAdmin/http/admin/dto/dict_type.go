package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqDictType struct {
	DictName string `form:"dictName" json:"dictName"` // 字典名称
	DictType string `form:"dictType" json:"dictType"` // 字典类型
	Status   string `form:"status" json:"status"`     // 状态(1-禁止;2-启用)
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
	DictName string `json:"dictName" form:"dictName"` // 字典名称
	DictType string `json:"dictType" form:"dictType"` // 字典类型
	Status   string `json:"status" form:"status"`     // 状态(1-禁止;2-启用)
	Remark   string `json:"remark" form:"remark"`     // 备注
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
