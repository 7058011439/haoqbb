package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	"github.com/7058011439/haoqbb/Http"
	"github.com/gin-gonic/gin"
)

type apiLoginLog struct {
	Http.Api
}

func init() {
	a := common.ServerAdmin.RegeditGroup("/api/sys_login_log", &apiLoginLog{}, common.CheckAdminToken).(*apiLoginLog)
	a.RegeditApi(Http.TypeGet, "", a.list)
}

func (a *apiLoginLog) list(c *gin.Context) {
	getList(c, &admin.LoginLog{}, &dto.QueryReqLoginLog{})
}
