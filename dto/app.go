package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/public"
	"time"
)

type AppAddInput struct {
	AppId    string `json:"app_id" form:"app_id" comment:"租户id" validate:"required"`
	Name     string `json:"name" form:"name" comment:"租户名称" validate:"required"`
	Secret   string `json:"secret" form:"secret" comment:"密钥" validate:""`
	WhiteIpS string `json:"white_ips" form:"white_ips" comment:"ip白名单，支持前缀匹配"`
	Qpd      int64  `json:"qpd" form:"qpd" comment:"日请求量限制" validate:""`
	Qps      int64  `json:"qps" form:"qps" comment:"每秒请求量限制" validate:""`
}

func (a *AppAddInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, a)
}


type AppUpdateInput struct {
	Id       int64  `json:"id" form:"id" gorm:"column:id" comment:"主键ID" validate:"required"`
	AppAddInput
}

func (a *AppUpdateInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, a)
}

type AppListInput struct {
	Info     string `json:"info" form:"info" comment:"查找信息" validate:""`
	PageSize int    `json:"page_size" form:"page_size" comment:"页数" validate:"required,min=1,max=999"`
	PageNo   int    `json:"page_no" form:"page_no" comment:"页码" validate:"required,min=1,max=999"`
}

func (a *AppListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, a)
}

type AppListItemOutput struct {
	Id        int64     `json:"id" gorm:"primary_key"`
	AppId     string    `json:"app_id" gorm:"column:app_id" description:"租户id	"`
	Name      string    `json:"name" gorm:"column:name" description:"租户名称	"`
	Secret    string    `json:"secret" gorm:"column:secret" description:"密钥"`
	WhiteIPS  string    `json:"white_ips" gorm:"column:white_ips" description:"ip白名单，支持前缀匹配		"`
	Qpd       int64     `json:"qpd" gorm:"column:qpd" description:"日请求量限制"`
	Qps       int64     `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
	RealQpd   int64       `json:"real_qpd" description:"日请求量限制"`
	RealQps   int64       `json:"real_qps" description:"每秒请求量限制"`
	UpdatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"添加时间	"`
	CreatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete  int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

type AppListOutput struct {
	Total    int64                   `json:"total" form:"total" comment:"总数" example:"20"`           //总数
	PageNo   int                     `json:"page_no" form:"page_no" comment:"页码" example:"1"`        //页码
	PageSize int                     `json:"page_size" form:"page_size" comment:"页数据量" example:"20"` //页数据量
	List     []AppListItemOutput 		  `json:"list" form:"list" comment:"数据列表" example:""`             //数据列表
}

type AppDetailInput struct {
	Id int64 `json:"id" form:"id" comment:"租户ID" validate:"required"`
}

func (a *AppDetailInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, a)
}


type AppStatOutput struct {
	Today     []int64 `json:"today" form:"today" comment:"今日统计" validate:"required"`
	Yesterday []int64 `json:"yesterday" form:"yesterday" comment:"昨日统计" validate:"required"`
}
