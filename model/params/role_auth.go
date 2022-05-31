package params

type GetRoleAuthList struct {
	PageInfo
	Name string `json:"name" form:"name"`
}

type UserRolePermission struct {
	Role       string `json:"role" form:"role"`
	Permission string `json:"permission" form:"permission"`
	UserID     string `json:"userID" form:"userID"`
}
