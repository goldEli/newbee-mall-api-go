package initialize

import (

	// "gologin/config"

	"github.com/go-redis/redis"

	// "github.com/sirupsen/logrus"
	"main.go/global"
)

func RedisInit() *redis.Client {
	// addr := fmt.Sprintf("%s:%d", config.Env.Redis.Host, config.Env.Redis.Port)
	addr := global.GVA_CONFIG.Redis.Addr()

	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
		// Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		// global.GVA_LOG.Error("redis error", zap.Error(err))
		panic(err)
	}
	global.GVA_LOG.Info("redis success: " + addr)

	return redisClient
}
