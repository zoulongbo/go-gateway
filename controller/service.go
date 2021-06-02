package controller

import (
	"errors"
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dao"
	"github.com/zoulongbo/go-gateway/dto/admin"
	"github.com/zoulongbo/go-gateway/middleware"
	"github.com/zoulongbo/go-gateway/public"
	"strconv"
	"strings"
	"time"
)

type ServiceController struct {
}

func RegisterServiceController(g *gin.RouterGroup) {
	s := ServiceController{}
	g.GET("/service_list", s.ServiceList)
	g.GET("/service_delete", s.ServiceDelete)
	g.POST("/service_add_http", s.ServiceAddHttp)
	g.POST("/service_update_http", s.ServiceUpdateHttp)
	g.GET("/service_detail", s.ServiceDetail)
	g.GET("/service_stat", s.ServiceStat)

	g.POST("/service_add_tcp", s.ServiceAddTcp)
	g.POST("/service_update_tcp", s.ServiceUpdateTcp)

	g.POST("/service_add_grpc", s.ServiceAddGRPC)
	g.POST("/service_update_grpc", s.ServiceUpdateGRPC)
}

//ServiceList godoc
//@Summary 服务列表
//@Description 服务列表
//@Tags 服务管理
//@ID /admin/service_list
//@Accept json
//@Produce json
//@Param query query admin.ServiceListInput true "body"
//@Success 200 {object} middleware.Response{data=admin.ServiceListOutput} "success"
//@Router /admin/service_list [get]
func (s *ServiceController) ServiceList(c *gin.Context) {
	params := &admin.ServiceListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &dao.ServiceInfo{}
	list, total, err := serviceInfo.PageList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	out := &admin.ServiceListOutput{
		Total:    total,
		PageNo:   params.PageNo,
		PageSize: params.PageSize,
	}
	var outList []admin.ServiceListItemOutput
	for _, listItem := range list {
		serviceDetail, err := listItem.ServiceDetail(c, tx, &listItem)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			return
		}
		clusterIp := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSslPort := lib.GetStringConf("base.cluster.cluster_ssl_port")
		serviceAddr := "unknow"

		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HttpRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HttpRule.NeedHttps == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIp, clusterSslPort, serviceDetail.HttpRule.Rule)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HttpRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HttpRule.NeedHttps == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIp, clusterPort, serviceDetail.HttpRule.Rule)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HttpRule.RuleType == public.HTTPRuleTypeDomain {
			serviceAddr = serviceDetail.HttpRule.Rule
		}
		if serviceDetail.Info.LoadType == public.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIp, serviceDetail.TcpRule.Port)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIp, serviceDetail.GRPCRule.Port)
		}
		ipList := serviceDetail.LoadBalance.GetIpListByModel()
		counter, err := public.FlowCountHandler.GetFlowCounter(public.FlowServicePrefix+listItem.ServiceName)
		if err != nil {
			middleware.ResponseError(c, 2004, err)
			return
		}
		outItem := admin.ServiceListItemOutput{
			Id:          listItem.Id,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			LoadType:    listItem.LoadType,
			ServiceAddr: serviceAddr,
			Qps:         counter.QPS,
			Qpd:         counter.TotalCount,
			TotalNode:   len(ipList),
		}
		outList = append(outList, outItem)
	}
	out.List = outList
	middleware.ResponseSuccess(c, out)
}

//ServiceDelete godoc
//@Summary 服务删除
//@Description 服务删除
//@Tags 服务管理
//@ID /admin/service_delete
//@Accept json
//@Produce json
//@Param id query string true "服务id"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /admin/service_delete [get]
func (s *ServiceController) ServiceDelete(c *gin.Context) {
	params := &admin.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &dao.ServiceInfo{Id: params.Id}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	serviceInfo.IsDelete = public.IsDeleteTrue
	if err = serviceInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
}

