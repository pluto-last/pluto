package service

import (
	"context"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"math"
	"pluto/global"
	"pluto/middleware/db"
	"pluto/model/constant"
	"pluto/model/table"
	"sync"
	"time"
)

var clientReady bool                      // websocket连接是否就绪
var callTaskClinet ClientPool             // 跟机器人bot的websocket连接
var taskCallingMap map[string]*TaskStatus // 正在拨打中的任务map
var taskCallingMapMutex sync.Mutex        // map锁， 防止map的并发问题

func init() {
	taskCallingMap = make(map[string]*TaskStatus)
}

type TaskStatus struct {
	TaskID     string
	Ctx        context.Context
	CancelFunc context.CancelFunc
	PlanChan   chan *table.CallPlan
	MapMutex   sync.Mutex
}

// DispatchTask 任务的分发
func DispatchTask(taskID string) {

	var task table.CallTask
	if global.GVA_DB.Where("id = ?", taskID).First(&task).RecordNotFound() {
		global.GVA_LOG.Info("任务信息没找到")
		return
	}

	// 1. 校验任务是否在允许拨打的时间段内，校验任务状态，校验用户钱包余额
	if !taskCanStart(task) {
		global.GVA_LOG.Info("任务无法正常开启")
		return
	}

	// 2. 避免任务重复执行，已经running中的task放在map中，同时通过锁来避免并发资源竞争
	taskCallingMapMutex.Lock()
	defer taskCallingMapMutex.Unlock() // 一定要记得释放锁
	_, ok := taskCallingMap[task.ID]
	if ok {
		global.GVA_LOG.Error("task has started:", zap.String("taskName", task.Name), zap.String("taskID", task.ID))
		return
	}

	// 3. 查询出该任务使用了哪些线路去进行拨打
	taskSips := []*table.TaskSip{}
	err := global.GVA_DB.Where("task_id = ? and concurrent > ?", task.ID, task.UserID, 0).
		Preload("UserSipInfo").
		Find(&taskSips).Error
	if err != nil {
		return
	}

	// 4.初始化context 以及plan通道，plan通道用来在电话号码的生产消费者模型里面共享数据
	ctx, CancelFunc := context.WithCancel(context.Background())
	planCh := make(chan *table.CallPlan)

	// 5. 启动消费者
	for _, taskSip := range taskSips {

		// 查询出线路相关的信息，ip,port，以及使用的并发数等等信息
		sip := new(table.Sip)
		global.GVA_DB.Where("id = ?", taskSip.UserSipInfo.SipID).First(sip)

		// 消费者线程启动， 用来处理生产者传过来的号码
		// 每个线路都启动一个消费者
		go callTaskClinet.DoPlan(ctx, planCh, taskSip.Concurrent, &task, sip)
	}

	taskStatus := &TaskStatus{
		TaskID:     task.ID,
		MapMutex:   sync.Mutex{},
		PlanChan:   planCh,
		Ctx:        ctx,
		CancelFunc: CancelFunc,
	}

	// 6. 启动生产者， 生产者只有一个，消费者有多个
	go DispatchPlan(taskStatus)

	// 7. 标记当前task为running状态
	taskCallingMap[task.ID] = taskStatus
}

// taskCanStart 是否能开始任务
func taskCanStart(task table.CallTask) bool {

	if task.Status != constant.CallTaskStatusCalling {
		return false
	}

	// 判断时间段
	now := time.Now().Local()
	todayMinute := now.Hour()*60 + now.Minute()

	if todayMinute < task.ExecuteStartTime || todayMinute > task.ExecuteEndTime {
		return false
	}

	// 校验余额
	var wallet table.Wallet
	if global.GVA_DB.Where("user_id = ?", task.UserID).First(&wallet).RecordNotFound() {
		return false
	}

	if wallet.Balance <= 0 {
		return false
	}

	return true
}

