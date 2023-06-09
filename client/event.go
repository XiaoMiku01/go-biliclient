package client

import (
	"encoding/json"
	"github.com/FishZe/go-libonebot/protocol"
	notifyapi "github.com/XiaoMiku01/bilibili-grpc-api-go/bilibili/broadcast/message/im"
	"github.com/XiaoMiku01/go-biliclient/onebot"
	"github.com/XiaoMiku01/go-biliclient/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"strconv"
	"sync"
)

var (
	// eventHandleFunc e.TypeUrl -> func(e *anypb.Any)
	eventHandleFunc sync.Map
)

var (
	BiliBot *onebot.BiliBot
)

func init() {
	eventHandleFunc.Store("type.googleapis.com/bilibili.broadcast.message.im.NotifyRsp", notifyHandler)
	eventHandleFunc.Store("type.googleapis.com/bilibili.broadcast.v1.HeartbeatResp", func(e *anypb.Any) {
		log.Debugln("收到心跳包回复:", e)
	})
}

func eventHandler(e *anypb.Any) {
	if f, ok := eventHandleFunc.Load(e.TypeUrl); ok {
		f.(func(e *anypb.Any))(e)
		return
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
				evt := protocol.NewMessageEventPrivate()
				msg := make(map[string]string)
				_ = json.Unmarshal([]byte(v.LastMsg.Content), &msg)
				evt.Message = append(evt.Message, protocol.GetSegmentText(msg["content"]))
				evt.UserId = strconv.FormatUint(v.LastMsg.SenderUid, 10)
				BiliBot.SendEvent(evt)
			}
		default:
			log.Debugln("收到未知通知回复:", utils.AnyToJSON(e))
		}
	} else {
		log.Debugln("收到未知通知回复:", utils.AnyToJSON(e))
	}
}
