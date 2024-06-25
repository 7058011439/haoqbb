package dto

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	commonHttp "github.com/7058011439/haoqbb/GoAdmin/http/common"
)

type QueryReqMenu struct {
	Title   string `form:"title" json:"title"`
	Visible int    `form:"visible" json:"visible"`
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
	MenuName   string  `json:"menuName"` // 菜单名
	Title      string  `json:"title"`    // 菜单描述
	Icon       string  `json:"icon"`     // 图标
	Path       string  `json:"path"`
	Paths      string  `json:"paths"`
	MenuType   string  `json:"menuType"`
	Action     string  `json:"action"`
	Permission string  `json:"permission"`
	ParentId   int64   `json:"parentId"`
	NoCache    bool    `json:"noCache"`
	Breadcrumb string  `json:"breadcrumb"`
	Component  string  `json:"component"`
	Sort       int     `json:"sort"`
	Visible    string  `json:"visible"`
	IsFrame    string  `json:"isFrame"`
	Apis       []int64 `json:"apis"`
}

func (u *UpdateReqMenu) Model() common.IUpdateData {
	model := &admin.Menu{}
	commonHttp.ReqToDB(u, model)
	return model
}

type InsertReqMenu struct {
	MenuName   string  `json:"menuName"` // 菜单名
	Title      string  `json:"title"`    // 菜单描述
	Icon       string  `json:"icon"`     // 图标
	Path       string  `json:"path"`
	Paths      string  `json:"paths"`
	MenuType   string  `json:"menuType"`
	Action     string  `json:"action"`
	Permission string  `json:"permission"`
	ParentId   int64   `json:"parentId"`
	NoCache    bool    `json:"noCache"`
	Breadcrumb string  `json:"breadcrumb"`
	Component  string  `json:"component"`
	Sort       int     `json:"sort"`
	Visible    string  `json:"visible"`
	IsFrame    string  `json:"isFrame"`
	Apis       []int64 `json:"apis"`
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
