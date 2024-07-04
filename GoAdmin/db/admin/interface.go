package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/jinzhu/gorm"
	"sort"
)

var dbManager *common.Manager

func LoadForumData(db common.IDataDB) common.IDataDB {
	return dbManager.LoadData(db)
}

func InsertForumData(db common.IDataDB) error {
	return dbManager.InsertData(db)
}

func UpdateForumData(db common.IDataDB) error {
	return dbManager.UpdateData(db)
}

func DeleteForumData(db common.IDataDB, ids ...int64) error {
	return dbManager.DeleteData(db, ids...)
}

func ClearMemData(db common.IDataDB, ids ...int64) error {
	return dbManager.ClearMemData(db, ids...)
}

func MysqlDB() *gorm.DB {
	return dbManager.MysqlDB()
}

func DBManager() *common.Manager {
	return dbManager
}

func GetAdmin(id int64) *User {
	return LoadForumData(&User{Model: common.Model{ID: id}}).(*User)
}

func GetAdminByUserName(userName string) *User {
	admin := &User{UserName: userName}
	dbManager.MysqlDB().Where(admin).First(admin)
	return admin
}

func GetRole(id int64) *Role {
	return LoadForumData(&Role{Model: common.Model{ID: id}}).(*Role)
}

func GetConfigByKey(key string) *Config {
	return LoadForumData(&Config{ConfigKey: key}).(*Config)
}

func GetMenuListSummary() (ret []*Menu) {
	dbManager.MysqlDB().Select("id, title, parent_id").Order("sort").Find(&ret)
	return
}

func GetRoleMenuId(roleId int64) (ret []int64) {
	menuList := GetRoleMenu(roleId)
	for _, menu := range menuList {
		ret = append(ret, menu.ID)
	}
	return
}

func GetRoleMenu(roleId int64) (ret []*Menu) {
	if roleId == common.AdminRoleId {
		dbManager.MysqlDB().Find(&ret)
	} else {
		var role Role
		dbManager.MysqlDB().Preload("SysMenu").Find(&role, roleId)
		ret = role.SysMenu
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Sort < ret[j].Sort
	})
	return
}

func GetMenuApi(menuId int64) []*Api {
	var menu Menu
	dbManager.MysqlDB().Preload("SysApi").Find(&menu, menuId)
	return menu.SysApi
}
