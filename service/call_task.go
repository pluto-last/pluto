package service

import (
	"errors"
	"pluto/global"
	"pluto/middleware/db"
	"pluto/middleware/xgorm"
	"pluto/model/constant"
	"pluto/model/params"
	"pluto/model/table"
	"time"
)

type CallTaskService struct{}

// GetMonthlyBillList 获取用户月账单列表
func (c *CallTaskService) GetMonthlyBillList(info params.GetMonthlyBillList) (list []table.MonthlyBill, total int64, err error) {
	db := global.GVA_DB.Model(&table.MonthlyBill{})
	if info.Date != "" {
		db = db.Where("date = ?", info.Date)
	}
	if info.Mobile != "" {
		db = db.Joins(` join sys_user on sys_user.id = m_recharge.user_id and sys_user.mobile like ?`, "%"+info.Mobile+"%")
	}
	err = db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = db.Limit(info.Limit).Offset(info.Offset).Order("created_at desc").Find(&list).Error
	return list, total, err
}

// CreateCallTask 创建呼叫推广任务
func (c *CallTaskService) CreateCallTask(callTask table.CallTask, taskSips []table.TaskSip) (err error) {
	tx := db.Begin(global.GVA_DB)
	defer tx.RollbackIfFailed()

	err = tx.Create(&callTask).Error
	if err != nil {
		return
	}

	for _, item := range taskSips {
		err = tx.Create(&item).Error
		if err != nil {
			return
		}
	}
	tx.Commit()
	return nil
}

// SetCallTask 编辑呼叫推广任务
func (c *CallTaskService) SetCallTask(callTask table.CallTask, taskSips []table.TaskSip) (err error) {
	if callTask.ID == "" {
		err = errors.New("id is null")
		return
	}

	tx := db.Begin(global.GVA_DB)
	defer tx.RollbackIfFailed()

	// 删除原来绑定的线路
	err = tx.Where("task_id = ?", callTask.ID).
		Delete(&table.TaskSip{}).Error
	if err != nil {
		return
	}
	// 修改
	err = global.GVA_DB.Table("m_call_task").Where("id = ?", callTask.ID).Updates(&callTask).Error

	// 更新綁定關系 简体
	for _, item := range taskSips {
		err = tx.Create(&item).Error
		if err != nil {
			return
		}
	}
	tx.Commit()
	return nil
}

