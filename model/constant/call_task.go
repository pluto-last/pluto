package constant

const (
	CallTaskStatusCreated             = "created" // 创建完
	CallTaskStatusCalling             = "calling"
	CallTaskStatusEnd                 = "end"         // 任务终止
	CallTaskStatusManualPause         = "manualPause" // 手工暂停
	CallTaskStatusFinished            = "finished"
	CallTaskStatusInsufficientBalance = "insufficientBalance" // 余额不足
)

const (
	PlanStatusReady    = "ready"
	PlanStatusCalling  = "calling"
	PlanStatusFinished = "finished"
)

const (
	RecordStatusCallerFailed = "callerFailed" // 线路不可用   线路原因导致呼叫失败
	RecordStatusFailed       = "failed"       // 失败        系统原因导致失败
	RecordStatusRingTimeout  = "ringTimeout"  // 呼叫超时
	RecordStatusRefused      = "refused"      // 拒接
	RecordStatusOK           = "ok"           // 正常
	RecordStatusBusy         = "busy"         // 忙
	RecordStatusUnavailable  = "unavailable"  // 号码不可用   被叫原因导致呼叫失败。除忙，拒接，呼叫超时外的其他呼叫失败问题。
	RecordStatusInBlacklist  = "inBlacklist"  // 黑名单
	RecordStatusRobotSpeaker = "robotSpeaker" // 程序自动应答
	RecordStatusCallLoss     = "callLoss"     // 呼损
)