// DispatchPlan 待呼叫号码的分发 （生产者）
func DispatchPlan(status *TaskStatus) {

	// 1. 方法结束时，代表task结束或者暂停，在defer方法中删除存放在map中的标记
	defer func() {
		taskCallingMapMutex.Lock()
		defer taskCallingMapMutex.Unlock()
		status.CancelFunc()
		global.GVA_LOG.Info("delete taskstatus ", zap.String("taskID", status.TaskID))
		delete(taskCallingMap, status.TaskID)
	}()

	global.GVA_LOG.Info("DispatchPlan taskid", zap.String("taskID", status.TaskID))

	// 2. 生产者开始生产
	for {
		// 查询当前task的状态
		var task table.CallTask
		if global.GVA_DB.Where("id = ?", status.TaskID).First(&task).RecordNotFound() {
			global.GVA_LOG.Error("呼叫任务未找到")
			return
		}
		if task.Status != constant.CallTaskStatusCalling {
			global.GVA_LOG.Error("呼叫任务已经暂停")
			return
		}

		// 查询当前task总并发数
		var concurrent int
		taskSips := []table.TaskSip{}
		if err := global.GVA_DB.Where("task_id = ?", task.ID).Find(&taskSips).Error; err != nil || len(taskSips) <= 0 {
			global.GVA_LOG.Error("task没找到对应的线路")
			return
		}
		for _, v := range taskSips {
			concurrent += v.Concurrent
		}

		// 判断项目是否完成,如果任务已经执行完毕，更改任务状态
		var planCount int64
		global.GVA_DB.Model(&table.CallPlan{}).Where(
			"task_id = ? and status in (?)",
			status.TaskID,
			[]string{constant.PlanStatusReady, constant.PlanStatusCalling}).
			Count(&planCount)

		// 修改任务为已经完成状态
		if planCount == 0 {
			global.GVA_LOG.Info("task中没有未开始,呼叫中的任务了,任务已经执行完毕")
			task.Status = constant.CallTaskStatusFinished
			global.GVA_DB.Save(task)
			return
		}

		// 根据并发数去获取电话号码
		plans := []*table.CallPlan{}
		global.GVA_DB.Where("task_id = ? and status in (?))",
			status.TaskID,
			[]string{constant.PlanStatusReady}).
			Limit(concurrent).
			Find(&plans)

		// 开始向消费者传递电话号码
		for _, plan := range plans {

			select {
			case <-status.Ctx.Done(): // ctx.Done 通知线程结束，直接return
				global.GVA_LOG.Info("task status context Done")
				return
			case status.PlanChan <- plan:
				plan.Status = constant.PlanStatusCalling // 修改此条号码的状态为呼叫中
				global.GVA_DB.Save(plan)
			case <-time.After(time.Second * 10):
				global.GVA_LOG.Info("生产者阻塞，等待10秒") // 没有缓冲区的chan是阻塞的，当所有消费着都是忙碌状态的时候，生产者先等待
			}
		}
	}

}

// PlanCallBack 呼叫成功后的回调
func PlanCallBack(info *CallBackResult, plan *table.CallPlan) error {

	tx := db.Begin(global.GVA_DB)
	defer tx.RollbackIfFailed()

	var userSip table.UserSip
	tx.Where("user_id = ? and sip_id = ? ", plan.UserID, plan.SipID).First(&userSip)

	// 1. 更新呼叫记录
	plan.DurationSec = info.DurationSec
	plan.IntentionTag = info.IntentionTag
	plan.Rounds = info.Rounds
	plan.Status = info.Status
	plan.Price = userSip.Price
	minute := float64(plan.DurationSec) / 60
	plan.Minute = int(math.Ceil(minute))
	plan.Cost = plan.Price * plan.Minute

	tx.Save(plan)

	// 2. 钱包扣费
	var wallet table.Wallet
	tx.Where("user_id = ?", plan.UserID).First(&wallet)
	wallet.Balance -= plan.Cost

	tx.Save(wallet)

	// 修改任务为钱包额度不足状态
	if wallet.Balance <= 0 {
		tx.Model(&table.CallTask{}).Where("id = ?", plan.TaskID).
			UpdateColumns(map[string]interface{}{
				"status":     constant.CallTaskStatusInsufficientBalance,
				"updated_at": time.Now(),
			})
	}

	// 3. 更新用户月账单
	ti := time.Now()
	date := time.Date(ti.Year(), ti.Month(), 1, 0, 0, 0, 0, ti.Location())
	bill := new(table.MonthlyBill)
	notFound := tx.Where("user_id = ? and date = ?", plan.UserID, date).First(bill).RecordNotFound()

	if notFound {
		bill.DurationMin = plan.Minute
		bill.Cost = plan.Cost
		bill.Count = 1
		bill.UserID = plan.UserID
		bill.Date = date
		tx.Create(bill)
	} else {
		err := tx.Model(bill).UpdateColumns(map[string]interface{}{
			"duration_min": gorm.Expr("duration_min + ?", plan.Minute),
			"cost":         gorm.Expr("cost + ?", plan.Cost),
			"count":        gorm.Expr("count + ?", 1),
			"updated_at":   time.Now(),
		}).Error

		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}
