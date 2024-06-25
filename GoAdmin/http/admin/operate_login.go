package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	"github.com/7058011439/haoqbb/Http"
	"github.com/gin-gonic/gin"
)

type apiOperateLog struct {
	Http.Api
}

func init() {
	a := common.ServerAdmin.RegeditGroup("/api/sys_operate_log", &apiOperateLog{}, common.CheckAdminToken).(*apiOperateLog)
	a.RegeditApi(Http.TypeGet, "", a.list)
}

func (a *apiOperateLog) list(c *gin.Context) {
	getList(c, &admin.OperateLog{}, &dto.QueryReqOperateLog{})
}
