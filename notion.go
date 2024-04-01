package main

import (
	"log"

	"github.com/jomei/notionapi"
)

var notionClient *notionapi.Client

func notionInit() {

	notionClient = notionapi.NewClient(notionapi.Token(conf.Notion_auth))
	log.Println("notion 初始化成功")
}
