package client

import (
	notifyapi "github.com/XiaoMiku01/bilibili-grpc-api-go/bilibili/broadcast/message/im"
	"github.com/XiaoMiku01/go-biliclient/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func eventHandler(e *anypb.Any) {
	//log.Println(e)
	switch e.TypeUrl {
	case "type.googleapis.com/bilibili.broadcast.v1.HeartbeatResp":
		// 心跳包回复
		log.Debugln("收到心跳包回复:", e)
	case "type.googleapis.com/bilibili.broadcast.message.im.NotifyRsp":
		// 通知回复
		notifyHandler(e)
	}
}

func notifyHandler(e *anypb.Any) {
	var notifyRsp notifyapi.NotifyRsp
	err := anypb.UnmarshalTo(e, &notifyRsp, proto.UnmarshalOptions{})
	if err != nil {
		log.Errorln("通知回复解析失败:", err)
		return
	}
	if notifyRsp.Cmd == uint64(notifyapi.CmdId_EN_CMD_ID_MSG_NOTIFY) {
		switch notifyRsp.PayloadType {
		case notifyapi.PLType_EN_PAYLOAD_BASE64:
			log.Debugln(utils.AnyToJSON(e))
			// 私信 or 应援团
			// TODO
			session := GetOneClient().GetNewSession()
			for _, v := range session.SessionList {
				log.Infoln("收到:", v)
			}
		default:
			log.Debugln("收到未知通知回复:", utils.AnyToJSON(e))
		}
	} else {
		log.Debugln("收到未知通知回复:", utils.AnyToJSON(e))
	}
}
