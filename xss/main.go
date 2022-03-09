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

type Article struct {
	Title template.HTML
	Detail template.HTML
}

func article(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	db, err := sql.Open("mysql", "acurd:acurd@tcp(127.0.0.1:3306)/test")
	if r.Method == "GET" {
		sqlStr:="SELECT title,detail FROM an_article WHERE id=19"
		var article Article
		var title string
		var detail string
		err = db.QueryRow(sqlStr).Scan(&title,&detail)
		article.Title=template.HTML(title)
		article.Detail=template.HTML(detail)
		log.Println(article)
		t, err := template.ParseFiles("xss/article.html")
		if err != nil{
			log.Fatal(err)
		}
		err = t.Execute(w, article)
		if err != nil{
			log.Fatal(err)
		}
		return
	}

	if r.Method=="POST"{
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		sqlStr := "INSERT INTO an_article (title,detail) VALUES (?,?)"

		r.ParseForm()
		fmt.Println("title:", r.Form["title"])
		fmt.Println("detail:", r.Form["detail"])
		var title = strings.Join(r.Form["title"], "")
		var detail = strings.Join(r.Form["detail"], "")
		res, _ := db.Exec(sqlStr, title, detail)
		id,_:=res.LastInsertId()
		fmt.Println(id)
	}

}

func main() {
	http.HandleFunc("/article", article)       //设置访问的路由     //设置访问的路由
	err := http.ListenAndServe(":80", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
