package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"one/infra/base_c"
	"time"
)

type Product struct {
	gorm.Model
	Title       string  `gorm:"type:varchar(255)"`
	Description string  `gorm:"type:text"`
	Image       string  `gorm:"type:varchar(255)"`
	OnSale      uint8   `gorm:"type:tinyint unsigned;default:1"`
	Rating      float32 `gorm:"type:float unsigned;default:5"`
	SoldCount   uint    `gorm:"type:int unsigned;default:0"`
	ReviewCount uint    `gorm:"type:int unsigned;default:0"`
	Price       float32 `gorm:"type:decimal(10,2)"`
}

type ProductSku struct {
	gorm.Model
	Title       string  `gorm:"type:varchar(255)" validator:"required"`
	Description string  `gorm:"type:text" validator:"required"`
	Price       float32 `gorm:"type:decimal(10,2)" validator:"required"`
	Stock       uint    `gorm:"type:int unsigned" validator:"required"`
	ProductId   uint
	Product     Product
}

type CreateProduct struct {
	Title       string `validator:"required"`
	Description string `validator:"required"`
	Image       string `validator:"required"`
	OnSale      uint8  `validator:"required"`
	Sku         []ProductSku
}

// 新增商品
func (p *Product) Create(tx *gorm.DB) error {
	return tx.Save(&p).Error
}

// 批量出入 sku
func (p *ProductSku) Create(tx *gorm.DB, sku []ProductSku, pid uint) error {
	sql := "INSERT INTO `product_skus` (`title`,`description`,`price`,`stock`,`product_id`,`created_at`,`updated_at`) VALUES "
	// 循环data数组,组合sql语句
	now := time.Now().Format("2006-01-02 15:04:05")
	for key, value := range sku {
		if len(sku)-1 == key {
			//最后一条数据 以分号结尾
			sql += fmt.Sprintf("('%s','%s',%f,%d,%d,'%v','%v');", value.Title, value.Description, value.Price, value.Stock, pid, now, now)
		} else {
			sql += fmt.Sprintf("('%s','%s',%f,%d,%d,'%v','%v'),", value.Title, value.Description, value.Price, value.Stock, pid, now, now)
		}
	}
	return tx.Exec(sql).Error
}

// 插入单个sku
func (p *ProductSku) CreateSku() error {
	return base_c.Database().Save(&p).Error
}

// 更新商品
func (p *Product) Update() error {
	return base_c.Database().Model(&Product{}).Updates(&p).Error
}

// 更新单品
func (p *ProductSku) Update() error {
	return base_c.Database().Model(&ProductSku{}).Updates(&p).Error
}

// 删除商品
func (p *Product) Delete() error {
	tx := base_c.Database().Begin()
	err := tx.Delete(&p).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Where("product_id = ?", p.ID).Delete(&ProductSku{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 删除单品
func (p *ProductSku) Delete() error {
	return base_c.Database().Delete(&p).Error
}

// 商品列表
func (p *Product) List(offset, limit string) (*[]Product, error) {
	var list = make([]Product, 10, 10)
	err := base_c.Database().Model(&Product{}).Offset(offset).Limit(limit).Find(&list).Error
	return &list, err
}

// 单品列表
func (p *ProductSku) List(pid uint, offset, limit string) (*[]ProductSku, error) {
	var list = make([]ProductSku, 10, 10)
	err := base_c.Database().Model(&ProductSku{}).Offset(offset).Limit(limit).Where("product_id = ?", pid).Find(&list).Error
	return &list, err
}
