package main

import (
	"controller"
	"log"
	"model"
	"net/http"
)

func main() {

	http.Handle("/css/", http.FileServer(http.Dir("style")))
	http.Handle("/blog/", http.FileServer(http.Dir("style")))
	http.Handle("/js/", http.FileServer(http.Dir("style")))

	// golang database/sql 中自带连接池，一个go进程执行一次sql.Open就好了
	db, err := model.Init()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// blog首页
	http.HandleFunc("/", controller.List)
	http.HandleFunc("/index", controller.List)
	http.HandleFunc("/list", controller.List)

	// 写blog
	http.HandleFunc("/blog", controller.Blog)
	http.HandleFunc("/blog/save", controller.Save)

	// blog详细页面
	http.HandleFunc("/blog/c", controller.Detail)

	// 修改
	http.HandleFunc("/blog/edit", controller.Edit)
	http.HandleFunc("/blog/doedit", controller.DoEdit)

	// 删除
	http.HandleFunc("/blog/delete", controller.Delete)

	// 根据Tag查询
	http.HandleFunc("/blog/findbytag", controller.FindByTag)

	// http.HandleFunc("/tag", controller.Tag)

	http.HandleFunc("/dologin", controller.DoLogin)
	http.HandleFunc("/login", controller.Login)

	err1 := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err1 != nil {
		log.Fatal("ListenAndServe: ", err1)
	}
}
