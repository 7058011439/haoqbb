package goAdmin

import (
	"github.com/7058011439/haoqbb/GoAdmin/config"
	// "github.com/7058011439/haoqbb/GoAdmin/docs"
	_ "github.com/7058011439/haoqbb/GoAdmin/http/admin"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	"github.com/7058011439/haoqbb/Log"
)

// @title         好奇宝宝后台
// @version       1.0
// @description   go-admin基础模块接口
// @contact.name  API Support
// @contact.url   http://www.swagger.io/support
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func init() {
	//common.ServerAdmin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	common.ServerAdmin.Start(config.HttpPort())
	//url := fmt.Sprintf("http://localhost:%v/swagger/index.html#", config.HttpPort())
	//Log.Debug(fmt.Sprintf("打开:%v,可进入API调试", url))
	Log.Log("基础服务启动完成")
}
