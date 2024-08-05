package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	"github.com/7058011439/haoqbb/Http"
	"github.com/gin-gonic/gin"
)

type apiDictData struct {
	Http.Api
}

func init() {
	a := common.ServerAdmin.RegeditGroup("/api/dict-data", &apiDictData{}, common.CheckAdminToken).(*apiDictData)
	a.RegeditApi(Http.TypeGet, "option-select", a.optionSelect)
	a.RegeditApi(Http.TypeGet, "", a.list)
	a.RegeditApi(Http.TypeGet, "/:id", a.info)
	a.RegeditApi(Http.TypePut, "", a.updateInfo)
	a.RegeditApi(Http.TypePost, "", a.addInfo)
	a.RegeditApi(Http.TypeDelete, "", a.delInfo)
}

// @Tags     好奇宝宝后台-字典数据管理
// @Summary  获取字典数据(下拉)列表
// @Param    token  header    string                true  "token"
// @Param    data   query     dto.QueryReqDictData  true  "字典类型"
// @Success  200    {object}  Http.WebResult
// @Router   /api/dict-data/option-select [get]
func (a *apiDictData) optionSelect(c *gin.Context) {
	ret := Http.NewResult(c)
	var requestData dto.QueryReqDictData
	if Http.Bind(c, &requestData) {
		_, listData := getDBList(&admin.DictData{}, &requestData)
		var retData []interface{}
		for _, data := range listData {
			retData = append(retData, map[string]interface{}{
				"label": data.DictLabel,
				"value": data.DictValue,
			})
		}
		ret.Success(common.ResponseSuccess, retData)
	}
}

// @Tags     好奇宝宝后台-字典数据管理
// @Summary  获取字典数据列表
// @Param    token  header    string                true  "token"
// @Param    data   query     dto.QueryReqDictData  true  "查询条件"
// @Success  200    {object}  Http.WebResult{data=[]admin.DictData}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/dict-data [get]
func (a *apiDictData) list(c *gin.Context) {
	getList(c, &admin.DictData{}, &dto.QueryReqDictData{})
}

// @Tags     好奇宝宝后台-字典数据管理
// @Summary  获取字典数据详情
// @Param    token  header    string  true  "token"
// @Param    id     path      int     true  "字典数据id"
// @Success  200    {object}  Http.WebResult{data=admin.DictData}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/dict-data/{id} [get]
func (a *apiDictData) info(c *gin.Context) {
	getItemById(c, &admin.DictData{}, nil)
}

// @Tags     好奇宝宝后台-字典数据管理
// @Summary  修改字典数据
// @Param    token  header    string                 true  "token"
// @Param    data   formData  dto.UpdateReqDictData  true  "字典数据信息"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/dict-data [put]
func (a *apiDictData) updateInfo(c *gin.Context) {
	updateItem(c, &dto.UpdateReqDictData{}, nil)
}

// @Tags     好奇宝宝后台-字典数据管理
// @Summary  新增字典数据
// @Param    token  header    string                 true  "token"
// @Param    data   formData  dto.InsertReqDictData  true  "字典数据信息"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/dict-data [post]
func (a *apiDictData) addInfo(c *gin.Context) {
	addItem(c, &dto.InsertReqDictData{})
}

// @Tags     好奇宝宝后台-字典数据管理
// @Summary  删除字典数据
// @Param    token  header    string                 true  "token"
// @Param    data   body      dto.DeleteReqDictData  true  "字典数据id"
// @Success  200    {object}  Http.WebResult
// @Failure  500    {object}  Http.WebResult
// @Router   /api/dict-data [delete]
func (a *apiDictData) delInfo(c *gin.Context) {
	deleteItem(c, &dto.DeleteReqDictData{}, nil)
}
