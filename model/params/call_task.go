package params

import (
	"pluto/model/table"
)

type GetMonthlyBillList struct {
	PageInfo
	Date   string `json:"date" form:"date"`     // 时间
	Mobile string `json:"mobile" form:"mobile"` // 用户手机号码
}

type SetCallTask struct {
	ID               string          `json:"id"  form:"id"`
	UserID           string          `json:"userID"  form:"date"`
	Name             string          `json:"name"  form:"name"`                         // 任务名称
	SceneID          string          `json:"sceneID"  form:"sceneID"  `                 // 话术ID
	Status           string          `json:"status"  form:"status"`                     // 状态
	ExecuteStartTime int             `json:"executeStartTime"  form:"executeStartTime"` // 任务允许拨打的开始时间（存分钟数）
	ExecuteEndTime   int             `json:"executeEndTime"  form:"executeEndTime"`     // 任务允许拨打的结束时间（存分钟数）
	Describe         string          `json:"describe"  form:"describe"`                 // 任务描述
	TaskSips         []table.TaskSip `json:"taskSips" form:"taskSips"`                  // 任务关联的线路
}

type GetCallTaskList struct {
	PageInfo
	Name   string `json:"name" form:"name"`     // 任务名称
	Status string `json:"status" form:"status"` // 任务状态
}

type RespCallTaskList struct {
	table.CallTask
	PlansTotalCount int `json:"plansTotalCount"` // 呼叫计划总数
	PlansDownCount  int `json:"plansDownCount"`  // 呼叫计划已完成数
}

type GetCallPlanList struct {
	PageInfo
	Status       string `json:"status" form:"status"`             // 任务状态
	CallStatus   string `json:"callStatus" form:"callStatus"`     // 呼叫状态
	IntentionTag string `json:"intentionTag" form:"intentionTag"` // 最终意向
	Mobile       string `json:"mobile" form:"mobile"`             // 手机号码
}

type ImportCallPlan struct {
	UserID string `json:"userID"  form:"userID"`
	TaskID string `json:"taskID" form:"userID"`
	Name   string `json:"name" form:"name"`
	Mobile string `json:"mobile" form:"mobile"`
	Status string `json:"status" form:"status"`
}

type ImportCallPlanData struct {
	Data        []table.CallPlan `json:"data"  form:"data"`
	IsRemoveDup bool             `json:"isRemoveDup"  form:"isRemoveDup"`
}

type CreateCallTask struct {
	CallTask table.CallTask  `json:"callTask" form:"callTask"`
	TaskSip  []table.TaskSip `json:"taskSip" form:"taskSip"`
}

type ChangeTaskStatus struct {
	ID     string `json:"id"  form:"id"`
	Status string `json:"status"  form:"status"`
}
