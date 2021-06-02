package dao

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dto"
	"github.com/zoulongbo/go-gateway/public"
	"net/http/httptest"
	"sync"
	"time"
)

type App struct {
	Id        int64     `json:"id" gorm:"primary_key"`
	AppId     string    `json:"app_id" gorm:"column:app_id" description:"租户id	"`
	Name      string    `json:"name" gorm:"column:name" description:"租户名称	"`
	Secret    string    `json:"secret" gorm:"column:secret" description:"密钥"`
	WhiteIPS  string    `json:"white_ips" gorm:"column:white_ips" description:"ip白名单，支持前缀匹配"`
	Qpd       int64     `json:"qpd" gorm:"column:qpd" description:"日请求量限制"`
	Qps       int64     `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"添加时间	"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete  int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

func (a *App) TableName() string {
	return "gateway_app"
}

func (a *App) Find(c *gin.Context, tx *gorm.DB, search *App) (*App, error) {
	model := &App{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

func (a *App) Save(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(public.GetGinTraceContext(c)).Save(a).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *App) AppList(c *gin.Context, tx *gorm.DB, params *dto.AppListInput) ([]App, int64, error) {
	var list []App
	var count int64
	pageNo := params.PageNo
	pageSize := params.PageSize

	//limit offset,pagesize
	offset := (pageNo - 1) * pageSize
	query := tx.SetCtx(public.GetGinTraceContext(c))
	query = query.Table(a.TableName()).Select("*")
	query = query.Where("is_delete=?", public.IsDeleteFalse)
	if params.Info != "" {
		query = query.Where(" (name like ? or app_id like ?)", "%"+params.Info+"%", "%"+params.Info+"%")
	}
	err := query.Limit(pageSize).Offset(offset).Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		return nil, 0, err
	}
	return list, count, nil
}

var AppHandler *AppManager

func init()  {
	AppHandler = NewAppManager()
}

type AppManager struct {
	AppMap map[string]*App
	AppSlice []*App
	Locker sync.RWMutex
	init sync.Once
	err error
}

func NewAppManager() *AppManager  {
	return &AppManager{
		AppMap:   map[string]*App{},
		AppSlice: []*App{},
		Locker:   sync.RWMutex{},
		init:     sync.Once{},
	}
}

func (am *AppManager)GetApp (app *App) (*App, error)  {
	for _, item := range am.AppSlice {
		if item.AppId == app.AppId && item.Secret == app.Secret {
			return app, nil
		}
	}
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	tx, err := lib.GetGormPool("default")
	if err != nil {
		return nil, err
	}
	app, err = app.Find(c, tx, app)
	if err != nil {
		return nil, err
	}
	am.AppSlice = append(am.AppSlice, app)
	am.Locker.Lock()
	defer am.Locker.Unlock()

	am.AppMap[app.AppId] = app
	return app, nil
}


func (am *AppManager) LoadOne() error  {
	am.init.Do(func() {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		tx, err := lib.GetGormPool("default")
		if err != nil {
			am.err = err
			return
		}
		app := &App{}
		appList, _, err := app.AppList(c, tx, &dto.AppListInput{PageNo:1,PageSize:9999})
		if err != nil {
			am.err = err
			return
		}
		am.Locker.Lock()
		defer am.Locker.Unlock()

		for _, item := range appList {
			tmpItem := item
			am.AppSlice = append(am.AppSlice, &tmpItem)
			am.AppMap[tmpItem.AppId] = &tmpItem
		}
	})
	return am.err
}