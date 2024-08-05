package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqLoginLog struct {
	UserName string `form:"userName" json:"userName"` // 账号
	Ip       string `form:"ip" json:"ip"`             // 登录ip
	Status   string `form:"status" json:"status"`     // 登录结果
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
