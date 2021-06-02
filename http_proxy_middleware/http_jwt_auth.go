package http_proxy_middleware

import (
	"errors"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dao"
	"github.com/zoulongbo/go-gateway/middleware"
	"github.com/zoulongbo/go-gateway/public"
	"strings"
	"time"
)

func HttpJwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 20001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := service.(*dao.ServiceDetail)
		if serviceDetail.AccessControl.OpenAuth == public.HTTPAccessControlAuth {
			var author []string
			if c.GetHeader("Authorization") != "" {
				author = strings.Split(c.GetHeader("Authorization"), " ")
			}
			if len(author) != 2 {
				middleware.ResponseError(c, 2001, errors.New("认证参数欠缺"))
				c.Abort()
				return
			}
			if author[0] != public.TokenType {
				middleware.ResponseError(c, 2001, errors.New("认证类型错误"))
				c.Abort()
				return
			}
			claims, err := public.JwtDecode(author[1])
			if err != nil {
				middleware.ResponseError(c, 2001, errors.New("认证解析失败"))
				c.Abort()
				return
			}
			//过期时间小于当前时间
			if claims.ExpiresAt < time.Now().In(lib.TimeLocation).Unix() {
				middleware.ResponseError(c, 2001, errors.New("身份认证已过期"))
				c.Abort()
				return
			}
			//if _, ok := dao.AppHandler.AppMap[claims.Issuer]; !ok {
			//	middleware.ResponseError(c, 2001, errors.New("认证信息不存在"))
			//	c.Abort()
			//	return
			//}
			app, err := dao.AppHandler.GetApp(&dao.App{AppId:claims.Issuer})
			if err != nil {
				middleware.ResponseError(c, 2001, errors.New("认证信息不存在"))
				c.Abort()
				return
			}
			c.Set("appDetail", app)
		}
		c.Next()
	}
}
