package dao

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/public"
)

type ServiceHttpRule struct {
	Id             int64  `json:"id" gorm:"primary_key" description:"自增主键"`
	ServiceId      int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	RuleType       int    `json:"rule_type" gorm:"column:rule_type" description:"匹配类型 0=url前缀url_prefix 1=域名domain"`
	Rule           string `json:"rule" gorm:"column:rule" description:"type=domain表示域名，type=url_prefix时表示url前缀"`
	NeedHttps      int    `json:"need_https" gorm:"column:need_https" description:"支持https 1=支持"`
	NeedStripUri   int    `json:"need_strip_uri" gorm:"column:need_strip_uri" description:"启用strip_uri 1=启用"`
	NeedWebsocket  int    `json:"need_websocket" gorm:"column:need_websocket" description:"是否支持websocket 1=支持"`
	UrlRewrite     string `json:"url_rewrite" gorm:"column:url_rewrite" description:"url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔"`
}

func (s *ServiceHttpRule) TableName() string {
	return "gateway_service_http_rule"
}

func (s *ServiceHttpRule) Find(c *gin.Context, tx *gorm.DB, search *ServiceHttpRule) (*ServiceHttpRule, error) {
	result := &ServiceHttpRule{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ServiceHttpRule) Save(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(public.GetGinTraceContext(c)).Save(s).Error
	if err != nil {
		return err
	}
	return nil
}
