package router

import (
	"github.com/gin-gonic/gin"
	"pluto/middleware/filter"
)

// InitRoleAuthRouter 角色权限相关
func InitRoleAuthRouter(Router *gin.RouterGroup) {
	roleAuthRouter := Router.Group("roleAuth").Use(filter.OperationRecord())
	{
		roleAuthRouter.GET("getRoleList", ApisCtl.GetRoleList)
		roleAuthRouter.POST("createRole", ApisCtl.CreateRole)
		roleAuthRouter.POST("setRole", ApisCtl.SetRole)
		roleAuthRouter.DELETE("deleteRole", ApisCtl.DeleteRole)
		roleAuthRouter.GET("getPermissionList", ApisCtl.GetPermissionList)
		roleAuthRouter.POST("createPermission", ApisCtl.CreatePermission)
		roleAuthRouter.POST("setPermission", ApisCtl.SetPermission)
		roleAuthRouter.DELETE("deletePermission", ApisCtl.DeletePermission)
		roleAuthRouter.POST("addRolePermission", ApisCtl.AddRolePermission)
		roleAuthRouter.DELETE("delRolePermission", ApisCtl.DelRolePermission)
		roleAuthRouter.GET("getRolePermissions", ApisCtl.GetRolePermissions)
		roleAuthRouter.POST("addUserRole", ApisCtl.AddUserRole)
		roleAuthRouter.DELETE("delUserRole", ApisCtl.DelUserRole)
		roleAuthRouter.GET("getUserRoles", ApisCtl.GetUserRoles)
		roleAuthRouter.GET("getUserAllPermissions", ApisCtl.GetUserAllPermissions)
	}
}
