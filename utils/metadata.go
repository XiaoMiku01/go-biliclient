package utils

import (
	bilimetadata "github.com/XiaoMiku01/bilibili-grpc-api-go/bilibili/metadata"
	"github.com/XiaoMiku01/bilibili-grpc-api-go/bilibili/metadata/device"
	"github.com/XiaoMiku01/bilibili-grpc-api-go/bilibili/metadata/locale"
	"github.com/XiaoMiku01/bilibili-grpc-api-go/bilibili/metadata/network"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func GetBiliMetaData(accessKey string) metadata.MD {
	device := &device.Device{
		MobiApp:  "android",
		Device:   "phone",
		Build:    6830300,
		Channel:  "bili",
		Buvid:    "XX82B818F96FB2F312B3A1BA44DB41892FF99",
		Platform: "android",
	}
	devicebin, _ := proto.Marshal(device)
	locale := &locale.Locale{
		Timezone: "Asia/Shanghai",
	}
	localebin, _ := proto.Marshal(locale)
	bilimetadata := &bilimetadata.Metadata{
		AccessKey: accessKey,
		MobiApp:   "android",
		Device:    "phone",
		Build:     6830300,
		Channel:   "bili",
		Buvid:     "XX82B818F96FB2F312B3A1BA44DB41892FF99",
		Platform:  "android",
	}
	bilimetadatabin, _ := proto.Marshal(bilimetadata)
	network := &network.Network{
		Type: network.NetworkType_WIFI,
	}
	networkbin, _ := proto.Marshal(network)
	md := metadata.Pairs(
		"x-bili-device-bin", string(devicebin),
		"x-bili-local-bin", string(localebin),
		"x-bili-metadata-bin", string(bilimetadatabin),
		"x-bili-network-bin", string(networkbin),
		"authorization", "identify_v1 "+accessKey,
	)
	return md
}
