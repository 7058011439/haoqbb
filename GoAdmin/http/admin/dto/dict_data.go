package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqDictData struct {
	DictType string `form:"dictType" json:"dictType"` // 字典类型
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
	DictLabel string `json:"dictLabel" form:"dictLabel"` // 数据标签
	DictType  string `json:"dictType" form:"dictType"`   // 字典类型
	DictValue string `json:"dictValue" form:"dictValue"` // 数据键值
	Sort      int    `json:"sort" form:"sort"`           // 显示顺序
	Status    string `json:"status" form:"status"`       // 状态(1-禁止;2-启用)
	Remark    string `json:"remark" form:"remark"`       // 备注
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
