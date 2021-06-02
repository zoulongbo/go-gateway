package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/public"
)

type OAuthInput struct {
	GrantType string `json:"grant_type" form:"grant_type" comment:"授权类型" example:"client_credentials" validate:"required"`  //授权类型
	Scope string `json:"scope" form:"scope" comment:"权限" example:"write_read" validate:"required"` //权限
}

func (oauth *OAuthInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, oauth)
}

type OAuthOutput struct {
	AccessToken string `json:"access_token" form:"access_token"` //access_token
	ExpireIn int `json:"expire_in" form:"expire_in"` //过期时间
	TokenType string `json:"token_type" form:"token_type"` //token类型
	Scope string `json:"scope" form:"scope"` //权限
}