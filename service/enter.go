package service

type ServiceGroup struct {
	OperationRecordService
	UserService
	RoleAuthService
	SipService
	RechargeService
	SceneService
	CallTaskService
	WSClientService
}

var ServiceGroupAPP = new(ServiceGroup)
