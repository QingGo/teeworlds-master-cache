package myconst

import "flag"

// ServerBrowseGetList token
var ServerBrowseGetList = []uint8{255, 255, 255, 255, 'r', 'e', 'q', '2'}

// ServerBrowseList token
var ServerBrowseList = []uint8{255, 255, 255, 255, 'l', 'i', 's', '2'}

// MaxServersPerPacket 一个包最多能带有多少个服务器地址
const MaxServersPerPacket = 75

// DataOffset 在数据包不是Extended类型时，数据包前面0xff的个数，这里用的不是Extended类型
const DataOffset = 6

// MaxPackets 最多发送多少个响应包，也许忽略这个限制也没关系？
const MaxPackets = 16

// 传输给teeworlds master server的字节报文，用来获取服务器列表信息
const (
	// ipv6列表
	PacketGetList1 = "\x20\x00\x00\x00\x00\x00\xff\xff\xff\xffreqt"
	// ipv4列表
	PacketGetList2 = "\x20\x00\x00\x00\x00\x00\xff\xff\xff\xffreq2"
)

// PostToken 用于验证POST请求
var PostToken = flag.String("PostToken", "", "-PostToken <your token>")

// ListenURL 用于rest的监听端口
var ListenURL = flag.String("ListenURL", "0.0.0.0:18080", "-ListenURL <your url>")

// ProxyURL 把在heroku上部署的应用作为反向代理，转发请求到实际服务器。
var ProxyURL = flag.String("ProxyURL", "", "-ProxyURL <your url>")
