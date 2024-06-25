package common

import (
	"encoding/json"
	"fmt"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
)

type IQueryParam interface {
	GetPageIndex() int
	GetPageSize() int
	Condition() interface{}
	ConditionTime() interface{}
	Preload() string
	Order() string
}

type IUpdateParam interface {
	Model() common.IUpdateData
}

type IInsertData interface {
	Model() common.IUpdateData
}

type IDeleteData interface {
	Model() common.IDataDB
	GetIds() []int64
}

type QueryParam struct {
	PageIndex int    `form:"pageIndex" json:"pageIndex"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
	BeginTime string `form:"beginTime" json:"beginTime"`
	EndTime   string `form:"endTime" json:"endTime"`
}

func (q *QueryParam) GetPageIndex() int {
	return q.PageIndex
}

func (q *QueryParam) GetPageSize() int {
	return q.PageSize
}

func (q *QueryParam) Condition() interface{} {
	return ""
}

func (q *QueryParam) ConditionTime() interface{} {
	condition := ""
	if q.BeginTime != "" {
		condition = fmt.Sprintf("created_at >= '%v'", q.BeginTime)
	}
	if q.EndTime != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf("created_at <= '%v'", q.EndTime)
	}
	return condition
}

func (q *QueryParam) Preload() string {
	return ""
}

func (q *QueryParam) Order() string {
	return ""
}

type UpdateId struct {
	ID int64 `json:"id" form:"id"` // 用户ID
}

func (u *UpdateId) GetId() int64 {
	return u.ID
}

type DeleteId struct {
	IDS []int64 `json:"ids" form:"ids"`
}

func (d *DeleteId) GetIds() []int64 {
	return d.IDS
}

func ReqToDB(req interface{}, model common.IDataDB) {
	d, _ := json.Marshal(req)
	json.Unmarshal(d, model)
}
