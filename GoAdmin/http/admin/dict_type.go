package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
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

func (a *apiDictType) optionSelect(c *gin.Context) {
	ret := Http.NewResult(c)
	_, listData := getDBList(&admin.DictType{}, &common.QueryParam{})
	ret.Success(common.ResponseSuccess, listData)
}

func (a *apiDictType) list(c *gin.Context) {
	getList(c, &admin.DictType{}, &dto.QueryReqDictType{})
}

func (a *apiDictType) info(c *gin.Context) {
	getItemById(c, &admin.DictType{}, nil)
}

func (a *apiDictType) updateInfo(c *gin.Context) {
	updateItem(c, &dto.UpdateReqDictType{}, nil)
}

func (a *apiDictType) addInfo(c *gin.Context) {
	addItem(c, &dto.InsertReqDictType{})
}

func (a *apiDictType) delInfo(c *gin.Context) {
	deleteItem(c, &dto.DeleteReqDictType{}, nil)
}
