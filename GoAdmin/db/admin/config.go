package admin

import "github.com/7058011439/haoqbb/GoAdmin/db/common"

type Config struct {
	common.Model
	ConfigName  string `json:"configName" gorm:"unique;"`            // 配置名称
	ConfigKey   string `json:"configKey" gorm:"unique;"`             // 配置key
	ConfigValue string `json:"configValue"`                          // 配置value
	ConfigType  string `json:"configType"`                           // 配置类型(系统内置(Y,N))
	IsFrontend  string `json:"isFrontend" gorm:"size:4;default:'1'"` // 是否前端配置(2,1)
	Remark      string `json:"remark"`                               // 备注
	common.ControlBy
}

func (c *Config) TableName() string {
	return "config"
}

func (c *Config) IsValid() bool {
	return c.Model.IsValid() && c.ConfigName != ""
}
