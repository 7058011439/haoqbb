package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	commondb "github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	_ "github.com/7058011439/haoqbb/GoAdmin/util"
	"github.com/7058011439/haoqbb/Http"
	"github.com/gin-gonic/gin"
)

type apiUser struct {
	Http.Api
}

func init() {
	a := common.ServerAdmin.RegeditGroup("/api/user", &apiUser{}, common.CheckAdminToken).(*apiUser)
	a.RegeditApi(Http.TypeGet, "", a.list)
	a.RegeditApi(Http.TypeGet, "/:id", a.info)
	a.RegeditApi(Http.TypePut, "", a.updateInfo)
	a.RegeditApi(Http.TypePost, "", a.addInfo)
	a.RegeditApi(Http.TypeDelete, "", a.delInfo)

	a.RegeditApi(Http.TypePut, "/status", a.updateStatus)
	a.RegeditApi(Http.TypePut, "/pwd", a.updatePassword)
}

func (a *apiUser) clearToken(userId int64) {
	common.UpdateToken(common.UserTypeAdmin, userId, "")
}

// @Tags     好奇宝宝后台-用户管理
// @Summary  获取用户列表
// @Param    token  header    string            true  "token"
// @Param    data   query     dto.QueryReqUser  true  "查询条件"
// @Success  200    {object}  Http.WebResult{data=util.WebResultCommonList{list=[]admin.User}}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/user [get]
func (a *apiUser) list(c *gin.Context) {
	getList(c, &admin.User{}, &dto.QueryReqUser{})
}

// @Tags     好奇宝宝后台-用户管理
// @Summary  获取用户详情
// @Param    token  header    string  true  "token"
// @Param    id     path      int     true  "用户id"
// @Success  200    {object}  Http.WebResult{data=admin.User}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/user/{id} [get]
func (a *apiUser) info(c *gin.Context) {
	getItemById(c, &admin.User{}, nil)
}

// @Tags     好奇宝宝后台-用户管理
// @Summary  修改用户
// @Param    token  header    string                   true  "token"
// @Param    data   formData  dto.UpdateReqUser  true  "用户信息"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/user [put]
func (a *apiUser) updateInfo(c *gin.Context) {
	updateItem(c, &dto.UpdateReqUser{}, func(data commondb.IDataDB) {
		user := data.(*admin.User)
		// 这个地方粗暴处理，因为不知道是否有改所属角色和状态，那就直接清理token，让用户重新登录
		a.clearToken(user.ID)
	})
}

// @Tags     好奇宝宝后台-用户管理
// @Summary  新增用户
// @Param    token  header    string                   true  "token"
// @Param    data   formData  dto.InsertReqUser  true  "用户信息"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/user [Post]
func (a *apiUser) addInfo(c *gin.Context) {
	addItem(c, &dto.InsertReqUser{})
}

// @Tags     好奇宝宝后台-用户管理
// @Summary  删除用户
// @Param    token  header    string             true  "token"
// @Param    data   body      dto.DeleteReqUser  true  "用户id"
// @Success  200    {object}  Http.WebResult
// @Failure  500    {object}  Http.WebResult
// @Router   /api/user [delete]
func (a *apiUser) delInfo(c *gin.Context) {
	deleteItem(c, &dto.DeleteReqUser{}, func(ids []int64) {
		// 用户删除，需要删除其token
		for _, id := range ids {
			a.clearToken(id)
		}
	})
}

// @Tags     好奇宝宝后台-用户管理
// @Summary  修改用户状态
// @Param    token  header    string             true  "token"
// @Param    data   formData  dto.UpdateReqUserStatus  true  "用户状态"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/user/status [put]
func (a *apiUser) updateStatus(c *gin.Context) {
	updateItem(c, &dto.UpdateReqUserStatus{}, func(data commondb.IDataDB) {
		user := data.(*admin.User)
		// 如果管理员禁用，需要删除其token
		if user.Status == commondb.StatusForbid {
			a.clearToken(user.ID)
		}
	})
}

// @Tags     好奇宝宝后台-用户管理
// @Summary  修改用户密码
// @Param    token  header    string             true  "token"
// @Param    data   formData  dto.UpdateReqUserStatus  true  "用户密码"
// @Success  200    {object}  Http.WebResult{data=int64}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/user/pwd [put]
func (a *apiUser) updatePassword(c *gin.Context) {
	updateItem(c, &dto.UpdateReqUserPassword{}, func(data commondb.IDataDB) {
		user := data.(*admin.User)
		// 修改管理员密码后，需要让该管理员token失效
		a.clearToken(user.ID)
	})
}
