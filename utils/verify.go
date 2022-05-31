package utils

var (
	RegisterVerify         = Rules{"Mobile": {NotEmpty()}, "Password": {NotEmpty()}, "Sms": {NotEmpty()}}
	LoginVerify            = Rules{"Mobile": {NotEmpty()}, "Password": {NotEmpty()}}
	ChangePasswordVerify   = Rules{"Mobile": {NotEmpty()}, "Password": {NotEmpty()}, "NewPassword": {NotEmpty()}}
	PageInfoVerify         = Rules{"Limit": {NotEmpty()}, "Offset": {NotEmpty()}}
	IdVerify               = Rules{"ID": {NotEmpty()}}
	RoleVerify             = Rules{"Name": {NotEmpty()}, "Describe": {NotEmpty()}}
	RolePermissionVerify   = Rules{"Role": {NotEmpty()}, "Permission": {NotEmpty()}}
	UserRoleVerify         = Rules{"Role": {NotEmpty()}, "UserID": {NotEmpty()}}
	SMSVerify              = Rules{"Mobile": {NotEmpty()}, "Captcha": {NotEmpty()}, "CaptchaID": {NotEmpty()}}
	CreateSipVerify        = Rules{"Mobile": {NotEmpty()}, "Name": {NotEmpty()}, "IntervalSec": {NotEmpty()}, "SipIP": {NotEmpty()}, "SipPort": {NotEmpty()}}
	GetUserSipList         = Rules{"UserID": {NotEmpty()}, "Limit": {NotEmpty()}, "Offset": {NotEmpty()}}
	UserAddSipVerify       = Rules{"UserID": {NotEmpty()}, "SipID": {NotEmpty()}, "ExpireAt": {NotEmpty()}, "Concurrent": {NotEmpty()}, "Price": {NotEmpty()}}
	Recharge               = Rules{"UserID": {NotEmpty()}, "Type": {NotEmpty()}, "Amount": {NotEmpty()}}
	ChangeTaskStatusVerify = Rules{"ID": {NotEmpty()}, "Status": {NotEmpty()}}
)
