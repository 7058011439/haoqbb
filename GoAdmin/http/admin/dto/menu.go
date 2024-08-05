package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqMenu struct {
	Title   string `form:"title" json:"title"`     // 菜单标题
	Visible int    `form:"visible" json:"visible"` // 是否显示
	commonHttp.QueryParam
}

func (q *QueryReqMenu) Condition() interface{} {
	model := &admin.Menu{}
	commonHttp.ReqToDB(q, model)
	return model
}

func (q *QueryReqMenu) Preload() string {
	return "SysApi"
}

func (q *QueryReqMenu) Order() string {
	return "sort asc"
}

type UpdateReqMenu struct {
	commonHttp.UpdateId
	MenuName   string  `json:"menuName"`   // 菜单名
	Title      string  `json:"title"`      // 菜单描述
	Icon       string  `json:"icon"`       // 图标
	Path       string  `json:"path"`       // url显示地址
	Paths      string  `json:"paths"`      // 没鸟用
	MenuType   string  `json:"menuType"`   // 菜单类型(M-目录, C-菜单, F-功能/按钮)
	Action     string  `json:"action"`     // 没鸟用
	Permission string  `json:"permission"` // 权限
	ParentId   int64   `json:"parentId"`   // 父节点id
	NoCache    bool    `json:"noCache"`    // 是否缓存
	Breadcrumb string  `json:"breadcrumb"` // 没鸟用
	Component  string  `json:"component"`  // 对应客户端代码(.vue)文件
	Sort       int     `json:"sort"`       // 排序
	Visible    string  `json:"visible"`    // 是否显示
	IsFrame    string  `json:"isFrame"`    // 没鸟用
	Apis       []int64 `json:"apis"`       // 该菜单拥有的api权限列表
}

func (u *UpdateReqMenu) Model() common.IUpdateData {
	model := &admin.Menu{}
	commonHttp.ReqToDB(u, model)
	return model
}

type InsertReqMenu struct {
	MenuName   string  `json:"menuName"`   // 菜单名
	Title      string  `json:"title"`      // 菜单描述
	Icon       string  `json:"icon"`       // 图标
	Path       string  `json:"path"`       // url显示地址
	Paths      string  `json:"paths"`      // 没鸟用
	MenuType   string  `json:"menuType"`   // 菜单类型(M-目录, C-菜单, F-功能/按钮)
	Action     string  `json:"action"`     // 没鸟用
	Permission string  `json:"permission"` // 权限
	ParentId   int64   `json:"parentId"`   // 父节点id
	NoCache    bool    `json:"noCache"`    // 是否缓存
	Breadcrumb string  `json:"breadcrumb"` // 没鸟用
	Component  string  `json:"component"`  // 对应客户端代码(.vue)文件
	Sort       int     `json:"sort"`       // 排序
	Visible    string  `json:"visible"`    // 是否显示
	IsFrame    string  `json:"isFrame"`    // 没鸟用
	Apis       []int64 `json:"apis"`       // 该菜单拥有的api权限列表
}

func (u *InsertReqMenu) Model() common.IUpdateData {
	model := &admin.Menu{}
	commonHttp.ReqToDB(u, model)
	admin.MysqlDB().Find(&model.SysApi, u.Apis)
	return model
}

type DeleteReqMenu struct {
	commonHttp.DeleteId
}

func (d *DeleteReqMenu) Model() common.IDataDB {
	return &admin.Menu{}
}
