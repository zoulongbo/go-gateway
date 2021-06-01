package http_proxy_middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dao"
	"github.com/zoulongbo/go-gateway/middleware"
	"github.com/zoulongbo/go-gateway/public"
	"strings"
)

func HttpWhiteListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 20001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := service.(*dao.ServiceDetail)
		var ipList []string
		if serviceDetail.AccessControl.WhiteList != "" {
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		if serviceDetail.AccessControl.OpenAuth == public.HTTPAccessControlAuth && len(ipList) > 0 {
			if ! public.InStringSlice(ipList, c.ClientIP()) {
				middleware.ResponseError(c, 30001, errors.New(fmt.Sprintf("client ip %s not in white iplist", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
