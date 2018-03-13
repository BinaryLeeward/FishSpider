package spider

import (
	"bufio"
	"conf"
	"database/sql"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
	"github.com/toqueteos/webbrowser"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var defaultUrls = []string{"http://fangzi.xmfish.com/web/search_hire.html"}
var preURL = "http://fangzi.xmfish.com"
var hasNewRecord = false

func Catch() {
	createTable()
	urlFile, err := os.Open(conf.URL_FILE)
	defer urlFile.Close()
	urls := make([]string, 0, 10)
	if err == nil {
		scanner := bufio.NewScanner(urlFile)
		for scanner.Scan() {
			log.Println("read from file.........")
			log.Println(scanner.Text()) // Println will add back the final '\n'
			urls = append(urls, scanner.Text())
		}
	}
	if len(urls) <= 0 {
		urls = defaultUrls
	}
	log.Printf(".......url length %d.....\n", len(urls))
	for {
		hasNewRecord = false
		log.Println(".......begin parse.....\n")

		for i := 0; i < len(urls); i++ {
			log.Println(".....parse url....." + urls[i] + "\n")
			parseURL(urls[i])
		}
		log.Println("........end  parse.....\n")
		if hasNewRecord {
			webbrowser.Open("http://localhost:" + strconv.Itoa(9999) + "/index.html")
		}
		time.Sleep(conf.SLEEP_SECONDS * 1000 * time.Millisecond)
	}
}

func parseURL(url string) {

	doc, err := goquery.NewDocument(url)
	checkErr(err)

	doc.Find(".list-word").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("h3 a").Attr("href")
		url := preURL + href
		title := s.Find("h3 a").Text()
		attr := s.Find(".list-attr").Text()
		addr := s.Find(".list-addr").Text()
		price := s.Find(".list-price").Text()
		publishTime := s.Find(".list-square").Next().Text()
		publishTime = strings.Replace(publishTime, "最近更新：", "", 1)
		insert(url, title, attr, addr, price, publishTime)
	})
}

func createTable() {
	db, err := sql.Open("sqlite3", conf.DB_FILE)
	defer db.Close()
	checkErr(err)
	_, err = db.Exec("create table if not exists xm_fish (id integer primary key autoincrement, url varchar(256),title varchar(256),attr varchar(256),addr varchar(256),price varchar(256),publish_time date,create_time date null)")
	checkErr(err)
}

func insert(url, title, attr, addr, price, publishTime string) {
	db, err := sql.Open("sqlite3", conf.DB_FILE)
	defer db.Close()
	checkErr(err)
	rows, err := db.Query("select * from xm_fish where url = ?", url)
	defer rows.Close()
	checkErr(err)
	if !rows.Next() {
		stmt, err := db.Prepare("INSERT INTO xm_fish(url,title,attr,addr,price,publish_time,create_time) values(?,?,?,?,?,?,?)")
		checkErr(err)
		_, err = stmt.Exec(url, title, attr, addr, price, publishTime, time.Now())
		defer stmt.Close()
		checkErr(err)
		hasNewRecord = true
	}
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
