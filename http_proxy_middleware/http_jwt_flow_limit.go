package http_proxy_middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dao"
	"github.com/zoulongbo/go-gateway/middleware"
	"github.com/zoulongbo/go-gateway/public"
)

func HttpJwtFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		app, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := app.(*dao.App)
		if appInfo.Qps > 0 {
			clientLimiter := public.FlowLimiterHandler.GetLimiter(public.FlowAppPrefix+appInfo.AppId+"_"+c.ClientIP(), float64(appInfo.Qps))
			if !clientLimiter.Allow() {
				middleware.ResponseError(c, 5002, errors.New(fmt.Sprintf("%v flow limit %v", c.ClientIP(), appInfo.Qps), ))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
