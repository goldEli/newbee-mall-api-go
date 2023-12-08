package global

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main.go/config"
)

var (
	GVA_DB     *gorm.DB
	GVA_VP     *viper.Viper
	GVA_LOG    *zap.Logger
	GVA_CONFIG config.Server
	GVA_REDIS  *redis.Client
)
