package models

import (
	"errors"
	"one/infra/base_c"
)

type CartItem struct {
	Id           uint `gorm:"primary_key"`
	UserId       uint `gorm:"type:int;not null"`
	ProductSkuId uint `gorm:"type:int;not null" validate:"required,onSale"`
	Amount       uint `gorm:"type:int unsigned;not null" validate:"required"`
	ProductSku   ProductSku
}

type CartInfo struct {
	Title string
	Image string
	Des string
	Price float32
	Amount uint
	OnSale uint8
	ProductSkuId uint
	ProductId uint
	Id uint
}

// 添加购物车
func (c *CartItem) Create() error {
	db := base_c.Database()
	// 如果商品已存在，累加 amount
	var item CartItem
	if db.Model(&CartItem{}).First(&item, "product_sku_id = ? and user_id = ?", c.ProductSkuId, c.UserId).RecordNotFound() {
		return db.Save(&c).Error
	} else {
		amount := c.Amount + item.Amount
		return db.Model(&CartItem{}).Where("product_sku_id = ? and user_id = ?", c.ProductSkuId, c.UserId).UpdateColumn("amount", amount).Error
	}
}

// 更新购物车
func (c *CartItem) Update() error {
	return base_c.Database().Model(&CartItem{}).Where("product_sku_id = ? and user_id = ?", c.ProductSkuId, c.UserId).UpdateColumn("amount", c.Amount).Error
}

// 删除购物车商品
func (c *CartItem) Delete() error {
	return base_c.Database().Delete(&c).Error
}

// 查看购物车
func (c *CartItem) List(uid uint, offset, limit string) (*[]CartInfo, error) {
	db := base_c.Database()
	var list []CartItem
	err := db.Preload("ProductSku").Preload("ProductSku.Product").Offset(offset).Limit(limit).Where("user_id = ?", uid).Find(&list).Error
	if err != nil {
		return nil, err
	}
	if list == nil {
		return nil, errors.New("购物车还是空的哦")
	}
	l := len(list)
	var row = make([]CartInfo, l, l)
	var del = make([]uint, 5, 5)
	var count,skip int
	for i, item := range list {
		// 跳过已删除的
		if item.ProductSku.ID == 0  {
			del = append(del, item.ProductSkuId)
			count++
			continue
		}
		row[i].Title = item.ProductSku.Product.Title
		row[i].Image = item.ProductSku.Product.Image
		row[i].Des = item.ProductSku.Description
		row[i].Price = item.ProductSku.Price
		row[i].Amount = item.Amount
		row[i].OnSale = item.ProductSku.Product.OnSale
		row[i].ProductSkuId = item.ProductSku.ID
		row[i].ProductId = item.ProductSku.Product.ID
		row[i].Id = item.Id
	}
	// 如果有删除的 sku，所有购物车清除该 sku
	if del != nil {
		err = db.Where("product_sku_id in (?)", del).Delete(&CartItem{}).Error
		if err != nil {
			return nil, err
		}
	}
	// 重新整理数据，去除空切片
	var data = make([]CartInfo, l-count, l-count)
	for i, v := range row {
		if v.Title == "" {
			skip++
			continue
		}
		i = i - skip
		data[i].Title = v.Title
		data[i].Image = v.Image
		data[i].Des = v.Des
		data[i].Price = v.Price
		data[i].Amount = v.Amount
		data[i].OnSale = v.OnSale
		data[i].ProductId = v.ProductId
		data[i].ProductSkuId = v.ProductSkuId
		data[i].Id = v.Id
	}
	return &data, err
}