package controller

import (
	"encoding/json"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dao"
	"github.com/zoulongbo/go-gateway/dto/admin"
	"github.com/zoulongbo/go-gateway/middleware"
	"github.com/zoulongbo/go-gateway/public"
	"time"
)

type LoginController struct {
}

func RegisterLoginController(group *gin.RouterGroup) {
	login := &LoginController{}
	group.POST("/login", login.Login)
	group.GET("/logout", login.Logout)
}

//Login godoc
//@Summary 后台登录
//@Description 后台登录
//@Tags 后台接口
//@ID /admin/login
//@Accept json
//@Produce json
//@Param body body admin.LoginInput true "body"
//@Success 200 {object} middleware.Response{data=admin.LoginOutput} "success"
//@Router /admin/login [post]
func (login *LoginController) Login(c *gin.Context) {
	params := &admin.LoginInput{}
	err := params.BindValidParam(c)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	adminInfo := &dao.Admin{}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//登录
	adminInfo, err = adminInfo.LoginCheck(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	//session
	sess := sessions.Default(c)
	sessionInfo := &admin.LoginSessionInfo{
		Id:        adminInfo.Id,
		Username:  adminInfo.Username,
		LoginTime: time.Now(),
	}
	sessionBts, err := json.Marshal(sessionInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	sess.Set(public.AdminLoginSessionKey, string(sessionBts))
	sess.Save()

	out := &admin.LoginOutput{Token: adminInfo.Username}
	middleware.ResponseSuccess(c, out)
}

//Logout godoc
//@Summary 退出登录
//@Description 退出登录
//@Tags 后台接口
//@ID /admin/logout
//@Accept json
//@Produce json
//@Success 200 {object} middleware.Response{} "success"
//@Router /admin/logout [get]
func (login *LoginController) Logout(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Delete(public.AdminLoginSessionKey)
	sess.Save()
	middleware.ResponseSuccess(c, "")
}
