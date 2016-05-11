package main

import (
	"server"
	"spider"
)

func main() {
	go spider.Catch()
	server.Run()
}
