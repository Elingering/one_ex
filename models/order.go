package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"one/infra/base_c"
	"strconv"
	"time"
)

type Order struct {
	gorm.Model
	No            string    `gorm:"type:varchar(255);not null;unique_index"`
	UserId        uint      `gorm:"type:int;not null"`
	Address       string    `gorm:"type:text"`
	TotalAmount   float32   `gorm:"type:decimal(10,2)" validator:"required"`
	Remark        string    `gorm:"type:text"`
	PaidAt        time.Time `gorm:"default:null"`
	PaymentMethod string    `gorm:"type:varchar(255);not null"`
	PaymentNo     string    `gorm:"type:varchar(255);not null"`
	RefundStatus  string    `gorm:"type:varchar(255)"`
	RefundNo      string    `gorm:"type:varchar(255);not null;unique_index"`
	Closed        uint8     `gorm:"type:tinyint unsigned;not null"`
	Reviewed      uint8     `gorm:"type:tinyint unsigned;not null"`
	ShipStatus    string    `gorm:"type:varchar(255);not null"`
	ShipData      string    `gorm:"type:text"`
	Extra         string    `gorm:"type:text"`
	User          User
	OrderItem     []OrderItem
}

type OrderItem struct {
	gorm.Model
	OrderId      uint      `gorm:"type:int;not null"`
	ProductId    uint      `gorm:"type:int;not null"`
	ProductSkuId uint      `gorm:"type:int;not null"`
	Amount       uint      `gorm:"type:int unsigned;not null" validate:"required"`
	Price        float32   `gorm:"type:decimal(10,2)" validator:"required"`
	Rating       float32   `gorm:"type:float unsigned;default:5"`
	Review       string    `gorm:"type:text"`
	ReviewedAt   time.Time `gorm:"default:null"`
	ProductSku   ProductSku
}

type CreateOrder struct {
	Address       uint
	Extra         string
	Item 		  []OrderItem
}

type OrderInfo struct {
	ItemInfo []ItemInfo
	Status string
	Total float32
	Address string
	Extra string
	No string
}

type ItemInfo struct {
	Image string
	Title string
	Des string
	Price float32
	Amount uint
	SubTotal float32
}

// 生产订单
func (o *Order) Create(tx *gorm.DB) error {
	return tx.Save(&o).Error
}

// 批量出入 sku
func (o *OrderItem) Create(tx *gorm.DB, item []OrderItem, oid uint) error {
	// 先计算价格
	l := len(item)
	var sku = make([]uint, l, l)
	var ids = make([]uint, l, l)
	for _, v := range item {
		sku = append(sku, v.ProductSkuId)
		ids = append(ids, v.ID)
	}
	var productSku []ProductSku
	tx.Where("id in (?)", sku).Find(&productSku)
	var tmp = make(map[uint]float32, l)
	for _, v := range productSku {
		tmp[v.ID] = v.Price
	}
	sql := "INSERT INTO `order_items` (`order_id`,`product_id`,`product_sku_id`,`amount`,`price`,`created_at`,`updated_at`) VALUES "
	// 循环data数组,组合sql语句
	now := time.Now().Format("2006-01-02 15:04:05")
	var total float32
	for key, value := range item {
		price := float32(value.Amount) * tmp[value.ProductSkuId]
		total += price
		if len(item)-1 == key {
			//最后一条数据 以分号结尾
			sql += fmt.Sprintf("(%d,%d,%d,%d,%f,'%v','%v');", oid, value.ProductId, value.ProductSkuId, value.Amount, price, now, now)
		} else {
			sql += fmt.Sprintf("(%d,%d,%d,%d,%f,'%v','%v'),", oid, value.ProductId, value.ProductSkuId, value.Amount, price, now, now)
		}
	}
	err := tx.Exec(sql).Error
	if err != nil {
		return err
	}
	// 更新订单信息
	no := GeneraOrderNo()
	err = tx.Model(&Order{}).Where("id = ?", oid).UpdateColumns(Order{No: no, TotalAmount: total, RefundNo:strconv.Itoa(int(oid))}).Error
	// 清空购物车
	err = tx.Where("id in (?)", ids).Delete(&CartItem{}).Error
	if err != nil {
		return err
	}
	return err
}

// 更新订单
func (o *Order) Update() error {
	return base_c.Database().Model(&Order{}).Where("id = ? and user_id = ?", o.ID, o.UserId).Updates(&o).Error
}

// 删除订单
func (o *Order) Delete() error {
	return base_c.Database().Delete(&o).Error
}

// 查看订单详情
func (o *Order) List(uid uint) (*OrderInfo, error) {
	db := base_c.Database()
	err := db.Preload("OrderItem").Preload("OrderItem.ProductSku").Preload("OrderItem.ProductSku.Product").Find(&o).Error
	if err != nil {
		return nil, err
	}
	if o.UserId != uid {
		return nil, errors.New("非法操作")
	}
	l := len(o.OrderItem)
	var row = make([]ItemInfo, l, l)
	for i, item := range o.OrderItem {
		row[i].Title = item.ProductSku.Product.Title
		row[i].Image = item.ProductSku.Product.Image
		row[i].Des = item.ProductSku.Description
		row[i].Price = item.ProductSku.Price
		row[i].Amount = item.Amount
		row[i].SubTotal = item.Price
	}
	var data OrderInfo
	data.ItemInfo = row
	data.Address = o.Address
	data.No = o.No
	data.Extra = o.Extra
	if o.PaidAt.IsZero() {
		data.Status = "未支付"
	} else {
		data.Status = "已支付"
	}
	data.Total = o.TotalAmount
	return &data, err
}