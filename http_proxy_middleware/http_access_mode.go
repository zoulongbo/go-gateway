package http_proxy_middleware


import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dao"
	"github.com/zoulongbo/go-gateway/middleware"
	"github.com/zoulongbo/go-gateway/public"
)

//拼配接入方式 基于请求信息
func HttpAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := dao.ServiceManagerHandle.HttpAccessMode(c)
		if err != nil {
			middleware.ResponseError(c, 10001, err)
			c.Abort()
			return
		}

		fmt.Println("marched service", public.Obj2Json(serviceDetail))
		c.Set("service", serviceDetail)
		c.Next()
	}
}

