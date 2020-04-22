package base

import (
	"gopkg.in/go-playground/validator.v9"
	"one/infra/base_c"
	"one/models"
)

// 自定义验证函数
func onSale(v validator.FieldLevel) bool {
	var sku models.ProductSku
	base_c.Database().Model(&models.ProductSku{}).First(&sku, "id = ?", v.Field().Uint())
	if sku.DeletedAt == nil && sku.ID > 0 {
		return true
	}
	return false
}