//ServiceAdd godoc
//@Summary HTTP服务新增
//@Description 服务新增
//@Tags 服务管理
//@ID /admin/service_add_http
//@Accept json
//@Produce json
//@Param body body admin.ServiceAddHttpInput true "body"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /admin/service_add_http [post]
func (s *ServiceController) ServiceAddHttp(c *gin.Context) {
	params := &admin.ServiceAddHttpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		middleware.ResponseError(c, 2004, errors.New("服务ip列表与权重不匹配"))
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//服务重复校验
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	if _, err = serviceInfo.Find(c, tx, serviceInfo); err == nil {
		middleware.ResponseError(c, 2002, errors.New("服务名已经存在"))
		return
	}
	//服务接入类型 + 前缀 重复校验
	httpRule := &dao.ServiceHttpRule{RuleType: params.RuleType, Rule: params.Rule}
	if _, err = httpRule.Find(c, tx, httpRule); err == nil {
		middleware.ResponseError(c, 2003, errors.New("服务接入前缀/域名已经存在"))
		return
	}
	//开启事务
	tx = tx.Begin()
	serviceModel := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
	}
	if err = serviceModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}
	httpRuleModel := &dao.ServiceHttpRule{
		ServiceId:      serviceModel.Id,
		RuleType:       params.RuleType,
		Rule:           params.Rule,
		NeedHttps:      params.NeedHttps,
		NeedStripUri:   params.NeedStripUri,
		NeedWebsocket:  params.NeedWebsocket,
		UrlRewrite:     params.UrlRewrite,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err = httpRuleModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}
	accessControl := &dao.ServiceAccessControl{
		ServiceId:         serviceModel.Id,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		ClientipFlowLimit: params.ClientipFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err = accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}
	loadBalance := &dao.ServiceLoadBalance{
		ServiceId:              serviceModel.Id,
		RoundType:              params.RoundType,
		IpList:                 params.IpList,
		WeightList:             params.WeightList,
		UpstreamConnectTimeout: params.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  params.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    params.UpstreamIdleTimeout,
		UpstreamMaxIdle:        params.UpstreamMaxIdle,
	}
	if err = loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
}

//ServiceUpdate godoc
//@Summary HTTP服务编辑
//@Description HTTP服务编辑
//@Tags 服务管理
//@ID /admin/service_update_http
//@Accept json
//@Produce json
//@Param body body admin.ServiceUpdateHttpInput true "body"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /admin/service_update_http [post]
func (s *ServiceController) ServiceUpdateHttp(c *gin.Context) {
	params := &admin.ServiceUpdateHttpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		middleware.ResponseError(c, 2001, errors.New("服务ip列表与权重不匹配"))
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	//服务校验
	serviceInfo := &dao.ServiceInfo{Id: params.Id}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, errors.New("服务不存在"))
		return
	}
	serviceDetail, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2004, errors.New("服务不存在"))
		return
	}
	//开启事务
	tx = tx.Begin()
	info := &dao.ServiceInfo{Id: params.Id}
	info, err = info.Find(c, tx, info)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, errors.New("服务不存在"))
		return
	}
	info.ServiceDesc = params.ServiceDesc
	info.ServiceDesc = params.ServiceName
	info.UpdateAt = time.Now()
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	httpRule := serviceDetail.HttpRule
	httpRule.NeedHttps = params.NeedHttps
	httpRule.NeedStripUri = params.NeedStripUri
	httpRule.NeedWebsocket = params.NeedWebsocket
	httpRule.UrlRewrite = params.UrlRewrite
	httpRule.HeaderTransfor = params.HeaderTransfor
	if err := httpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.ClientipFlowLimit = params.ClientipFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}

	loadBalance := serviceDetail.LoadBalance
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	loadBalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	loadBalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	loadBalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2009, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
}

//ServiceDetail godoc
//@Summary 服务详情
//@Description 服务详情
//@Tags 服务管理
//@ID /admin/service_detail
//@Accept json
//@Produce json
//@Param id query  string true "服务id" default(62)
//@Success 200 {object} dao.ServiceDetail "success"
//@Router /admin/service_detail [get]
func (s *ServiceController) ServiceDetail(c *gin.Context) {
	idStr := c.Query("id")
	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &dao.ServiceInfo{Id: idInt}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, errors.New("服务不存在"))
		return
	}
	info, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, info)
}

