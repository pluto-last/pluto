package auth

import (
	"pluto/model/table"
	"pluto/service"
)

var roleAuthService = service.ServiceGroupAPP.RoleAuthService

var Roles = []*table.Role{
	{Name: "Admin", Describe: "管理员角色", FullDescribe: "具有管理员操作权限"},
}

const (
	PermissionGetUserList = "GetUserList"
)

var Permissions = []*table.Permission{
	{Name: PermissionGetUserList, Describe: "获取用户列表", FullDescribe: "获取用户列表"},
}

// InitPermission 初始化权限
func InitPermission() {
	for _, v := range Roles {
		roleAuthService.AddRole(v)
	}

	for _, v := range Permissions {
		roleAuthService.AddPermission(v)
		roleAuthService.AddRolePermission("Admin", v.Name)
	}
}
