package router

import (
	"github.com/gin-gonic/gin"
)

// InitWSRouter websocket相关
func InitWSRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("ws").Use()
	{
		userRouter.GET("client", ApisCtl.BotWebsocket) // 分页获取用户列表
	}
}
