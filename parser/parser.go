package parser

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/QingGo/teeworlds-master-cache/datatype"
)

// ParseServerInfo 把服务器返回的信息解析成ip端口列表
func ParseServerInfo(inforaw []byte) []datatype.ServerAddr {
	inforaw = inforaw[14:]
	numServers := len(inforaw) / 18
	iplist := make([]datatype.ServerAddr, 0, 10)
	for i := 0; i < numServers; i++ {
		var ip net.IP
		if string(inforaw[i*18:i*18+12]) == "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xff\xff" {
			// ipv4
			ip = net.IP(inforaw[i*18+12 : i*18+16])
		} else {
			ip = net.IP(inforaw[i*18 : i*18+16])
		}
		port := int(inforaw[i*18+16])*256 + int(inforaw[i*18+17])
		fmt.Println(ip.String(), port)
		iplist = append(iplist, datatype.ServerAddr{
			IP:   ip.String(),
			Port: port})
	}
	return iplist
}

// ParseIPListToBytes 把ip端口列表解析成此cache服务器返回的信息
func ParseIPListToBytes(iplist []datatype.ServerAddr, clientIP net.IP) []byte {
	responseBytes := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	responseBytes = append(responseBytes, net.IP{108, 105, 115, 50}...)
	for _, ipport := range iplist {
		ipString := ipport.IP
		responseBytes = append(responseBytes, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255}...)
		for _, ipNumString := range strings.Split(ipString, ".") {
			ipNumInt, err := strconv.Atoi(ipNumString)
			if err != nil {
				log.Debugf("ip解析错误：%s", err)
			}
			responseBytes = append(responseBytes, byte(ipNumInt))
		}
		portInt := ipport.Port
		portBytes := make([]byte, 2)
		portuInt16 := uint16(portInt)
		binary.BigEndian.PutUint16(portBytes, portuInt16)
		responseBytes = append(responseBytes, portBytes...)
	}
	return responseBytes
}
