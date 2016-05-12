package server

import (
	"conf"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	_ "log"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type Info struct {
	Id          int
	Url         string
	Title       string
	Attr        string
	Addr        string
	Price       string
	PublishTime time.Time
	CreateTime  time.Time
}

var pageHTML = `<html>{{range .}} {{.Id}} <a href="{{.Url}}" target="_blank">链接</a> {{.Price}} —  {{.PublishTime}} <br/> &nbsp;&nbsp;&nbsp;&nbsp;{{.Title}} <br/> &nbsp;&nbsp;&nbsp;&nbsp;{{.Attr}} {{.Addr}} <hr/><p/> {{end}}</html>`

func handler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", conf.DB_FILE)
	defer db.Close()
	checkErr(err)
	rows, err := db.Query("select * from xm_fish order by publish_time desc limit ?", conf.SHOW_COUNT)
	defer rows.Close()
	list := make([]Info, 0, conf.SHOW_COUNT)
	for rows.Next() {
		info := new(Info)
		err = rows.Scan(&info.Id, &info.Url, &info.Title, &info.Attr, &info.Addr, &info.Price, &info.PublishTime, &info.CreateTime)
		list = append(list, *info)
		checkErr(err)
	}
	t := template.New("page") //创建一个模板
	t, err = t.Parse(pageHTML)
	checkErr(err)
	t.Execute(w, list)
}

func Run() {
	http.HandleFunc("/index.html", handler)
	http.ListenAndServe(":"+strconv.Itoa(9999), nil)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
