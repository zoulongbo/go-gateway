package dao

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/public"
)

type ServiceAccessControl struct {
	Id                int64  `json:"id" gorm:"primary_key" description:"自增主键"`
	ServiceId         int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	OpenAuth          int    `json:"open_auth" gorm:"column:open_auth" description:"是否开启权限 1=开启"`
	BlackList         string `json:"black_list" gorm:"column:black_list" description:"黑名单"`
	WhiteList         string `json:"white_list" gorm:"column:white_list" description:"白名单"`
	WhiteHostName     string `json:"white_host_name" gorm:"column:white_host_name" description:"白名单主机"`
	ClientipFlowLimit int    `json:"clientip_flow_limit" gorm:"column:clientip_flow_limit" description:"客户端ip限流"`
	ServiceFlowLimit  int    `json:"service_flow_limit" gorm:"column:service_flow_limit" description:"服务端限流"`
}

func (s *ServiceAccessControl) TableName() string {
	return "gateway_service_access_control"
}

func (s *ServiceAccessControl) Find(c *gin.Context, tx *gorm.DB, search *ServiceAccessControl) (*ServiceAccessControl, error) {
	result := &ServiceAccessControl{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ServiceAccessControl) Save(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(public.GetGinTraceContext(c)).Save(s).Error
	if err != nil {
		return err
	}
	return nil
}
