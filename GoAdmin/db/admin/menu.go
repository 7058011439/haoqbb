package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Menu struct {
	common.Model
	MenuName   string  `json:"menuName"`   // 菜单名
	Title      string  `json:"title"`      // 菜单描述
	Icon       string  `json:"icon"`       // 图标
	Path       string  `json:"path"`       // 浏览器Url地址
	MenuType   string  `json:"menuType"`   // 类型(M-目录, C-菜单, F-功能/按钮)
	Action     string  `json:"action"`     // todo
	Permission string  `json:"permission"` // 权限
	ParentId   int64   `json:"parentId"`   // 父节点id
	NoCache    bool    `json:"noCache"`    // todo
	Breadcrumb string  `json:"breadcrumb"` // todo
	Component  string  `json:"component"`  // 对应客户端代码(.vue)文件
	Sort       int     `json:"sort"`       // 排序
	Visible    string  `json:"visible"`    // 是否显示
	IsFrame    string  `json:"isFrame"`    // todo
	SysApi     []*Api  `json:"sysApi" gorm:"many2many:menu_api_rule"`
	SysRole    []*Role `json:"-" gorm:"many2many:role_menu_rule"`
	Apis       []int64 `json:"apis" gorm:"-"`
	DataScope  string  `json:"dataScope" gorm:"-"`
	Params     string  `json:"params" gorm:"-"`
	IsSelect   bool    `json:"is_select" gorm:"-"`
	Children   []*Menu `json:"children" gorm:"-"`
	common.ControlBy
}

func (m *Menu) TableName() string {
	return "menu"
}

func (m *Menu) IsValid() bool {
	return m.Model.IsValid() && m.MenuName != ""
}

func (m *Menu) GetParentId() int64 {
	return m.ParentId
}

func (m *Menu) AddChild(child common.IChild) {
	m.Children = append(m.Children, child.(*Menu))
}

func (m *Menu) Summary() map[string]interface{} {
	ret := map[string]interface{}{
		"id":    m.ID,
		"label": m.Title,
	}
	if len(m.Children) > 0 {
		var child []interface{}
		for _, c := range m.Children {
			child = append(child, c.Summary())
		}
		ret["children"] = child
	}
	return ret
}

func (m *Menu) BeforeDelete(tx *gorm.DB) (err error) {
	childCount := int64(0)
	tx.Model(&Menu{}).Where(&Menu{ParentId: m.ID}).Count(&childCount)
	if childCount > 0 {
		return fmt.Errorf("当前菜单还有子菜单，请处理子菜单")
	}
	if err := tx.Model(m).Association("SysApi").Clear().Error; err == nil {
		return tx.Model(m).Association("SysRole").Clear().Error
	} else {
		return err
	}
}
