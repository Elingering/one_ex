package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"net/http"
	"one/infra"
	"one/infra/base"
	"one/infra/base_c"
	"one/infra/helper"
	"one/infra/middleware"
	"one/models"
	"strings"
	"time"
)

func init() {
	infra.RegisterApi(new(AccountApi))
}

type AccountApi struct {}

func (a *AccountApi) Init() {
	groupRouter := base.Gin().Group("/v1")
	groupRouter.POST("/register", a.register)
	groupRouter.POST("/login", a.login)
	m := groupRouter.Use(middleware.Jwt())
	m.DELETE("/logout", a.logout)
	m.GET("/verify", a.verify)
}

// 账户创建
func (a *AccountApi) register(c *gin.Context) {
	// 参数检验
	d, _ := c.GetRawData()
	var newUser models.User
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
	token := time.Now().Unix()
	base_c.Redis().Set(models.VerifyPre+newUser.Email, token, time.Duration(30)*time.Minute)
	singe := fmt.Sprintf("%d", token)
	//发送邮件
	go func() {
		err := helper.SendEmail(
			"imlzqiang@163.com",
			"",
			"",
			"尊敬的用户您好",
			"感谢您的注册，点击链接:<a href=\"localhost:8080/v1/verify?singe="+singe+"\">激活账号！</a>",
			"",
			newUser.Email)
		if err != nil {
			logrus.Error("Email send error: " + newUser.Email)
		}
	}()
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "注册成功",
		"data":   "",
	})
}

// 账户登录
func (a *AccountApi) login(c *gin.Context) {
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
	user := userParams.Login()
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
		Email: user.Email,
	}, expiresAt, issuer)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"msg":    err.Error(),
			"data":   "",
		})
	}
	type loginInfo struct {
		models.User
		Token string
	}
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "登录成功",
		"data":   loginInfo{*user, token},
	})
}

// 邮箱验证
func (a *AccountApi) verify(c *gin.Context) {
	claim := helper.GetJwtClaimsByContext(c)
	email := claim.Email
	singe, _ := c.GetQuery("singe")
	str, err := base_c.Redis().Get(models.VerifyPre + email).Result()
	if str != singe {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    "验证失败（链接过期/重复操作）",
			"data":   "",
		})
		return
	}
	err = models.Verify(email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	base_c.Redis().Del(models.VerifyPre + email)
	c.JSON(http.StatusOK, gin.H{
		"status": false,
		"msg":    "验证成功",
		"data":   "",
	})
}

// 账户退出
func (a *AccountApi) logout(c *gin.Context) {
	// 参数校验
	str := c.GetHeader("Authorization")
	str = strings.Replace(str, "Bearer ", "", 1)
	claim := helper.GetJwtClaimsByContext(c)
	//redis 黑名单
	ok, _ := base_c.Redis().SetNX(str, 1, time.Duration(claim.ExpiresAt-time.Now().Unix())*time.Second).Result()
	if ok {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"msg":    "操作成功",
			"data":   "",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    "操作失败",
			"data":   "",
		})
	}
}
