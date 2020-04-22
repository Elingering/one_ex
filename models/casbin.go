package models

import (
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	"github.com/tietang/props/kvs"
	"one/infra/base_c"
)

//权限结构
type CasbinModel struct {
	Ptype    string `gorm:"type:varchar(255);not null"`
	RoleName string `gorm:"type:varchar(255);not null"`
	Path     string `gorm:"type:varchar(255);not null"`
	Method   string `gorm:"type:varchar(255);not null"`
}

// 添加权限
func (c *CasbinModel) AddCasbin(cm CasbinModel) bool {
	e := Casbin()
	return e.AddPolicy(cm.RoleName, cm.Path, cm.Method)
}

// 用户关联角色
func (c *CasbinModel) LinkCasbin(user, role string) bool {
	e := Casbin()
	return e.AddRoleForUser(user, role)
}

// 持久化到数据库
func Casbin() *casbin.Enforcer {
	a := gormadapter.NewAdapterByDB(base_c.Database())
	file := kvs.GetCurrentFilePath("rbac_model.conf", 1)
	e := casbin.NewEnforcer(file, a)
	e.LoadPolicy()
	return e
}

//// 权限列表
//func (c *CasbinModel) List(offset, limit string) (*[]CasbinModel, error) {
//	var list = make([]CasbinModel, 10, 10)
//	err := base_c.Database().Model(&CasbinModel{}).Offset(offset).Limit(limit).Find(&list).Error
//	return &list, err
//}
