package api

import (
	"pluto/global"
	"pluto/model/params"
	"pluto/model/reply"
	"pluto/model/table"
	"pluto/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RoleAuthCtl struct {
}

// @Tags RoleAuth
// @Summary 分页获取角色列表
// @Produce  application/json
// @Param data query params.GetRoleAuthList true "分页获取角色列表"
// @Success 200
// @Router  /roleAuth/getRoleList [get]
func (ra *RoleAuthCtl) GetRoleList(c *gin.Context) {
	var pageInfo params.GetRoleAuthList
	_ = c.Bind(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := roleAuthService.GetRoleList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取角色列表失败!", zap.Any("err", err))
		reply.FailWithMessage("获取角色列表失败", c)
	} else {
		reply.OkWithDetailed(params.PageResult{
			List:   list,
			Total:  total,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
		}, "获取成功", c)
	}
}

// @Tags RoleAuth
// @Summary 创建角色
// @Produce  application/json
// @Param data body table.Role true "创建角色"
// @Success 200
// @Router  /roleAuth/createRole [post]
func (ra *RoleAuthCtl) CreateRole(c *gin.Context) {
	var role table.Role
	_ = c.Bind(&role)
	if err := utils.Verify(role, utils.RoleVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if err := roleAuthService.CreateRole(role); err != nil {
		global.GVA_LOG.Error("创建角色失败!", zap.Any("err", err))
		reply.FailWithMessage("创建角色失败", c)
	} else {
		reply.OkWithMessage("创建角色成功", c)
	}
}

// @Tags RoleAuth
// @Summary 修改角色信息
// @Produce  application/json
// @Param data body table.Role true "修改角色信息"
// @Success 200
// @Router  /roleAuth/setRole [post]
func (ra *RoleAuthCtl) SetRole(c *gin.Context) {
	var role table.Role
	_ = c.Bind(&role)
	if err := utils.Verify(role, utils.RoleVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if err := roleAuthService.SetRole(role); err != nil {
		global.GVA_LOG.Error("修改角色失败!", zap.Any("err", err))
		reply.FailWithMessage("修改角色失败", c)
	} else {
		reply.OkWithMessage("修改角色成功", c)
	}
}

// @Tags RoleAuth
// @Summary 删除角色信息
// @Produce  application/json
// @Param roleName query string true "删除角色信息"
// @Success 200
// @Router  /roleAuth/deleteRole [delete]
func (ra *RoleAuthCtl) DeleteRole(c *gin.Context) {
	roleName := c.Query("roleName")
	if roleName == "" {
		reply.FailWithMessage("roleName 不能为空", c)
	}
	if err := roleAuthService.DeleteRole(roleName); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		reply.FailWithMessage("删除失败", c)
	} else {
		reply.OkWithMessage("删除成功", c)
	}
}

// @Tags RoleAuth
// @Summary 查询所有权限
// @Produce  application/json
// @Param data query params.GetRoleAuthList true "查询所有权限"
// @Success 200
// @Router  /roleAuth/getPermissionList [get]
func (ra *RoleAuthCtl) GetPermissionList(c *gin.Context) {
	var pageInfo params.GetRoleAuthList
	_ = c.Bind(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := roleAuthService.GetPermissionList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取权限列表失败!", zap.Any("err", err))
		reply.FailWithMessage("获取权限列表失败", c)
	} else {
		reply.OkWithDetailed(params.PageResult{
			List:   list,
			Total:  total,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
		}, "获取成功", c)
	}
}

// @Tags RoleAuth
// @Summary 创建权限
// @Produce  application/json
// @Param data body table.Permission true "创建权限"
// @Success 200
// @Router  /roleAuth/createPermission [post]
func (ra *RoleAuthCtl) CreatePermission(c *gin.Context) {
	var permission table.Permission
	_ = c.Bind(&permission)
	if err := utils.Verify(permission, utils.RoleVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if err := roleAuthService.CreatePermission(permission); err != nil {
		global.GVA_LOG.Error("创建权限失败!", zap.Any("err", err))
		reply.FailWithMessage("创建权限失败", c)
	} else {
		reply.OkWithMessage("操作成功", c)
	}
}

// @Tags RoleAuth
// @Summary 修改权限信息
// @Produce  application/json
// @Param data body table.Permission true "修改权限信息"
// @Success 200
// @Router  /roleAuth/setPermission [post]
func (ra *RoleAuthCtl) SetPermission(c *gin.Context) {
	var permission table.Permission
	_ = c.Bind(&permission)
	if err := utils.Verify(permission, utils.RoleVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if err := roleAuthService.SetPermission(permission); err != nil {
		global.GVA_LOG.Error("修改权限失败!", zap.Any("err", err))
		reply.FailWithMessage("修改权限失败", c)
	} else {
		reply.OkWithMessage("操作成功", c)
	}
}

// @Tags RoleAuth
// @Summary 删除权限
// @Produce  application/json
// @Param permissionName query string true "删除权限"
// @Success 200
// @Router  /roleAuth/deletePermission [delete]
func (ra *RoleAuthCtl) DeletePermission(c *gin.Context) {
	permissionName := c.Query("permissionName")
	if permissionName == "" {
		reply.FailWithMessage("permissionName 不能为空", c)
	}
	if err := roleAuthService.DeletePermission(permissionName); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		reply.FailWithMessage("删除失败", c)
	} else {
		reply.OkWithMessage("删除成功", c)
	}
}

// @Tags RoleAuth
// @Summary 角色添加权限
// @Produce  application/json
// @Param data body params.UserRolePermission true "角色添加权限"
// @Success 200
// @Router  /roleAuth/addRolePermission [post]
func (ra *RoleAuthCtl) AddRolePermission(c *gin.Context) {
	var req params.UserRolePermission
	_ = c.Bind(&req)
	if err := utils.Verify(req, utils.RolePermissionVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if err := roleAuthService.AddRolePermission(req.Role, req.Permission); err != nil {
		global.GVA_LOG.Error("角色添加权限失败!", zap.Any("err", err))
		reply.FailWithMessage("角色添加权限失败", c)
	} else {
		reply.OkWithMessage("操作成功", c)
	}
}

// @Tags RoleAuth
// @Summary 角色删除权限的
// @Produce  application/json
// @Param data query params.UserRolePermission true "角色删除权限的"
// @Success 200
// @Router  /roleAuth/delRolePermission [delete]
func (ra *RoleAuthCtl) DelRolePermission(c *gin.Context) {
	var req params.UserRolePermission
	_ = c.Bind(&req)
	if err := utils.Verify(req, utils.RolePermissionVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if err := roleAuthService.DelRolePermission(req.Role, req.Permission); err != nil {
		global.GVA_LOG.Error("角色删除权限失败!", zap.Any("err", err))
		reply.FailWithMessage("角色删除权限失败", c)
	} else {
		reply.OkWithMessage("操作成功", c)
	}
}

// @Tags RoleAuth
// @Summary 获取全部角色和权限的关联关系
// @Produce  application/json
// @Param role query string true "获取全部角色和权限的关联关系"
// @Success 200
// @Router  /roleAuth/getRolePermissions [get]
func (ra *RoleAuthCtl) GetRolePermissions(c *gin.Context) {
	role := c.Query("role")
	if role == "" {
		reply.FailWithMessage("role不能为空", c)
	}
	permissions, err := roleAuthService.GetRolePermissions(role)
	if err != nil {
		global.GVA_LOG.Error("获取全部角色和权限的关联关系失败!", zap.Any("err", err))
		reply.FailWithMessage("获取全部角色和权限的关联关系失败", c)
		return
	}

	reply.OkWithDetailed(permissions, "获取成功", c)
}

// @Tags RoleAuth
// @Summary 用户添加角色
// @Produce  application/json
// @Param data body params.UserRolePermission true "角色添加权限"
// @Success 200
// @Router  /roleAuth/addUserRole [post]
func (ra *RoleAuthCtl) AddUserRole(c *gin.Context) {
	var req params.UserRolePermission
	_ = c.Bind(&req)
	if err := utils.Verify(req, utils.UserRoleVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if err := roleAuthService.AddUserRole(req.Role, req.UserID); err != nil {
		global.GVA_LOG.Error("用户添加到组失败!", zap.Any("err", err))
		reply.FailWithMessage("用户添加到组失败", c)
	} else {
		reply.OkWithMessage("操作成功", c)
	}
}

// @Tags RoleAuth
// @Summary 用户移除组
// @Produce  application/json
// @Param data query params.UserRolePermission true "用户移除组"
// @Success 200
// @Router  /roleAuth/delUserRole [delete]
func (ra *RoleAuthCtl) DelUserRole(c *gin.Context) {
	var req params.UserRolePermission
	_ = c.Bind(&req)
	if err := utils.Verify(req, utils.UserRoleVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if err := roleAuthService.DelUserRole(req.Role, req.Permission); err != nil {
		global.GVA_LOG.Error("用户移除组失败!", zap.Any("err", err))
		reply.FailWithMessage("用户移除组失败", c)
	} else {
		reply.OkWithMessage("操作成功", c)
	}
}

// @Tags RoleAuth
// @Summary 获取指定用户全部角色
// @Produce  application/json
// @Param userID query string true "获取指定用户全部角色"
// @Success 200
// @Router  /roleAuth/getUserRoles [get]
func (ra *RoleAuthCtl) GetUserRoles(c *gin.Context) {
	userID := c.Query("userID")
	if userID == "" {
		reply.FailWithMessage("userID不能为空", c)
	}
	roles, err := roleAuthService.GetUserRoles(userID)
	if err != nil {
		global.GVA_LOG.Error("获取指定用户全部角色失败!", zap.Any("err", err))
		reply.FailWithMessage("获取指定用户全部角色失败", c)
		return
	}

	reply.OkWithDetailed(roles, "获取成功", c)
}

// @Tags RoleAuth
// @Summary 获取用户的所有权限
// @Produce  application/json
// @Param userID query string true "获取用户的所有权限"
// @Success 200
// @Router  /roleAuth/getUserAllPermissions [get]
func (ra *RoleAuthCtl) GetUserAllPermissions(c *gin.Context) {
	userID := c.Query("userID")
	if userID == "" {
		reply.FailWithMessage("userID不能为空", c)
	}
	userPermissions, err := roleAuthService.GetUserAllPermissions(userID)
	if err != nil {
		global.GVA_LOG.Error("获取用户的所有权限失败!", zap.Any("err", err))
		reply.FailWithMessage("获取用户的所有权限失败", c)
		return
	}

	reply.OkWithDetailed(userPermissions, "获取成功", c)
}
