package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type BountyList struct {
	Id          string
	Name        string
	Description string
}

func dbInit() {
	// 连接到SQLite数据库
	tavernDb, _ := sql.Open("sqlite3", "./bounty_data.db")

	// 创建悬赏表
	_, err := tavernDb.Exec(`CREATE TABLE IF NOT EXISTS bountylist (
						id TEXT PRIMARY KEY,
						name TEXT,
						description TEXT
					)`)
	if err != nil {
		log.Println(err)
	}

	if err := tavernDb.Ping(); err != nil {
		log.Fatalf("open database fail")
		return
	}

	log.Println("酒馆数据库初始化成功")

	// 单独功能 单独创建数据库
	// xxxDb, _ := sql.Open("sqlite3", "./xxx.db")

}

// 插入悬赏数据到悬赏表
func insertBountyList(bounty BountyList) {
	db, _ := sql.Open("sqlite3", "./bounty_data.db")

	_, err := db.Exec("INSERT INTO bountylist(id, name, description) VALUES(?, ?, ?)", bounty.Id, bounty.Name, bounty.Description)
	if err != nil {
		log.Println(err)
	}
	log.Println("bounty inserted successfully.")
}

func selectBountyList() []string {
	var idList []string
	db, _ := sql.Open("sqlite3", "./bounty_data.db")

	// 从悬赏表中查询数据并打印
	sql := "SELECT id FROM bountylist"
	log.Println(db.Ping())

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		idList = append(idList, id)
	}
	// 检查是否有任何迭代过程中的错误
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return idList
}
