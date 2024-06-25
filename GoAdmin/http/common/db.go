package common

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/7058011439/haoqbb/Http"
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

func GetList[T any](dbManager *common.Manager, c *gin.Context, model T, requestCondition IQueryParam) {
	if Http.Bind(c, requestCondition) {
		ret := Http.NewResult(c)
		count, list := GetDBList(dbManager, model, requestCondition)
		ret.Success(ResponseSuccess,
			map[string]interface{}{
				"count":     count,
				"list":      list,
				"pageIndex": requestCondition.GetPageIndex(),
				"pageSize":  requestCondition.GetPageSize(),
			})
	}
}

func GetItemById(dbManager *common.Manager, c *gin.Context, data common.IDataDB, afterGet func(data common.IDataDB)) {
	ret := Http.NewResult(c)
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		data.SetId(id)
		data = dbManager.LoadData(data)
		if afterGet != nil {
			afterGet(data)
		}
		ret.Success(ResponseSuccess, data)
	} else {
		ret.Fail("参数错误", nil)
	}
}

func UpdateItem(dbManager *common.Manager, c *gin.Context, data IUpdateParam, afterUpdate func(data common.IDataDB)) {
	if Http.Bind(c, data) {
		ret := Http.NewResult(c)
		model := data.Model()
		model.SetUpdateBy(GetAdminId(c))
		if err := dbManager.UpdateData(model); err == nil {
			if afterUpdate != nil {
				afterUpdate(model)
			}
			ret.Success(ResponseSuccess, model.GetId())
		} else {
			ret.Fail(err.Error(), nil)
		}
	}
}

func AddItem(dbManager *common.Manager, c *gin.Context, data IInsertData) {
	if Http.Bind(c, data) {
		ret := Http.NewResult(c)
		model := data.Model()
		model.SetCreateBy(GetAdminId(c))
		if err := dbManager.InsertData(model); err == nil {
			ret.Success(ResponseSuccess, model.GetId())
		} else {
			ret.Fail(err.Error(), nil)
		}
	}
}

func DeleteItem(dbManager *common.Manager, c *gin.Context, data IDeleteData, afterDelete func(ids []int64)) {
	if Http.Bind(c, data) {
		ret := Http.NewResult(c)
		if err := dbManager.DeleteData(data.Model(), data.GetIds()...); err == nil {
			if afterDelete != nil {
				afterDelete(data.GetIds())
			}
			ret.Success(ResponseSuccess, data.GetIds())
		} else {
			ret.Fail(err.Error(), nil)
		}
	}
}

func fillChild(list []common.IChild, parent common.IChild, fun func(item common.IChild) bool) {
	for _, d := range list {
		if d.GetParentId() == parent.GetId() && fun(d) {
			fillChild(list, d, fun)
			parent.AddChild(d)
		}
	}
}

func GetDBList[T any](dbManager *common.Manager, model T, queryParam IQueryParam) (int64, []T) {
	var count int64
	var list []T // 使用 []*T 来声明一个指向指针切片的变量
	pageSize := queryParam.GetPageSize()
	pageIndex := queryParam.GetPageIndex()

	db := dbManager.MysqlDB().Model(model)
	if queryParam.Preload() != "" {
		db = db.Preload(queryParam.Preload())
	}
	db = db.Where(queryParam.Condition())
	db = db.Where(queryParam.ConditionTime())
	db.Count(&count)
	if queryParam.Order() != "" {
		db = db.Order(queryParam.Order())
	}
	if pageSize > 0 || pageIndex > 0 {
		db = db.Limit(pageSize).Offset((pageIndex - 1) * pageSize)
	}
	db.Find(&list)

	return count, list
}
