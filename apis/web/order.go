package web

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"one/infra"
	"one/infra/base"
	"one/infra/base_c"
	"one/infra/helper"
	"one/infra/middleware"
	"one/models"
	"strconv"
	"time"
)

func init() {
	infra.RegisterApi(new(OrderApi))
}

type OrderApi struct {}

func (a *OrderApi) Init() {
	r := base.Gin().Group("/v1").Use(middleware.Jwt())
	r.POST("/order", a.create)
	r.PUT("/order", a.update)
	r.DELETE("/order", a.delete)
	r.GET("/order", a.detail)
}

func (a *OrderApi) create(c *gin.Context) {
	d, _ := c.GetRawData()
	var row models.CreateOrder
	jsoniter.Unmarshal(d, &row)
	// 参数校验
	if err := base.ValidateStruct(&row); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	// 执行逻辑
	claim := helper.GetJwtClaimsByContext(c)
	var order models.Order
	order.Extra = row.Extra
	// 查询用户地址
	addr, err := models.GetByIdAndUserId(row.Address, claim.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	address := addr.Province + addr.City + addr.District + addr.Address + " " + strconv.Itoa(addr.Zip) + " " + addr.ContactName + " " + addr.ContactPhone
	order.Address = address
	order.UserId = claim.Id
	tx := base_c.Database().Begin()
	err = order.Create(tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	var item models.OrderItem
	err = item.Create(tx, row.Item, order.ID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "创建成功",
		"data":   "",
	})
}

func (a *OrderApi) update(c *gin.Context) {
	// 参数校验
	d, _ := c.GetRawData()
	var order models.Order
	jsoniter.Unmarshal(d, &order)
	// 执行逻辑
	if order.PaymentNo != "" {
		order.PaidAt = time.Now()
	}
	claim := helper.GetJwtClaimsByContext(c)
	order.UserId = claim.Id
	err := order.Update()
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

func (a *OrderApi) delete(c *gin.Context) {
	// 参数校验
	d, _ := c.GetRawData()
	var order models.Order
	jsoniter.Unmarshal(d, &order)
	// 执行逻辑
	err := order.Delete()
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

func (a *OrderApi) detail(c *gin.Context) {
	var order models.Order
	d, _ := c.GetRawData()
	jsoniter.Unmarshal(d, &order)
	// 执行逻辑
	claim := helper.GetJwtClaimsByContext(c)
	list, err := order.List(claim.Id)
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
