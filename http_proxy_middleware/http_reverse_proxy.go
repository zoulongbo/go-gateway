package http_proxy_middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dao"
	"github.com/zoulongbo/go-gateway/middleware"
	"github.com/zoulongbo/go-gateway/reverse_proxy"
)

//拼配接入方式 基于请求信息
func HttpReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 20001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := service.(*dao.ServiceDetail)
		lb, err := dao.LoadBalanceHandler.GetLoadBalance(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 20002, err)
			c.Abort()
			return
		}
		trans, err := dao.TransporterHandler.GetTransporter(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 20003, err)
			c.Abort()
			return
		}

		proxy := reverse_proxy.NewLoadBalanceReverseProxy(c, lb, trans)
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
		return
	}
}
