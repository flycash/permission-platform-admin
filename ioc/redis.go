package ioc

import (
	"github.com/gotomicro/ego/core/econf"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	type Config struct {
		Addr string
	}
	var cfg Config
	err := econf.UnmarshalKey("redis", &cfg)
	if err != nil {
		panic(err)
	}
	cmd := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})
	return cmd
}
