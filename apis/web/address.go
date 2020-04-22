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
)

func init() {
	infra.RegisterApi(new(AddressApi))
}

type AddressApi struct {}

func (a *AddressApi) Init() {
	r := base.Gin().Group("/v1").Use(middleware.Jwt())
	r.POST("/address", a.create)
	r.PUT("/address", a.update)
	r.DELETE("/address", a.delete)
	r.GET("/address", a.list)
}

func (a *AddressApi) create(c *gin.Context) {
	// 参数校验
	d, _ := c.GetRawData()
	var address models.UserAddress
	jsoniter.Unmarshal(d, &address)
	if err := base.ValidateStruct(&address); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	// 执行逻辑
	claim := helper.GetJwtClaimsByContext(c)
	address.UserId = claim.Id
	err := address.Create()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "创建成功",
		"data":   "",
	})
}

func (a *AddressApi) update(c *gin.Context) {
	// 参数校验
	d, _ := c.GetRawData()
	var address models.UserAddress
	jsoniter.Unmarshal(d, &address)
	if err := base.ValidateStruct(&address); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	// 执行逻辑
	err := address.Update()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "更新成功",
		"data":   "",
	})
}

func (a *AddressApi) delete(c *gin.Context) {
	// 参数校验
	d, _ := c.GetRawData()
	var address models.UserAddress
	jsoniter.Unmarshal(d, &address)
	// 执行逻辑
	err := address.Delete()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "删除成功",
		"data":   "",
	})
}

func (a *AddressApi) list(c *gin.Context) {
	var address models.UserAddress
	offset, _ := c.GetQuery("offset")
	limit, _ := c.GetQuery("limit")
	// 执行逻辑
	claim := helper.GetJwtClaimsByContext(c)
	address.UserId = claim.Id
	list, err := address.List(offset, limit)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "",
		"data":   list,
	})
}
