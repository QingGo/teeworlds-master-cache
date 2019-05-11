package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"
)

// 传输给teeworlds的字节报文，用来获取服务器列表信息
const (
	// ipv6列表
	PacketGetList1 = "\x20\x00\x00\x00\x00\x00\xff\xff\xff\xffreqt"
	// ipv4列表
	PacketGetList2 = "\x20\x00\x00\x00\x00\x00\xff\xff\xff\xffreq2"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

type ipport struct {
	// 这里不用大写的变量在json序列化时无法导出
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

// GetListClient 用来获取服务器列表的类
type GetListClient struct {
}

func (p *GetListClient) getServerInfo(masterURL string) []byte {
	conn, err := net.Dial("udp4", masterURL)
	checkError(err)
	defer conn.Close()
	timeout, _ := time.ParseDuration("10s")
	conn.SetReadDeadline(time.Now().Add(timeout))
	_, err = conn.Write([]byte(PacketGetList2))
	checkError(err)
	var buf [20480]byte
	n, err := conn.Read(buf[0:])
	inforaw := buf[0:n]
	checkError(err)
	return inforaw
}

func (p *GetListClient) parseServerInfo(inforaw []byte) []ipport {
	inforaw = inforaw[14:]
	numServers := len(inforaw) / 18
	iplist := make([]ipport, 0, 10)
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
		iplist = append(iplist, ipport{ip.String(), port})
	}
	return iplist
}

func (p *GetListClient) savetojsonfile(iplist []ipport) {
	iplistjson, err := json.Marshal(iplist)
	checkError(err)
	err = ioutil.WriteFile("fetched_ip_list.json", iplistjson, 0644)
	checkError(err)
}

func test1() {
	s := GetListClient{}
	exampleInfo := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 108, 105, 115, 50, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 164, 132, 46, 180, 32, 112, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 164, 132, 46, 180, 32, 111}
	iplist := s.parseServerInfo(exampleInfo)
	fmt.Println(iplist)
	s.savetojsonfile(iplist)
}

func test2() {
	s := GetListClient{}
	rawInfo := s.getServerInfo("31.186.251.128:8300")
	fmt.Println(rawInfo)
	iplist := s.parseServerInfo(rawInfo)
	fmt.Println(rawInfo)
	fmt.Println(len(iplist))
	s.savetojsonfile(iplist)
}

// func main() {
// 	test1()
// }
