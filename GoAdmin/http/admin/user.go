package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	commondb "github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/7058011439/haoqbb/GoAdmin/http/admin/dto"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
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

func (a *apiUser) list(c *gin.Context) {
	getList(c, &admin.User{}, &dto.QueryReqUser{})
}

func (a *apiUser) info(c *gin.Context) {
	getItemById(c, &admin.User{}, nil)
}

func (a *apiUser) updateInfo(c *gin.Context) {
	updateItem(c, &dto.UpdateReqUser{}, func(data commondb.IDataDB) {
		user := data.(*admin.User)
		// 这个地方粗暴处理，因为不知道是否有改所属角色和状态，那就直接清理token，让用户重新登录
		a.clearToken(user.ID)
	})
}

func (a *apiUser) addInfo(c *gin.Context) {
	addItem(c, &dto.InsertReqUser{})
}

func (a *apiUser) delInfo(c *gin.Context) {
	deleteItem(c, &dto.DeleteReqUser{}, func(ids []int64) {
		// 用户删除，需要删除其token
		for _, id := range ids {
			a.clearToken(id)
		}
	})
}

func (a *apiUser) updateStatus(c *gin.Context) {
	updateItem(c, &dto.UpdateReqUserStatus{}, func(data commondb.IDataDB) {
		user := data.(*admin.User)
		// 如果管理员禁用，需要删除其token
		if user.Status == commondb.StatusForbid {
			a.clearToken(user.ID)
		}
	})
}

func (a *apiUser) updatePassword(c *gin.Context) {
	updateItem(c, &dto.UpdateReqUserPassword{}, func(data commondb.IDataDB) {
		user := data.(*admin.User)
		// 修改管理员密码后，需要让该管理员token失效
		a.clearToken(user.ID)
	})
}
