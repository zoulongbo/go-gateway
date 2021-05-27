package dao

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dto/admin"
	"github.com/zoulongbo/go-gateway/public"
	"time"
)

type ServiceInfo struct {
	Id          int64     `json:"id" gorm:"primary_key" description:"自增主键"`
	LoadType    int       `json:"load_type" gorm:"column:load_type" description:"负载类型 0=http 1=tcp 2=grpc"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	CreateAt    time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	UpdateAt    time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete    int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (s *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

func (s *ServiceInfo) PageList(c *gin.Context, tx *gorm.DB, params *admin.ServiceListInput) ([]ServiceInfo, int64, error) {
	total := int64(0)
	var list []ServiceInfo
	offset := (params.PageNo - 1) * params.PageSize
	query := tx.SetCtx(public.GetGinTraceContext(c)).Table(s.TableName()).Where("is_delete = ?", public.IsDeleteFalse)
	if params.Info != "" {
		query = query.Where("service_name like ? or  service_desc like ?", "%"+params.Info+"%", "%"+params.Info+"%")
	}

	if err := query.Limit(params.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, nil
	}

	query.Count(&total)

	return list, total, nil
}

func (s *ServiceInfo) ServiceDetail(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceDetail, error) {
	httpRule := &ServiceHttpRule{ServiceId: search.Id}
	httpRule, err := httpRule.Find(c, tx, httpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	tcpRule := &ServiceTcpRule{ServiceId: search.Id}
	tcpRule, err = tcpRule.Find(c, tx, tcpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	grpcRule := &ServiceGRPCRule{ServiceId: search.Id}
	grpcRule, err = grpcRule.Find(c, tx, grpcRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	accessControl := &ServiceAccessControl{ServiceId: search.Id}
	accessControl, err = accessControl.Find(c, tx, accessControl)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	loadBalance := &ServiceLoadBalance{ServiceId: search.Id}
	loadBalance, err = loadBalance.Find(c, tx, loadBalance)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &ServiceDetail{
		Info:          search,
		HttpRule:      httpRule,
		TcpRule:       tcpRule,
		GRPCRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}, nil
}

func (s *ServiceInfo) Find(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	result := &ServiceInfo{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ServiceInfo) Save(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(public.GetGinTraceContext(c)).Save(s).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceInfo) GroupByLoadType (c *gin.Context, tx *gorm.DB) ([]admin.DashServiceStatItemOutput, error)  {
	var list []admin.DashServiceStatItemOutput
	query := tx.SetCtx(public.GetGinTraceContext(c))
	if err := query.Table(s.TableName()).Where("is_delete=0").Select("load_type, count(*) as value").Group("load_type").Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
