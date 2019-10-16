package config

import (
	"encoding/json"
	"io/ioutil"
)

/**
 加载配置文件
 @auth liukelin
*/
var Conf *Config

type Config struct {
	ServerPort     		string `json:"server_port"`
	RedisAddr      		string `json:"redis_addr"` // redis配置
	Queuename 	   		string `json:"queue_name"` // 消息通道队列名, 缺省值 queue_job_task, 
	QueueControlname	string `json:"queue_control_name"` // 任务控制队列名  缺省值 {queue_name}_control 
	Logs		   		bool	`json:"logs"`			// 开启执行日志
	LogsChannel			string	`json:"logs_channel"`	// 日志上报 队列通道   缺省值 {queue_name}_logs
	LogsHttp			string	`json:"logs_http"`		// 日志上报 web接口
	CmdDir				string	`json:"cmd_dir"`		// 脚本执行目录，用于限制越权  缺省值 /tmp
}

func init() {
	Conf = &Config{Logs: true}
}

func LoadConfig(confPath string) error {
	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &Conf)
	return err
}
