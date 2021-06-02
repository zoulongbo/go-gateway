package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/public"
	"time"
)

type LoginInput struct {
	Username string `json:"username" form:"username" comment:"姓名" example:"admin" validate:"required"`  //用户名
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"` //密码
}

func (params *LoginInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type LoginOutput struct {
	Token string `json:"token" form:"token" comment:"token" example:"token" validate:""`
}

type LoginSessionInfo struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	LoginTime time.Time `json:"login_time"`
}
