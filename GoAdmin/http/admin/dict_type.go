package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	_ "github.com/7058011439/haoqbb/GoAdmin/util"
	"github.com/7058011439/haoqbb/Http"
	"github.com/gin-gonic/gin"
)

type apiDictType struct {
	Http.Api
}

func init() {
	a := common.ServerAdmin.RegeditGroup("/api/dict-type", &apiDictType{}, common.CheckAdminToken).(*apiDictType)
	a.RegeditApi(Http.TypeGet, "option-select", a.optionSelect)
	a.RegeditApi(Http.TypeGet, "", a.list)
	a.RegeditApi(Http.TypeGet, "/:id", a.info)
	a.RegeditApi(Http.TypePut, "", a.updateInfo)
	a.RegeditApi(Http.TypePost, "", a.addInfo)
	a.RegeditApi(Http.TypeDelete, "", a.delInfo)
}

// @Tags     好奇宝宝后台-字典类型管理
// @Summary  获取字典类型(下拉)列表
// @Param    token  header    string  true  "token"
// @Success  200    {object}  Http.WebResult{data=[]admin.DictType}
// @Router   /api/dict-type/option-select [get]
func (a *apiDictType) optionSelect(c *gin.Context) {
	ret := Http.NewResult(c)
	_, listData := getDBList(&admin.DictType{}, &common.QueryParam{})
	ret.Success(common.ResponseSuccess, listData)
}

// @Tags     好奇宝宝后台-字典类型管理
// @Summary  获取字典类型列表
// @Param    token  header    string                true  "token"
// @Param    data   query     dto.QueryReqDictType  true  "查询条件"
// @Success  200    {object}  Http.WebResult{data=util.WebResultCommonList{list=[]admin.DictType}}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/dict-type [get]
func (a *apiDictType) list(c *gin.Context) {
	getList(c, &admin.DictType{}, &dto.QueryReqDictType{})
}

// @Tags     好奇宝宝后台-字典类型管理
// @Summary  获取字典类型详情
// @Param    token  header    string  true  "token"
// @Param    id     path      int     true  "字典类型id"
// @Success  200    {object}  Http.WebResult{data=admin.DictType}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/dict-type/{id} [get]
func (a *apiDictType) info(c *gin.Context) {
	getItemById(c, &admin.DictType{}, nil)
}

// @Tags     好奇宝宝后台-字典类型管理
// @Summary  修改字典类型
// @Param    token  header    string                 true  "token"
// @Param    data   formData  dto.UpdateReqDictType  true  "字典类型信息"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/dict-type [put]
func (a *apiDictType) updateInfo(c *gin.Context) {
	updateItem(c, &dto.UpdateReqDictType{}, nil)
}

// @Tags     好奇宝宝后台-字典类型管理
// @Summary  新增字典类型
// @Param    token  header    string                 true  "token"
// @Param    data   formData  dto.InsertReqDictType  true  "字典类型信息"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/dict-type [post]
func (a *apiDictType) addInfo(c *gin.Context) {
	addItem(c, &dto.InsertReqDictType{})
}

// @Tags     好奇宝宝后台-字典类型管理
// @Summary  删除字典类型
// @Param    token  header    string                 true  "token"
// @Param    data   body      dto.DeleteReqDictType  true  "字典类型id"
// @Success  200    {object}  Http.WebResult
// @Failure  500    {object}  Http.WebResult
// @Router   /api/dict-type [delete]
func (a *apiDictType) delInfo(c *gin.Context) {
	deleteItem(c, &dto.DeleteReqDictType{}, nil)
}
