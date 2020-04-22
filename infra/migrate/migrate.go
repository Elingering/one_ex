package migrate

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"one/infra"
	"one/infra/base_c"
	"one/models"
)

type MigrateStarter struct {
	infra.BaseStarter
}

func (s *MigrateStarter) Setup(ctx infra.StarterContext) {
	base_c.Database().AutoMigrate(
		&models.User{},
		&models.UserAddress{},
		&models.Admin{},
		&models.Role{},
		&models.Product{},
		&models.ProductSku{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		)
}