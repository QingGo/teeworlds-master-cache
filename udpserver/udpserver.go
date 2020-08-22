package udpserver

import (
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/QingGo/teeworlds-master-cache/cache"
	"github.com/QingGo/teeworlds-master-cache/parser"
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
	n, addr, err := server.conn.ReadFromUDP(buf)
	if err != nil {
		return
	}
	request := buf[:n]
	if request != nil {
		log.Debugf("收到请求%s(%s)", request, string(request))
	}

	if len(request) > 6 && string(request[6:]) == "\xff\xff\xff\xffreq2" {
		response := parser.ParseIPListToBytes(cache.ServerAddrList, addr.IP)
		if response != nil {
			log.Debugf("返回响应%s(%s)", response, string(response))
			server.conn.WriteToUDP(response, addr)
		}
	}
}

// Run 运行udp cache server
func (server UDPServer) Run() {
	log.Infof("运行cache server：%+v", server.listerAddr)
	go func() {
		for {
			server.handleUDPClient()
		}
	}()
}
