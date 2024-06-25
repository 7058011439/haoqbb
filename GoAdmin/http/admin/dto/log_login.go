package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqLoginLog struct {
	UserName string `form:"userName" json:"userName"`
	Ip       string `form:"ip" json:"ip"`
	Status   string `form:"status" json:"status"`
	commonHttp.QueryParam
}

func (q *QueryReqLoginLog) Condition() interface{} {
	model := &admin.LoginLog{}
	commonHttp.ReqToDB(q, model)
	return model
}

func (q *QueryReqLoginLog) Order() string {
	return "created_at desc"
}
