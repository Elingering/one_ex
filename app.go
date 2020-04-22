package one

import (
	_ "one/apis/web"
	"one/infra"
	"one/infra/base"
	"one/infra/base_c"
	"one/infra/migrate"
)

func init() {
	// 自定义 logrus
	base.InitLogrus()
	infra.Register(&base.PropsStarter{})
	infra.Register(&base_c.EmailStarter{})
	infra.Register(&base_c.DatabaseStarter{})
	infra.Register(&migrate.MigrateStarter{})
	infra.Register(&base_c.RedisStarter{})
	infra.Register(&base.ValidatorStarter{})
	//infra.Register(&base.GoRPCStarter{})
	//infra.Register(&gorpc.GoRpcApiStarter{})
	//infra.Register(&jobs.RefundExpiredJobStarter{})
	infra.Register(&base.GinStarter{})
	infra.Register(&infra.WebApiStarter{})
}
