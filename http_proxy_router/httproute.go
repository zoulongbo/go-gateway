package http_proxy_router

import (
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/http_proxy_middleware"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Use(
		http_proxy_middleware.HttpAccessModeMiddleware(),
		http_proxy_middleware.HttpBlackListMiddleware(),
		http_proxy_middleware.HttpWhiteListMiddleware(),
		http_proxy_middleware.HttpHeaderTransferMiddleware(),
		http_proxy_middleware.HttpStripUriMiddleware(),
		http_proxy_middleware.HttpUrlWriteMiddleware(),
		http_proxy_middleware.HttpReverseProxyMiddleware(),
		)
	return router
}
