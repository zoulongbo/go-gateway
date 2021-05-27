package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/public"
)

type ServiceListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" example:"" validate:""`                      //关键词
	PageNo   int    `json:"page_no" form:"page_no" comment:"页码" example:"1" validate:""`                //页码
	PageSize int    `json:"page_size" form:"page_size" comment:"页数据量" example:"20" validate:"required"` //页数据量
}

type ServiceListOutput struct {
	Total    int64                   `json:"total" form:"total" comment:"总数" example:"20"`           //总数
	PageNo   int                     `json:"page_no" form:"page_no" comment:"页码" example:"1"`        //页码
	PageSize int                     `json:"page_size" form:"page_size" comment:"页数据量" example:"20"` //页数据量
	List     []ServiceListItemOutput `json:"list" form:"list" comment:"数据列表" example:""`             //数据列表
}

type ServiceListItemOutput struct {
	Id          int64  `json:"id" form:"id" comment:"id" example:"1" validate:""`                         //id
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" example:"名称"`              //服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"描述"`              //服务描述
	LoadType    int    `json:"load_type" form:"load_type" comment:"负载类型 0=http;1=tcp;2=grpc" example:"1"` //服务类型
	ServiceAddr string `json:"service_addr" form:"service_addr" comment:"服务地址" example:"127.0.0.1"`       //服务地址
	Qps         int64  `json:"qps" form:"qps" comment:"qps" example:"1"`                                  //qps
	Qpd         int64  `json:"qpd" form:"qpd" comment:"qpd" example:"1"`                                  //qpd
	TotalNode   int    `json:"total_node" form:"total_node" comment:"节点数" example:"1"`                    //节点数
}

func (s *ServiceListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, s)
}

type ServiceDeleteInput struct {
	Id int64 `json:"id" form:"id" comment:"服务id" example:"1" validate:"required"` //服务id
}

func (s *ServiceDeleteInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, s)
}

type ServiceAddHttpInput struct {
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" example:"i_am_test1" validate:"required,valid_service_name"` //服务名
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"test1的desc" validate:"required,max=255,min=1"`      //服务描述

	RuleType       int    `json:"rule_type" form:"rule_type" comment:"匹配类型" example:"1" validate:"max=1,min=0"`                                   //接入类型
	Rule           string `json:"rule" form:"rule" comment:"t域名或者前缀"  example:"/test_http_service" validate:"required,valid_rule"`                //域名或者前缀
	NeedHttps      int    `json:"need_https" form:"need_https" comment:"是否支持https" example:"1" validate:"max=1,min=0"`                            //支持https
	NeedStripUri   int    `json:"need_strip_uri" form:"need_strip_uri" comment:"是否启用strip_uri" example:"1" validate:"max=1,min=0"`                //启用strip_uri
	NeedWebsocket  int    `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket" example:"1" validate:"max=1,min=0"`                //是否支持websocket
	UrlRewrite     string `json:"url_rewrite" form:"url_rewrite" description:"url重写功能" example:"$url" validate:"valid_url_rewrite"`               //url重写功能
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"header转换支持增加" example:"agent" validate:"valid_header_transfor"` //header转换

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限" example:"1" validate:"max=1,min=0"`                 //关键词
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"127.0.0.1" validate:""`                   //黑名单ip
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单ip" example:"127.0.0.1" validate:""`                   //白名单ip
	ClientipFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端ip限流" example:"1" validate:"min=0"` //客户端ip限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" example:"1" validate:"min=0"`      //服务端限流

	RoundType              int    `json:"round_type" form:"round_type" comment:"轮询方式" example:"1" validate:"max=3,min=0"`                                //轮询方式
	IpList                 string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"127.0.0.1:80" validate:"required,valid_ipportlist"`             //ip列表
	WeightList             string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"1,2,3" validate:"required,valid_weightlist"`           //权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时, 单位s" example:"1" validate:"min=0"`   //建立连接超时, 单位s
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时, 单位s" example:"1" validate:"min=0"` //获取header超时, 单位s
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"链接最大空闲时间, 单位s" example:"20" validate:"min=0"`      //链接最大空闲时间, 单位s
	UpstreamMaxIdle        int    `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大空闲链接数" example:"5" validate:"min=0"`                     //最大空闲链接数
}

type ServiceUpdateHttpInput struct {
	Id int64 `json:"id" form:"id" comment:"服务Id" example:"62" validate:"required,min=1"` //服务id
	ServiceAddHttpInput
}

func (s *ServiceAddHttpInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, s)
}

type ServiceStatOutput struct {
	Today     []int64 `json:"today" form:"today" comment:"当日统计数据" example:""`
	Yesterday []int64 `json:"yesterday" form:"yesterday" comment:"前一天统计数据" example:""`
}

type ServiceAddGRPCInput struct {
	ServiceName       string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`
	Port              int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`
	HeaderTransfor    string `json:"header_transfor" form:"header_transfor" comment:"metadata转换" validate:"valid_header_transfor"`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`
	IpList            string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"`
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`
}

func (params *ServiceAddGRPCInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ServiceUpdateGRPCInput struct {
	Id int64 `json:"id" form:"id" comment:"服务Id" example:"62" validate:"required,min=1"` //服务id
	ServiceAddGRPCInput
}
type ServiceAddTcpInput struct {
	ServiceName       string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`
	Port              int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`
	HeaderTransfor    string `json:"header_transfor" form:"header_transfor" comment:"header头转换" validate:""`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`
	IpList            string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"`
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`
}

type ServiceUpdateTcpInput struct {
	Id int64 `json:"id" form:"id" comment:"服务Id" example:"62" validate:"required,min=1"` //服务id
	ServiceAddTcpInput
}

func (params *ServiceAddTcpInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
