
# Swagger 调试工具生成说明
这里是后台管理核心，只有：后台管理(公共)接口

## 生成后台管理(公共)接口
```shell
swag i -dir http/admin --instanceName core --parseDependency
```

## 程序调用代码
```code
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

func init() {
	RegeditSwagger(common.ServerAdmin, "core", config.HttpPort())
	common.ServerAdmin.Start(config.HttpPort())
	Log.Log("基础服务启动完成")
}

func RegeditSwagger(server *Http.Server, path string, port int) {
	server.GET(fmt.Sprintf("/swagger/%v/*any", path), ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName(path)))
	Log.Debug(fmt.Sprintf("打开:http://localhost:%v/swagger/%v/index.html#,可进入API调试", port, path))
}
```

## 编译运行程序
```shell
    只是模块不单独运行
```

## 测试调试地址
后台管理(公共)接口 [http://localhost:corePort/swagger/core/index.html](http://localhost:8000/swagger/core/index.html)
