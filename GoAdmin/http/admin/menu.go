package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	commondb "github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	"github.com/7058011439/haoqbb/GoAdmin/util"
	"github.com/7058011439/haoqbb/Http"
	"github.com/gin-gonic/gin"
)

type apiMenu struct {
	Http.Api
}

const (
	menuTypeMain  = "M" // 一级主菜单
	menuTypeChild = "C" // 子菜单
	menuTypeFunc  = "F" // 功能性菜单
)

func init() {
	a := common.ServerAdmin.RegeditGroup("/api/menu", &apiMenu{}, common.CheckAdminToken).(*apiMenu)
	a.RegeditApi(Http.TypeGet, "/menurole", a.menuRole)
	a.RegeditApi(Http.TypeGet, "", a.list)
	a.RegeditApi(Http.TypeGet, "/:id", a.info)
	a.RegeditApi(Http.TypePut, "", a.updateInfo)
	a.RegeditApi(Http.TypePost, "", a.addInfo)
	a.RegeditApi(Http.TypeDelete, "", a.deleteInfo)

	a.RegeditApi(Http.TypeGet, "/menuTree/:id", a.menuTree)
}

func (a *apiMenu) getRootMenu(listData []*admin.Menu, fun func(menu *admin.Menu) bool) (rootMenu []*admin.Menu) {
	var listInterface []commondb.IChild
	for _, data := range listData {
		listInterface = append(listInterface, data)
	}
	for _, data := range listData {
		if data.ParentId == 0 {
			fillChild(listInterface, data, func(item commondb.IChild) bool {
				return fun(item.(*admin.Menu))
			})
			rootMenu = append(rootMenu, data)
		}
	}
	return
}

func (a *apiMenu) findFamilyMenu(menu *admin.Menu, allMenu map[int64]*admin.Menu) []*admin.Menu {
	if menu.ParentId == 0 {
		return []*admin.Menu{menu}
	} else {
		if m, ok := allMenu[menu.ParentId]; ok {
			familyMenu := a.findFamilyMenu(m, allMenu)
			familyMenu = append(familyMenu, menu)
			return familyMenu
		} else {
			return nil
		}
	}
}

func (a *apiMenu) menuRole(c *gin.Context) {
	ret := Http.NewResult(c)
	menuList := admin.GetRoleMenu(common.GetAdminRoleId(c))
	if common.GetAdminRoleId(c) != commondb.AdminRoleId {
		allMenu := admin.GetRoleMenu(commondb.AdminRoleId)
		menuMap := map[int64]*admin.Menu{}
		for _, m := range allMenu {
			menuMap[m.ID] = m
		}
		for _, menu := range menuList {
			menuList = append(menuList, a.findFamilyMenu(menu, menuMap)...)
		}
		menuList = util.RemoveDuplicates(menuList, func(t interface{}) interface{} {
			return t.(*admin.Menu).ID
		})
	}
	ret.Success(common.ResponseSuccess, a.getRootMenu(menuList, func(menu *admin.Menu) bool {
		return menu.MenuType == menuTypeMain || menu.MenuType == menuTypeChild
	}))
}

func (a *apiMenu) list(c *gin.Context) {
	var requestData dto.QueryReqMenu
	if Http.Bind(c, &requestData) {
		ret := Http.NewResult(c)
		_, allMenu := getDBList(&admin.Menu{}, &requestData)
		ret.Success(common.ResponseSuccess, a.getRootMenu(allMenu, func(menu *admin.Menu) bool {
			return true
		}))
	}
}

func (a *apiMenu) menuTree(c *gin.Context) {
	ret := Http.NewResult(c)
	var retData []map[string]interface{}

	allMenu := admin.GetMenuListSummary()
	rootMenu := a.getRootMenu(allMenu, func(menu *admin.Menu) bool {
		return true
	})

	for _, menu := range rootMenu {
		retData = append(retData, menu.Summary())
	}

	ret.Success(common.ResponseSuccess, map[string]interface{}{
		"menus":       retData,
		"checkedKeys": admin.GetRoleMenuId(getReqId(c)),
	})
}

func (a *apiMenu) info(c *gin.Context) {
	getItemById(c, &admin.Menu{}, func(data commondb.IDataDB) {
		menu := data.(*admin.Menu)
		menu.SysApi = admin.GetMenuApi(menu.GetId())
		menu.Apis = []int64{}
		for _, api := range menu.SysApi {
			menu.Apis = append(menu.Apis, api.ID)
		}
	})
}

func (a *apiMenu) updateInfo(c *gin.Context) {
	updateItem(c, &dto.UpdateReqMenu{}, func(data commondb.IDataDB) {
		menu := data.(*admin.Menu)

		// 修改菜单对应api关系
		admin.MysqlDB().Model(menu).Association("SysApi").Clear()
		var apiList []*admin.Api
		admin.MysqlDB().Find(&apiList, menu.Apis)
		admin.MysqlDB().Model(menu).Association("SysApi").Append(apiList)

		// 修改菜单可能权限变更，需刷新角色对应权限列表
		common.RefreshPermission()
	})
}

func (a *apiMenu) addInfo(c *gin.Context) {
	addItem(c, &dto.InsertReqMenu{})
}

func (a *apiMenu) deleteInfo(c *gin.Context) {
	deleteItem(c, &dto.DeleteReqMenu{}, func(ids []int64) {
		// 删除菜单可能权限变更，需刷新角色对应权限列表
		common.RefreshPermission()
	})
}
