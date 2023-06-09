package onebot

import (
	onebot "github.com/FishZe/go-libonebot"
	v12 "github.com/FishZe/go-libonebot/connect/v12"
	log "github.com/sirupsen/logrus"
)

type BiliBot struct {
	bot *onebot.Bot
}

// NewBot
// TODO: 传入配置文件
func NewBot() *BiliBot {
	onebotConfig := onebot.OneBotConfig{PlatForm: "BiliBili", Version: "1.0.0", Implementation: "BiliClient"}
	// 创建一个bot
	bot := onebot.NewOneBot(onebotConfig, "123456789")
	// 创建一个连接
	conn, err := v12.NewOneBotV12Connect(v12.OneBotV12Config{
		// 连接类型
		// Http
		ConnectType: v12.ConnectTypeHttp,
		HttpConfig: v12.OneBotV12HttpConfig{
			Host:            "127.0.0.1",
			Port:            20003,
			EventEnable:     true,
			EventBufferSize: 500,
		},
	})
	if err != nil {
		log.Println(err)
	}
	// 把连接加入到bot
	err = bot.AddConnection(conn)
	if err != nil {
		log.Println(err)
	}
	conn2, err := v12.NewOneBotV12Connect(v12.OneBotV12Config{
		ConnectType: v12.ConnectTypeWebSocketReverse,
		WebsocketReverseConfig: v12.OneBotV12WebsocketReverseConfig{
			Url:               "ws://192.168.81.137:20001",
			ReconnectInterval: 5000,
		},
	})
	if err != nil {
		log.Println(err)
	}
	err = bot.AddConnection(conn2)
	if err != nil {
		log.Println(err)
	}
	return &BiliBot{bot: bot}
}

func (b *BiliBot) SendEvent(e any) {
	// 发送事件
	b.bot.SendEvent(e)
}
