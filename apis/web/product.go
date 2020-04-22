package web

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"net/http"
	"one/infra"
	"one/infra/base"
	"one/infra/base_c"
	"one/infra/middleware"
	"one/models"
	"strconv"
)

func init() {
	infra.RegisterApi(new(ProductApi))
}

type ProductApi struct {}

func (p *ProductApi) Init() {
	r := base.Gin().Group("/v1/admin").Use(middleware.Jwt()).Use(middleware.Auth())
	r.POST("/product", p.create)
	r.POST("/sku", p.skuCreate)
	r.PUT("/product", p.update)
	r.PUT("/sku", p.skuUpdate)
	r.DELETE("/product", p.delete)
	r.DELETE("/sku", p.skuDelete)
	r.GET("/product", p.list)
	r.GET("/sku", p.skuList)
}

func (p *ProductApi) create(c *gin.Context) {
	d, _ := c.GetRawData()
	var data models.CreateProduct
	jsoniter.Unmarshal(d, &data)
	// 参数校验
	if err := base.ValidateStruct(&data); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	// 执行逻辑
	var product models.Product
	product.Title = data.Title
	product.Image = data.Image
	product.Description = data.Description
	product.OnSale = data.OnSale
	tx := base_c.Database().Begin()
	err := product.Create(tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	var sku models.ProductSku
	err = sku.Create(tx, data.Sku, product.ID)
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
		"data":   data,
	})
}

func (p *ProductApi) skuCreate(c *gin.Context) {
	d, _ := c.GetRawData()
	var sku models.ProductSku
	jsoniter.Unmarshal(d, &sku)
	// 参数校验
	if err := base.ValidateStruct(&sku); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	// 执行逻辑
	err := sku.CreateSku()
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

func (p *ProductApi) update(c *gin.Context) {
	d, _ := c.GetRawData()
	var product models.Product
	jsoniter.Unmarshal(d, &product)
	logrus.Info(product)
	err := product.Update()
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
		"msg":    "操作成功",
		"data":   "",
	})
}

func (p *ProductApi) skuUpdate(c *gin.Context) {
	d, _ := c.GetRawData()
	var sku models.ProductSku
	jsoniter.Unmarshal(d, &sku)
	err := sku.Update()
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
		"msg":    "操作成功",
		"data":   "",
	})
}

func (p *ProductApi) delete(c *gin.Context) {
	d, _ := c.GetRawData()
	var product models.Product
	jsoniter.Unmarshal(d, &product)
	err := product.Delete()
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
		"msg":    "操作成功",
		"data":   "",
	})
}

func (p *ProductApi) skuDelete(c *gin.Context) {
	d, _ := c.GetRawData()
	var sku models.ProductSku
	jsoniter.Unmarshal(d, &sku)
	err := sku.Delete()
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
		"msg":    "操作成功",
		"data":   "",
	})
}

func (p *ProductApi) list(c *gin.Context) {
	var product models.Product
	offset, _ := c.GetQuery("offset")
	limit, _ := c.GetQuery("limit")
	// 执行逻辑
	list, err := product.List(offset, limit)
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

func (p *ProductApi) skuList(c *gin.Context) {
	var sku models.ProductSku
	offset, _ := c.GetQuery("offset")
	limit, _ := c.GetQuery("limit")
	pid, _ := c.GetQuery("productId")
	_pid, _ := strconv.Atoi(pid)
	// 执行逻辑
	list, err := sku.List(uint(_pid), offset, limit)
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