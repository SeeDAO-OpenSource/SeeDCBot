package main

// 当 notion database有新的page id后，将page内容同步至Dc论坛
import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jomei/notionapi"
)

// 判断ID是否在列表中
func isIDInList(id string, idList []string) bool {
	for _, item := range idList {
		if id == item {
			return true
		}
	}
	return false
}

// 格式化成MD
func descriptionformat(description string) string {
	descriptionMarkdown := "> " + strings.ReplaceAll(description, "\n", "\n> ")
	return descriptionMarkdown
}

func monitorNotionTavern() {
	// // 连接 notion
	// client := notionapi.NewClient(notionapi.Token(conf.Notion_auth))

	// 获取 notion database 数据
	notionBountyDb, err := notionClient.Database.Query(context.Background(), notionapi.DatabaseID(conf.TavernSync_NotionDb_id), &notionapi.DatabaseQueryRequest{})
	if err != nil {
		log.Println("err", err)
		return
	}
	// 获取本地数据库已经同步的列表
	idList := selectBountyList()

	// 对比本地数据库和notion database进行同步新增
	length := len(notionBountyDb.Results)
	for i := length - 1; i >= 0; i-- {
		// 悬赏帖子
		value := notionBountyDb.Results[i]
		bountyId := string(value.ID)

		// 不同步已归档与已认领的悬赏
		bountyState, ok := value.Properties["悬赏状态"].(*notionapi.SelectProperty)
		if !ok || len(bountyState.Select.Name) == 0 {
		} else {
			if bountyState.Select.Name == "已归档" || bountyState.Select.Name == "已认领" {
				continue
			}
		}

		// 判断是否已被同步
		if isIDInList(bountyId, idList) {
			log.Println("ID %s 在列表中\n", bountyId)
			continue
		}

		// 提取属性值
		bountyName := "未填写"
		bountyDescription := "未填写"
		bountyType := "未填写"
		bountSkillRequire := "未填写"
		bountyReward := "未填写"
		bountyContactPerson := "未填写"
		bountyContactWechat := "未填写"
		bountyContactDc := "未填写"
		bountExpireDate := "未填写"

		// 判断各属性的值是否存在，存在则赋值
		bountyNameProp, ok := value.Properties["悬赏名称"].(*notionapi.TitleProperty)
		if !ok || len(bountyNameProp.Title) == 0 {
			log.Println("悬赏名称属性不存在或为空")
		} else {
			bountyName = bountyNameProp.Title[0].PlainText
			log.Println("悬赏名称:", bountyName)
		}

		bountyDescriptionProp, ok := value.Properties["任务说明"].(*notionapi.RichTextProperty)
		if !ok || len(bountyDescriptionProp.RichText) == 0 {
			log.Println("任务说明属性不存在或为空")
		} else {
			bountyDescription = bountyDescriptionProp.RichText[0].PlainText
			// bountyDescriptionFormat = descriptionformat(bountyDescription)
			log.Println("任务说明:", bountyDescription)
		}

		bountyTypeProp, ok := value.Properties["悬赏类型"].(*notionapi.MultiSelectProperty)
		if !ok || len(bountyTypeProp.MultiSelect) == 0 {
			log.Println("悬赏类型属性不存在或为空")
		} else {
			bountyType = bountyTypeProp.MultiSelect[0].Name
			// bountyTypeFormat = descriptionformat(bountyType)
			log.Println("悬赏类型:", bountyType)
		}

		bountSkillRequireProp, ok := value.Properties["技能要求"].(*notionapi.RichTextProperty)
		if !ok || len(bountSkillRequireProp.RichText) == 0 {
			log.Println("技能要求属性不存在或为空")
		} else {
			bountSkillRequire = bountSkillRequireProp.RichText[0].PlainText
			log.Println("技能要求:", bountSkillRequireProp)
		}

		bountyRewardProp, ok := value.Properties["贡献报酬"].(*notionapi.RichTextProperty)
		if !ok || len(bountyRewardProp.RichText) == 0 {
			log.Println("贡献报酬属性不存在或为空")
		} else {
			bountyReward = bountyRewardProp.RichText[0].PlainText
			// bountyRewardFormat = descriptionformat(bountyReward)
			log.Println("贡献报酬:", bountyReward)
		}

		bountyContactPersonProp, ok := value.Properties["对接人"].(*notionapi.RichTextProperty)
		if !ok || len(bountyContactPersonProp.RichText) == 0 {
			log.Println("对接人属性不存在或为空")
		} else {
			bountyContactPerson = bountyContactPersonProp.RichText[0].PlainText
			// bountyContactPersonFormat = descriptionformat(bountyContactPerson)
			log.Println("对接人:", bountyContactPerson)
		}

		bountyContactWechatProp, ok := value.Properties["联系方式：微信"].(*notionapi.RichTextProperty)
		if !ok || len(bountyContactWechatProp.RichText) == 0 {
			log.Println("联系方式：微信属性不存在或为空")
		} else {
			bountyContactWechat = bountyContactWechatProp.RichText[0].PlainText
			// bountyContactWechatFormat = descriptionformat(bountyContactWechat)
			log.Println("联系方式：微信:", bountyContactWechat)
		}

		bountyContactDcProp, ok := value.Properties["联系方式：Discord"].(*notionapi.RichTextProperty)
		if !ok || len(bountyContactDcProp.RichText) == 0 {
			log.Println("联系方式：Discord属性不存在或为空")
		} else {
			bountyContactDc = bountyContactDcProp.RichText[0].PlainText
			// bountyContactDcFormat = descriptionformat(bountyContactDc)
			log.Println("联系方式：Discord:", bountyContactDc)
		}

		bountExpireDateProp, ok := value.Properties["招募截止时间"].(*notionapi.RichTextProperty)
		if !ok || len(bountExpireDateProp.RichText) == 0 {
			log.Println("招募截止时间属性不存在或为空")
		} else {
			bountExpireDate = bountExpireDateProp.RichText[0].PlainText
			log.Println("招募截止时间:", bountExpireDateProp)
		}

		descriptionFormat := `### 悬赏类型：
%s
### 技能要求:
%s
### 贡献报酬： 
%s
### 对接人：
%s
### 微信：
%s
### Discord：
%s
### 招募截止时间:
%s`
		description := fmt.Sprintf(descriptionFormat, bountyType, bountSkillRequire, bountyReward, bountyContactPerson, bountyContactWechat, bountyContactDc, bountExpireDate)

		embedContent := []*discordgo.MessageEmbed{{
			Description: description,
			Color:       0x00ff00, // Green color

		}}

		bounty := BountyList{
			Id:          bountyId,
			Name:        bountyName,
			Description: description,
		}

		// DC 创建帖子
		// 根据 Discord API 文档，auto_archive_duration 参数的有效值是 60、1440、4320 或 10080。这些值分别对应着 1 小时、1 天、3 天和 7 天。
		log.Println(discordSession)
		forumMessage, err := discordSession.ForumThreadStart(conf.TavernSync_DcChannel_id, bounty.Name, 10080, "### 任务说明：\n"+bountyDescription+"\n"+"@everyone")
		if err != nil {
			log.Println("DC插入失败：", err)
			return
		}
		// 发送联络信息
		discordSession.ChannelMessageSendEmbeds(forumMessage.ID, embedContent)

		// 将同步数据插入本地数据库
		insertBountyList(bounty)
	}
}

func tavernSyncInit() {
	// 启动时直接执行一次
	monitorNotionTavern()
	log.Println("酒馆同步初始化成功")

	// 定时器执行
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Println(time.Now())
			start := time.Now()
			monitorNotionTavern()
			end := time.Now()
			log.Println("同步完成，执行耗时:", end.Sub(start))
		}
	}

}
