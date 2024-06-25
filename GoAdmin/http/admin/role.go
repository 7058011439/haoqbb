package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	commondb "github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
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

func (a *apiRole) list(c *gin.Context) {
	getList(c, &admin.Role{}, &dto.QueryReqRole{})
}

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

func (a *apiRole) addInfo(c *gin.Context) {
	addItem(c, &dto.InsertReqRole{})
}

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
