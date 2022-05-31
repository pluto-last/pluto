package filter

import (
	"fmt"
	"pluto/model/reply"
	"pluto/service"

	"github.com/gin-gonic/gin"
)

var roleAuthService = service.ServiceGroupAPP.RoleAuthService

// HasPermission 是否拥有指定权限的中间件
func HasPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, permission := range permissions {
			userIDStr, exists := c.Get("userID")
			if !exists {
				reply.FailWithMessage("鉴权获取用户信息失败", c)
				c.Abort()
				return
			}
			userID := fmt.Sprint(userIDStr)
			if !roleAuthService.HasPermission(userID, permission) {
				reply.FailWithMessage("您没有权限访问", c)
				c.Abort()
				return
			}
		}
		// 处理请求
		c.Next()
	}
}

// HasRole 是否拥有指定角色的中间件
func HasRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, role := range roles {
			userIDStr, exists := c.Get("userID")
			if !exists {
				reply.FailWithMessage("鉴权获取用户信息失败", c)
				c.Abort()
				return
			}
			userID := fmt.Sprint(userIDStr)
			if !roleAuthService.HasRole(userID, role) {
				reply.FailWithMessage("您没有权限访问", c)
				c.Abort()
				return
			}
		}
		// 处理请求
		c.Next()
	}
}
