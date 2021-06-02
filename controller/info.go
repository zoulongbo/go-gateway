package controller

import (
	"encoding/json"
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dao"
	"github.com/zoulongbo/go-gateway/dto/admin"
	"github.com/zoulongbo/go-gateway/middleware"
	"github.com/zoulongbo/go-gateway/public"
)

type InfoController struct {
}

func RegisterInfoController(group *gin.RouterGroup) {
	info := &InfoController{}
	group.GET("/info", info.Info)
	group.POST("/change_pwd", info.ChangePwd)
}

//Info godoc
//@Summary 登录信息
//@Description 登录信息
//@Tags 后台接口
//@ID /admin/info
//@Accept json
//@Produce json
//@Success 200 {object} middleware.Response{data=admin.InfoOutput} "success"
//@Router /admin/info [get]
func (i *InfoController) Info(c *gin.Context) {
	sess := sessions.Default(c)
	sessionInfo := sess.Get(public.AdminLoginSessionKey)
	adminSessionInfo := &admin.LoginSessionInfo{}
	err := json.Unmarshal([]byte(fmt.Sprint(sessionInfo)), adminSessionInfo)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
	}
	out := admin.InfoOutput{
		Id:           adminSessionInfo.Id,
		Name:         adminSessionInfo.Username,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "",
		Introduction: "",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(c, out)
}

//ChangePwd godoc
//@Summary 修改密码
//@Description 修改密码
//@Tags 后台接口
//@ID /admin/change_pwd
//@Accept json
//@Produce json
//@Param body body admin.ChangePwdInput true "body"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /admin/change_pwd [post]
func (i *InfoController) ChangePwd(c *gin.Context) {
	params := &admin.ChangePwdInput{}
	err := params.BindValidParam(c)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	sess := sessions.Default(c)
	sessionInfo := sess.Get(public.AdminLoginSessionKey)
	adminSessionInfo := &admin.LoginSessionInfo{}
	err = json.Unmarshal([]byte(fmt.Sprint(sessionInfo)), adminSessionInfo)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
	}
	adminInfo := &dao.Admin{}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	adminInfo, err = adminInfo.Find(c, tx, &dao.Admin{Username: adminSessionInfo.Username})
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	saltPwd := public.GenSaltPassword(adminInfo.Salt, params.Password)
	adminInfo.Password = saltPwd
	err = adminInfo.Save(c, tx)
	if err != nil {
		middleware.ResponseError(c, 2004, err)
		return
	}

	middleware.ResponseSuccess(c, "")
}
