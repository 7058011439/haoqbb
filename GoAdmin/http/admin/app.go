package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	"github.com/7058011439/haoqbb/Http"
	"github.com/7058011439/haoqbb/Log"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

type apiApp struct {
	Http.Api
}

func init() {
	a := common.ServerAdmin.RegeditGroup("/api/app", &apiApp{}).(*apiApp)
	a.RegeditApi(Http.TypeGet, "/captcha", a.captcha)
	a.RegeditApi(Http.TypeGet, "/captchaTest", a.captchaTest)
	a.RegeditApi(Http.TypeGet, "/app-config", a.config)
}

// @Tags     好奇宝宝后台-系统工具
// @Summary  获取验证码
// @Success  200  {object}  Http.WebResult
// @Router   /api/app/captcha [get]
func (a *apiApp) captcha(c *gin.Context) {
	ret := Http.NewResult(c)
	if ids, codeUrl, _, err := base64Captcha.NewCaptcha(base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80), base64Captcha.DefaultMemStore).Generate(); err == nil {
		ret.Success(common.ResponseSuccess, map[string]interface{}{
			"codeUrl": codeUrl,
			"uuid":    ids,
		})
	} else {
		Log.Error("DriverDigitFunc error, %s", err.Error())
		ret.Fail("验证码获取失败", err)
	}
}

// @Tags     好奇宝宝后台-系统工具
// @Summary  获取验证码(测试)
// @Success  200  {object}  Http.WebResult
// @Router   /api/app/captchaTest [get]
func (a *apiApp) captchaTest(c *gin.Context) {
	ret := Http.NewResult(c)
	if ids, codeUrl, answer, err := base64Captcha.NewCaptcha(base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80), base64Captcha.DefaultMemStore).Generate(); err == nil {
		ret.Success(common.ResponseSuccess, map[string]interface{}{
			"codeUrl": codeUrl,
			"uuid":    ids,
			"answer":  answer,
		})
	} else {
		Log.Error("DriverDigitFunc error, %s", err.Error())
		ret.Fail("验证码获取失败", err)
	}
}

// @Tags     好奇宝宝后台-系统工具
// @Summary  获取配置信息
// @Success  200  {object}  Http.WebResult
// @Router   /api/app/app-config [get]
func (a *apiApp) config(c *gin.Context) {
	ret := Http.NewResult(c)
	ret.Success(common.ResponseSuccess, map[string]interface{}{
		"sys_app_logo": admin.GetConfigByKey("sys_app_logo").ConfigValue,
		"sys_app_name": admin.GetConfigByKey("sys_app_name").ConfigValue,
	})
}
