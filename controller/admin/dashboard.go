package admin

import (
	"errors"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dao"
	"github.com/zoulongbo/go-gateway/dto/admin"
	"github.com/zoulongbo/go-gateway/middleware"
	"github.com/zoulongbo/go-gateway/public"
	"math/rand"
	"time"
)

type DashboardController struct {

}

func RegisterDashboardController(g *gin.RouterGroup) {
	dashboard := &DashboardController{}
	g.GET("/panel_group_data", dashboard.PanelGroupData)
	g.GET("/flow_stat", dashboard.FlowStat)
	g.GET("/service_stat", dashboard.ServiceStat)
}

// PanelGroupData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags 首页大盘
// @ID /dashboard/panel_group_data
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=admin.PanelGroupDataOutput} "success"
// @Router /dashboard/panel_group_data [get]
func (d *DashboardController) PanelGroupData (c *gin.Context)  {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &dao.ServiceInfo{}
	_, serviceNum, err := serviceInfo.PageList(c, tx, &admin.ServiceListInput{PageSize: 1, PageNo: 1})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	app := &dao.App{}
	_, appNum, err := app.AppList(c, tx, &admin.AppListInput{PageNo: 1, PageSize: 1})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	out := &admin.PanelGroupDataOutput{
		ServiceNum:      serviceNum,
		AppNum:          appNum,
	}
	middleware.ResponseSuccess(c, out)
}

// FlowStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 首页大盘
// @ID /dashboard/flow_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=admin.ServiceStatOutput} "success"
// @Router /dashboard/flow_stat [get]
func (d *DashboardController) FlowStat (c *gin.Context)  {
	//今日流量全天小时级访问统计
	var todayList, yesterdayList []int64
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < time.Now().Hour(); i++ {
		todayList = append(todayList, rand.Int63n(100))
	}
	for i := 0; i < 23; i++ {
		yesterdayList = append(yesterdayList, rand.Int63n(100))
	}
	middleware.ResponseSuccess(c, &admin.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 首页大盘
// @ID /dashboard/service_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=admin.DashServiceStatOutput} "success"
// @Router /dashboard/service_stat [get]
func (d *DashboardController) ServiceStat (c *gin.Context)  {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &dao.ServiceInfo{}
	list, err := serviceInfo.GroupByLoadType(c, tx)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	var legend []string
	for index, item := range list {
		name, ok := public.LoadTypeMap[item.LoadType]
		if !ok {
			middleware.ResponseError(c, 2003, errors.New("load_type not found"))
			return
		}
		list[index].Name = name
		legend = append(legend, name)
	}
	out := &admin.DashServiceStatOutput{
		Legend: legend,
		Data:   list,
	}
	middleware.ResponseSuccess(c, out)
}