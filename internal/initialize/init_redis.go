package initialize

import (
	"fmt"
	"id-maker/config"

	"github.com/go-redis/redis"
)

var redisDb *redis.Client

// 初始化日志
func Init_Redis(cfg *config.RedisConfig) (err error) {

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Password: fmt.Sprintf("%s", cfg.Password), // no password set
		DB:       cfg.Db,                          // use default DB
		PoolSize: cfg.PoolSize,
	})

	_, err = client.Ping().Result()
	return err
}

func Close_Redis() {
	_ = redisDb.Close()
}
