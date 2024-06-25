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

func (a *apiDictData) list(c *gin.Context) {
	getList(c, &admin.DictData{}, &dto.QueryReqDictData{})
}

func (a *apiDictData) info(c *gin.Context) {
	getItemById(c, &admin.DictData{}, nil)
}

func (a *apiDictData) updateInfo(c *gin.Context) {
	updateItem(c, &dto.UpdateReqDictData{}, nil)
}

func (a *apiDictData) addInfo(c *gin.Context) {
	addItem(c, &dto.InsertReqDictData{})
}

func (a *apiDictData) delInfo(c *gin.Context) {
	deleteItem(c, &dto.DeleteReqDictData{}, nil)
}
