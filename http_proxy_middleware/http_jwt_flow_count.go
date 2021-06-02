package http_proxy_middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dao"
	"github.com/zoulongbo/go-gateway/middleware"
	"github.com/zoulongbo/go-gateway/public"
)

func HttpJwtFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		app, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := app.(*dao.App)
		//3、 租户统计
		appCounter,err := public.FlowCountHandler.GetFlowCounter(public.FlowAppPrefix+appInfo.AppId)
		if err != nil {
			middleware.ResponseError(c, 40001, err)
			c.Abort()
			return
		}
		appCounter.Increase()

		if appInfo.Qpd >0 && appCounter.TotalCount > appInfo.Qpd {
			middleware.ResponseError(c, 40002, errors.New(fmt.Sprintf("qpd limit %v - curr:%v", appInfo.Qpd, appCounter.TotalCount)))
			c.Abort()
			return
		}

		c.Next()
		return
	}
}
