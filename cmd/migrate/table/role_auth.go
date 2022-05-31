package table

type Role struct {
	Name         string `gorm:"primary_key"`
	Describe     string
	FullDescribe string
}

func (Role) TableName() string {
	return "sys_role"
}

type Permission struct {
	Name         string `gorm:"primary_key"`
	Describe     string
	FullDescribe string
}

func (Permission) TableName() string {
	return "sys_permission"
}

type UserRoles struct {
	ID   uint `gorm:"primary_key"`
	User string
	Role string
}

func (UserRoles) TableName() string {
	return "sys_user_roles"
}

type RolePermissions struct {
	ID         uint `gorm:"primary_key"`
	Role       string
	Permission string
}

func (RolePermissions) TableName() string {
	return "sys_role_permissions"
}
