package client

import (
	"context"
	broadcastapi "github.com/XiaoMiku01/bilibili-grpc-api-go/bilibili/broadcast/v1"
	imapi "github.com/XiaoMiku01/bilibili-grpc-api-go/bilibili/im/interfaces/v1"
	"github.com/XiaoMiku01/go-biliclient/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/anypb"
	"sync"
	"time"
)

type BiliClient struct {
	User                   *User
	biliGrpcClient         *grpc.ClientConn
	biliBroadcastClient    *grpc.ClientConn
	stopBroadcastSyncWg    sync.WaitGroup
	stopBroadcastHeartBeat bool
	MsgId                  int64
	Seq                    int64
}

var BiliClients []*BiliClient

func NewBiliClient(u User) *BiliClient {
	var b *BiliClient
	defer func() {
		BiliClients = append(BiliClients, b)
	}()
	biliGrpcClient := NewGrpcClient("grpc.biliapi.net:443")
	biliBroadcastClient := NewGrpcClient("broadcast.chat.bilibili.com:7824")
	b = &BiliClient{
		User:                &u,
		biliGrpcClient:      biliGrpcClient,
		biliBroadcastClient: biliBroadcastClient,
		stopBroadcastSyncWg: sync.WaitGroup{},
		MsgId:               1,
		Seq:                 1,
	}
	return b
}

func GetOneClient() *BiliClient {
	return BiliClients[0]
}

func (c *BiliClient) BroadcastConnect() {
	md := utils.GetBiliMetaData(c.User.LoginInfo.AccessKey)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	broadcastClient := broadcastapi.NewBroadcastTunnelClient(c.biliBroadcastClient)
	stream, err := broadcastClient.CreateTunnel(ctx)
	if err != nil {
		log.Fatalln("创建双向 rpc 通道失败:", err)
	}
	c.BroadcastAuth(stream)
	c.stopBroadcastSyncWg.Add(1)
	go c.broadcastHeartBeat(stream)
	c.stopBroadcastSyncWg.Add(1)
	go c.BroadcastRead(stream)
}

func (c *BiliClient) BroadcastReConnect() {
	c.stopBroadcastHeartBeat = true
	c.stopBroadcastSyncWg.Wait()
	c.stopBroadcastHeartBeat = false
	c.BroadcastConnect()
}

func (c *BiliClient) BroadcastAuth(stream broadcastapi.BroadcastTunnel_CreateTunnelClient) {
	authReq := &broadcastapi.AuthReq{
		ConnId:    utils.GetUUid(),
		Guid:      utils.GetUUid(),
		LastMsgId: c.MsgId,
	}
	anyReq, _ := anypb.New(authReq)
	frame := &broadcastapi.BroadcastFrame{
		Options: &broadcastapi.FrameOption{
			MessageId: c.MsgId,
			Sequence:  c.Seq,
		},
		TargetPath: "/bilibili.broadcast.v1.Broadcast/Auth",
		Body:       anyReq,
	}
	err := stream.Send(frame)
	if err != nil {
		log.Fatalln("发送鉴权包失败:")
		grpcErr(err)
		return
	}
	log.Println("鉴权成功")
	c.Seq++
	subReq := &broadcastapi.TargetPath{
		TargetPaths: []string{"/bilibili.broadcast.message.im.Notify/WatchNotify"},
	}
	anyReq, _ = anypb.New(subReq)
	frame = &broadcastapi.BroadcastFrame{
		Options: &broadcastapi.FrameOption{
			MessageId: c.MsgId,
			Sequence:  c.Seq,
		},
		TargetPath: "/bilibili.broadcast.v1.Broadcast/Subscribe",
		Body:       anyReq,
	}
	err = stream.Send(frame)
	if err != nil {
		log.Fatalln("发送订阅包失败:")
		grpcErr(err)
		return
	}
	c.Seq++
	log.Println("订阅通知成功")
}

func (c *BiliClient) broadcastHeartBeat(stream broadcastapi.BroadcastTunnel_CreateTunnelClient) {
	defer c.stopBroadcastSyncWg.Done()
	for !c.stopBroadcastHeartBeat {
		heartbeatReq := &broadcastapi.HeartbeatReq{}
		anyReq, _ := anypb.New(heartbeatReq)
		frame := &broadcastapi.BroadcastFrame{
			Options: &broadcastapi.FrameOption{
				MessageId: c.MsgId,
				Sequence:  c.Seq,
			},
			TargetPath: "/bilibili.broadcast.v1.Broadcast/Heartbeat",
			Body:       anyReq,
		}
		err := stream.Send(frame)
		if err != nil {
			log.Errorln("发送心跳包失败:")
			grpcErr(err)
			return
		}
		c.Seq++
		log.Debugln("发送心跳包成功", frame)
		time.Sleep(time.Second * 10)
	}
}

func (c *BiliClient) BroadcastRead(stream broadcastapi.BroadcastTunnel_CreateTunnelClient) {
	defer c.stopBroadcastSyncWg.Done()
	for !c.stopBroadcastHeartBeat {
		in, err := stream.Recv()
		if err != nil {
			log.Errorln("接收消息失败:", err)
			grpcErr(err)
			continue
		}
		if in.Options.Status != nil {
			log.Errorf("接收到错误消息,错误码%d,错误信息%s:", in.Options.Status.Code, in.Options.Status.Message)
			continue
		}
		eventHandler(in.Body)
	}
}

func (c *BiliClient) GetNewSession() *imapi.RspSessions {
	md := utils.GetBiliMetaData(c.User.LoginInfo.AccessKey)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	imClient := imapi.NewImInterfaceClient(c.biliGrpcClient)
	imReq := &imapi.ReqNewSessions{
		Size: 1,
	}
	imResp, err := imClient.NewSessions(ctx, imReq)
	if err != nil {
		log.Println("err")
		grpcErr(err)
		return nil
	}
	return imResp
}
