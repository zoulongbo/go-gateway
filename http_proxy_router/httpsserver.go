package http_proxy_router

import (
	"context"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/cert_file"
	"github.com/zoulongbo/go-gateway/middleware"
	"log"
	"net/http"
	"time"
)

var (
	HttpsSrvHandler *http.Server
)

func HttpsServerRun() {
	gin.SetMode(lib.ConfBase.DebugMode)
	r := InitRouter(middleware.RecoveryMiddleware(), middleware.RequestLog())
	HttpsSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.https.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.https.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.https.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.https.max_header_bytes")),
	}
	log.Printf(" [INFO] HttpProxyRun:%s\n", lib.GetStringConf("proxy.https.addr"))
	if err := HttpsSrvHandler.ListenAndServeTLS(cert_file.Path("server.crt"), cert_file.Path("server.key")); err != nil {
		log.Fatalf(" [ERROR] HttpProxyRun:%s err:%v\n", lib.GetStringConf("proxy.https.addr"), err)
	}
}

func HttpsServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpsSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpsProxyStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpsProxyStop stopped\n")
}
