package main

import (
	"encoding/json"
	"os"
)

type Conf struct {
	Calendar_id             string
	Guild_id                string
	Discord_bot_auth        string
	Notion_auth             string
	TavernSync_NotionDb_id  string
	TavernSync_DcChannel_id string
}

var conf *Conf

func configLoad() (*Conf, error) {
	var conf Conf
	// 打开文件
	config_file, _ := os.Open("config.json")

	// 关闭文件
	defer config_file.Close()

	//NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	decoder := json.NewDecoder(config_file)

	conf = Conf{}

	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	decoder.Decode(&conf)
	return &conf, nil
}

func main() {
	conf, _ = configLoad()

	// sqlite 初始化
	dbInit()
	// dc 初始化
	dcInit()
	// notion 初始化
	notionInit()
	// 谷歌初始化
	googleInit()

	// 酒馆同步初始化
	go tavernSyncInit()
	// google日历同步初始化
	go calendarsync()
	<-make(chan struct{})

}
