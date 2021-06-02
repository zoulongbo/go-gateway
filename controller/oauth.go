package controller

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dao"
	"github.com/zoulongbo/go-gateway/dto"
	"github.com/zoulongbo/go-gateway/middleware"
	"github.com/zoulongbo/go-gateway/public"
	"strings"
	"time"
)

type OAuthController struct {
}

func RegisterOAuthController(group *gin.RouterGroup) {
	oauth := &OAuthController{}
	group.POST("/tokens", oauth.Tokens)
}

//Tokens godoc
//@Summary 获取token
//@Description 获取token
//@Tags OAUTH
//@ID /oauth/tokens
//@Accept json
//@Produce json
//@Param body body dto.OAuthInput true "body"
//@Success 200 {object} middleware.Response{data=dto.OAuthOutput} "success"
//@Router /oauth/tokens [post]
func (oauth *OAuthController) Tokens(c *gin.Context) {
	params := &dto.OAuthInput{}
	err := params.BindValidParam(c)
	fmt.Println(params)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	var author []string
	if c.GetHeader("Authorization") != "" {
		author = strings.Split(c.GetHeader("Authorization"), " ")
	}
	if len(author) != 2 {
		middleware.ResponseError(c, 2001, errors.New("参数欠缺,认证失败"))
		return
	}
	authorByte, err := base64.StdEncoding.DecodeString(author[1])
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	authorArr := strings.Split(string(authorByte), ":")
	if len(authorArr) != 2 {
		middleware.ResponseError(c, 2003, errors.New("账号密码必传"))
		return
	}
	app := &dao.App{
		AppId:  authorArr[0],
		Secret: authorArr[1],
	}
	app, err = dao.AppHandler.GetApp(app)
	if err != nil {
		middleware.ResponseError(c, 2004, err)
		return
	}
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(public.JwtExpires * time.Second).In(lib.TimeLocation).Unix(),
		Issuer:    app.AppId,
		NotBefore: 0,
		Subject:   "",
	}
	token, err := public.JwtEncode(claims)
	if err != nil {
		middleware.ResponseError(c, 2005, err)
		return
	}

	out := &dto.OAuthOutput{
		AccessToken: token,
		ExpireIn:    public.JwtExpires,
		TokenType:   public.TokenType,
		Scope:       "write_read",
	}
	middleware.ResponseSuccess(c, out)
}
