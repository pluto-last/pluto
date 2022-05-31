package service

import (
	"errors"
	"pluto/global"
	"pluto/middleware/db"
	"pluto/model/params"
	"pluto/model/table"
	"strings"
)

type RoleAuthService struct{}

// GetAllRole 查询所有角色
func (roleAuthService *RoleAuthService) GetRoleList(info params.GetRoleAuthList) (roles []table.Role, total int64, err error) {
	db := global.GVA_DB.Model(&table.Role{})
	if info.Name != "" {
		db = db.Where("name like ?", "%"+info.Name+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(info.Limit).Offset(info.Offset).Find(&roles).Error
	return roles, total, err

}

// CreateRole 创建角色
func (roleAuthService *RoleAuthService) CreateRole(role table.Role) (err error) {
	existRole := new(table.Role)
	notFound := global.GVA_DB.Where(" name = ?", role.Name).First(existRole).RecordNotFound()
	if !notFound {
		return errors.New("role existed")
	}
	return global.GVA_DB.Create(&role).Error
}

// SetRole 修改角色信息
func (roleAuthService *RoleAuthService) SetRole(role table.Role) (err error) {
	existRole := &table.Role{Name: role.Name}
	if err = global.GVA_DB.First(existRole).Error; err != nil {
		return err
	}
	if err = global.GVA_DB.Save(&role).Error; err != nil {
		return err
	}
	return nil
}

// DeleteRole 删除角色
func (roleAuthService *RoleAuthService) DeleteRole(roleName string) (err error) {
	if roleName == "" {
		return errors.New("role should not be empty string")
	}
	if err = global.GVA_DB.Where(table.UserRoles{Role: roleName}).Delete(&table.UserRoles{}).Error; err != nil {
		return err
	}
	if err := global.GVA_DB.Where(table.RolePermissions{Role: roleName}).Delete(&table.RolePermissions{}).Error; err != nil {
		return err
	}
	if err := global.GVA_DB.Delete(table.Role{Name: roleName}).Error; err != nil {
		return err
	}
	return nil
}

// GetPermissionList 查询所有权限
func (roleAuthService *RoleAuthService) GetPermissionList(info params.GetRoleAuthList) (roles []table.Permission, total int64, err error) {
	db := global.GVA_DB.Model(&table.Permission{})
	if info.Name != "" {
		db = db.Where("name like ?", "%"+info.Name+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(info.Limit).Offset(info.Offset).Find(&roles).Error
	return roles, total, err

}

// CreatePermission 创建权限
func (roleAuthService *RoleAuthService) CreatePermission(permission table.Permission) (err error) {
	existPermission := new(table.Permission)
	notFound := global.GVA_DB.Where(" name = ?", permission.Name).First(existPermission).RecordNotFound()
	if !notFound {
		return errors.New("role existed")
	}
	return global.GVA_DB.Create(&permission).Error
}

// SetPermission 修改权限信息
func (roleAuthService *RoleAuthService) SetPermission(permission table.Permission) (err error) {
	existPermission := &table.Permission{Name: permission.Name}
	if err = global.GVA_DB.First(existPermission).Error; err != nil {
		return err
	}
	if err = global.GVA_DB.Save(&permission).Error; err != nil {
		return err
	}
	return nil
}

// DeletePermission 删除权限
func (roleAuthService *RoleAuthService) DeletePermission(permissionName string) (err error) {
	if permissionName == "" {
		return errors.New("role should not be empty string")
	}
	if err := global.GVA_DB.Where(table.RolePermissions{Role: permissionName}).Delete(&table.RolePermissions{}).Error; err != nil {
		return err
	}
	if err := global.GVA_DB.Delete(table.Permission{Name: permissionName}).Error; err != nil {
		return err
	}
	return nil
}

// AddRolePermission 角色添加权限
func (roleAuthService *RoleAuthService) AddRolePermission(role, permission string) (err error) {

	// 先删除，后增加
	tx := db.Begin(global.GVA_DB)
	defer tx.RollbackIfFailed()

	err = tx.Exec("delete from sys_role_permissions where role = ?", role).Error
	if err != nil {
		return
	}

	list := strings.Split(permission, ",")
	for _, item := range list {

		err = tx.Create(&table.RolePermissions{Role: role, Permission: item}).Error
		if err != nil {
			return
		}
	}

	tx.Commit()
	return nil
}

// DelRolePermission 角色删除权限
func (roleAuthService *RoleAuthService) DelRolePermission(role, permission string) (err error) {
	if role == "" || permission == "" {
		return errors.New("role or permission should not be empty string")
	}
	err = global.GVA_DB.Where(table.RolePermissions{Role: role, Permission: permission}).Delete(&table.RolePermissions{}).Error
	return err
}

// GetRolePermissions 获取全部角色和权限的关联关系
func (roleAuthService *RoleAuthService) GetRolePermissions(role string) (permissions []table.Permission, err error) {
	sql := "JOIN sys_role_permissions ON sys_role_permissions.permission = sys_permission.name and sys_role_permissions.role = ?"
	err = global.GVA_DB.Joins(sql, role).Find(&permissions).Error
	return permissions, err
}

// GetUserRoles 获取指定用户全部角色
func (roleAuthService *RoleAuthService) GetUserRoles(user string) (roles []table.Role, err error) {
	sql := "JOIN sys_user_roles ON sys_user_roles.role = sys_role.name and sys_user_roles.user = ?"
	err = global.GVA_DB.Joins(sql, user).Find(&roles).Error
	return roles, err
}

// AddUserRole 用户添加到组
func (roleAuthService *RoleAuthService) AddUserRole(role, user string) (err error) {
	// 先删除，后增加
	tx := db.Begin(global.GVA_DB)
	defer tx.RollbackIfFailed()

	err = tx.Exec(`delete from sys_user_roles where "user" = ?`, user).Error
	if err != nil {
		return
	}

	list := strings.Split(role, ",")
	for _, item := range list {
		err = tx.Create(&table.UserRoles{Role: item, User: user}).Error
		if err != nil {
			return
		}

	}

	tx.Commit()
	return nil
}

// DelUserRole 用户移除组
func (roleAuthService *RoleAuthService) DelUserRole(role, user string) (err error) {
	if user == "" || role == "" {
		return errors.New("user or role should not be empty string")
	}
	userRoles := new(table.UserRoles)
	userRoles.Role = role
	userRoles.User = user
	return global.GVA_DB.Where(userRoles).Delete(userRoles).Error
}

// GetUserAllPermissions 获取用户的所有权限
func (roleAuthService *RoleAuthService) GetUserAllPermissions(user string) (userPermissions []table.Permission, err error) {

	roles, err := roleAuthService.GetUserRoles(user)
	if err != nil {
		return userPermissions, err
	}
	for _, role := range roles {
		rolePermissions, err := roleAuthService.GetRolePermissions(role.Name)
		if err != nil {
			continue
		}
		for _, permission := range rolePermissions {
			userPermissions = append(userPermissions, permission)
		}
	}
	return userPermissions, err
}

// HasRole 用户是否拥有指定角色
func (roleAuthService *RoleAuthService) HasRole(user, role string) bool {
	userRoles := new(table.UserRoles)
	if global.GVA_DB.Where(table.UserRoles{User: user, Role: role}).First(userRoles).RecordNotFound() {
		return false
	}
	return true
}

func (roleAuthService *RoleAuthService) HasPermission(user, permission string) bool {
	sql := "JOIN sys_user_roles ON sys_user_roles.role = sys_role_permissions.role AND sys_user_roles.user = ? "
	rolePermissions := new(table.RolePermissions)
	query := global.GVA_DB.Joins(sql, user).
		Where(table.RolePermissions{Permission: permission}).Find(rolePermissions)
	err := query.Error
	if err != nil {
		return false
	}

	notFound := query.RecordNotFound()
	if !notFound {
		return true
	}
	return false
}

func (roleAuthService *RoleAuthService) AddPermission(permission *table.Permission) error {
	return global.GVA_DB.FirstOrCreate(permission).Error
}

func (roleAuthService *RoleAuthService) AddRole(role *table.Role) error {
	return global.GVA_DB.FirstOrCreate(role).Error
}