// GetCallTaskList 获取任务列表
func (c *CallTaskService) GetCallTaskList(info params.GetCallTaskList) (list []params.RespCallTaskList, total int64, err error) {
	db := global.GVA_DB.Model(&table.CallTask{})
	if info.Name != "" {
		db = db.Where("name like ?", "%"+info.Name+"%")
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	err = db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = db.Limit(info.Limit).Offset(info.Offset).Order("created_at desc").Find(&list).Error

	// 补充呼叫任务的总数和已完成数
	for i, item := range list {
		global.GVA_DB.Model(&table.CallPlan{}).
			Where("task_id = ? and status = ? ", item.ID, []string{constant.PlanStatusFinished}).Count(list[i].PlansDownCount)

		global.GVA_DB.Model(&table.CallPlan{}).
			Where("task_id = ?", item.ID).Count(list[i].PlansTotalCount)
	}

	return list, total, err
}

// GetCallTaskByID 获取任务详情
func (c *CallTaskService) GetCallTaskByID(id string) (callTask table.CallTask, err error) {
	err = global.GVA_DB.Where("id = ?", id).Preload("Sips").First(&callTask).Error
	return
}

// DeleteCallTask 删除任务
func (c *CallTaskService) DeleteCallTask(id string) (err error) {
	tx := db.Begin(global.GVA_DB)
	defer tx.RollbackIfFailed()
	err = global.GVA_DB.Where("id = ?", id).Delete(&table.CallTask{}).Error
	if err != nil {
		return
	}

	err = global.GVA_DB.Where("task_id = ?", id).Delete(&table.TaskSip{}).Error
	if err != nil {
		return
	}
	tx.Commit()
	return
}

// GetCallPlanList 获取呼叫计划列表
func (c *CallTaskService) GetCallPlanList(info params.GetCallPlanList) (list []params.RespCallTaskList, total int64, err error) {
	db := global.GVA_DB.Model(&table.CallPlan{})
	if info.Mobile != "" {
		db = db.Where("mobile like ?", "%"+info.Mobile+"%")
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}

	if info.CallStatus != "" {
		db = db.Where("call_status = ?", info.CallStatus)
	}

	if info.IntentionTag != "" {
		db = db.Where("intention_tag = ?", info.IntentionTag)
	}

	err = db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = db.Limit(info.Limit).Offset(info.Offset).Order("update_at desc").Find(&list).Error
	return list, total, err
}

// ImportCallPlan 导入呼叫任务
func (c *CallTaskService) ImportCallPlan(plans []table.CallPlan, isRemoveDup bool) (err error) {

	var realInsertPlans []table.CallPlan
	// 是否去重
	if isRemoveDup {
		var tempPlans []table.CallPlan
		// 本次导入的号码集合
		importPlansMp := make(map[string]struct{}, len(plans))
		for _, item := range plans {
			importPlansMp[item.Mobile] = struct{}{}
		}

		// 本次导入的号码去重
		// 访问 map 中不存在的 key 时，Go 则会返回元素对应数据类型的零值，比如 nil、’’ 、false 和 0，取值操作总有值返回
		for i, item := range plans {
			if _, ok := importPlansMp[item.Mobile]; !ok {
				tempPlans = append(tempPlans, plans[i])
			}
		}

		// 数据库中已存在的号码集合
		dbPlansMp := make(map[string]struct{}, len(plans))
		var dbPlans []table.CallPlan
		global.GVA_DB.Where("task_id = ?", plans[0].TaskID).Find(&dbPlans)
		for _, item := range plans {
			dbPlansMp[item.Mobile] = struct{}{}
		}

		// 与数据库中已存在的号码去重
		for i, item := range tempPlans {
			if _, ok := dbPlansMp[item.Mobile]; !ok {
				realInsertPlans = append(realInsertPlans, tempPlans[i])
			}
		}
	} else {
		realInsertPlans = plans
	}

	// 使用封装好的方法，传入接口的实现
	batchArgs := make([]xgorm.SQLFieldsGetter, len(realInsertPlans))
	for i := range realInsertPlans {
		batchArgs[i] = realInsertPlans[i]
	}

	tx := global.GVA_DB.Begin()

	// 这里可控制每次插入一百条数据
	err = xgorm.BatchInsert(tx, 100, batchArgs)
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	return
}

// ChangeTaskStatus 开始/暂停任务
func (c *CallTaskService) ChangeTaskStatus(id, status string) (err error) {

	// 查询线路信息
	var callTask table.CallTask
	if global.GVA_DB.Where("id = ?").First(&callTask).RecordNotFound() {
		return errors.New("查询线路信息失败")
	}

	if status == constant.CallTaskStatusCalling {
		if callTask.Status == constant.CallTaskStatusCalling || callTask.Status == constant.CallTaskStatusFinished {
			return errors.New("任务无法开始")
		}

		var taskSip []table.TaskSip
		err := global.GVA_DB.Where("task_id = ?", id).Preload("UserSipInfo").Find(&taskSip).Error
		if err != nil || len(taskSip) <= 0 {
			return errors.New("请先给任务勾选线路")
		}

		now := time.Now()
		// 校验是否有过期的线路
		for _, item := range taskSip {
			if item.UserSipInfo.ExpireAt.Before(now) {
				return errors.New("任务中存在已过期的线路，无法开始")
			}
		}

		// 校验用户钱包
		var wallet table.Wallet
		if global.GVA_DB.Where("user_id = ?", callTask.UserID).First(&wallet).RecordNotFound() {
			return errors.New("查询用户钱包失败")
		}

		if wallet.Balance <= 0 {
			callTask.Status = constant.CallTaskStatusInsufficientBalance
			global.GVA_DB.Save(&callTask)
			return errors.New("用户钱包余额不足")
		}

		// 任务状态更改为进行中
		callTask.Status = constant.CallTaskStatusCalling
		global.GVA_DB.Save(&callTask)

		// 开始任务
		go DispatchTask(callTask.ID)

	} else if status == constant.CallTaskStatusManualPause {

		if callTask.Status != constant.CallTaskStatusCalling {
			return errors.New("任务不需要暂停")
		}
		// 任务状态更改为暂停
		callTask.Status = constant.CallTaskStatusManualPause
		global.GVA_DB.Save(&callTask)
	} else {
		return errors.New("任务状态错误")
	}

	return nil
}
