package router

import (
	"github.com/gin-gonic/gin"
	"pluto/middleware/filter"
)

// InitRechargeRouter 充值相关路由
func InitRechargeRouter(Router *gin.RouterGroup) {
	rechargeRouter := Router.Group("recharge").Use(filter.OperationRecord())
	{
		rechargeRouter.GET("getRechargeList", ApisCtl.GetRechargeList)
		rechargeRouter.POST("createRecharge", ApisCtl.CreateRecharge)
	}
}
