package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
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

func (a *apiConfig) list(c *gin.Context) {
	getList(c, &admin.Config{}, &dto.QueryReqConfig{})
}

func (a *apiConfig) info(c *gin.Context) {
	getItemById(c, &admin.Config{}, nil)
}

func (a *apiConfig) updateInfo(c *gin.Context) {
	updateItem(c, &dto.UpdateReqConfig{}, nil)
}

func (a *apiConfig) addInfo(c *gin.Context) {
	addItem(c, &dto.InsertReqConfig{})
}

func (a *apiConfig) delInfo(c *gin.Context) {
	deleteItem(c, &dto.DeleteReqConfig{}, nil)
}

func (a *apiConfig) infoByKey(c *gin.Context) {
	ret := Http.NewResult(c)
	data := admin.GetConfigByKey(c.Param("key"))
	ret.Success(common.ResponseSuccess, map[string]interface{}{
		"configKey":   data.ConfigKey,
		"configValue": data.ConfigValue,
	})
}
