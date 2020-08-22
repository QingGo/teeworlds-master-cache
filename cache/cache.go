package cache

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/QingGo/teeworlds-master-cache/datatype"
	"github.com/QingGo/teeworlds-master-cache/myconst"
)

// ServerAddrList 用来缓存可用的ip列表
var ServerAddrList []datatype.ServerAddr

// Init 初始化缓存模块
func Init() {
	ServerAddrList = []datatype.ServerAddr{
		{
			IP:   "164.132.46.180",
			Port: 8304,
		},
	}
}

func getServerInfo(masterURL string) []byte {
	conn, err := net.Dial("udp4", masterURL)
	if err != nil {
		log.Warnf("尝试从master获取信息失败：%s\n", err)
	}
	defer conn.Close()
	timeout, _ := time.ParseDuration("10s")
	conn.SetReadDeadline(time.Now().Add(timeout))
	_, err = conn.Write([]byte(myconst.PacketGetList2))
	if err != nil {
		log.Warnf("尝试从master获取信息失败：%s\n", err)
	}
	var buf [20480]byte
	n, err := conn.Read(buf[0:])
	inforaw := buf[0:n]
	if err != nil {
		log.Warnf("尝试从master获取信息失败：%s\n", err)
	}
	return inforaw
}
