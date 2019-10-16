package httpserver

import (
	"fmt"
	"net/http"
	"io/ioutil"
	// "encoding/json"
	"time"
)

/**
  接收任务请求
  POST
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
 */
 func TaskCreateAction(w http.ResponseWriter, r *http.Request) {
	// 解析参数, 默认是不会解析的
	r.ParseForm()

	body, _ := ioutil.ReadAll(r.Body)
	// 解析

	fmt.Println(time.Now(), "args:", r.Form)
	fmt.Println(time.Now(), "body:", string(body[:]))
	fmt.Fprintf(w, "success.")
	// io.WriteString(w, "success.")
}


/**
 任务控制
*/
func TaskControlAction(w http.ResponseWriter, r *http.Request) {
	// 解析参数, 默认是不会解析的
	r.ParseForm()

	// d := r.Form["d"]
	// d := r.FormValue("d")
	// key := r.FormValue("sign")
	fmt.Println(time.Now(), "args:", r.Form)
	fmt.Fprintf(w, "success.")
	// io.WriteString(w, "success.")
}
