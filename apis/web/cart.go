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
	infra.RegisterApi(new(CartApi))
}

type CartApi struct {}

func (a *CartApi) Init() {
	r := base.Gin().Group("/v1").Use(middleware.Jwt())
	r.POST("/cart", a.create)
	r.PUT("/cart", a.update)
	r.DELETE("/cart", a.delete)
	r.GET("/cart", a.list)
}

func (a *CartApi) create(c *gin.Context) {
	// 参数校验
	d, _ := c.GetRawData()
	var item models.CartItem
	jsoniter.Unmarshal(d, &item)
	if err := base.ValidateStruct(&item); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	// 执行逻辑
	claim := helper.GetJwtClaimsByContext(c)
	item.UserId = claim.Id
	err := item.Create()
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
		"msg":    "添加成功",
		"data":   "",
	})
}

func (a *CartApi) update(c *gin.Context) {
	// 参数校验
	d, _ := c.GetRawData()
	var item models.CartItem
	jsoniter.Unmarshal(d, &item)
	if err := base.ValidateStruct(&item); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	// 执行逻辑
	claim := helper.GetJwtClaimsByContext(c)
	item.UserId = claim.Id
	err := item.Update()
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

func (a *CartApi) delete(c *gin.Context) {
	// 参数校验
	d, _ := c.GetRawData()
	var item models.CartItem
	jsoniter.Unmarshal(d, &item)
	// 执行逻辑
	err := item.Delete()
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

func (a *CartApi) list(c *gin.Context) {
	var item models.CartItem
	offset, _ := c.GetQuery("offset")
	limit, _ := c.GetQuery("limit")
	// 执行逻辑
	claim := helper.GetJwtClaimsByContext(c)
	list, err := item.List(claim.Id, offset, limit)
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
