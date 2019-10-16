package main

/**
 入口程序
 @auth liukelin
*/

import (
	"os"
	// "log"
	"fmt"
	"flag"
	"path/filepath"
	"taskjob/config"
	"taskjob/server"
	"github.com/golang/glog"
)

var ( 
	ConfName      = "config.json"
	TasksTubeName = "queue_job_task"
	CmdDir 		  = "/tmp"
	PoolSize      = 10
	ChPool = make(chan int64, PoolSize) // 最多可同时执行几个任务
)

// 获取项目当前路径
func getCurrentPath(path ...string) string {
	selfDir, err := filepath.Abs(os.Args[0])
	if err != nil {
		panic(err)
	}
	return filepath.Join(filepath.Dir(selfDir), filepath.Join(path...))
}

func main() {
	
	confName := flag.String("config", "", "配置文件路径")
	flag.Parse()

	glog.Infoln("-----------Start TaskSched----------")

	if confName == nil {
		ConfName = getCurrentPath(ConfName) // 默认为跟目录
	}else {
		ConfName = *confName
	}

	err := config.LoadConfig(ConfName)
	if err != nil {
		glog.Fatal(err)
	}

	// 组织缺省参数
	if config.Conf.Queuename == "" {
		config.Conf.Queuename = TasksTubeName
	}
	if config.Conf.QueueControlname == "" {
		config.Conf.QueueControlname = fmt.Sprintf("%s_control", config.Conf.Queuename) 
	}
	if config.Conf.LogsChannel == "" && config.Conf.LogsHttp == ""{
		config.Conf.LogsChannel = fmt.Sprintf("%s_logs", config.Conf.Queuename)
	}
	if config.Conf.ServerPort == "" {
		config.Conf.ServerPort = "9527"
	} 

	glog.Infof("%#v", config.Conf)

	// run
	server.Server(ChPool)
}
