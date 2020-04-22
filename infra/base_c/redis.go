package base_c

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"one/infra"
	"time"
)

// 定义redis链接池
var client *redis.Client

func Redis() *redis.Client {
	return client
}

// redis starter， 并且设置为全局
type RedisStarter struct {
	infra.BaseStarter
}

func (s *RedisStarter) Setup(ctx infra.StarterContext) {
	addr := ctx.Props().GetDefault("redis.addr", "homestead")
	password := ctx.Props().GetDefault("redis.password", "secret")
	db := ctx.Props().GetIntDefault("redis.db", 0)
	poolsize := ctx.Props().GetIntDefault("redis.poolsize", 5)
	maxretries := ctx.Props().GetIntDefault("redis.maxretries", 3)
	idletimeout := ctx.Props().GetIntDefault("redis.idletimeout", 10)
	client = redis.NewClient(&redis.Options{
		Addr:        addr,                                     // Redis地址
		Password:    password,                                 // Redis账号
		DB:          db,                                       // Redis库
		PoolSize:    poolsize,                                 // Redis连接池大小
		MaxRetries:  maxretries,                               // 最大重试次数
		IdleTimeout: time.Duration(idletimeout) * time.Second, // 空闲链接超时时间
	})
	pong, err := client.Ping().Result()
	if err == redis.Nil {
		logrus.Info("Redis异常")
	} else if err != nil {
		logrus.Info("失败:", err)
	} else {
		logrus.Info(pong)
	}
}
