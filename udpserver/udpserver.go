package udpserver

import (
	"net"
	"reflect"

	log "github.com/sirupsen/logrus"

	"github.com/QingGo/teeworlds-master-cache/cache"
	"github.com/QingGo/teeworlds-master-cache/myconst"
)

// UDPServer 用来响应客户端的ip列表请求
type UDPServer struct {
	listerAddr net.UDPAddr
	conn       *net.UDPConn
}

// NewUDPServer 生成新的UDPServer实例
func NewUDPServer(ip string, port int) (*UDPServer, error) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}
	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		return nil, err
	}
	server := UDPServer{
		listerAddr: addr,
		conn:       conn,
	}
	return &server, nil
}

func (server UDPServer) handleUDPClient() {
	buf := make([]byte, 2048)
	for {
		n, addr, err := server.conn.ReadFromUDP(buf)
		if err != nil {
			return
		}
		request := buf[:n]
		if request != nil {
			log.Debugf("收到请求%s(%s)", request, string(request))
		}
		go server.handle(request, addr)
	}
}

func (server UDPServer) handle(request []byte, addr *net.UDPAddr) {
	// 确认udp请求的token
	if len(request) > myconst.DataOffset && reflect.DeepEqual(request[myconst.DataOffset:myconst.DataOffset+len(myconst.ServerBrowseGetList)], myconst.ServerBrowseGetList) {
		// 因为发送估计耗时较长，因此先加锁把数据取出，再发送
		cache.RWLock.RLock()
		responseList := make([][]byte, len(cache.UDPResponseList))
		for i := range cache.UDPResponseList {
			responseList[i] = make([]byte, len(cache.UDPResponseList[i]))
			copy(responseList[i], cache.UDPResponseList[i])
		}
		cache.RWLock.RUnlock()

		if len(responseList) != 0 {
			log.Debugf("返回响应长度%+v，首记录长度%+v", len(responseList), len(responseList[0]))
			for _, response := range responseList {
				server.conn.WriteToUDP(response, addr)
			}

		}
	}
}

// Run 运行udp cache server
func (server UDPServer) Run() {
	log.Infof("运行cache server：%+v", server.listerAddr)
	server.handleUDPClient()
}
