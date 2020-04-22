package models

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"one/infra/base_c"
	"time"
)

// jwt
var JwtSecret = []byte("lzq2020")

const (
	JwtClaimsKey = "JwtClaim"
	PwdSalt      = "lzq2020"
	VerifyPre    = "verify_"
)

type Claim struct {
	Id    uint
	Name  string
	Email string
	Role  string
	jwt.StandardClaims
}

type User struct {
	gorm.Model
	Name            string    `gorm:"type:varchar(255);not null" validate:"required,max=12"`
	Email           string    `gorm:"type:varchar(255);unique_index;not null" validate:"required,email"`
	EmailVerifiedAt time.Time `gorm:"default:null"`
	Password        string    `gorm:"type:varchar(255);not null" validate:"required,min=6,max=12"`
	RememberToken   string    `gorm:"default:null"`
}

type UserLogin struct {
	Email           string    `validate:"required,email"`
	Password        string    `validate:"required,min=6,max=12"`
}

// 新增用户
func (u *User) Register() error {
	if !base_c.Database().Where("email = ?", u.Email).First(&User{}).RecordNotFound() {
		return errors.New("邮箱已经存在")
	}
	u.Password = DigestString(PwdSalt + u.Password)
	return base_c.Database().Save(&u).Error
}

// 用户登录
func (u *UserLogin) Login() *User {
	var user User
	u.Password = DigestString(PwdSalt + u.Password)
	if !base_c.Database().Where("email = ? and password = ?", u.Email, u.Password).First(&user).RecordNotFound() {
		return &user
	}
	return nil
}

// 通过邮箱获取用户
//func (u *User) GetUserByEmail(email string) *User {
//	var user User
//	if !base_c.Database().Where("email = ?", email).First(&user).RecordNotFound() {
//		return &user
//	}
//	return nil
//}

// 验证邮箱
func Verify(email string) error {
	return base_c.Database().Model(&User{}).Where("email = ?", email).UpdateColumn("email_verified_at", time.Now()).Error
}