package config

const (
	// 日志返回类型
	LOG_TYPE_VALIDATE = "validate" // 任务内容不合法
	LOG_TYPE_STDOUT = "stdout"	// 日志类型  正常输出
	LOG_TYPE_STDERR = "stderr"	// 日志类型  执行错误
	LOG_TYPE_DONE 	= "done"	// 日志类型  执行完成
	// 任务控制类型
	CONTROL_TYPE_KILL= "kill"	// 任务控制，终止任务
)

// 任务格式
type TaskInfo struct {
	Id       	string		`json:"id"`		// 任务id 标示同一个任务唯一
	Type		string		`json:"type"`	// 任务类型
	Tube	 	string		`json:"tube"`	// 管道名称，可以为空, 比如优先任务，指定目标任务
	TaskData    TaskData	`json:"task"`	// 任务内容，详细内容。需要执行的命令
}

// 任务执行内容
type TaskData struct {
	Cmdline	string		`json:"cmdline"`  // 运行命令行
	CmdArgs	[]string	`json:"cmd_args"` // 命令行参数，可为空
	Timeout	int			`json:"timeout"`  // 超时时间, 这个利用 context 来控制
}

// 任务控制内容
type TaskControl struct {
	Id		string		`json:"id"`		// 任务id 标示同一个任务唯一
	Type	string		`json:"type"`	// 类型 kill 终止任务
}

// 任务执行日志
type TaskLog struct {
	Id			string	`json:"id"`			// 任务id 标示同一个任务唯一
	Type		string	`json:"type"`		// 任务类型
	ChildNum 	int		`json:"child_num"` 	// 子任务编号
	MsgType 	string	`json:"msg_type"`	// 消息类型  stdout 正常执行输出, stderr 执行错误， done 执行完成
	Data 		string	`json:"data"`		// 日志内容
}