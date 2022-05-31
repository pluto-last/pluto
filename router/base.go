package router

import (
	"github.com/gin-gonic/gin"
	"pluto/middleware/filter"
)

// InitBaseRouter 初始化公共路由
func InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("base").Use(filter.OperationRecord())
	userRouter := Router.Group("user").Use(filter.OperationRecord())
	{
		baseRouter.GET("captcha", ApisCtl.Captcha)
		baseRouter.POST("sms", ApisCtl.SMS)
		userRouter.POST("login", ApisCtl.Login)
		userRouter.POST("register", ApisCtl.Register)
		userRouter.POST("resetPassword", ApisCtl.ResetPassword)
	}
	return baseRouter
}
