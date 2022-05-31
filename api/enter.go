package api

import "pluto/service"

type ApisGroup struct {
	UserCtl
	CaptchaCtl
	RoleAuthCtl
	SipCtl
	RechargeCtl
	SceneCtl
	CallTaskCtl
	WSClientCtl
}

var ApisGroupsAPP = new(ApisGroup)

// service 实例
var operationRecordService = service.ServiceGroupAPP.OperationRecordService
var userService = service.ServiceGroupAPP.UserService
var roleAuthService = service.ServiceGroupAPP.RoleAuthService
var SipService = service.ServiceGroupAPP.SipService
var RechargeService = service.ServiceGroupAPP.RechargeService
var SceneService = service.ServiceGroupAPP.SceneService
var CallTaskService = service.ServiceGroupAPP.CallTaskService
var WSClientService = service.ServiceGroupAPP.WSClientService
