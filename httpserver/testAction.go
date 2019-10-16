package httpserver

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"time"
)

/**
* [testAction 请求业务处理]
* @param  {[type]} w http.ResponseWriter [description]
* @param  {[type]} r *http.Request       [description]
* @return {[type]}   [description]
*/
func TestAction(w http.ResponseWriter, r *http.Request) {
	// 解析参数, 默认是不会解析的
	r.ParseForm()

	// d := r.Form["d"]
	// d := r.FormValue("d")
	// key := r.FormValue("sign")
	s, _ := ioutil.ReadAll(r.Body) //把  body 内容读入字符串 s
	fmt.Println(time.Now(), "body:", s)
	fmt.Println(time.Now(), "form:", r.Form)
	fmt.Fprintf(w, "success.")
	// io.WriteString(w, "success.")
}