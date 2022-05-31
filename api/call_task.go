package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"pluto/global"
	"pluto/model/constant"
	"pluto/model/params"
	"pluto/model/reply"
	"pluto/utils"
)

type CallTaskCtl struct{}

// @Tags Bill
// @Summary 获取用户月账单列表
// @Produce  application/json
// @Param body body params.GetMonthlyBillList true  "获取用户月账单列表"
// @Success 200
// @Router /bill/getMonthlyBillList [get]
func (call *CallTaskCtl) GetMonthlyBillList(c *gin.Context) {
	var pageInfo params.GetMonthlyBillList
	_ = c.Bind(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := CallTaskService.GetMonthlyBillList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取用户月账单列表失败!", zap.Any("err", err))
		reply.FailWithMessage("获取用户月账单列表失败", c)
	} else {
		reply.OkWithDetailed(params.PageResult{
			List:   list,
			Total:  total,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
		}, "获取成功", c)
	}
}

// @Tags Call
// @Summary 创建呼叫推广任务
// @Produce  application/json
// @Param body body params.CreateCallTask true  "创建呼叫推广任务"
// @Success 200
// @Router /call/createCallTask [post]
func (call *CallTaskCtl) CreateCallTask(c *gin.Context) {
	var r params.CreateCallTask
	_ = c.BindJSON(&r)

	r.CallTask.Status = constant.CallTaskStatusCreated
	err := CallTaskService.CreateCallTask(r.CallTask, r.TaskSip)
	if err != nil {
		global.GVA_LOG.Error("创建呼叫推广任务失败!", zap.Any("err", err))
		reply.FailWithMessage("创建呼叫推广任务失败", c)
		return
	}
	reply.Ok(c)
}

// @Tags Call
// @Summary 编辑呼叫推广任务
// @Produce  application/json
// @Param body body params.CreateCallTask true  "编辑呼叫推广任务"
// @Success 200
// @Router /call/setCallTask [put]
func (call *CallTaskCtl) SetCallTask(c *gin.Context) {
	var r params.CreateCallTask
	_ = c.BindJSON(&r)

	err := CallTaskService.SetCallTask(r.CallTask, r.TaskSip)
	if err != nil {
		global.GVA_LOG.Error("编辑呼叫推广任务失败!", zap.Any("err", err))
		reply.FailWithMessage("编辑呼叫推广任务失败", c)
		return
	}
	reply.Ok(c)
}

// @Tags Call
// @Summary 获取呼叫任务列表
// @Produce  application/json
// @Param body body params.GetCallTaskList true  "获取呼叫任务列表"
// @Success 200
// @Router /call/getCallTaskList [get]
func (call *CallTaskCtl) GetCallTaskList(c *gin.Context) {
	var pageInfo params.GetCallTaskList
	_ = c.Bind(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := CallTaskService.GetCallTaskList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取呼叫任务列表失败!", zap.Any("err", err))
		reply.FailWithMessage("获取呼叫任务列表失败", c)
	} else {
		reply.OkWithDetailed(params.PageResult{
			List:   list,
			Total:  total,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
		}, "获取成功", c)
	}
}

// @Tags Call
// @Summary 获取任务详情
// @Produce  application/json
// @Param id query string true "id"
// @Success 200
// @Router /call/GetCallTaskByID [get]
func (call *CallTaskCtl) GetCallTaskByID(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		reply.FailWithMessage("ID不能为空", c)
		return
	}
	if callTask, err := CallTaskService.GetCallTaskByID(id); err != nil {
		global.GVA_LOG.Error("获取任务详情失败!", zap.Any("err", err))
		reply.FailWithMessage("获取任务详情失败", c)
	} else {
		reply.OkWithData(callTask, c)
	}
}

// @Tags Call
// @Summary 删除任务
// @Produce  application/json
// @Param id query string true "id"
// @Success 200
// @Router /call/deleteCallTask [delete]
func (call *CallTaskCtl) DeleteCallTask(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		reply.FailWithMessage("ID不能为空", c)
		return
	}
	if err := CallTaskService.DeleteCallTask(id); err != nil {
		global.GVA_LOG.Error("删除任务失败!", zap.Any("err", err))
		reply.FailWithMessage("删除任务失败", c)
	} else {
		reply.Ok(c)
	}
}

// @Tags Call
// @Summary 导入呼叫任务
// @Produce  application/json
// @Param body body params.ImportCallPlanData true  "导入呼叫任务"
// @Success 200
// @Router /call/importCallPlan [POST]
func (call *CallTaskCtl) ImportCallPlan(c *gin.Context) {
	var r params.ImportCallPlanData
	_ = c.BindJSON(&r)

	err := CallTaskService.ImportCallPlan(r.Data, r.IsRemoveDup)
	if err != nil {
		global.GVA_LOG.Error("导入呼叫任务失败!", zap.Any("err", err))
		reply.FailWithMessage("导入呼叫任务失败", c)
		return
	}
	reply.Ok(c)
}

// @Tags Call
// @Summary 开始/暂停任务
// @Produce  application/json
// @Param body body params.ChangeTaskStatus true  "导入呼叫任务"
// @Success 200
// @Router /call/changeTaskStatus [POST]
func (call *CallTaskCtl) ChangeTaskStatus(c *gin.Context) {
	var r params.ChangeTaskStatus
	_ = c.Bind(&r)
	if err := utils.Verify(r, utils.ChangeTaskStatusVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	err := CallTaskService.ChangeTaskStatus(r.ID, r.Status)
	if err != nil {
		global.GVA_LOG.Error("导入呼叫任务失败!", zap.Any("err", err))
		reply.FailWithMessage(err.Error(), c)
		return
	}
	reply.Ok(c)
}
