package table

import (
	"pluto/global"
	"time"
)

type Task struct {
	global.UUID
	Name          string
	Desc          string
	Detail        string
	TaskType      string
	LimitTimes    int
	PrizePoint    int
	Status        string
	MaxExcuteTime int64 // 任务最大运行时间（ms）
	ExcuteStart   int
	ExcuteEnd     int
}

func (Task) TableName() string {
	return "t_task"
}

type TaskRecord struct {
	global.UUID
	TaskID        string
	UserID        string
	RobotID       string
	Status        string    // 运行中，已完成，超时，异常中止
	MaxExcuteTime time.Time // 超时时间
	StartTime     time.Time
	EndTime       time.Time
	Remark        string
	PrizePoint    int
}

func (TaskRecord) TableName() string {
	return "t_task_record"
}

type TaskResult struct {
	global.UUID
	TaskID       string
	TaskRecordID string
	Data         string
}

func (TaskResult) TableName() string {
	return "t_task_result"
}
