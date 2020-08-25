package cache

import (
	"encoding/json"
	"math/rand"
	"net"
	"net/http"
	"strconv"
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
			// go getServerInfoFromRestAPI()
			go getServerInfoFromMasterList()
			<-t.C
		}
	}()

	go func() {
		t := time.NewTicker(5 * time.Minute)
		defer t.Stop()
		for {
			// 首次运行先拉取一次
			// go getServerInfoFromRestAPI()
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

	recordMap := make(map[datatype.ServerAddr]bool)

	for _, newServerAddr := range rspStruct.Servers {
		port, err := strconv.ParseInt(newServerAddr.ServerPort, 10, 32)
		if err != nil {
			log.Warnf("解析ip端口失败：%s", newServerAddr.ServerPort)
			continue
		}
		newServerAddr2 := datatype.ServerAddr{
			IP:   newServerAddr.ServerIP,
			Port: int(port),
		}
		recordMap[newServerAddr2] = true
	}
	newServerAddrList := make([]datatype.ServerAddr, 0)
	for key := range recordMap {
		newServerAddrList = append(newServerAddrList, key)
	}

	RWLock.Lock()
	for _, serverAddr := range ServerAddrList {
		recordMap[serverAddr] = true
	}

	ServerAddrList = newServerAddrList
	UDPResponseList = parser.ParseIPListToBytes(ServerAddrList)
	RWLock.Unlock()
}

var masterURLList = []string{
	"master1.teeworlds.com:8300",
	"master2.teeworlds.com:8300",
	"master3.teeworlds.com:8300",
	"master4.teeworlds.com:8300",
}
var packageChan = make(chan []byte, 1000)

func getServerInfoFromMasterList() {
	// 防止程序挂掉
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("未知的panic：%s", r)
		}
	}()

	for _, masterURL := range masterURLList {
		go getServerInfoFromMaster(masterURL)
	}

	recordMap := make(map[datatype.ServerAddr]bool)

	// 处理packageChan接收到的数据包
	for data := range packageChan {
		serverAddrList := parser.ParseServerInfo(data)
		for _, newServerAddr := range serverAddrList {
			recordMap[newServerAddr] = true
		}
	}
	newServerAddrList := make([]datatype.ServerAddr, 0)
	for key := range recordMap {
		newServerAddrList = append(newServerAddrList, key)
	}

	RWLock.Lock()
	for _, serverAddr := range ServerAddrList {
		recordMap[serverAddr] = true
	}
	ServerAddrList = newServerAddrList
	UDPResponseList = parser.ParseIPListToBytes(ServerAddrList)
	RWLock.Unlock()

}

func getServerInfoFromMaster(masterURL string) {
	// 防止程序挂掉
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("未知的panic：%s", r)
		}
	}()

	tryCount := 0
	var clientConn *net.UDPConn
	var err error
	for {
		clientPort := rand.Intn(22801-22223) + 22223
		clientAddr := net.UDPAddr{
			Port: clientPort,
			IP:   net.ParseIP("0.0.0.0"),
		}
		clientConn, err = net.ListenUDP("udp", &clientAddr)
		if err == nil {
			break
		} else if tryCount >= 3 {
			log.Warnf("尝试了%d次分配端口都失败，中断。", tryCount)
			break
		}
		tryCount++
	}
	defer clientConn.Close()

	serverAddr, _ := net.ResolveUDPAddr("udp", masterURL)
	log.Debugf("%+v\n", serverAddr)
	_, err = clientConn.WriteTo([]byte(myconst.PacketGetList2), serverAddr)

	buf := make([]byte, 32768)
	timeout := time.After(time.Second * 30)
	clientConn.SetDeadline(time.Now().Add(time.Second * 30))
	for {
		select {
		case <-timeout:
			log.Debug("等待结束")
			close(packageChan)
			return
		default:
			log.Debug("尝试读取响应")
			n, _, err := clientConn.ReadFromUDP(buf)
			if err != nil {
				log.Warnf("读取响应错误：%s", err.Error())
				break
			}
			if n > 14 {
				// 先复制再传入channel
				inforaw := make([]byte, n)
				copy(inforaw, buf[0:n])
				packageChan <- inforaw
			}
		}

	}

}
