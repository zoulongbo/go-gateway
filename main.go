package main

import (
	"flag"
	"github.com/e421083458/golang_common/lib"
	"github.com/zoulongbo/go-gateway/dao"
	"github.com/zoulongbo/go-gateway/http_proxy_router"
	"github.com/zoulongbo/go-gateway/router"
	"os"
	"os/signal"
	"syscall"
)

var (
	endpoint = flag.String("endpoint", "", "input endpoint like dashboard or server")
	config = flag.String("config", "", "input config like ./conf/dev/")
)
func main() {
	flag.Parse()
	if *endpoint == "" || *config == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *endpoint == "dashboard" {
		lib.InitModule(*config, []string{"base", "mysql", "redis"})
		defer lib.Destroy()
		router.HttpServerRun()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		router.HttpServerStop()
	} else {
		lib.InitModule(*config, []string{"base", "mysql", "redis"})
		defer lib.Destroy()

		//服务数据初始化
		dao.ServiceManagerHandle.LoadOnce()

		//http代理服务
		go func() {
			http_proxy_router.HttpServerRun()
		}()
		//https代理服务
		go func() {
			http_proxy_router.HttpsServerRun()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		http_proxy_router.HttpServerStop()
		http_proxy_router.HttpsServerStop()
	}
}
