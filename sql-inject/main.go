package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("sql-inject/login.html")
		t.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		r.ParseForm()
		fmt.Println("name:", r.Form["name"])
		fmt.Println("password:", r.Form["password"])
		var name = strings.Join(r.Form["name"], "")
		var password = strings.Join(r.Form["password"], "")
		var id int
		db, err := sql.Open("mysql", "acurd:acurd@tcp(127.0.0.1:3306)/test")
		if err != nil {
			log.Fatal(err)
		}
		sqlStr:="SELECT id FROM user WHERE name = " + name + " AND password =" + password + ""
		log.Println(sqlStr)
		err = db.QueryRow(sqlStr).Scan(&id)
		err = db.QueryRow(sqlStr,name,password).Scan(&id)
		if err == sql.ErrNoRows {
			w.Write([]byte("登录失败!"))
			return
		}

		if id != 0 {
			w.Write([]byte("登录成功!"))
			return
		} else {
			w.Write([]byte("登录失败!"))
			return
		}
		return
	}
}

func main() {
	http.HandleFunc("/login", login)       //设置访问的路由     //设置访问的路由
	err := http.ListenAndServe(":80", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
