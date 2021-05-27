package dao

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/public"
)

type ServiceTcpRule struct {
	Id        int64 `json:"id" gorm:"primary_key" description:"自增主键"`
	ServiceId int64 `json:"service_id" gorm:"column:service_id" description:"服务id"`
	Port      int   `json:"port" gorm:"column:port" description:"端口"`
}

func (s *ServiceTcpRule) TableName() string {
	return "gateway_service_tcp_rule"
}

func (s *ServiceTcpRule) Find(c *gin.Context, tx *gorm.DB, search *ServiceTcpRule) (*ServiceTcpRule, error) {
	result := &ServiceTcpRule{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ServiceTcpRule) Save(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(public.GetGinTraceContext(c)).Save(s).Error
	if err != nil {
		return err
	}
	return nil
}
