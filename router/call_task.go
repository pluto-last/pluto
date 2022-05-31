package router

import (
	"github.com/gin-gonic/gin"
	"pluto/middleware/filter"
)

// InitCallTaskRouter 初始化任务相关路由
func InitCallTaskRouter(Router *gin.RouterGroup) {
	sipRouter := Router.Group("call").Use(filter.OperationRecord())
	{
		sipRouter.GET("getMonthlyBillList", ApisCtl.GetMonthlyBillList)
		sipRouter.POST("createCallTask", ApisCtl.CreateCallTask)
		sipRouter.PUT("setCallTask", ApisCtl.SetCallTask)
		sipRouter.GET("getCallTaskList", ApisCtl.GetCallTaskList)
		sipRouter.GET("getCallTaskByID", ApisCtl.GetCallTaskByID)
		sipRouter.DELETE("deleteCallTask", ApisCtl.DeleteCallTask)
		sipRouter.POST("importCallPlan", ApisCtl.ImportCallPlan)
		sipRouter.POST("changeTaskStatus", ApisCtl.ChangeTaskStatus)
	}
}
