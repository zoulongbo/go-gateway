package dao

import (
	"errors"
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dto/admin"
	"github.com/zoulongbo/go-gateway/public"
	"net/http/httptest"
	"strings"
	"sync"
)

type ServiceDetail struct {
	Info          *ServiceInfo          `json:"info" description:"服务信息"`
	HttpRule      *ServiceHttpRule      `json:"http" description:"http 规则"`
	TcpRule       *ServiceTcpRule       `json:"tcp" description:"tcp 规则"`
	GRPCRule      *ServiceGRPCRule      `json:"grpc" description:"grpc 规则"`
	LoadBalance   *ServiceLoadBalance   `json:"load_balance" description:"负载"`
	AccessControl *ServiceAccessControl `json:"access_control" description:"access_control"`
}

var ServiceManagerHandle *ServiceManager

//自动加载方法
func init()  {
	ServiceManagerHandle = NewServiceManager()
}

type ServiceManager struct {
	ServiceMap   map[string]*ServiceDetail
	ServiceSlice []*ServiceDetail
	Locker       sync.RWMutex
	init         sync.Once
	err          error
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		ServiceMap:   map[string]*ServiceDetail{},
		ServiceSlice: []*ServiceDetail{},
		Locker:       sync.RWMutex{},
		init:         sync.Once{},
		err:          nil,
	}
}

func (sm *ServiceManager) HttpAccessMode(c *gin.Context) (*ServiceDetail, error) {
	//1、前缀匹配      2、域名匹配
	host := c.Request.Host
	host = host[0:strings.Index(host, ":")]
	fmt.Println("host: ", host)
	path := c.Request.URL.Path

	for _, item :=range sm.ServiceSlice {
		if item.Info.LoadType != public.LoadTypeHTTP {
			continue
		}
		//域名
		if item.HttpRule.RuleType == public.HTTPRuleTypeDomain {
			if item.HttpRule.Rule == host {
				return item, nil
			}
		}
		//前缀
		if item.HttpRule.RuleType == public.HTTPRuleTypePrefixURL {
			if strings.HasPrefix(path, item.HttpRule.Rule) {
				return item, nil
			}
		}
	}
	return nil, errors.New("not matched service")
}


func (sm *ServiceManager) LoadOnce() error  {
	sm.init.Do(func() {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		tx, err := lib.GetGormPool("default")
		if err != nil {
			sm.err = err
			return
		}
		serviceInfo := &ServiceInfo{}
		params := &admin.ServiceListInput{PageNo:1, PageSize:999999}
		list, _, err := serviceInfo.PageList(c, tx, params)
		if err != nil {
			sm.err = err
			return
		}
		//map读取时写入可能内存溢出  加锁
		sm.Locker.Lock()
		defer sm.Locker.Unlock()

		for _, listItem := range list {
			tmpItem := listItem
			serviceDetail, err := tmpItem.ServiceDetail(c, tx, &tmpItem)
			if err != nil {
				sm.err = err
				return
			}
			sm.ServiceMap[listItem.ServiceName] = serviceDetail
			sm.ServiceSlice = append(sm.ServiceSlice, serviceDetail)
		}
	})
	return sm.err
}
