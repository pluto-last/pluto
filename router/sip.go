package router

import (
	"github.com/gin-gonic/gin"
	"pluto/middleware/filter"
)

// InitSipRouter 线路相关
func InitSipRouter(Router *gin.RouterGroup) {
	sipRouter := Router.Group("sip").Use(filter.OperationRecord())
	{
		sipRouter.POST("createSip", ApisCtl.CreateSip)
		sipRouter.GET("getSipList", ApisCtl.GetSipList)
		sipRouter.GET("getSipInfo", ApisCtl.GetSipInfo)
		sipRouter.PUT("setSipInfo", ApisCtl.SetSipInfo)
		sipRouter.DELETE("deleteSip", ApisCtl.DeleteSip)
		sipRouter.GET("getUserSips", ApisCtl.GetUserSips)
		sipRouter.POST("userAddSip", ApisCtl.UserAddSip)
		sipRouter.PUT("setUserSip", ApisCtl.SetUserSip)
		sipRouter.DELETE("userDelSip", ApisCtl.UserDelSip)
	}
}
