package web

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"one/infra"
	"one/infra/base"
	"one/infra/helper"
	"one/infra/middleware"
	"one/models"
	"time"
)

func init() {
	infra.RegisterApi(new(AdminApi))
}

type AdminApi struct {}

func (a *AdminApi) Init() {
	api := AccountApi{}
	groupRouter := base.Gin().Group("/v1/admin")
	groupRouter.POST("/register", a.register)
	groupRouter.POST("/login", a.login)
	m := groupRouter.Use(middleware.Jwt())
	m.DELETE("/logout", api.logout)
}

// 账户创建
func (a *AdminApi) register(c *gin.Context) {
	// 参数检验
	d, _ := c.GetRawData()
	var newUser models.Admin
	jsoniter.Unmarshal(d, &newUser)
	if err := base.ValidateStruct(&newUser); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	// 注册逻辑
	if err := newUser.Register(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	//token := time.Now().Unix()
	//base_c.Redis().Set(models.VerifyPre+newUser.Email, token, time.Duration(30)*time.Minute)
	//singe := fmt.Sprintf("%d", token)
	////发送邮件
	//go func() {
	//	err := helper.SendEmail(
	//		"imlzqiang@163.com",
	//		"",
	//		"",
	//		"尊敬的用户您好",
	//		"感谢您的注册，点击链接:<a href=\"localhost:8080/v1/verify?singe="+singe+"\">激活账号！</a>",
	//		"",
	//		newUser.Email)
	//	if err != nil {
	//		logrus.Error("Email send error: " + newUser.Email)
	//	}
	//}()
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "注册成功",
		"data":   "",
	})
}

// 账户登录
func (a *AdminApi) login(c *gin.Context) {
	// 参数校验
	d, _ := c.GetRawData()
	var userParams models.UserLogin
	jsoniter.Unmarshal(d, &userParams)
	if err := base.ValidateStruct(&userParams); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	// 登录逻辑
	user := userParams.AdminLogin()
	if user == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    "邮箱或密码错误",
			"data":   "",
		})
		return
	}
	// redis 记录登录信息
	expiresAt := time.Now().Add(base.Props().GetDurationDefault("jwt.expires", time.Duration(7*24)*time.Hour)).Unix()
	// jwt
	issuer := base.Props().GetDefault("app.name", "one")
	token, err := helper.GenToken(&models.Claim{
		Id:    user.ID,
		Name:  user.Name,
		Role:  user.Role,
	}, expiresAt, issuer)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"msg":    err.Error(),
			"data":   "",
		})
	}
	type loginInfo struct {
		models.Admin
		Token string
	}
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "登录成功",
		"data":   loginInfo{*user, token},
	})
}
