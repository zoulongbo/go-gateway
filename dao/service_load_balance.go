package dao

import (
	"fmt"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/public"
	"github.com/zoulongbo/go-gateway/reverse_proxy/load_balance"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ServiceLoadBalance struct {
	Id            int64  `json:"id" gorm:"primary_key" description:"自增主键"`
	ServiceId     int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	CheckMethod   int    `json:"check_method" gorm:"column:check_method" description:"检查方法 tcpchk=检测端口是否握手成功	"`
	CheckTimeout  int    `json:"check_timeout" gorm:"column:check_timeout" description:"check超时时间	"`
	CheckInterval int    `json:"check_interval" gorm:"column:check_interval" description:"检查间隔, 单位s		"`
	RoundType     int    `json:"round_type" gorm:"column:round_type" description:"轮询方式 round/weight_round/random/ip_hash"`
	IpList        string `json:"ip_list" gorm:"column:ip_list" description:"ip列表"`
	WeightList    string `json:"weight_list" gorm:"column:weight_list" description:"权重列表"`
	ForbidList    string `json:"forbid_list" gorm:"column:forbid_list" description:"禁用ip列表"`

	UpstreamConnectTimeout int `json:"upstream_connect_timeout" gorm:"column:upstream_connect_timeout" description:"下游建立连接超时, 单位s"`
	UpstreamHeaderTimeout  int `json:"upstream_header_timeout" gorm:"column:upstream_header_timeout" description:"下游获取header超时, 单位s	"`
	UpstreamIdleTimeout    int `json:"upstream_idle_timeout" gorm:"column:upstream_idle_timeout" description:"下游链接最大空闲时间, 单位s	"`
	UpstreamMaxIdle        int `json:"upstream_max_idle" gorm:"column:upstream_max_idle" description:"下游最大空闲链接数"`
}

func (s *ServiceLoadBalance) TableName() string {
	return "gateway_service_load_balance"
}

func (s *ServiceLoadBalance) Find(c *gin.Context, tx *gorm.DB, search *ServiceLoadBalance) (*ServiceLoadBalance, error) {
	result := &ServiceLoadBalance{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ServiceLoadBalance) GetIpListByModel() []string {
	return strings.Split(s.IpList, ",")
}

func (s *ServiceLoadBalance) GetWeightListByModel() []string {
	return strings.Split(s.WeightList, ",")
}

func (s *ServiceLoadBalance) Save(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(public.GetGinTraceContext(c)).Save(s).Error
	if err != nil {
		return err
	}
	return nil
}

//复杂均衡策略
func init()  {
	LoadBalanceHandler = NewLoadBalanceManager()
	TransporterHandler = NewTransporter()
}
var LoadBalanceHandler  *LoadBalanceManager

type LoadBalanceManager struct {
	LoadBalanceMap map[string]LoadBalanceItem
	LoadBalanceSlice []LoadBalanceItem
	Locker sync.RWMutex
}

type LoadBalanceItem struct {
	LoadBalance load_balance.LoadBalance
	ServiceName string
}

func NewLoadBalanceManager() *LoadBalanceManager  {
	return &LoadBalanceManager{
		LoadBalanceMap:   map[string]LoadBalanceItem{},
		LoadBalanceSlice: []LoadBalanceItem{},
		Locker:           sync.RWMutex{},
	}
}

func (l *LoadBalanceManager) GetLoadBalance (service *ServiceDetail) (load_balance.LoadBalance, error) {
	for  _, loadItem := range l.LoadBalanceSlice{
		if loadItem.ServiceName == service.Info.ServiceName {
			return loadItem.LoadBalance, nil
		}
	}
	schema := "http"
	if service.HttpRule.NeedHttps == public.HTTPRuleNeedHttps {
		schema = "https"
	}
	ipList := service.LoadBalance.GetIpListByModel()
	weightList := service.LoadBalance.GetWeightListByModel()
	ipConf := map[string]string{}
	for index, ip := range ipList {
		ipConf[ip] = weightList[index]
	}

	mConf ,err := load_balance.NewLoadBalanceCheckConf(fmt.Sprintf("%s://%s", schema, "%s"), ipConf)
	if err != nil {
		return nil, err
	}
	loadBalance := load_balance.LoadBalanceFactorWithConf(load_balance.LbType(service.LoadBalance.RoundType), mConf)
	loadBalanceItem := LoadBalanceItem{
		LoadBalance: loadBalance,
		ServiceName: service.Info.ServiceName,
	}
	l.LoadBalanceSlice = append(l.LoadBalanceSlice, loadBalanceItem)

	l.Locker.Lock()
	defer l.Locker.Unlock()

	l.LoadBalanceMap[service.Info.ServiceName] = loadBalanceItem
	return loadBalance, nil
}

var TransporterHandler * Transporter

type Transporter struct {
	TransporterMap map[string]*TransporterItem
	TransporterSlice []*TransporterItem
	Locker sync.RWMutex
}

type TransporterItem struct {
	Trans *http.Transport
	ServiceName string
}

func NewTransporter() *Transporter  {
	return &Transporter{
		TransporterMap:   map[string]*TransporterItem{},
		TransporterSlice: []*TransporterItem{},
		Locker:           sync.RWMutex{},
	}
}

func (t *Transporter) GetTransporter (service *ServiceDetail) (*http.Transport, error) {
	for  _, transItem := range t.TransporterSlice{
		if transItem.ServiceName == service.Info.ServiceName {
			return transItem.Trans, nil
		}
	}

	trans := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(service.LoadBalance.UpstreamConnectTimeout) * time.Second,
		}).DialContext,
		MaxIdleConns:          service.LoadBalance.UpstreamMaxIdle,
		IdleConnTimeout:       time.Duration(service.LoadBalance.UpstreamIdleTimeout) * time.Second,
		ResponseHeaderTimeout: time.Duration(service.LoadBalance.UpstreamHeaderTimeout) * time.Second,
	}
	transporterItem := &TransporterItem{
		Trans:       trans,
		ServiceName: service.Info.ServiceName,
	}
	t.TransporterSlice = append(t.TransporterSlice, transporterItem)
	t.Locker.Lock()
	defer  t.Locker.Unlock()
	t.TransporterMap[service.Info.ServiceName] = transporterItem
	return trans, nil
}
