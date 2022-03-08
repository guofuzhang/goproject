package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
)

func article(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		r.ParseForm()
		title := r.FormValue("title")
		log.Println(title)
		titleHtml:=template.HTML(title)
		t, _ := template.ParseFiles("xss/article.html")
		t.Execute(w, titleHtml)
		return
	}

}

func main() {
	http.HandleFunc("/article", article)       //设置访问的路由     //设置访问的路由
	err := http.ListenAndServe(":80", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
