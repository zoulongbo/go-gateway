package dao

import (
	"errors"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/zoulongbo/go-gateway/dto"
	"github.com/zoulongbo/go-gateway/public"
	"time"
)

type Admin struct {
	Id       int       `json:"id" gorm:"primary_key" description:"自增主键"`
	Username string    `json:"username" gorm:"column:username" description:"用户名"`
	Password string    `json:"password" gorm:"column:password" description:"密码"`
	Salt     string    `json:"salt" gorm:"column:salt" description:"密码盐"`
	CreateAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	UpdateAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *Admin) TableName() string {
	return "gateway_admin"
}

func (t *Admin) Find(c *gin.Context, tx *gorm.DB, search *Admin) (*Admin, error) {
	result := &Admin{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *Admin) LoginCheck(c *gin.Context, tx *gorm.DB, param *dto.LoginInput) (*Admin, error) {
	adminInfo, err := t.Find(c, tx, &Admin{Username: param.Username, IsDelete: 0})
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if adminInfo.Password != public.GenSaltPassword(adminInfo.Salt, param.Password) {
		return nil, errors.New("密码错误")
	}
	return adminInfo, nil
}

func (t *Admin) Save(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
	if err != nil {
		return err
	}
	return nil
}