//ServiceStat godoc
//@Summary 服务统计
//@Description 服务统计
//@Tags 服务管理
//@ID /admin/service_stat
//@Accept json
//@Produce json
//@Param id query  string true "服务id" default(62)
//@Success 200 {object} admin.ServiceStatOutput "success"
//@Router /admin/service_stat [get]
func (s *ServiceController) ServiceStat(c *gin.Context) {
	idStr := c.Query("id")
	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &dao.ServiceInfo{Id: idInt}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, errors.New("服务不存在"))
		return
	}

	var todayList, yesterdayList []int64
	counter, err := public.FlowCountHandler.GetFlowCounter(public.FlowServicePrefix+serviceInfo.ServiceName)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	curr := time.Now()
	for i := 0; i <= curr.Hour(); i++ {
		currTime := time.Date(curr.Year(), curr.Month(), curr.Day(), i, 0, 0, 0 , lib.TimeLocation)
		count, _ := counter.GetHourData(currTime)
		todayList = append(todayList, count)
	}
	yes := curr.Add(-1 * 24 * time.Hour)
	for i := 0; i <= 23; i++ {
		yesTime := time.Date(yes.Year(), yes.Month(), yes.Day(), i, 0, 0, 0 , lib.TimeLocation)
		count, _ := counter.GetHourData(yesTime)
		yesterdayList = append(yesterdayList, count)
	}
	middleware.ResponseSuccess(c, &admin.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}

