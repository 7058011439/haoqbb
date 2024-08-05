package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqConfig struct {
	ConfigName string `form:"configName" json:"configName"` // 配置名称
	ConfigKey  string `form:"configKey" json:"configKey"`   // 配置key
	ConfigType string `form:"configType" json:"configType"` // 配置类型(系统内置(Y,N))
	commonHttp.QueryParam
}

func (q *QueryReqConfig) Condition() interface{} {
	model := &admin.Config{}
	commonHttp.ReqToDB(q, model)
	return model
}

type UpdateReqConfig struct {
	commonHttp.UpdateId
	InsertReqConfig
}

func (u *UpdateReqConfig) Model() common.IUpdateData {
	model := &admin.Config{}
	commonHttp.ReqToDB(u, model)
	return model
}

type InsertReqConfig struct {
	ConfigName  string `json:"configName"`  // 配置名称
	ConfigKey   string `json:"configKey"`   // 配置key
	ConfigValue string `json:"configValue"` // 配置value
	ConfigType  string `json:"configType"`  // 配置类型(系统内置(Y,N))
	IsFrontend  string `json:"isFrontend"`  // 是否前端配置(2-否,1-是)
	Remark      string `json:"remark"`      // 备注
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
