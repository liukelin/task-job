package server

/**
 task相关
 @auth liukelin
*/

import (
	"fmt"
	"time"
	"encoding/json"
	"path/filepath"
	"github.com/golang/glog"
	"github.com/go-cmd/cmd"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"taskjob/config"
)

type Task struct {
	TaskInfo    config.TaskInfo		// 任务内容，详细内容。需要执行的命令
	cancel      func()				// 停止任务，开关
	errG        *errgroup.Group		
	Qconn		Queue				// mq连接
}

// 实际执行job
func TaskRun(t *Task, ChPool chan int64){
	defer func(){
		// 释放池
		<-ChPool
		
		// 释放map
		TaskPool.Delete(t.TaskInfo.Id)
	}()

	// 记录任务
	TaskPool.Store(t.TaskInfo.Id, t)

	// 执行任务
	t.Execute()
}

// 任务控制器
func TaskController(t *config.TaskControl) {
	// 判断任务是否在执行
	if v,ok := TaskPool.Load(t.Id);ok{
		// 断言取出
		vv,ok := v.(*Task)
		if !ok {
			glog.Infof("Task map val Type:%s, ->err.", t.Id)
			return
		}

		// 操作类型
		switch t.Type {
		case config.CONTROL_TYPE_KILL:  // 中止任务
			
			err := vv.Stop()
			if err != nil {
				glog.Infof("Task Control:%s, Id:%s, ->err:%v", config.CONTROL_TYPE_KILL, t.Id, err)
				return
			}
			glog.Infof("Task Control:%s, Id:%s ->success.", config.CONTROL_TYPE_KILL, t.Id)
		}
	}
	glog.Infof("TaskPool Id:%s, ->None.", t.Id)
	return
}

// 执行shell
func (task *Task) Execute() (err error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	task.cancel = cancelFunc

	if task.TaskInfo.TaskData.Timeout > 0 {
		ctx, _ = context.WithTimeout(ctx, time.Duration(task.TaskInfo.TaskData.Timeout)*time.Second)
	}
	task.errG, ctx = errgroup.WithContext(ctx)

	/**
	// 用于拆分子任务
	childTaskPropertyList, err := task.splitTask()
	if err != nil {
		return err
	}
	glog.Infoln("childTaskPropertyList->", childTaskPropertyList)
	**/
	childTaskNum := 1 // 任务拆分, 冗余，后期考虑子任务场景
	for i := 0; i < childTaskNum; i++ {
		index := i
		task.errG.Go(func() error {
			childTaskID := fmt.Sprintf("%s_%d", task.TaskInfo.Id, index)

			glog.Infof("Run TaskChildTasId->%s", childTaskID)

			// Disable output buffering, enable streaming
			cmdOptions := cmd.Options{
				Buffered:  false,
				Streaming: true,
			}

			CmdArgs := []string{filepath.Join(config.Conf.CmdDir, task.TaskInfo.TaskData.Cmdline)}
			CmdArgs = append(CmdArgs, task.TaskInfo.TaskData.CmdArgs...)

			// Create Cmd with options
			childCmd := cmd.NewCmdOptions(cmdOptions, "bash", CmdArgs...)

			// Print STDOUT and STDERR lines streaming from Cmd
			// 获取执行输出
			go func() {
				for {
					select {
					case line := <-childCmd.Stdout:
						
						glog.Infof("childCmd.Stdout:%s.", line)
						// 执行输出
						if config.Conf.Logs {
							err_ := task.PushTaskLog(line, config.LOG_TYPE_STDOUT, index)
							if err_ != nil {
								glog.Error("PushTaskLog error:%s, %s, %v", childTaskID, config.LOG_TYPE_STDOUT, err_)
							}
						}

					case line := <-childCmd.Stderr:
						
						glog.Infof("childCmd.Stderr:%s, %s", childTaskID, line)
						//执行报错
						err_ := task.PushTaskLog(line, config.LOG_TYPE_STDERR, index)
						if err_ != nil {
							glog.Error("PushTaskLog error:%s, %s, %v", childTaskID, config.LOG_TYPE_STDERR, err_)
						}
						
					case line := <-childCmd.Done():
						// 子任务完成调用
						glog.Infof("TaskLchildCmd:%s:%v:Done.", childTaskID, line)
						err_ := task.PushTaskLog(config.LOG_TYPE_DONE, config.LOG_TYPE_DONE, index)
						if err_ != nil {
							glog.Error("PushTaskLog error:%s,%v", config.LOG_TYPE_DONE, err_)
						}
						return
					}
				}
			}()

			select {
			case status := <-childCmd.Start():
				// Cmd has finished but wait for goroutine to print all lines
				for len(childCmd.Stdout) > 0 || len(childCmd.Stderr) > 0 {
					time.Sleep(10 * time.Millisecond)
				}
				glog.Infof("childTaskID->%s,executed status->%v", childTaskID, status)
				if status.Exit != 0 && status.Error == nil {
					return fmt.Errorf("childTaskID->%s execute failed, EXIT CODE:%d", childTaskID, status.Exit)
				}
				return status.Error
			case <-ctx.Done(): // 强制结束, 整个任务
				glog.Infof("childTaskID->%s, kill has been CANCELED", childTaskID)
				err := childCmd.Stop()
				if err != nil {
					glog.Error("failed to stop childTaskID->%s, err: %v", childTaskID, err)
				}
				glog.Infof("childTaskID->%s,ctx.Done.", childTaskID)
				return ctx.Err()
			}
		})
	}

	// glog.Infof("wait->Runing...") 
	// 阻塞等待所有字任务完成
	if err := task.errG.Wait(); err != nil {
		return err
	}
	// glog.Infof("wait->Done") 
	return nil
}

/**
  结束任务
 **/
func (task *Task) Stop() error {
	if task.cancel != nil {
		task.cancel()
	}
	return nil
}

// 上报日志
func (task *Task) PushTaskLog(line string, logType string, childTaskNum int) error {
	if (line == ""){
		return nil
	}
	log := &config.TaskLog{
		Id: 		task.TaskInfo.Id,
		Type:		task.TaskInfo.Type,
		ChildNum: 	childTaskNum, 
		MsgType:	logType,
		Data:		line,
	}
	buf, err := json.MarshalIndent(log, "", "    ")
    if err != nil {
        return err
	}
	err_ := task.Qconn.Push(string(buf[:]), config.Conf.LogsChannel)
	if err_ != nil {
		return err_
	}
	return nil
}

// 解析任务内容
func ParsTaskJson(data string) (Task, error) {
	var taskInfo config.TaskInfo
	
	task := &Task{}
	// 解析json
	err := json.Unmarshal([]byte(data), &taskInfo)
	if err != nil {
		return *task, err
	}
	// 校验任务内容是否合法
	err_ := ValidateTask(taskInfo)
	if err_ != nil {
		return *task, err_
	}
	task = &Task{
		TaskInfo: taskInfo,
	}

	return *task, nil
}

// 解析任务控制 内容
func ParsTaskControlJson(data string) (config.TaskControl, error) {
	var t config.TaskControl
	// 解析json
	err := json.Unmarshal([]byte(data), &t)
	if err != nil {
		return t, err
	}
	// 校验任务控制参数是否合法

	return t, nil
}
