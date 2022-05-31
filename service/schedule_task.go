package service

import (
	"go.uber.org/zap"
	"pluto/global"

	"github.com/robfig/cron/v3"
)

func init() {
	ScheduleTask()
}

// ScheduleTask 定时任务调度中心
func ScheduleTask() {

	// 初始化任务调度中心
	c := cron.New()

	// 呼叫超过10分钟没反应的号码设置为超时，一分钟调用一次
	c.AddFunc("@every 1m", func() {
		err := taskTimeOut()
		if err != nil {
			global.GVA_LOG.Error("taskTimeOut run err :", zap.Any("error", err))
		}
	})

	// 调用中心开始运行
	c.Start()
}

func taskTimeOut() error {

	global.GVA_LOG.Error("定时任务执行一次了")

	return nil
}
