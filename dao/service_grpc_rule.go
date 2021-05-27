package dao

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/public"
)

type ServiceGRPCRule struct {
	Id             int64  `json:"id" gorm:"primary_key" description:"自增主键"`
	ServiceId      int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	Port           int    `json:"port" gorm:"column:port" description:"端口"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔"`
}

func (s *ServiceGRPCRule) TableName() string {
	return "gateway_service_grpc_rule"
}

func (s *ServiceGRPCRule) Find(c *gin.Context, tx *gorm.DB, search *ServiceGRPCRule) (*ServiceGRPCRule, error) {
	result := &ServiceGRPCRule{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ServiceGRPCRule) Save(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(public.GetGinTraceContext(c)).Save(s).Error
	if err != nil {
		return err
	}
	return nil
}
