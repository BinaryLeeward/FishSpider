package spider

import (
	"conf"
	"database/sql"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
	"github.com/toqueteos/webbrowser"
	"log"
	"strconv"
	"time"
)

var urls = []string{"http://fangzi.xmfish.com/web/search_hire.html?h=&hf=1&ca=59201&r=5920112&s=103&a=&rm=&f=&d=&tp=&l=0&tg=&hw=1&o=&ot=0&tst=0", "http://fangzi.xmfish.com/web/search_hire.html?h=&hf=1&ca=59201&r=5920113&s=103&a=&rm=&f=&d=&tp=&l=0&tg=&hw=1&o=&ot=1&tst=0", "http://fangzi.xmfish.com/web/search_hire.html?h=&hf=1&ca=59201&r=5920114&s=103&a=&rm=&f=&d=&tp=&l=0&tg=&hw=1&o=&ot=0&tst=0"}
var preURL = "http://fangzi.xmfish.com"
var hasNewRecord = false

func Catch() {
	createTable()
	for {
		hasNewRecord = false
		log.Println(".......begin parse.....")
		for i := 0; i < len(urls); i++ {
			parseURL(urls[i])
		}
		log.Println("........end  parse.....")
		if hasNewRecord {
			webbrowser.Open("http://localhost:" + strconv.Itoa(9999) + "/index.html")
		}
		time.Sleep(conf.SLEEP_SECONDS * 1000 * time.Millisecond)
	}
}

func parseURL(url string) {

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".list-word").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("h3 a").Attr("href")
		url := preURL + href
		title := s.Find("h3 a").Text()
		attr := s.Find(".list-attr").Text()
		addr := s.Find(".list-addr").Text()
		price := s.Find(".list-price").Text()

		insert(url, title, attr, addr, price)
	})
}

func createTable() {
	db, err := sql.Open("sqlite3", conf.DB_FILE)
	checkErr(err)
	_, err = db.Exec("create table if not exists xm_fish (id integer primary key autoincrement, url varchar(256),title varchar(256),attr varchar(256),addr varchar(256),price varchar(256),create_time date null)")
	checkErr(err)
}

func insert(url, title, attr, addr, price string) {
	db, err := sql.Open("sqlite3", conf.DB_FILE)
	defer db.Close()
	checkErr(err)
	rows, err := db.Query("select * from xm_fish where url = ?", url)
	checkErr(err)
	if !rows.Next() {
		stmt, err := db.Prepare("INSERT INTO xm_fish(url,title,attr,addr,price,create_time) values(?,?,?,?,?,?)")
		checkErr(err)
		_, err = stmt.Exec(url, title, attr, addr, price, time.Now())
		checkErr(err)
		err = stmt.Close()
		checkErr(err)
		hasNewRecord = true
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
