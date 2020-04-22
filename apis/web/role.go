package web

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"one/infra"
	"one/infra/base"
	"one/infra/base_c"
	"one/infra/middleware"
	"one/models"
)

func init() {
	infra.RegisterApi(new(RoleApi))
}

type RoleApi struct {}

func (r *RoleApi) Init() {
	groupRouter := base.Gin().Group("/v1/admin")
	m := groupRouter.Use(middleware.Jwt()).Use(middleware.Auth())
	m.POST("/role", r.create)
	m.GET("/role", r.list)
	m.PUT("/role", r.update)
	m.DELETE("/role", r.delete)
	m.POST("/role/auth", r.auth)
	m.POST("/role/assign", r.assign)
}

// 新增角色
func (r *RoleApi) create(c *gin.Context) {
	// 参数校验
	d, _ := c.GetRawData()
	var role models.Role
	jsoniter.Unmarshal(d, &role)
	if err := base.ValidateStruct(&role); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	// 执行逻辑
	err := role.Create()
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

// 更新角色
func (r *RoleApi) update(c *gin.Context) {
	// 参数校验
	d, _ := c.GetRawData()
	var role models.Role
	jsoniter.Unmarshal(d, &role)
	if err := base.ValidateStruct(&role); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	// 执行逻辑
	// 角色的唯一标识，不可变更。接收后设为空值，struct更新自动忽略
	role.RoleName = ""
	err := role.Update()
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

// 删除角色
func (r *RoleApi) delete(c *gin.Context) {
	// 参数校验
	d, _ := c.GetRawData()
	var role models.Role
	jsoniter.Unmarshal(d, &role)
	// 执行逻辑
	err := role.Delete()
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

// 角色列表
func (r *RoleApi) list(c *gin.Context) {
	var role models.Role
	offset, _ := c.GetQuery("offset")
	limit, _ := c.GetQuery("limit")
	// 执行逻辑
	list, err := role.List(offset, limit)
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

// 角色授权
func (r *RoleApi) auth(c *gin.Context) {
	d, _ := c.GetRawData()
	var casbin models.CasbinModel
	jsoniter.Unmarshal(d, &casbin)
	ok := casbin.AddCasbin(casbin)
	if ok {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"msg":     "保存成功",
			"data":"",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"msg":     "保存失败",
			"data":"",
		})
	}
}

// 分配角色
func (r *RoleApi) assign(c *gin.Context) {
	// 参数校验
	d, _ := c.GetRawData()
	var role models.AdminRole
	jsoniter.Unmarshal(d, &role)
	if err := base.ValidateStruct(&role); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    err.Error(),
			"data":   "",
		})
		return
	}
	// 执行逻辑
	tx := base_c.Database().Begin()
	err := role.Update(tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":     "保存失败",
			"data":"",
		})
		return
	}
	var casbin models.CasbinModel
	ok := casbin.LinkCasbin(role.AdminRole, role.RoleName)
	if !ok {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":     "保存失败",
			"data":"",
		})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":     "保存成功",
		"data":"",
	})
}
