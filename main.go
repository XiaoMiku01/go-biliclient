package main

import (
	"github.com/XiaoMiku01/go-biliclient/client"
	_ "github.com/XiaoMiku01/go-biliclient/logger"
)

func main() {
	u := client.User{
		LoginInfo: &client.LoginInfo{
			AccessKey: "",
		},
	}
	c := client.NewBiliClient(u)
	c.BroadcastConnect()
	select {}
}
