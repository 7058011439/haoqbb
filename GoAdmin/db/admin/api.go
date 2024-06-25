package admin

import (
	"fmt"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/jinzhu/gorm"
)

type Api struct {
	common.Model
	Title   string  `json:"title"`                            // 标题
	Path    string  `json:"path"`                             // 地址
	Type    string  `json:"type"`                             // 类型(BUS, SYS)
	Method  string  `json:"method"`                           // 请求方式(GET, PUT, DELETE, POST)
	Log     string  `json:"log" gorm:"default:'Y';size:4;"`   // 是否记录日志(Y, N)
	SysMenu []*Menu `json:"-" gorm:"many2many:menu_api_rule"` // 这是一个无用的字段,主要方便删除api的时候(通过BeforeDelete)快捷清理menu_api_rule数据
	common.ControlBy
}

func (a *Api) TableName() string {
	return "api"
}

func (a *Api) IsValid() bool {
	return a.Model.IsValid() && a.Title != ""
}

func (a *Api) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Model(a).Association("SysMenu").Clear().Error
}

func (a *Api) FormatPath() string {
	return fmt.Sprintf("%-10s%s", a.Method, a.Path)
}
