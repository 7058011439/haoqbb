package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqOperateLog struct {
	UserName string `form:"userName" json:"userName"`
	FullPath string `form:"fullPath" json:"fullPath"`
	Method   string `form:"method" json:"method"`
	commonHttp.QueryParam
}

func (q *QueryReqOperateLog) Condition() interface{} {
	model := &admin.OperateLog{}
	commonHttp.ReqToDB(q, model)
	return model
}

func (q *QueryReqOperateLog) Order() string {
	return "created_at desc"
}
