package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// 返回第二个 Event 数组中与第一个数组ID不同的元素
func difference(a, b []Event) []Event {
	m := make(map[string]bool)
	result := []Event{}

	// 将数组 b 中的元素添加到映射 m 中
	for _, element := range b {
		m[element.ID] = true
	}

	// 遍历数组 a，将不在映射 m 中的元素添加到结果数组中
	for _, element := range a {
		if _, ok := m[element.ID]; !ok {
			result = append(result, element)
		}
	}

	return result
}

func monitorDcEvents() {
	guildID := conf.Guild_id

	// 获取已同步列表
	calendarDbList := selectCalendarList()

	calendarDcDiffList := []Event{}

	// 特殊DC活动排除 [社区大会 翻译公会周会 研发公会周会 S7市政厅通气会]
	specialEventsID := []string{"1242826196899729519", "1237972510612520991", "1248649014858092634", "1254452696967020677"}

	// 获取DC当前列表
	events, err := discordSession.GuildScheduledEvents(guildID, true)
	if err != nil {
		log.Println("获取服务器活动失败:", err)
		return
	}

	// 遍历活动
	for _, event := range events {
		calendarDcDiff := Event{
			ID: event.ID,
		}

		calendarDcData := calendarEvent{
			Id:                 event.ID,
			Name:               event.Name,
			Creator:            event.Creator.Username,
			Description:        event.Description,
			ScheduledStartTime: event.ScheduledStartTime,
			ScheduledEndTime:   event.ScheduledEndTime,
			Location:           event.EntityMetadata.Location,
			GeventId:           "",
		}

		// 添加到对比列表
		calendarDcDiffList = append(calendarDcDiffList, calendarDcDiff)

		// 同步判断键
		shouldPush := true

		// 判断是否同步过
		for _, cacalendar := range calendarDbList {
			// 判断是否需要同步
			if calendarDcData.Id == cacalendar.ID && calendarDcData.ScheduledStartTime == cacalendar.ScheduledStartTime {
				log.Println("已同步过", calendarDcData.Name)
				shouldPush = false
				break
			}
		}
		// 判断是否在特殊列表中
		for _, specialID := range specialEventsID {
			if calendarDcData.Id == specialID {
				shouldPush = false
				log.Println("特殊活动", calendarDcData.Id, calendarDcData.Name)
				break
			}
		}

		// 执行同步
		if shouldPush {
			googleEventPush(calendarDcData)
		}
	}

	// 对比数据库中的活动与DC活动的差别，进行删除
	diff := difference(calendarDbList, calendarDcDiffList)
	for _, eventDiff := range diff {
		srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(googleClient))
		if err != nil {
			log.Println("创建Calendar服务失败：%v", err)
		}
		calendarID := conf.Calendar_id
		err = srv.Events.Delete(calendarID, eventDiff.GEventID).Do()
		if err != nil {
			log.Println("删除日历失败：%v", err)
			return
		}
		log.Println("日历删除成功", eventDiff.Name, eventDiff.ID, eventDiff.GEventID)
		updateCalendarRmMark(eventDiff.GEventID)
	}

}

func googleEventPush(calendarData calendarEvent) {

	// 计算结束时间为开始时间的一个小时后，解决周期性时间过大结束时间问题
	start := calendarData.ScheduledStartTime
	end := start.Add(time.Hour)

	// 创建日历
	event := &calendar.Event{
		Summary: calendarData.Name,
		Start: &calendar.EventDateTime{
			DateTime: start.Format("2006-01-02T15:04:05-07:00"), // 设置日历开始时间
		},
		End: &calendar.EventDateTime{
			DateTime: end.Format("2006-01-02T15:04:05-07:00"), // 设置日历结束时间
		},
		Description: calendarData.Description,
		Location:    calendarData.Location,
	}

	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(googleClient))
	if err != nil {
		log.Println("创建Calendar服务失败：%v", err)
	}

	// 添加到google日历
	// calendarID := "primary" // 默认 primary
	calendarID := conf.Calendar_id
	event, err = srv.Events.Insert(calendarID, event).Do()
	if err != nil {
		log.Println("插入日历失败：%v", err)
	}

	// 将信息插入数据库
	calendarData.GeventId = event.Id
	insertCalendarList(calendarData)
	log.Printf("日历已添加：%s\n", event.Summary, event.Id, event.HtmlLink)

}

func calendarsync() {
	// 启动时直接执行一次
	monitorDcEvents()
	log.Println("日历同步初始化成功")

	// 定时器执行
	ticker := time.NewTicker(6 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Println(time.Now())
			start := time.Now()
			monitorDcEvents()
			end := time.Now()
			log.Println("同步完成，执行耗时:", end.Sub(start))
		}
	}
}
