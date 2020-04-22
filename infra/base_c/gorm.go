package base_c

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"one/infra"
	"time"
)

// db 数据库实例
var db *gorm.DB

func Database() *gorm.DB {
	return db
}

// dbx 数据库starter， 并且设置为全局
type DatabaseStarter struct {
	infra.BaseStarter
}

func (s *DatabaseStarter) Setup(ctx infra.StarterContext) {
	conf := ctx.Props()
	user := conf.GetDefault("mysql.user", "homestead")
	pwd := conf.GetDefault("mysql.password", "secret")
	host := conf.GetDefault("mysql.host", "192.016.10.10")
	port := conf.GetDefault("mysql.port", "3306")
	database := conf.GetDefault("mysql.database", "one")
	handle, err := gorm.Open("mysql", user+":"+pwd+"@tcp("+host+":"+port+")/"+database+"?parseTime=true&charset=utf8mb4")
	//检查数据库是否连接成功
	if handle.Error != nil {
		panic(err)
		return
	}
	maxLifetime := conf.GetIntDefault("connMaxLifetime", 12)
	maxIdleConns := conf.GetIntDefault("maxIdleConns", 1)
	maxOpenConns := conf.GetIntDefault("maxOpenConns", 3)
	// SetConnMaxLifetime 设置可重用连接的最长时间
	handle.DB().SetConnMaxLifetime(time.Duration(maxLifetime) * time.Hour)
	// SetMaxIdleConns 设置空闲连接池中的最大连接数。
	handle.DB().SetMaxIdleConns(maxIdleConns)
	// SetMaxOpenConns 设置数据库连接最大打开数。
	handle.DB().SetMaxOpenConns(maxOpenConns)
	//handle.LogMode(true)打印执行sql方便调试
	logrus.Info(handle.DB().Ping())
	db = handle
}