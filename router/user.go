package router

import (
	"pluto/middleware/auth"
	"pluto/middleware/filter"

	"github.com/gin-gonic/gin"
)

// InitUserRouter 用户相关
func InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user").Use(filter.OperationRecord())
	{
		userRouter.DELETE("deleteUser", ApisCtl.DeleteUser)                                                  // 删除用户
		userRouter.PUT("setUserInfo", ApisCtl.SetUserInfo)                                                   // 设置用户信息
		userRouter.GET("getUserList", filter.HasPermission(auth.PermissionGetUserList), ApisCtl.GetUserList) // 分页获取用户列表
		userRouter.GET("getUserInfo", ApisCtl.GetUserInfo)                                                   // 获取自身信息
		userRouter.POST("batchCreateUser", ApisCtl.BatchCreateUser)
		userRouter.POST("resetPasswordByID", ApisCtl.ResetPasswordByID)
	}
}
