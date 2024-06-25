package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	commondb "github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	"github.com/gin-gonic/gin"
	"strconv"
)

func getReqId(c *gin.Context) int64 {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		return id
	} else {
		return 0
	}
}

func getList[T any](c *gin.Context, model T, requestCondition common.IQueryParam) {
	common.GetList(admin.DBManager(), c, model, requestCondition)
}

func getItemById(c *gin.Context, data commondb.IDataDB, afterGet func(data commondb.IDataDB)) {
	common.GetItemById(admin.DBManager(), c, data, afterGet)
}

func updateItem(c *gin.Context, data common.IUpdateParam, afterUpdate func(data commondb.IDataDB)) {
	common.UpdateItem(admin.DBManager(), c, data, afterUpdate)
}

func addItem(c *gin.Context, data common.IInsertData) {
	common.AddItem(admin.DBManager(), c, data)
}

func deleteItem(c *gin.Context, data common.IDeleteData, afterDelete func(ids []int64)) {
	common.DeleteItem(admin.DBManager(), c, data, afterDelete)
}

func fillChild(list []commondb.IChild, parent commondb.IChild, fun func(item commondb.IChild) bool) {
	for _, d := range list {
		if d.GetParentId() == parent.GetId() && fun(d) {
			fillChild(list, d, fun)
			parent.AddChild(d)
		}
	}
}

func getDBList[T any](model T, queryParam common.IQueryParam) (int64, []T) {
	return common.GetDBList(admin.DBManager(), model, queryParam)
}
