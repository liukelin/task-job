package server

/**
 主要服务启动
 @auth liukelin
*/

import (
	"time"
	"sync"
	"github.com/golang/glog"
	"taskjob/config"
	// "taskjob/httpserver"
	mq "taskjob/queue/redis" // 引入redis mq
	// mq "taskjob/queue/kafka" // 引入kafka mq
)

// MQ
type Queue mq.QueueInterface

// 用于存储当前正在运行的任务
var TaskPool sync.Map

// 消费服务 常驻
// 判断当前正在处理的任务，是否需要继续接收处理任务
func Server(ChPool chan int64){

	// 连接MQ
	q := &mq.Queue{
		Conf: &mq.Conf{
			Host:config.Conf.RedisAddr,
			Qname:config.Conf.Queuename,
		},
	}
	q.Connect()

	// 从公共任务队列获取任务并消费
	go func(){
		for {
			ret, err := q.Get("")
			if err == nil {
				if ret != "" {	
					task, _err := ParsTaskJson(ret)
					if _err == nil {
						ChPool<-1  		// 如果pool满了，此处将阻塞
						task.Qconn = q
						go TaskRun(&task, ChPool)
					}else{
						glog.Error("parsJobJson error:%s, %v", ret, _err)
					}
				}
			}else{
				glog.Error("queue get error:%s, err: %v", ret, err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// 从专属任务队列获取任务并消费，用于插队优先任务


	// 从控制器队列获取控制任务 并消费
	go func(){
		for {
			ret, err := q.Get(config.Conf.QueueControlname)
			if err == nil {
				if ret != "" {	
					taskControl, _err := ParsTaskControlJson(ret)
					if _err == nil {
						go TaskController(&taskControl)
					}else{
						glog.Error("parsTaskControllerJson error:%s, %v", ret, _err)
					}
				}
			}else{
				glog.Error("queue get error:%s, %s, %v", config.Conf.QueueControlname, ret, err)
			}
			time.Sleep(1 * time.Second)
		}
	}()


	// 启动API服务
	// go func(){
	// 	httpserver.HttpServer()
	// }()

	// 阻塞主进程，用于控制整个程序暂停
	IsStop := make(chan bool)
	for {
		select {
		case stop := <-IsStop: // 收到停止程序命令
			if stop {
				glog.Warning("Stop...\r\n")
				break
			}
		}
		time.Sleep(1 * time.Second)
	}
}

