package table

type Role struct {
	Name         string `gorm:"primary_key" json:"name" form:"name"`
	Describe     string `json:"describe" form:"describe"`
	FullDescribe string `json:"fullDescribe" form:"fullDescribe"`
}

func (Role) TableName() string {
	return "sys_role"
}

type Permission struct {
	Name         string `gorm:"primary_key" json:"name" form:"name"`
	Describe     string `json:"describe" form:"describe"`
	FullDescribe string `json:"fullDescribe" form:"fullDescribe"`
}

func (Permission) TableName() string {
	return "sys_permission"
}

type UserRoles struct {
	ID   uint   `gorm:"primary_key"  json:"id" form:"id"`
	User string `json:"user" form:"user"`
	Role string `json:"role" form:"role"`
}

func (UserRoles) TableName() string {
	return "sys_user_roles"
}

type RolePermissions struct {
	ID         uint   `gorm:"primary_key"  json:"id" form:"id"`
	Role       string `json:"role" form:"role"`
	Permission string `json:"permission" form:"permission"`
}

func (RolePermissions) TableName() string {
	return "sys_role_permissions"
}
