package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	"github.com/7058011439/haoqbb/Http"
	"github.com/gin-gonic/gin"
)

type apiPost struct {
	Http.Api
}

func init() {
	a := common.ServerAdmin.RegeditGroup("/api/post", &apiPost{}, common.CheckAdminToken).(*apiPost)
	a.RegeditApi(Http.TypeGet, "", a.list)
	a.RegeditApi(Http.TypeGet, "/:id", a.info)
	a.RegeditApi(Http.TypePut, "", a.updateInfo)
	a.RegeditApi(Http.TypePost, "", a.addInfo)
	a.RegeditApi(Http.TypeDelete, "", a.delInfo)
}

func (a *apiPost) list(c *gin.Context) {
	getList(c, &admin.Post{}, &dto.QueryReqPost{})
}

func (a *apiPost) info(c *gin.Context) {
	getItemById(c, &admin.Post{}, nil)
}

func (a *apiPost) updateInfo(c *gin.Context) {
	updateItem(c, &dto.UpdateReqPost{}, nil)
}

func (a *apiPost) addInfo(c *gin.Context) {
	addItem(c, &dto.InsertReqPost{})
}

func (a *apiPost) delInfo(c *gin.Context) {
	deleteItem(c, &dto.DeleteReqPost{}, nil)
}
