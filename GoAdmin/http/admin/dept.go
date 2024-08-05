package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	commondb "github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	"github.com/7058011439/haoqbb/Http"
	"github.com/gin-gonic/gin"
)

type apiDept struct {
	Http.Api
}

func init() {
	a := common.ServerAdmin.RegeditGroup("/api/dept", &apiDept{}, common.CheckAdminToken).(*apiDept)
	a.RegeditApi(Http.TypeGet, "", a.list)
	a.RegeditApi(Http.TypeGet, "/:id", a.info)
	a.RegeditApi(Http.TypePut, "", a.updateInfo)
	a.RegeditApi(Http.TypePost, "", a.addInfo)
	a.RegeditApi(Http.TypeDelete, "", a.delInfo)

	a.RegeditApi(Http.TypeGet, "/deptTree", a.deptTree)
}

func (a *apiDept) getRootDept(listData []*admin.Dept, fun func(dept *admin.Dept) bool) (rootDept []*admin.Dept) {
	var listInterface []commondb.IChild
	for _, data := range listData {
		listInterface = append(listInterface, data)
	}
	for _, data := range listData {
		if data.ParentId == 0 {
			fillChild(listInterface, data, func(item commondb.IChild) bool {
				return fun(item.(*admin.Dept))
			})
			rootDept = append(rootDept, data)
		}
	}
	return
}

// @Tags     好奇宝宝后台-部门管理
// @Summary  获取部门列表
// @Param    token  header    string  true  "token"
// @Param    data   query     dto.QueryReqDept  true  "查询条件"
// @Success  200    {object}  Http.WebResult{data=[]admin.Dept}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/dept [get]
func (a *apiDept) list(c *gin.Context) {
	ret := Http.NewResult(c)
	var requestData dto.QueryReqDept
	if Http.Bind(c, &requestData) {
		_, allDept := getDBList(&admin.Dept{}, &requestData)
		ret.Success(common.ResponseSuccess, a.getRootDept(allDept, func(dept *admin.Dept) bool {
			return true
		}))
	}
}

// @Tags     好奇宝宝后台-部门管理
// @Summary  获取部门详情
// @Param    token  header    string  true  "token"
// @Param    id     path      int     true  "部门id"
// @Success  200    {object}  Http.WebResult{data=admin.Dept}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/dept/{id} [get]
func (a *apiDept) info(c *gin.Context) {
	getItemById(c, &admin.Dept{}, nil)
}

// @Tags     好奇宝宝后台-部门管理
// @Summary  修改部门
// @Param    token  header    string             true  "token"
// @Param    data   formData  dto.UpdateReqDept  true  "部门信息"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/dept [put]
func (a *apiDept) updateInfo(c *gin.Context) {
	updateItem(c, &dto.UpdateReqDept{}, nil)
}

// @Tags     好奇宝宝后台-部门管理
// @Summary  新增部门
// @Param    token  header    string             true  "token"
// @Param    data   formData  dto.InsertReqDept  true  "部门信息"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/dept [post]
func (a *apiDept) addInfo(c *gin.Context) {
	addItem(c, &dto.InsertReqDept{})
}

// @Tags     好奇宝宝后台-部门管理
// @Summary  删除部门
// @Param    token  header    string            true  "token"
// @Param    data   body      dto.DeleteReqDept  true  "部门id"
// @Success  200    {object}  Http.WebResult
// @Failure  500    {object}  Http.WebResult
// @Router   /api/dept [delete]
func (a *apiDept) delInfo(c *gin.Context) {
	deleteItem(c, &dto.DeleteReqDept{}, nil)
}

// @Tags     好奇宝宝后台-部门管理
// @Summary  获取部门(简介)列表
// @Param    token  header    string             true  "token"
// @Success  200    {object}  Http.WebResult{data=[]admin.Dept}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/dept/deptTree [get]
func (a *apiDept) deptTree(c *gin.Context) {
	ret := Http.NewResult(c)
	var retData []map[string]interface{}

	_, allDept := getDBList(&admin.Dept{}, &dto.QueryReqDept{})
	rootDept := a.getRootDept(allDept, func(dept *admin.Dept) bool {
		return true
	})

	for _, dept := range rootDept {
		retData = append(retData, dept.Summary())
	}

	ret.Success(common.ResponseSuccess, retData)
}
