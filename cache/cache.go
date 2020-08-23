package cache

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
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

	// 每隔一段时间拉取一次服务器列表
	go func() {
		t := time.NewTicker(10 * time.Minute)
		defer t.Stop()
		for {
			// 首次运行先拉取一次
			go getServerInfoFromRestAPI()
			<-t.C
		}
	}()
}

var serverInfoRestAPI = "https://api.status.tw/2.0/server/list"
var myClient = &http.Client{Timeout: 30 * time.Second}

func getServerInfoFromRestAPI() {
	// 防止程序挂掉
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("未知的panic：%s", r)
		}
	}()

	r, err := myClient.Get(serverInfoRestAPI)
	if err != nil {
		log.Warnf("从status.tw获取服务器信息失败：%s", err)
		return
	}
	defer r.Body.Close()

	var rspStruct datatype.ServerInfoRestAPIResponse
	err = json.NewDecoder(r.Body).Decode(&rspStruct)
	if err != nil {
		log.Warnf("解析json响应失败：%s", err)
		return
	}
	// 合并新获取的结果和原来的结果，而不是覆盖
	RWLock.Lock()
	recordMap := make(map[string]bool)
	for _, serverAddr := range ServerAddrList {
		recordMap[fmt.Sprintf("%s:%d", serverAddr.IP, serverAddr.Port)] = true
	}
	for _, newServerAddr := range rspStruct.Servers {
		recordMap[fmt.Sprintf("%s:%s", newServerAddr.ServerIP, newServerAddr.ServerPort)] = true
	}
	newServerAddrList := make([]datatype.ServerAddr, 0)
	for key := range recordMap {
		addrSlice := strings.Split(key, ":")
		if len(addrSlice) == 2 {
			port, err := strconv.ParseInt(addrSlice[1], 10, 32)
			if err != nil {
				log.Warnf("解析ip端口失败：%s", key)
				continue
			}
			serverAddr := datatype.ServerAddr{
				IP:   addrSlice[0],
				Port: int(port),
			}
			newServerAddrList = append(newServerAddrList, serverAddr)
		}
	}
	ServerAddrList = newServerAddrList
	UDPResponseList = parser.ParseIPListToBytes(ServerAddrList)
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
