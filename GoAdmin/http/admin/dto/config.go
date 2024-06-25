package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqConfig struct {
	ConfigName string `form:"configName" json:"configName"`
	ConfigKey  string `form:"configKey" json:"configKey"`
	ConfigType string `form:"configType" json:"configType"`
	commonHttp.QueryParam
}

func (q *QueryReqConfig) Condition() interface{} {
	model := &admin.Config{}
	commonHttp.ReqToDB(q, model)
	return model
}

type UpdateReqConfig struct {
	commonHttp.UpdateId
	ConfigName  string `json:"configName"`
	ConfigKey   string `json:"configKey"`
	ConfigValue string `json:"configValue"`
	ConfigType  string `json:"configType"`
	IsFrontend  string `json:"isFrontend"`
	Remark      string `json:"remark"`
}

func (u *UpdateReqConfig) Model() common.IUpdateData {
	model := &admin.Config{}
	commonHttp.ReqToDB(u, model)
	return model
}

type InsertReqConfig struct {
	ConfigName  string `json:"configName"`
	ConfigKey   string `json:"configKey"`
	ConfigValue string `json:"configValue"`
	ConfigType  string `json:"configType"`
	IsFrontend  string `json:"isFrontend"`
	Remark      string `json:"remark"`
}

func (i *InsertReqConfig) Model() common.IUpdateData {
	model := &admin.Config{}
	commonHttp.ReqToDB(i, model)
	return model
}

type DeleteReqConfig struct {
	commonHttp.DeleteId
}

func (d *DeleteReqConfig) Model() common.IDataDB {
	return &admin.Config{}
}
