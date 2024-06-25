package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	commondb "github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
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

func (a *apiApi) list(c *gin.Context) {
	getList(c, &admin.Api{}, &dto.QueryReqApi{})
}

func (a *apiApi) info(c *gin.Context) {
	getItemById(c, &admin.Api{}, nil)
}

func (a *apiApi) updateInfo(c *gin.Context) {
	updateItem(c, &dto.UpdateReqApi{}, func(data commondb.IDataDB) {
		// 修改api后，需刷新角色对应权限列表
		common.RefreshPermission()
	})
}

func (a *apiApi) addInfo(c *gin.Context) {
	addItem(c, &dto.InsertReqApi{})
}

func (a *apiApi) delInfo(c *gin.Context) {
	deleteItem(c, &dto.DeleteReqApi{}, func(ids []int64) {
		// 删除api后，需刷新角色对应权限列表
		common.RefreshPermission()
	})
}
