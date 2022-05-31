package main

import (
	"fmt"
	"go.uber.org/zap"
	"log"
	"pluto/config"
	"pluto/global"
	"pluto/middleware/auth"
	"pluto/middleware/cache"
	"pluto/middleware/db"
	myzap "pluto/middleware/zap"
	"pluto/router"
)

func InitServer() {
	var err error
	global.GVA_CONFIG, err = config.Get(configPath) // 初始化配置
	if err != nil {
		log.Fatalln("read conf file : ", err)
	}

	global.GVA_LOG = myzap.Zap() // 初始化zap日志库

	global.GVA_DB, err = db.Instance(global.GVA_CONFIG) // 初始化数据库库
	if err != nil {
		log.Fatalln("init db err : ", err)
	}

	global.GVA_REDIS, err = cache.InitRedis(global.GVA_CONFIG.Cache.Addr, global.GVA_CONFIG.Cache.Password, 0) // 初始化redis服务
	if err != nil {
		log.Fatalln("redis init err : ", err)
	}

	auth.InitPermission() // 初始化权限
}

// StartServer 启动服务
func StartServer() {
	// 初始化路由
	Router := router.Routers()

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.Port)

	// 启动服务
	Router.Run(address)
	global.GVA_LOG.Info("server run success on ", zap.String("address", address))
}
