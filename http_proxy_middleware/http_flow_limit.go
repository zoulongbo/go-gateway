package http_proxy_middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dao"
	"github.com/zoulongbo/go-gateway/middleware"
	"github.com/zoulongbo/go-gateway/public"
)

func HttpFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 50001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := service.(*dao.ServiceDetail)
		if serviceDetail.AccessControl.ServiceFlowLimit > 0 {
			qps := float64(serviceDetail.AccessControl.ServiceFlowLimit)
			serviceLimiter := public.FlowLimiterHandler.GetLimiter(public.FlowServicePrefix+serviceDetail.Info.ServiceName, qps)
			if ! serviceLimiter.Allow() {
				middleware.ResponseError(c, 50002, errors.New(fmt.Sprintf("service flow limit %v", qps)))
				c.Abort()
				return
			}
		}
		if serviceDetail.AccessControl.ClientipFlowLimit > 0 {
			clientQps := float64(serviceDetail.AccessControl.ClientipFlowLimit)
			clientLimiter := public.FlowLimiterHandler.GetLimiter(public.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+c.ClientIP(), clientQps)
			if ! clientLimiter.Allow() {
				middleware.ResponseError(c, 50002, errors.New(fmt.Sprintf("client flow limit %v", clientQps)))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
