package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	_ "github.com/7058011439/haoqbb/GoAdmin/util"
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

// @Tags     好奇宝宝后台-职位管理
// @Summary  获取职位列表
// @Param    token  header    string            true  "token"
// @Param    data   query     dto.QueryReqPost  true  "查询条件"
// @Success  200    {object}  Http.WebResult{data=util.WebResultCommonList{list=[]admin.Post}}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/post [get]
func (a *apiPost) list(c *gin.Context) {
	getList(c, &admin.Post{}, &dto.QueryReqPost{})
}

// @Tags     好奇宝宝后台-职位管理
// @Summary  获取职位详情
// @Param    token  header    string  true  "token"
// @Param    id     path      int     true  "职位id"
// @Success  200    {object}  Http.WebResult{data=admin.Post}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/post/{id} [get]
func (a *apiPost) info(c *gin.Context) {
	getItemById(c, &admin.Post{}, nil)
}

// @Tags     好奇宝宝后台-职位管理
// @Summary  修改职位
// @Param    token  header    string             true  "token"
// @Param    data   formData  dto.UpdateReqPost  true  "职位信息"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/post [put]
func (a *apiPost) updateInfo(c *gin.Context) {
	updateItem(c, &dto.UpdateReqPost{}, nil)
}

// @Tags     好奇宝宝后台-职位管理
// @Summary  新增职位
// @Param    token  header    string             true  "token"
// @Param    data   formData  dto.InsertReqPost  true  "职位信息"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/post [post]
func (a *apiPost) addInfo(c *gin.Context) {
	addItem(c, &dto.InsertReqPost{})
}

// @Tags     好奇宝宝后台-职位管理
// @Summary  删除职位
// @Param    token  header    string             true  "token"
// @Param    data   body      dto.DeleteReqPost  true  "职位id"
// @Success  200    {object}  Http.WebResult
// @Failure  500    {object}  Http.WebResult
// @Router   /api/post [delete]
func (a *apiPost) delInfo(c *gin.Context) {
	deleteItem(c, &dto.DeleteReqPost{}, nil)
}
