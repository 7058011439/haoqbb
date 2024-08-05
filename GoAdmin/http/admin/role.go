package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	commondb "github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	_ "github.com/7058011439/haoqbb/GoAdmin/util"
	"github.com/7058011439/haoqbb/Http"
	"github.com/gin-gonic/gin"
)

type apiRole struct {
	Http.Api
}

func init() {
	a := common.ServerAdmin.RegeditGroup("/api/role", &apiRole{}, common.CheckAdminToken).(*apiRole)
	a.RegeditApi(Http.TypeGet, "", a.list)
	a.RegeditApi(Http.TypeGet, "/:id", a.info)
	a.RegeditApi(Http.TypePut, "", a.updateInfo)
	a.RegeditApi(Http.TypePost, "", a.addInfo)
	a.RegeditApi(Http.TypeDelete, "", a.delInfo)

	a.RegeditApi(Http.TypePut, "/status", a.updateStatus)
}

// @Tags     好奇宝宝后台-角色管理
// @Summary  获取角色列表
// @Param    token  header    string            true  "token"
// @Param    data   query     dto.QueryReqRole  true  "查询条件"
// @Success  200    {object}  Http.WebResult{data=util.WebResultCommonList{list=[]admin.Role}}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/role [get]
func (a *apiRole) list(c *gin.Context) {
	getList(c, &admin.Role{}, &dto.QueryReqRole{})
}

// @Tags     好奇宝宝后台-角色管理
// @Summary  获取角色详情
// @Param    token  header    string  true  "token"
// @Param    id     path      int     true  "角色id"
// @Success  200    {object}  Http.WebResult{data=admin.Role}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/role/{id} [get]
func (a *apiRole) info(c *gin.Context) {
	getItemById(c, &admin.Role{}, func(data commondb.IDataDB) {
		role := data.(*admin.Role)
		role.MenuIds = admin.GetRoleMenuId(role.GetId())
	})
}

func (a *apiRole) clearToken(roleId int64) {
	var userList []*admin.User
	admin.MysqlDB().Where(&admin.User{RoleId: roleId}).Find(&userList)
	for _, user := range userList {
		common.UpdateToken(common.UserTypeAdmin, user.ID, "")
		admin.ClearMemData(user)
	}
}

// @Tags     好奇宝宝后台-角色管理
// @Summary  修改角色
// @Param    token  header    string             true  "token"
// @Param    data   formData  dto.UpdateReqRole  true  "角色信息"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/role [put]
func (a *apiRole) updateInfo(c *gin.Context) {
	updateItem(c, &dto.UpdateReqRole{}, func(data commondb.IDataDB) {
		role := data.(*admin.Role)

		// 修改角色对应菜单关系
		admin.MysqlDB().Model(role).Association("SysMenu").Clear()
		var menuList []*admin.Menu
		admin.MysqlDB().Find(&menuList, role.MenuIds)
		admin.MysqlDB().Model(role).Association("SysMenu").Append(menuList)

		// 这里粗暴处理，修改角色的时候可能会修改权限和状态，所以直接将角色对应的用户token全部清理
		a.clearToken(role.ID)

		// 修改角色后，需刷新角色对应权限列表
		common.RefreshPermission()
	})
}

// @Tags     好奇宝宝后台-角色管理
// @Summary  新增角色
// @Param    token  header    string             true  "token"
// @Param    data   formData  dto.InsertReqRole  true  "角色信息"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/role [Post]
func (a *apiRole) addInfo(c *gin.Context) {
	addItem(c, &dto.InsertReqRole{})

	// 新增角色后，需刷新角色对应权限列表
	common.RefreshPermission()
}

// @Tags     好奇宝宝后台-角色管理
// @Summary  删除角色
// @Param    token  header    string             true  "token"
// @Param    data   body      dto.DeleteReqRole  true  "角色id"
// @Success  200    {object}  Http.WebResult
// @Failure  500    {object}  Http.WebResult
// @Router   /api/role [delete]
func (a *apiRole) delInfo(c *gin.Context) {
	deleteItem(c, &dto.DeleteReqRole{}, func(ids []int64) {
		for _, id := range ids {
			// 当前角色下属用户token全部清理
			a.clearToken(id)
		}

		// 删除角色后，需刷新角色对应权限列表
		common.RefreshPermission()
	})
}

// @Tags     好奇宝宝后台-角色管理
// @Summary  修改角色状态
// @Param    token  header    string                   true  "token"
// @Param    data   formData  dto.UpdateReqRoleStatus  true  "角色状态"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/role/status [put]
func (a *apiRole) updateStatus(c *gin.Context) {
	updateItem(c, &dto.UpdateReqRoleStatus{}, func(data commondb.IDataDB) {
		role := data.(*admin.Role)
		// 角色禁用,角色对应用户token需要清理
		if role.Status == commondb.StatusForbid {
			a.clearToken(role.ID)
		}

		// 修改角色状态后，需刷新角色对应权限列表
		common.RefreshPermission()
	})
}
