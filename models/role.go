package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"one/infra/base_c"
)

type Role struct {
	gorm.Model
	ShowName            string    `gorm:"type:varchar(255);not null" validate:"required"`
	RoleName        string    `gorm:"type:varchar(255);unique_index;not null" validate:"required"`
	Sort			uint	  `gorm:"default:0"`
}

type AdminRole struct {
	AdminId   uint  `validate:"required"`
	AdminRole   string  `validate:"required"`
	RoleName string `validate:"required"`
}

// 新增角色
func (r *Role) Create() error {
	if !base_c.Database().Where("role_name = ?", r.RoleName).First(&Role{}).RecordNotFound() {
		return errors.New("角色已经存在")
	}
	if !base_c.Database().Where("show_name = ?", r.ShowName).First(&Role{}).RecordNotFound() {
		return errors.New("角色名已经存在")
	}
	return base_c.Database().Save(&r).Error
}

// 更新角色
func (r *Role) Update() error {
	var role Role
	base_c.Database().Where("show_name = ?", r.ShowName).Find(&role)
	if role.ID != 0 && role.ID != r.ID {
		return errors.New("角色名已经存在")
	}
	return base_c.Database().Model(&Role{}).Where("id = ?",r.ID).Updates(&r).Error
}

// 删除角色
func (r *Role) Delete() error {
	return base_c.Database().Delete(&r).Error
}

// 角色列表
func (r *Role) List(offset, limit string) (*[]Role, error) {
	var list = make([]Role, 10, 10)
	err := base_c.Database().Model(&Role{}).Offset(offset).Limit(limit).Find(&list).Error
	return &list, err
}
