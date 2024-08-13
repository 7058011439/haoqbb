package GoAdmin

import (
	"fmt"
	"github.com/7058011439/haoqbb/GoAdmin/config"
	_ "github.com/7058011439/haoqbb/GoAdmin/docs"
	_ "github.com/7058011439/haoqbb/GoAdmin/http/admin"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	"github.com/7058011439/haoqbb/Http"
	"github.com/7058011439/haoqbb/Log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title         好奇宝宝后台
// @version       1.0
// @description   go-admin基础模块接口
// @contact.name  API Support
// @contact.url   http://www.swagger.io/support
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func init() {
	RegeditSwagger(common.ServerAdmin, "core", config.HttpPort())
	common.ServerAdmin.Start(config.HttpPort())
	Log.Log("基础服务启动完成")
}

func RegeditSwagger(server *Http.Server, path string, port int) {
	server.GET(fmt.Sprintf("/swagger/%v/*any", path), ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName(path)))
	Log.Debug(fmt.Sprintf("打开:http://localhost:%v/%v/index.html#,可进入API调试", port, path))
}
