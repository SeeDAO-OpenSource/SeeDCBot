package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

var googleClient *http.Client

func googleInit() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "googlekey.json")

	ctx := context.Background()
	client, err := google.DefaultClient(ctx, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("创建HTTP客户端失败：%v", err)
	}
	log.Println("Google 初始化成功")
	googleClient = client

}
