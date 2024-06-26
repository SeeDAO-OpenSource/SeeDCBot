package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var discordSession *discordgo.Session

func dcInit() {
	// 初始化连接

	discordSession, _ = discordgo.New("Bot " + conf.Discord_bot_auth)

	err := discordSession.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Dc 初始化成功")

}
