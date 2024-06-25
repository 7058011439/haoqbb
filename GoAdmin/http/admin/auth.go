package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	common2 "github.com/7058011439/haoqbb/GoAdmin/db/common"
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	"github.com/7058011439/haoqbb/Http"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

type apiAuth struct {
	Http.Api
}

func init() {
	a := common.ServerAdmin.RegeditGroup("/api/auth", &apiAuth{}).(*apiAuth)
	a.RegeditApi(Http.TypePost, "/login", a.login)
	a.RegeditApi(Http.TypeGet, "/info", a.info, common.CheckAdminToken)
	a.RegeditApi(Http.TypePost, "/refresh_token", a.refreshToken, common.CheckAdminToken)
	a.RegeditApi(Http.TypePost, "/logout", a.logout, common.CheckAdminToken)
	a.RegeditApi(Http.TypePut, "/setpwd", a.setPwd, common.CheckAdminToken)

	a.RegeditApi(Http.TypeGet, "/profile", a.profile, common.CheckAdminToken)
}

func (a apiAuth) loginRet(c *gin.Context, userName string, msg string, data interface{}) {
	ret := Http.NewResult(c)
	status := "2"
	if msg == common.ResponseSuccess {
		ret.Success(msg, data)
	} else {
		ret.Fail(msg, data)
		status = "1"
	}

	admin.InsertForumData(&admin.LoginLog{
		UserName: userName,
		Status:   status,
		Remark:   msg,
		Ip:       c.RemoteIP(),
	})
}

// @Tags     MiniGame 后台相关
// @Summary  管理员登录
// @Param    username  formData  string  true  "账号"
// @Param    password  formData  string  true  "密码"
// @Success  200       {object}  util.WebResult
// @Failure  500       {object}  util.WebResult
// @Router   /api/admin/login [post]
func (a *apiAuth) login(c *gin.Context) {
	requestData := struct {
		UserName string `form:"username" json:"username" binding:"required"`
		PassWord string `form:"password" json:"password" binding:"required"`
		Code     string `form:"code" json:"code"`
		UUID     string `form:"uuid" json:"uuid" binding:"required"`
	}{}
	if Http.Bind(c, &requestData) {
		if base64Captcha.DefaultMemStore.Verify(requestData.UUID, requestData.Code, true) {
			if manager := admin.GetAdminByUserName(requestData.UserName); manager.IsValid() && manager.PassWord == requestData.PassWord {
				if manager.Status != common2.StatusForbid {
					if role := admin.GetRole(manager.RoleId); role.IsValid() && role.Status != common2.StatusForbid {
						token, _ := common.NewToken(map[string]interface{}{
							common.TokenKeyId:            manager.ID,
							common.TokenKeyRoleId:        manager.RoleId,
							common.TokenKeyAdminUserName: requestData.UserName,
						}, 24*7)
						a.loginRet(c, requestData.UserName, common.ResponseSuccess, map[string]interface{}{"token": token})
						common.UpdateToken(common.UserTypeAdmin, manager.ID, token)
					} else {
						a.loginRet(c, requestData.UserName, "角色状态错误", nil)
					}
				} else {
					a.loginRet(c, requestData.UserName, "账号禁用中", nil)
				}
			} else {
				a.loginRet(c, requestData.UserName, "账号或者密码错误", nil)
			}
		} else {
			a.loginRet(c, requestData.UserName, "验证码错误", nil)
		}
	}
}

func (a *apiAuth) getPermissions(c *gin.Context) (ret []string) {
	roleId := common.GetAdminRoleId(c)
	if roleId == common2.AdminRoleId {
		ret = append(ret, "*:*:*")
	} else {
		menuList := admin.GetRoleMenu(roleId)
		for _, menu := range menuList {
			if menu.Permission != "" {
				ret = append(ret, menu.Permission)
			}
		}
	}
	return ret
}

func (a *apiAuth) info(c *gin.Context) {
	ret := Http.NewResult(c)
	currAdmin := common.GetCurrAdmin(c)
	var roles []string
	role := &admin.Role{}
	for _, id := range currAdmin.RoleIds {
		role.SetId(id)
		roles = append(roles, admin.LoadForumData(role).(*admin.Role).RoleKey)
	}
	permissions := a.getPermissions(c)
	ret.Success(common.ResponseSuccess, map[string]interface{}{
		"userId":       currAdmin.ID,
		"userName":     currAdmin.UserName,
		"name":         currAdmin.NickName,
		"buttons":      permissions,
		"roles":        roles,
		"avatar":       currAdmin.Avatar,
		"introduction": "I'm superman",
		"deptId":       currAdmin.DeptId,
		"permissions":  permissions,
	})
}

func (a *apiAuth) refreshToken(c *gin.Context) {
	ret := Http.NewResult(c)
	admin := common.GetCurrAdmin(c)
	// todo
	ret.Success(common.ResponseSuccess, map[string]interface{}{
		"userId":       admin.ID,
		"userName":     admin.UserName,
		"name":         admin.NickName,
		"buttons":      []string{"*:*:*"},
		"roles":        []string{"admin"},
		"avatar":       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		"introduction": "妈卖麻花儿",
		"deptId":       1,
		"permissions":  []string{"*:*:*"},
	})
}

func (a *apiAuth) logout(c *gin.Context) {
	common.UpdateToken(common.UserTypeAdmin, common.GetAdminId(c), "")
	Http.NewResult(c).Success("success", nil)
}

// 个人中心
func (a *apiAuth) profile(c *gin.Context) {
	ret := Http.NewResult(c)
	currAdmin := common.GetCurrAdmin(c)

	var roles []*admin.Role
	for _, id := range currAdmin.RoleIds {
		role := &admin.Role{}
		role.SetId(id)
		roles = append(roles, admin.LoadForumData(role).(*admin.Role))
	}

	var posts []*admin.Post
	for _, id := range currAdmin.PostIds {
		post := &admin.Post{}
		post.SetId(id)
		posts = append(posts, admin.LoadForumData(post).(*admin.Post))
	}

	ret.Success(common.ResponseSuccess, map[string]interface{}{
		"posts": posts,
		"roles": roles,
		"user":  currAdmin,
	})
}

// 重置密码
func (a *apiAuth) setPwd(c *gin.Context) {
	ret := Http.NewResult(c)
	requestData := struct {
		OldPassword string `form:"oldPassword" json:"oldPassword" binding:"required"`
		NewPassword string `form:"newPassword" json:"newPassword" binding:"required"`
	}{}

	if Http.Bind(c, &requestData) {
		currAdmin := common.GetCurrAdmin(c)
		if currAdmin.PassWord != requestData.OldPassword {
			ret.Fail("修改失败:旧密码错误", nil)
		} else {
			currAdmin.PassWord = requestData.NewPassword
			admin.UpdateForumData(currAdmin)
			// 清理当前玩家token
			common.UpdateToken(common.UserTypeAdmin, currAdmin.ID, "")
			ret.Success("修改成功:请重新登录", nil)
		}
	}
}
