package router

import (
	"github.com/gin-gonic/gin"
	"pluto/middleware/filter"
)

// InitSceneRouter 话术相关
func InitSceneRouter(Router *gin.RouterGroup) {
	sipRouter := Router.Group("scene").Use(filter.OperationRecord())
	{
		sipRouter.POST("createScene", ApisCtl.CreateScene)
		sipRouter.PUT("setSceneNode", ApisCtl.SetSceneNode)
		sipRouter.PUT("setScene", ApisCtl.SetScene)
		sipRouter.GET("getSceneList", ApisCtl.GetSceneList)
		sipRouter.GET("getSceneInfoByID", ApisCtl.GetSceneInfoByID)
		sipRouter.DELETE("deleteScene", ApisCtl.DeleteScene)
	}
}
