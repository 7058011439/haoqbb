package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	_ "github.com/7058011439/haoqbb/GoAdmin/util"
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

// @Tags     好奇宝宝后台-日志管理
// @Summary  查询登录日志
// @Param    token  header    string                true  "token"
// @Param    data   query     dto.QueryReqLoginLog  true  "查询条件"
// @Success  200    {object}  Http.WebResult{data=util.WebResultCommonList{list=[]admin.LoginLog}}
// @Router   /api/sys_login_log [get]
func (a *apiLoginLog) list(c *gin.Context) {
	getList(c, &admin.LoginLog{}, &dto.QueryReqLoginLog{})
}
