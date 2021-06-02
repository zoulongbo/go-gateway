package http_proxy_router

import (
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/controller"
	"github.com/zoulongbo/go-gateway/http_proxy_middleware"
	"github.com/zoulongbo/go-gateway/middleware"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.New()
	//todo router := gin.Default()    new不会打印日志  default会 有io消耗
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//
	oauth := router.Group("/oauth")
	oauth.Use(
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.TranslationMiddleware(),
	)
	{
		controller.RegisterOAuthController(oauth)
	}

	//根
	router.Use(
		http_proxy_middleware.HttpAccessModeMiddleware(),
		http_proxy_middleware.HttpFlowCountMiddleware(),
		http_proxy_middleware.HttpFlowLimitMiddleware(),

		http_proxy_middleware.HttpJwtAuthMiddleware(),

		http_proxy_middleware.HttpBlackListMiddleware(),
		http_proxy_middleware.HttpWhiteListMiddleware(),
		http_proxy_middleware.HttpHeaderTransferMiddleware(),
		http_proxy_middleware.HttpStripUriMiddleware(),
		http_proxy_middleware.HttpUrlWriteMiddleware(),
		http_proxy_middleware.HttpReverseProxyMiddleware(),
		)

	return router
}
