package test

import (
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	"one/infra"
	"one/infra/base"
	"one/infra/base_c"
	"one/infra/migrate"
)

func init() {
	// 获取程序运行文件所在路径
	file := kvs.GetCurrentFilePath("../brun/config.ini", 1)
	// 加载和解析配置文件
	conf := ini.NewIniFileCompositeConfigSource(file)
	// 自定义 logrus
	base.InitLogrus()
	infra.Register(&base.PropsStarter{})
	infra.Register(&base_c.EmailStarter{})
	infra.Register(&base_c.DatabaseStarter{})
	infra.Register(&migrate.MigrateStarter{})
	infra.Register(&base_c.RedisStarter{})
	infra.Register(&base.ValidatorStarter{})
	app := infra.New(conf)
	app.Start()
}
