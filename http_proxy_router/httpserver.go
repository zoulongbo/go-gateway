package http_proxy_router

import (
	"context"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/middleware"
	"log"
	"net/http"
	"time"
)

var (
	HttpSrvHandler *http.Server
)

func HttpServerRun() {
	gin.SetMode(lib.ConfBase.DebugMode)
	r := InitRouter(middleware.RecoveryMiddleware(), middleware.RequestLog())
	HttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.http.max_header_bytes")),
	}
	log.Printf(" [INFO] HttpProxyRun:%s\n", lib.GetStringConf("proxy.http.addr"))
	if err := HttpSrvHandler.ListenAndServe(); err != nil {
		log.Fatalf(" [ERROR] HttpProxyRun:%s err:%v\n", lib.GetStringConf("proxy.http.addr"), err)
	}
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpProxyStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpProxyStop stopped\n")
}