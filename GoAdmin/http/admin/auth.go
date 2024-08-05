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

type loginRet struct {
	Token string `json:"token"`
}

func newToken(id int64, roleId int64, userName string) string {
	token, _ := common.NewToken(map[string]interface{}{
		common.TokenKeyId:            id,
		common.TokenKeyRoleId:        roleId,
		common.TokenKeyAdminUserName: userName,
	}, 24*7)
	return token
}

// @Tags     好奇宝宝后台-鉴权相关
// @Summary  管理员登录
// @Param    username  formData  string  true  "账号"
// @Param    password  formData  string  true  "密码"
// @Param    code      formData  string  true  "验证码"
// @Param    uuid      formData  string  true  "验证码id"
// @Success  200          {object}  Http.WebResult{data=loginRet}
// @Failure  500          {object}  Http.WebResult
// @Router   /api/auth/login [post]
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
						token := newToken(manager.ID, manager.RoleId, requestData.UserName)
						a.loginRet(c, requestData.UserName, common.ResponseSuccess, &loginRet{Token: token})
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

type adminInfo struct {
	UserId       int64    `json:"userId"`       // 用户id
	UserName     string   `json:"userName"`     // 用户名
	Name         string   `json:"name"`         // 昵称
	Buttons      []string `json:"buttons"`      // 权限
	Roles        []string `json:"roles"`        // 角色
	Avatar       string   `json:"avatar"`       // 头像
	Introduction string   `json:"introduction"` // 介绍
	DeptId       int64    `json:"deptId"`       // 部门id
	Permissions  []string `json:"permissions"`  // 权限
}

// @Tags     好奇宝宝后台-鉴权相关
// @Summary  获取用户信息
// @Param    token  header    string  true  "token"
// @Success  200    {object}  Http.WebResult{data=adminInfo}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/auth/info [get]
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
	ret.Success(common.ResponseSuccess, &adminInfo{
		UserId:       currAdmin.ID,
		UserName:     currAdmin.UserName,
		Name:         currAdmin.NickName,
		Buttons:      permissions,
		Roles:        roles,
		Avatar:       currAdmin.Avatar,
		Introduction: "I'm superman",
		DeptId:       currAdmin.DeptId,
		Permissions:  permissions,
	})
}

// @Tags     好奇宝宝后台-鉴权相关
// @Summary  刷新token
// @Param    token  header    string  true  "token"
// @Success  200    {object}  Http.WebResult{data=loginRet}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/auth/refresh_token [post]
func (a *apiAuth) refreshToken(c *gin.Context) {
	ret := Http.NewResult(c)
	admin := common.GetCurrAdmin(c)
	token := newToken(admin.ID, admin.RoleId, admin.UserName)
	ret.Success(common.ResponseSuccess, &loginRet{Token: token})
	common.UpdateToken(common.UserTypeAdmin, admin.ID, token)
}

// @Tags     好奇宝宝后台-鉴权相关
// @Summary  管理员登出
// @Param    token  header    string  true  "token"
// @Success  200    {object}  Http.WebResult
// @Failure  500    {object}  Http.WebResult
// @Router   /api/auth/logout [post]
func (a *apiAuth) logout(c *gin.Context) {
	common.UpdateToken(common.UserTypeAdmin, common.GetAdminId(c), "")
	Http.NewResult(c).Success("success", nil)
}

type profileInfo struct {
	Posts []*admin.Post `json:"posts"` // 职位列表
	Roles []*admin.Role `json:"roles"` // 角色列表
	User  *admin.User   `json:"user"`  // 用户详情
}

// @Tags     好奇宝宝后台-鉴权相关
// @Summary  个人中心
// @Param    token  header    string  true  "token"
// @Success  200    {object}  Http.WebResult{data=profileInfo}
// @Failure  500    {object}  Http.WebResult
// @Router   /api/auth/profile [get]
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

	ret.Success(common.ResponseSuccess, &profileInfo{
		Posts: posts,
		Roles: roles,
		User:  currAdmin,
	})
}

// @Tags     好奇宝宝后台-鉴权相关
// @Summary  管理员重置密码
// @Param    token        header    string  true  "token"
// @Param    oldPassword  formData  string  true  "旧密码"
// @Param    newPassword  formData  string  true  "新密码"
// @Success  200       {object}  Http.WebResult
// @Failure  500       {object}  Http.WebResult
// @Router   /api/auth/setpwd [put]
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
