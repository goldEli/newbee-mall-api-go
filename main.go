package main

import (
	"time"

	"main.go/core"
	"main.go/global"
	"main.go/initialize"
)

// @title Mall
// @version 1.0
// @description Mall description
// @termsOfService http://swagger.io/terms/

// @contact.name miaoyu2009@qq.com
// @contact.url miaoyu2009@qq.com
// @contact.email miaoyu2009@qq.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8888
// @BasePath /manage-api/v1/
func main() {

	global.GVA_VP = core.Viper()              // 初始化Viper
	global.GVA_LOG = core.Zap()               // 初始化zap日志库
	global.GVA_REDIS = initialize.RedisInit() // redis
	global.GVA_DB = initialize.Gorm()         // gorm连接数据库
	global.GVA_REDIS.Set("server_start_time", time.Now(), time.Hour).Err()

	core.RunWindowsServer()
}
