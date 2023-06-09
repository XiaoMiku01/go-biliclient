package main

import (
	"github.com/XiaoMiku01/go-biliclient/client"
	_ "github.com/XiaoMiku01/go-biliclient/logger"
	"github.com/XiaoMiku01/go-biliclient/onebot"
)

func main() {
	bot := onebot.NewBot()
	client.BiliBot = bot
	u := client.User{
		LoginInfo: &client.LoginInfo{
			AccessKey: "",
		},
	}
	c := client.NewBiliClient(u)
	c.BroadcastConnect()
	select {}
}
