package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func loadIPFromJSON() []ipport {
	var iplist []ipport
	f, err := os.Open("fetched_ip_list.json")
	checkError(err)
	defer f.Close()
	dec := json.NewDecoder(f)
	dec.Decode(&iplist)
	return iplist
}

func parseIPListToBytes(iplist []ipport, clientIP net.IP) []byte {
	responseBytes := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	responseBytes = append(responseBytes, net.IP{108, 105, 115, 50}...)
	for _, ipport := range iplist {
		ipString := ipport.IP
		responseBytes = append(responseBytes, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255}...)
		for _, ipNumString := range strings.Split(ipString, ".") {
			ipNumInt, err := strconv.Atoi(ipNumString)
			checkError(err)
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

func handleClient(conn *net.UDPConn, iplist []ipport) {
	buf := make([]byte, 2048)
	n, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		return
	}
	request := buf[:n]
	if request != nil {
		fmt.Println(request, string(request))
		// fmt.Println([]byte("\x20\x00\x00\x00\x00\x00\xff\xff\xff\xffreq2"))
	}

	if string(request[6:]) == "\xff\xff\xff\xffreq2" {
		fmt.Println("return response:")
		response := parseIPListToBytes(iplist, addr.IP)
		fmt.Println(response)
		conn.WriteToUDP(response, addr)
	}
}

func test3() {
	iplist := loadIPFromJSON()
	clientIP := net.IP{108, 105, 115, 50}
	response := parseIPListToBytes(iplist, clientIP)
	fmt.Println(response)
	exampleInfo := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 108, 105, 115, 50, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 164, 132, 46, 180, 32, 112, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 164, 132, 46, 180, 32, 111}
	fmt.Println(exampleInfo)
	fmt.Println(string(exampleInfo) == string(response))
}

func main() {
	// test3()
	iplist := loadIPFromJSON()
	addr := net.UDPAddr{
		Port: 8300,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	checkError(err)
	for {
		handleClient(conn, iplist)
	}
}
