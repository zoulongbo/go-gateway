package admin

import (
	"errors"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dao"
	"github.com/zoulongbo/go-gateway/dto/admin"
	"github.com/zoulongbo/go-gateway/middleware"
	"github.com/zoulongbo/go-gateway/public"
	"time"
)

type AppController struct {
}

func RegisterAppController(g *gin.RouterGroup) {
	app := &AppController{}
	g.POST("/app_add", app.AppAdd)
	g.POST("/app_update", app.AppUpdate)
	g.GET("/app_detail", app.AppDetail)
	g.GET("/app_delete", app.AppDelete)
	g.GET("/app_list", app.AppList)
	g.GET("/app_stat", app.AppStat)
}

// AppAdd godoc
// @Summary 租户添加
// @Description 租户添加
// @Tags 租户管理
// @ID /app/app_add
// @Accept  json
// @Produce  json
// @Param body body admin.AppAddInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_add [post]
func (app *AppController) AppAdd(c *gin.Context) {
	params := &admin.AppAddInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//验证app_id是否被占用
	search := &dao.App{
		AppId: params.AppId,
	}
	if _, err := search.Find(c, lib.GORMDefaultPool, search); err == nil {
		middleware.ResponseError(c, 2002, errors.New("租户ID被占用，请重新输入"))
		return
	}
	if params.Secret == "" {
		params.Secret = public.MD5(params.AppId)
	}
	tx := lib.GORMDefaultPool
	info := &dao.App{
		AppId:    params.AppId,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIPS: params.WhiteIpS,
		Qps:      params.Qps,
		Qpd:      params.Qpd,
	}
	if err := info.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}

// AppUpdate godoc
// @Summary 租户更新
// @Description 租户更新
// @Tags 租户管理
// @ID /app/app_update
// @Accept  json
// @Produce  json
// @Param body body admin.AppUpdateInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_update [post]
func (app *AppController) AppUpdate(c *gin.Context) {
	params := &admin.AppUpdateInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.App{
		Id: params.Id,
	}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	if params.Secret == "" {
		params.Secret = public.MD5(params.AppId)
	}
	info.Name = params.Name
	info.Secret = params.Secret
	info.WhiteIPS = params.WhiteIpS
	info.Qps = params.Qps
	info.Qpd = params.Qpd
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}

// AppDetail godoc
// @Summary 租户详情
// @Description 租户详情
// @Tags 租户管理
// @ID /app/app_detail
// @Accept  json
// @Produce  json
// @Param id query admin.AppDetailInput true "租户id"
// @Success 200 {object} middleware.Response{data=dao.App} "success"
// @Router /app/app_detail [get]
func (app *AppController) AppDetail(c *gin.Context) {
	params := &admin.AppDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.App{
		Id: params.Id,
	}
	detail, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	middleware.ResponseSuccess(c, detail)
	return
}

// AppDelete godoc
// @Summary 租户删除
// @Description 租户删除
// @Tags 租户管理
// @ID /app/app_delete
// @Accept  json
// @Produce  json
// @Param id query admin.AppDetailInput true "租户id"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_delete [get]
func (app *AppController) AppDelete(c *gin.Context) {
	params := &admin.AppDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.App{
		Id: params.Id,
	}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	info.IsDelete = 1
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}


// AppList godoc
// @Summary 租户列表
// @Description 租户列表
// @Tags 租户管理
// @ID /app/app_list
// @Accept  json
// @Produce  json
// @Param info query admin.AppListInput true "body"
// @Success 200 {object} middleware.Response{data=admin.AppListOutput} "success"
// @Router /app/app_list [get]
func (app *AppController) AppList(c *gin.Context) {
	params := &admin.AppListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	info := &dao.App{}
	list, total, err := info.AppList(c, lib.GORMDefaultPool, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	var outputList []admin.AppListItemOutput
	for _, item := range list {
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			c.Abort()
			return
		}
		outputList = append(outputList, admin.AppListItemOutput{
			Id:       item.Id,
			AppId:    item.AppId,
			Name:     item.Name,
			Secret:   item.Secret,
			WhiteIPS: item.WhiteIPS,
			Qpd:      item.Qpd,
			Qps:      item.Qps,
		})
	}
	output := &admin.AppListOutput{
		PageSize: params.PageSize,
		PageNo:   params.PageNo,
		List:     outputList,
		Total:    total,
	}
	middleware.ResponseSuccess(c, output)
	return
}

// AppStat godoc
// @Summary 租户统计
// @Description 租户统计
// @Tags 租户管理
// @ID /app/app_stat
// @Accept  json
// @Produce  json
// @Param id query admin.AppDetailInput true "租户id"
// @Success 200 {object} middleware.Response{data=admin.AppStatOutput} "success"
// @Router /app/app_stat [get]
func (app *AppController) AppStat(c *gin.Context) {
	params := &admin.AppDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	detail := &dao.App{Id: params.Id,}
	detail, err := detail.Find(c, lib.GORMDefaultPool, detail)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	//今日流量全天小时级访问统计
	var todayStat, yesterdayStat []int64
	counter, err := public.FlowCountHandler.GetFlowCounter(public.FlowAppPrefix+detail.AppId)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	curr := time.Now()
	for i := 0; i <= curr.Hour(); i++ {
		currTime := time.Date(curr.Year(), curr.Month(), curr.Day(), i, 0, 0, 0 , lib.TimeLocation)
		count, _ := counter.GetHourData(currTime)
		todayStat = append(todayStat, count)
	}
	yes := curr.Add(-1 * 24 * time.Hour)
	for i := 0; i <= 23; i++ {
		yesTime := time.Date(yes.Year(), yes.Month(), yes.Day(), i, 0, 0, 0 , lib.TimeLocation)
		count, _ := counter.GetHourData(yesTime)
		yesterdayStat = append(yesterdayStat, count)
	}
	middleware.ResponseSuccess(c, &admin.AppStatOutput{
		Today:     todayStat,
		Yesterday: yesterdayStat,
	})
}
