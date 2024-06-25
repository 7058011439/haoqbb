package common

import (
	"github.com/7058011439/haoqbb/GoAdmin/config"
	"github.com/7058011439/haoqbb/Http"
)

var ServerAdmin = Http.NewHttpServer(config.HttpVersion(), true)
