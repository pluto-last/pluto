package router

import (
	"pluto/api"
	"pluto/global"
	"pluto/middleware/filter"
	"pluto/middleware/jwt"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "pluto/docs"
)

var ApisCtl = api.ApisGroupsAPP

func Routers() *gin.Engine {
	var Router = gin.Default()

	//Router.Use(middleware.LoadTls())  //https了

	Router.Use(filter.Cors()) // 跨域

	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//获取路由组实例
	PublicGroup := Router.Group("v1")
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
	}
	{
		InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
		InitWSRouter(PublicGroup)   // ws_client连接的请求
	}
	PrivateGroup := Router.Group("v1")
	PrivateGroup.Use(jwt.JWTAuth()) // 私有路由都需要JWT鉴权
	{
		InitUserRouter(PrivateGroup)     // 注册用户路由
		InitRoleAuthRouter(PrivateGroup) // 角色权限相关的路由
		InitSipRouter(PrivateGroup)      // 注册线路相关路由
		InitRechargeRouter(PrivateGroup) // 注册充值相关的路由
		InitCallTaskRouter(PrivateGroup) // 推广任务相关的路由
		InitSceneRouter(PrivateGroup)    // 话术相关的路由

	}

	global.GVA_LOG.Info("router register success")
	return Router
}
