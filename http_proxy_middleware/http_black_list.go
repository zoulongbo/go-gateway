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

func HttpBlackListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 20001, errors.New("service not found"))
			c.Abort()
			return
		}
		var blackIpList, whiteIpList []string
		serviceDetail := service.(*dao.ServiceDetail)
		if serviceDetail.AccessControl.BlackList != "" {
			blackIpList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}
		if serviceDetail.AccessControl.WhiteList != "" {
			whiteIpList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		if serviceDetail.AccessControl.OpenAuth == public.HTTPAccessControlAuth && len(blackIpList) > 0 && len(whiteIpList) == 0 {
			if public.InStringSlice(blackIpList, c.ClientIP()) {
				middleware.ResponseError(c, 30001, errors.New(fmt.Sprintf("client ip %s in black iplist", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
