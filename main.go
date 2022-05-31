package main

import (
	"flag"
	"fmt"
)

var buildTime = "no time info"
var buildVersion = "no git version info"
var buildGoVersion = "no go version info"
var buildBy = "no builder info"

var configPath string

// @title Swagger pluto API
// @version 1.0.0
// @description pluto 永远滴神!!!
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name token
// @BasePath
func main() {
	fmt.Println("BuildTime:", buildTime)
	fmt.Println("BuildVersion:", buildVersion)
	fmt.Println("BuildGoVersion:", buildGoVersion)
	fmt.Println("BuildBy:", buildBy)

	flag.StringVar(&configPath, "config", "./conf/conf.yaml", "config file path")
	flag.Parse()

	InitServer()

	// 启动服务
	StartServer()

}
