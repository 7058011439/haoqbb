package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	commondb "github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	_ "github.com/7058011439/haoqbb/GoAdmin/util"
	"github.com/7058011439/haoqbb/Http"
	"github.com/gin-gonic/gin"
)

type apiApi struct {
	Http.Api
}

func init() {
	a := common.ServerAdmin.RegeditGroup("/api/sys-api", &apiApi{}, common.CheckAdminToken).(*apiApi)
	a.RegeditApi(Http.TypeGet, "", a.list)
	a.RegeditApi(Http.TypeGet, "/:id", a.info)
	a.RegeditApi(Http.TypePut, "", a.updateInfo)
	a.RegeditApi(Http.TypePost, "", a.addInfo)
	a.RegeditApi(Http.TypeDelete, "", a.delInfo)
}

// @Tags     好奇宝宝后台-接口管理
// @Summary  获取接口列表
// @Param    token  header    string           true  "token"
// @Param    data   query     dto.QueryReqApi  true  "查询条件"
// @Success  200    {object}  Http.WebResult{data=util.WebResultCommonList{list=[]admin.Api}}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/sys-api [get]
func (a *apiApi) list(c *gin.Context) {
	getList(c, &admin.Api{}, &dto.QueryReqApi{})
}

// @Tags     好奇宝宝后台-接口管理
// @Summary  获取接口详情
// @Param    token  header    string  true  "token"
// @Param    id     path      int     true  "接口id"
// @Success  200    {object}  Http.WebResult{data=admin.Api}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/sys-api/{id} [get]
func (a *apiApi) info(c *gin.Context) {
	getItemById(c, &admin.Api{}, nil)
}

// @Tags     好奇宝宝后台-接口管理
// @Summary  修改接口
// @Param    token  header    string            true  "token"
// @Param    data   formData  dto.UpdateReqApi  true  "接口信息"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/sys-api [put]
func (a *apiApi) updateInfo(c *gin.Context) {
	updateItem(c, &dto.UpdateReqApi{}, func(data commondb.IDataDB) {
		// 修改api后，需刷新角色对应权限列表
		common.RefreshPermission()
	})
}

// @Tags     好奇宝宝后台-接口管理
// @Summary  新增接口
// @Param    token  header    string            true  "token"
// @Param    data   formData  dto.InsertReqApi  true  "接口信息"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/sys-api [post]
func (a *apiApi) addInfo(c *gin.Context) {
	addItem(c, &dto.InsertReqApi{})
}

// @Tags     好奇宝宝后台-接口管理
// @Summary  删除接口
// @Param    token  header    string            true  "token"
// @Param    data   body      dto.DeleteReqApi  true  "接口id"
// @Success  200    {object}  Http.WebResult
// @Failure  500    {object}  Http.WebResult
// @Router   /api/sys-api [delete]
func (a *apiApi) delInfo(c *gin.Context) {
	deleteItem(c, &dto.DeleteReqApi{}, func(ids []int64) {
		// 删除api后，需刷新角色对应权限列表
		common.RefreshPermission()
	})
}
