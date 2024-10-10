package common

import (
	"fmt"
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/7058011439/haoqbb/Log"
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

var permissionMutex sync.RWMutex
var permissionList = map[int64]map[string]bool{} // 权限列表 map[roleId]map[api]bool
var sysApi = map[string]bool{}                   // 系统api(不需要验证权限)
var logApi = map[string]bool{}                   // 需要记录日志api

func init() {
	RefreshPermission()
}

func RefreshPermission() {
	start := time.Now()
	permissionMutex.Lock()
	permissionMutex.Unlock()
	permissionList = map[int64]map[string]bool{}
	sysApi = map[string]bool{}
	logApi = map[string]bool{}

	var roleList []*admin.Role
	admin.MysqlDB().Preload("SysMenu.SysApi").Find(&roleList)
	for _, role := range roleList {
		refreshPermission(role)
	}

	var apiList []*admin.Api
	admin.MysqlDB().Find(&apiList)
	for _, api := range apiList {
		path := api.FormatPath()
		if api.Type == "SYS" {
			sysApi[path] = true
		}
		if api.Log == "Y" {
			logApi[path] = true
		}
	}

	Log.Debug("初始化权限列表耗时:%v ms", time.Since(start).Milliseconds())
	/*
		for id, p := range permissionList {
			fmt.Println(id, len(p))
			for a := range p {
				fmt.Printf("\t%v\n", a)
			}
		}

		for api := range sysApi {
			fmt.Println(api)
		}
	*/
}

func refreshPermission(role *admin.Role) {
	// 系统管理员有所有权限, 无需添加权限列表
	if role.ID == common.AdminRoleId {
		return
	}
	// 非禁用状态添加角色权限
	if role.Status != common.StatusForbid {
		permissionList[role.ID] = map[string]bool{}
		for _, menu := range role.SysMenu {
			for _, api := range menu.SysApi {
				if api.Type == "BUS" {
					permissionList[role.ID][api.FormatPath()] = true
				}
			}
		}
	}
}

func checkPermission(c *gin.Context, path string) error {
	roleId := GetAdminRoleId(c)
	if roleId == common.AdminRoleId {
		return nil
	}
	if _, ok := sysApi[path]; ok {
		return nil
	}
	permissionMutex.RLock()
	defer permissionMutex.RUnlock()
	if apiList, ok := permissionList[roleId]; ok {
		if _, ok := apiList[path]; ok {
			return nil
		} else {
			Log.Error("%v: 用户没有该权限, path = %v", ResponsePermissionError, path)
			return fmt.Errorf("%v: 用户没有该权限", ResponsePermissionError)
		}
	} else {
		return fmt.Errorf("%v: 用户组错误", ResponsePermissionError)
	}
}
