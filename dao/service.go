package dao

type ServiceDetail struct {
	Info          *ServiceInfo          `json:"info" description:"服务信息"`
	HttpRule      *ServiceHttpRule      `json:"http" description:"http 规则"`
	TcpRule       *ServiceTcpRule       `json:"tcp" description:"tcp 规则"`
	GRPCRule      *ServiceGRPCRule      `json:"grpc" description:"grpc 规则"`
	LoadBalance   *ServiceLoadBalance   `json:"load_balance" description:"负载"`
	AccessControl *ServiceAccessControl `json:"access_control" description:"access_control"`
}
