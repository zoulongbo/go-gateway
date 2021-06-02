package http_proxy_middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dao"
	"github.com/zoulongbo/go-gateway/middleware"
	"github.com/zoulongbo/go-gateway/public"
)

func HttpFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 20001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := service.(*dao.ServiceDetail)
		//1、 全站统计
		//2、 服务统计
		//3、 租户统计
		totalCounter,err := public.FlowCountHandler.GetFlowCounter(public.FlowTotal)
		serviceCounter,err := public.FlowCountHandler.GetFlowCounter(public.FlowServicePrefix+serviceDetail.Info.ServiceName)
		if err != nil {
			middleware.ResponseError(c, 40002, err)
			c.Abort()
			return
		}
		totalCounter.Increase()
		serviceCounter.Increase()
		c.Next()
	}
}
