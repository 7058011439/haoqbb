package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	_ "github.com/7058011439/haoqbb/GoAdmin/util"
	"github.com/7058011439/haoqbb/Http"
	"github.com/gin-gonic/gin"
)

type apiConfig struct {
	Http.Api
}

func init() {
	a := common.ServerAdmin.RegeditGroup("/api/config", &apiConfig{}, common.CheckAdminToken).(*apiConfig)
	a.RegeditApi(Http.TypeGet, "", a.list)
	a.RegeditApi(Http.TypeGet, "/:id", a.info)
	a.RegeditApi(Http.TypePut, "", a.updateInfo)
	a.RegeditApi(Http.TypePost, "", a.addInfo)
	a.RegeditApi(Http.TypeDelete, "", a.delInfo)

	a.RegeditApi(Http.TypeGet, "/key/:key", a.infoByKey)
}

// @Tags     好奇宝宝后台-参数管理
// @Summary  获取参数列表
// @Param    token  header    string              true  "token"
// @Param    data   query     dto.QueryReqConfig  true  "查询条件"
// @Success  200    {object}  Http.WebResult{data=util.WebResultCommonList{list=[]admin.Config}}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/config [get]
func (a *apiConfig) list(c *gin.Context) {
	getList(c, &admin.Config{}, &dto.QueryReqConfig{})
}

// @Tags     好奇宝宝后台-参数管理
// @Summary  获取参数详情
// @Param    token  header    string  true  "token"
// @Param    id     path      int     true  "参数id"
// @Success  200    {object}  Http.WebResult{data=admin.Config}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/config/{id} [get]
func (a *apiConfig) info(c *gin.Context) {
	getItemById(c, &admin.Config{}, nil)
}

// @Tags     好奇宝宝后台-参数管理
// @Summary  修改参数
// @Param    token  header    string               true  "token"
// @Param    data   formData  dto.UpdateReqConfig  true  "参数信息"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/config [put]
func (a *apiConfig) updateInfo(c *gin.Context) {
	updateItem(c, &dto.UpdateReqConfig{}, nil)
}

// @Tags     好奇宝宝后台-参数管理
// @Summary  新增参数
// @Param    token  header    string               true  "token"
// @Param    data   formData  dto.InsertReqConfig  true  "参数信息"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/config [post]
func (a *apiConfig) addInfo(c *gin.Context) {
	addItem(c, &dto.InsertReqConfig{})
}

// @Tags     好奇宝宝后台-参数管理
// @Summary  删除参数
// @Param    token  header    string               true  "token"
// @Param    data   body      dto.DeleteReqConfig  true  "参数id"
// @Success  200    {object}  Http.WebResult
// @Failure  500    {object}  Http.WebResult
// @Router   /api/config [delete]
func (a *apiConfig) delInfo(c *gin.Context) {
	deleteItem(c, &dto.DeleteReqConfig{}, nil)
}

// @Tags     好奇宝宝后台-参数管理
// @Summary  获取参数(byKey)详情
// @Param    token  header    string  true  "token"
// @Param    key    path      string  true  "参数key"
// @Success  200    {object}  Http.WebResult{data=admin.Config}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/config/key/{key} [get]
func (a *apiConfig) infoByKey(c *gin.Context) {
	ret := Http.NewResult(c)
	data := admin.GetConfigByKey(c.Param("key"))
	ret.Success(common.ResponseSuccess, map[string]interface{}{
		"configKey":   data.ConfigKey,
		"configValue": data.ConfigValue,
	})
}
