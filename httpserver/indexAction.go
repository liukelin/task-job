package httpserver

import (
	"fmt"
	"net/http"
	"time"
)

/**
* [indexAction 请求业务处理]
* @param  {[type]} w http.ResponseWriter [description]
* @param  {[type]} r *http.Request       [description]
* @return {[type]}   [description]
*/
func IndexAction(w http.ResponseWriter, r *http.Request) {
	// 解析参数, 默认是不会解析的
	r.ParseForm()
	// s, _ := ioutil.ReadAll(r.Body) //把  body 内容读入字符串 s
	// d := r.Form["d"]
	// d := r.FormValue("d")
	// key := r.FormValue("sign")
	fmt.Println(time.Now(), "args:", r.Form)
	fmt.Fprintf(w, "success.")
	// io.WriteString(w, "success.")
}
