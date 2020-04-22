package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"one/infra/base_c"
	"time"
)

type Admin struct {
	gorm.Model
	Name            string    `gorm:"type:varchar(255);not null" validate:"required,max=12"`
	Email           string    `gorm:"type:varchar(255);unique_index;not null" validate:"required,email"`
	Role    		string    `gorm:"type:varchar(255);not null"`
	EmailVerifiedAt time.Time `gorm:"default:null"`
	Password        string    `gorm:"type:varchar(255);not null" validate:"required,min=6,max=12"`
	RememberToken   string    `gorm:"default:null"`
}

// 新增用户
func (u *Admin) Register() error {
	if !base_c.Database().Where("email = ?", u.Email).First(&Admin{}).RecordNotFound() {
		return errors.New("邮箱已经存在")
	}
	u.Password = DigestString(PwdSalt + u.Password)
	return base_c.Database().Save(&u).Error
}

// 用户登录
func (u *UserLogin) AdminLogin() *Admin {
	var user Admin
	u.Password = DigestString(PwdSalt + u.Password)
	if !base_c.Database().Where("email = ? and password = ?", u.Email, u.Password).First(&user).RecordNotFound() {
		return &user
	}
	return nil
}

// 分配角色
func (u *AdminRole) Update(tx *gorm.DB) error {
	return tx.Model(&Admin{}).Where("id = ?", u.AdminId).UpdateColumn("role", u.AdminRole).Error
}
