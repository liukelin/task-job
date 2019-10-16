package httpserver
/*
* http server
*/

import (
	"log"
	"net/http"
	"taskjob/config"
)

/**
* web server 入口
*/
func HttpServer() {

	http.HandleFunc("/", IndexAction)
	http.HandleFunc("/test", TestAction)
	// 任务创建
	http.HandleFunc("/task/create", TaskCreateAction)
	// 控制任务
	http.HandleFunc("/task/control", TaskControlAction)
	

	portStr := ":" + config.Conf.ServerPort
	log.Println("http Server Port:", portStr)
	// mux := http.NewServeMux()
	// err := http.ListenAndServe(portStr, mux)
	err := http.ListenAndServe(portStr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
