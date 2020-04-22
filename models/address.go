package models

import (
	"github.com/jinzhu/gorm"
	"one/infra/base_c"
)

type UserAddress struct {
	gorm.Model
	UserId       uint  `gorm:"type:int;not null"`
	Province     string `gorm:"type:varchar(255);not null" validate:"required"`
	City         string `gorm:"type:varchar(255);not null" validate:"required"`
	District     string `gorm:"type:varchar(255);not null" validate:"required"`
	Address      string `gorm:"type:varchar(255);not null" validate:"required"`
	Zip          int   `gorm:"type:varchar(255);not null" validate:"required"`
	ContactName  string `gorm:"type:varchar(255);not null" validate:"required"`
	ContactPhone string `gorm:"type:varchar(255);not null" validate:"required"`
}

// 新增地址
func (u *UserAddress) Create() error {
	return base_c.Database().Save(&u).Error
}

// 更新地址
func (u *UserAddress) Update() error {
	return base_c.Database().Model(&UserAddress{}).Where("user_id = ?", u.UserId).Updates(&u).Error
}

// 删除地址
func (u *UserAddress) Delete() error {
	return base_c.Database().Delete(&u).Error
}

// 地址列表
func (u *UserAddress) List(offset, limit string) (*[]UserAddress, error) {
	var list = make([]UserAddress, 10, 10)
	err := base_c.Database().Model(&UserAddress{}).Offset(offset).Limit(limit).Where("user_id = ?", u.UserId).Find(&list).Error
	return &list, err
}

// 获取单个地址
func GetByIdAndUserId(id, uid uint) (*UserAddress, error) {
	var address UserAddress
	address.ID = id
	err := base_c.Database().Where("user_id = ?", uid).First(&address).Error
	return &address, err
}
