package cache

import (
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/QingGo/teeworlds-master-cache/datatype"
	"github.com/QingGo/teeworlds-master-cache/myconst"
	"github.com/QingGo/teeworlds-master-cache/parser"
)

// ServerAddrList 用来缓存可用的ip列表
var ServerAddrList []datatype.ServerAddr

// UDPResponseList 用来缓存应该udp服务器应该返回的ip列表的列表
var UDPResponseList [][]byte

// RWLock 对缓存使用全局的读写锁，场景应该是写少读多，所以每次写ServerAddrList的时候同时构建UDPResponseList
var RWLock sync.RWMutex

// Init 初始化缓存模块
func Init() {
	RWLock.Lock()
	ServerAddrList = []datatype.ServerAddr{
		{
			IP:   "164.132.46.180",
			Port: 8304,
		},
		{
			IP:   "127.0.1.1",
			Port: 8304,
		},
	}
	UDPResponseList = parser.ParseIPListToBytes(ServerAddrList)
	log.Debugf("初始化的ip列表为：%+v", ServerAddrList)
	log.Debugf("初始化的udp响应为：%+v", UDPResponseList)
	RWLock.Unlock()
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
