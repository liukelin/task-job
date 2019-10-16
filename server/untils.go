package server

import (
	"fmt"
	// "github.com/golang/glog"
	"taskjob/config"
)


// 校验任务参数格式
func ValidateTask(t config.TaskInfo) error {
	
	// 必要参数校验
	if t.Id == "" {
		return fmt.Errorf("ValidateTask Id is null:")
	}

	return nil
}

