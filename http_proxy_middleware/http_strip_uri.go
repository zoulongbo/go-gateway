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

func HttpStripUriMiddleware() gin.HandlerFunc  {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 20001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := service.(*dao.ServiceDetail)
		if serviceDetail.HttpRule.RuleType == public.HTTPRuleTypePrefixURL && serviceDetail.HttpRule.NeedStripUri == public.HTTPRuleNeedHttps {
			fmt.Println("replace before:", c.Request.URL.Path)
			c.Request.URL.Path = strings.Replace(c.Request.URL.Path, serviceDetail.HttpRule.Rule, "", 1)
			fmt.Println("replace after:", c.Request.URL.Path)
		}
		c.Next()
	}
}
