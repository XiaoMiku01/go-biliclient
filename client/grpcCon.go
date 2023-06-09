package client

import (
	"crypto/tls"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	"github.com/XiaoMiku01/bilibili-grpc-api-go/bilibili/rpc"
)

func NewGrpcClient(addr string) *grpc.ClientConn {
	creds := grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		MinVersion: tls.VersionTLS11,
	}))
	kacp := keepalive.ClientParameters{
		PermitWithoutStream: true,
	}
	var err error
	grpcClient, err := grpc.Dial(addr, creds, grpc.WithKeepaliveParams(kacp))
	if err != nil {
		log.Fatalf("gRpc (%s) 连接失败: %v", addr, err)
		return nil
	}
	log.Debugf("gRpc (%s) 连接成功", addr)
	return grpcClient
}

// 格式化grpc错误信息
func grpcErr(err error) {
	status, ok := status.FromError(err)
	if !ok {
		log.Errorf("gRpc 未知报错: %v", err)
		return
	}
	// B站的grpc接口返回的错误码 例如鉴权错误
	if status.Code() == codes.Unknown && len(status.Details()) > 0 {
		rpc, ok := status.Details()[0].(*rpc.Status)
		if !ok {
			log.Errorf("B站 gRpc 接口未知报错: %v", err)
			return
		}
		log.Errorf("B站 gRpc 接口报错,错误码: %d, 错误信息: %s", rpc.Code, rpc.Message)
		return
	} else {
		log.Errorf("B站 gRpc 接口未知报错: %v", err)
		return
	}
}
