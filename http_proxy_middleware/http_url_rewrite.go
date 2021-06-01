package http_proxy_middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dao"
	"github.com/zoulongbo/go-gateway/middleware"
	"regexp"
	"strings"
)

func HttpUrlWriteMiddleware() gin.HandlerFunc  {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 20001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := service.(*dao.ServiceDetail)
		for _, item := range strings.Split(serviceDetail.HttpRule.UrlRewrite, ",") {
			items := strings.Split(item, " ")
			if len(items) != 2 {
				continue
			}
			reg, err := regexp.Compile(items[0])
			if err != nil {
				fmt.Println("regexp.Compile err:", err)
				continue
			}
			replacePath := reg.ReplaceAllString(c.Request.URL.Path, items[1])
			c.Request.URL.Path = replacePath

		}
		c.Next()
	}
}