//ServiceAddTcp godoc
//@Summary Tcp服务新增
//@Description Tcp服务新增
//@Tags 服务管理
//@ID /admin/service_add_tcp
//@Accept json
//@Produce json
//@Param body body admin.ServiceAddTcpInput true "body"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /admin/service_add_tcp [post]
func (s *ServiceController) ServiceAddTcp(c *gin.Context) {
	params := &admin.ServiceAddTcpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//验证 service_name 是否被占用
	infoSearch := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if _, err := infoSearch.Find(c, lib.GORMDefaultPool, infoSearch); err == nil {
		middleware.ResponseError(c, 2002, errors.New("服务名被占用，请重新输入"))
		return
	}

	//验证端口是否被占用?
	tcpRuleSearch := &dao.ServiceTcpRule{
		Port: params.Port,
	}
	if _, err := tcpRuleSearch.Find(c, lib.GORMDefaultPool, tcpRuleSearch); err == nil {
		middleware.ResponseError(c, 2003, errors.New("服务端口被占用，请重新输入"))
		return
	}
	grpcRuleSearch := &dao.ServiceGRPCRule{
		Port: params.Port,
	}
	if _, err := grpcRuleSearch.Find(c, lib.GORMDefaultPool, grpcRuleSearch); err == nil {
		middleware.ResponseError(c, 2004, errors.New("服务端口被占用，请重新输入"))
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2005, errors.New("ip列表与权重设置不匹配"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()
	info := &dao.ServiceInfo{
		LoadType:    public.LoadTypeTCP,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}
	loadBalance := &dao.ServiceLoadBalance{
		ServiceId:  info.Id,
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	httpRule := &dao.ServiceTcpRule{
		ServiceId: info.Id,
		Port:      params.Port,
	}
	if err := httpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}

	accessControl := &dao.ServiceAccessControl{
		ServiceId:         info.Id,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientipFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2009, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
	return
}

//ServiceUpdateTcp godoc
//@Summary Tcp服务编辑
//@Description Tcp服务编辑
//@Tags 服务管理
//@ID /admin/service_update_tcp
//@Accept json
//@Produce json
//@Param body body admin.ServiceUpdateTcpInput true "body"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /admin/service_update_tcp [post]
func (s *ServiceController) ServiceUpdateTcp(c *gin.Context) {
	params := &admin.ServiceUpdateTcpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2002, errors.New("ip列表与权重设置不匹配"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()

	service := &dao.ServiceInfo{
		Id: params.Id,
	}
	detail, err := service.ServiceDetail(c, lib.GORMDefaultPool, service)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}

	loadBalance := &dao.ServiceLoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	loadBalance.ServiceId = info.Id
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}

	tcpRule := &dao.ServiceTcpRule{}
	if detail.TcpRule != nil {
		tcpRule = detail.TcpRule
	}
	tcpRule.ServiceId = info.Id
	tcpRule.Port = params.Port
	if err := tcpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}

	accessControl := &dao.ServiceAccessControl{}
	if detail.AccessControl != nil {
		accessControl = detail.AccessControl
	}
	accessControl.ServiceId = info.Id
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientipFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
	return
}

//ServiceAddGRPC godoc
//@Summary GRPC服务新增
//@Description GRPC服务新增
//@Tags 服务管理
//@ID /admin/service_add_grpc
//@Accept json
//@Produce json
//@Param body body admin.ServiceAddGRPCInput true "body"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /admin/service_add_grpc [post]
func (s *ServiceController) ServiceAddGRPC(c *gin.Context) {
	params := &admin.ServiceAddGRPCInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//验证 service_name 是否被占用
	infoSearch := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if _, err := infoSearch.Find(c, lib.GORMDefaultPool, infoSearch); err == nil {
		middleware.ResponseError(c, 2002, errors.New("服务名被占用，请重新输入"))
		return
	}

	//验证端口是否被占用?
	tcpRuleSearch := &dao.ServiceTcpRule{
		Port: params.Port,
	}
	if _, err := tcpRuleSearch.Find(c, lib.GORMDefaultPool, tcpRuleSearch); err == nil {
		middleware.ResponseError(c, 2003, errors.New("服务端口被占用，请重新输入"))
		return
	}
	grpcRuleSearch := &dao.ServiceGRPCRule{
		Port: params.Port,
	}
	if _, err := grpcRuleSearch.Find(c, lib.GORMDefaultPool, grpcRuleSearch); err == nil {
		middleware.ResponseError(c, 2004, errors.New("服务端口被占用，请重新输入"))
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2005, errors.New("ip列表与权重设置不匹配"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()
	info := &dao.ServiceInfo{
		LoadType:    public.LoadTypeGRPC,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	loadBalance := &dao.ServiceLoadBalance{
		ServiceId:  info.Id,
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	grpcRule := &dao.ServiceGRPCRule{
		ServiceId:      info.Id,
		Port:           params.Port,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err := grpcRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}

	accessControl := &dao.ServiceAccessControl{
		ServiceId:         info.Id,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientipFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2009, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
	return
}

//ServiceUpdateGRPC godoc
//@Summary GRPC服务编辑
//@Description GRPC服务编辑
//@Tags 服务管理
//@ID /admin/service_update_grpc
//@Accept json
//@Produce json
//@Param body body admin.ServiceUpdateGRPCInput true "body"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /admin/service_update_grpc [post]
func (s *ServiceController) ServiceUpdateGRPC(c *gin.Context) {
	params := &admin.ServiceUpdateGRPCInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2002, errors.New("ip列表与权重设置不匹配"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()

	service := &dao.ServiceInfo{
		Id: params.Id,
	}
	detail, err := service.ServiceDetail(c, lib.GORMDefaultPool, service)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}

	loadBalance := &dao.ServiceLoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	loadBalance.ServiceId = info.Id
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}

	grpcRule := &dao.ServiceGRPCRule{}
	if detail.GRPCRule != nil {
		grpcRule = detail.GRPCRule
	}
	grpcRule.ServiceId = info.Id
	//grpcRule.Port = params.Port
	grpcRule.HeaderTransfor = params.HeaderTransfor
	if err := grpcRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	accessControl := &dao.ServiceAccessControl{}
	if detail.AccessControl != nil {
		accessControl = detail.AccessControl
	}
	accessControl.ServiceId = info.Id
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientipFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
	return
}
