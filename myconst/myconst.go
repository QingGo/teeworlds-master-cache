package myconst

import "flag"

// ServerBrowseGetList token
var ServerBrowseGetList = string([]uint8{255, 255, 255, 255, 'r', 'e', 'q', '2'})

// ServerBrowseList token
var ServerBrowseList = string([]uint8{255, 255, 255, 255, 'l', 'i', 's', '2'})

// 传输给teeworlds master server的字节报文，用来获取服务器列表信息
const (
	// ipv6列表
	PacketGetList1 = "\x20\x00\x00\x00\x00\x00\xff\xff\xff\xffreqt"
	// ipv4列表
	PacketGetList2 = "\x20\x00\x00\x00\x00\x00\xff\xff\xff\xffreq2"
)

// PostToken 用于验证POST请求
var PostToken = flag.String("-PostToken", "66666", "-PostToken <your token>")
