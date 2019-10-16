# 任务执行器 task-job
@auth liukelin

说明
```
 这是一个简单的抢夺任务的任务执行器。 任务是由执行器自己从queue获取并执行的。
后续在前面加上一个任务分配器，任务分配到执行器各自的队列，就可以实现简单的任务调度了。 
```

背景
```
 初始原因是容器内需要跑一些宿主机上的任务，所以在宿主机需要有这样一个接收任务并执行的工具。这里就封装成一个 基于queue的异步任务调度工具。下发和状态日志上报走消息队列。
```

目标
```
 该系统承接的是一个分布式的任务调度系统。支持 任务的下发、执行、任务状态、执行日志、结果反馈，后续 可以通过http接口管理。
```

依赖组件 
```
redis 或 kafka
```

配置文件 config.json
```
{
	"server_port":"",           // 冗余， http端口
	"redis_addr":"127.0.0.1:6379", // redis配置路径，单点或者sentinel
	"queue_name":"",            // 消息通道队列名, 缺省值 queue_job_task, 
    	"queue_control_name":"",    // 任务控制队列名  缺省值 {queue_name}_control 
	"logs":false,               // 开启执行日志上报
	"logs_channel":"queue"      // 日志上报 队列通道   缺省值 {queue_name}_logs
    	"logs_http":""              // 日志上报 web接口
    	"cmd_dir:""                 // 脚本执行目录，用于限制越权  缺省值 "/tmp"
}
```

任务格式 Task
```
{
    "id":"101",                // 消息ID全局唯一
    "type":"start",            // 任务类型  冗余
    "tube":"task_list",
    "task":{                    
        "cmdline":"bash xxx",  // 执行命令行
        "cmd_args":["a","b"],  // 执行参数 最终是拼接的方式
        "timeout":1            // 任务执行超时时间, 小于0 不超时，默认不超时
    }
}
```

任务控制格式 TaskControl
```
{
    "id":"101",               // 消息ID全局唯一
    "type":"kill",            // kill 终止任务
}
```

任务执行输出日志返回格式 TaskLog
```
{
    "id":"101",             // 消息ID
    "type":"start",         // 任务类型
    "msg_type":"stdout",    // 消息类型  stdout 正常执行输出, stderr 执行错误， done 执行完成
    "data":"xxx"            // 日志内容
}
```

运行
```
go run main.go -config=/tmp/config.json -alsologtostderr

    -log_dir=log -alsologtostderr
    以上打印日志将会同时打印在 log/ 目录和标准错误输出中（-alsologtostderr）
```


下发一个任务
```
rpush queue_job_task '{"id":"101","type":"test","task":{"cmdline":"install.sh","cmd_args":["feature/homestead"],"timeout":3600}}'
```


终止一个任务
```
rpush queue_job_task_control '{"id":"101", "type":"kill"}'
```

二阶段
```
 利用etcd做注册中心，服务发现与治理.
```


编译
```
    windows
    ./build.sh windows

    linux
    ./build.sh linux

    mac
    ./build.sh darwin

    freebsd
    ./build.sh freebsd

```
