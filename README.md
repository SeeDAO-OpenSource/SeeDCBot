# <p align="center"> SeeDcBot </p>
<p align="center"> 将SeeDAO生态连接至Discord社区 </p>

## 支持功能
* Notion 悬赏酒馆同步至Dc论坛

## 安装部署
1. 下载依赖包：
    ```shell
    go mod init main
    go mod tidy
    ```
2. 配置文件 
    ```go
    // 配置 config.json
    {
        // discord token
        "Discord_bot_auth": "",
        // notion token
        "Notion_auth": "",

        // 酒馆同步配置
        // database id
        "TavernSync_NotionDb_id": "",
        // dc频道id
        "TavernSync_DcChannel_id": ""
    }
    ```
3. 构建项目：
    ```shell
    go build
    ```
4. 运行项目:
    ```
    chmod 700 ./main
    ./main
    ```

## 开发指南
1. Notion 调用
    ```
    notionClient.
    ```
2. Discord 调用
    ```
    discordSession.
    ```
3. Sqlite 调用
    ```
    sql.Open("sqlite3", "./....db")
    ```