package parser

import (
	"encoding/binary"
	"net"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/QingGo/teeworlds-master-cache/datatype"
	"github.com/QingGo/teeworlds-master-cache/myconst"
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
		// log.Debugf("%s:%d", ip.String(), port)
		iplist = append(iplist, datatype.ServerAddr{
			IP:   ip.String(),
			Port: port})
	}
	return iplist
}

// ParseIPListToBytes 把ip端口列表解析成此cache服务器返回的信息
func ParseIPListToBytes(iplist []datatype.ServerAddr) (responseList [][]byte) {
	responseBytesHeader := make([]byte, myconst.DataOffset+len(myconst.ServerBrowseList))
	for i := 0; i < myconst.DataOffset; i++ {
		responseBytesHeader[i] = 0XFF
	}
	for i, singleByte := range myconst.ServerBrowseList {
		responseBytesHeader[i+myconst.DataOffset] = singleByte
	}

	responseList = make([][]byte, 0)
	ipCounter := 0
	for _, ipport := range iplist {
		if ipCounter%myconst.MaxServersPerPacket == 0 {
			// 更新引用
			responseBytes := make([]byte, len(responseBytesHeader))
			copy(responseBytes, responseBytesHeader)
			// 提前放进去，存的只是slice的长度和背后的数组，后面直接append responseBytes不会影响这里
			responseList = append(responseList, responseBytes)
		}
		ipCounter++
		// 这里只考虑ipv4
		responseList[len(responseList)-1] = append(responseList[len(responseList)-1], 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255)
		for _, ipNumString := range strings.Split(ipport.IP, ".") {
			ipNumInt, err := strconv.Atoi(ipNumString)
			if err != nil {
				log.Debugf("ip解析错误：%s", err)
			}
			responseList[len(responseList)-1] = append(responseList[len(responseList)-1], byte(ipNumInt))
		}
		portBytes := make([]byte, 2)
		portuInt16 := uint16(ipport.Port)
		binary.BigEndian.PutUint16(portBytes, portuInt16)
		responseList[len(responseList)-1] = append(responseList[len(responseList)-1], portBytes...)
	}
	return responseList
}
