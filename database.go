package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type BountyList struct {
	Id          string
	Name        string
	Description string
}

type calendarEvent struct {
	Id                 string
	Name               string
	Creator            string
	Description        string
	ScheduledStartTime time.Time
	ScheduledEndTime   *time.Time
	Location           string
	GeventId           string
}

func dbInit() {
	// 酒馆同步数据库
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
		log.Fatalf("tavernDb open database fail")
		return
	}

	log.Println("酒馆数据库初始化成功")

	// 日历同步数据库
	calendarDb, _ := sql.Open("sqlite3", "./calendar_data.db")

	// 创建悬赏表
	_, err = calendarDb.Exec(`CREATE TABLE IF NOT EXISTS calendarlist (
						id TEXT,
						name TEXT,
						creator TEXT,
						description TEXT,
						scheduledstarttime TIME,
						scheduledendtime TIME,
						location TEXT,
						geventid TEXT,
						rm  BOOLEAN DEFAULT false
						)`)
	if err != nil {
		log.Println(err)
	}

	if err := calendarDb.Ping(); err != nil {
		log.Fatalf("calendarDb open database fail")
		return
	}

	log.Println("日历数据库初始化成功")

	// 单独功能 单独创建数据库
	// xxxDb, _ := sql.Open("sqlite3", "./xxx.db")

}

// 插入悬赏数据到悬赏表
func insertBountyList(bounty BountyList) {
	db, _ := sql.Open("sqlite3", "./bounty_data.db")

	_, err := db.Exec("INSERT INTO bountylist(id, name, description) VALUES(?, ?, ?)", bounty.Id, bounty.Name, bounty.Description)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("bounty inserted successfully.")
}

// 插入数据到日历同步表
func insertCalendarList(calendar calendarEvent) {
	db, _ := sql.Open("sqlite3", "./calendar_data.db")

	_, err := db.Exec("INSERT INTO calendarlist(id, name, creator, description, scheduledstarttime, scheduledendtime, location, geventid) VALUES(?, ?, ?, ?, ?, ?, ?, ?)", calendar.Id, calendar.Name, calendar.Creator, calendar.Description, calendar.ScheduledStartTime, calendar.ScheduledEndTime, calendar.Location, calendar.GeventId)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("calendar inserted successfully.")
}

// 获取悬赏表内容
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

type Event struct {
	ID                 string
	Name               string
	GEventID           string
	ScheduledStartTime time.Time
}

// 获取日历表内容
func selectCalendarList() []Event {
	var events []Event
	var scheduledStartTimeStr string

	db, _ := sql.Open("sqlite3", "./calendar_data.db")

	// 从日历表中查询数据并打印
	timeNow := time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02 15:04:05")
	sql := "SELECT id, name, scheduledstarttime, geventid FROM calendarlist where scheduledstarttime > ? AND rm == false"
	log.Println(db.Ping())

	rows, err := db.Query(sql, timeNow)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &scheduledStartTimeStr, &event.GEventID)
		if err != nil {
			log.Fatal(err)
		}
		scheduledStartTime, _ := time.Parse("2006-01-02 15:04:05-07:00", scheduledStartTimeStr)
		event.ScheduledStartTime = scheduledStartTime
		events = append(events, event)
	}
	// 检查是否有任何迭代过程中的错误
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return events
}

func updateCalendarRmMark(geventid string) {
	db, _ := sql.Open("sqlite3", "./calendar_data.db")

	// 修改已删除日历为false
	sql := "update calendarlist set rm = true where geventid = ?"
	_, err := db.Exec(sql, geventid)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("calendar updated successfully.")
}
