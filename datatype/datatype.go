package datatype

import "fmt"

// ServerAddr 用来存放服务器的ip和端口信息
type ServerAddr struct {
	// 这里不用大写的变量在json序列化时无法导出
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

func (ipport *ServerAddr) String() string {
	return fmt.Sprintf("%s:%d", ipport.IP, ipport.Port)
}

// GetAddrListRespone 返回获取服务器列表的restful api响应
type GetAddrListRespone struct {
	Code    int
	Message string
	Data    []ServerAddr
}

// ChangeRespone post, put, delete的响应
type ChangeRespone struct {
	Code    int
	Message string
}

// PostAddrListRequest 全量更新服务器列表的restful api请求
type PostAddrListRequest struct {
	Token string
	Data  []ServerAddr
}

// PutAddrRequest 增加一个服务器记录的restful api请求
type PutAddrRequest struct {
	Token string
	Data  ServerAddr
}

// DeleteAddrRequest 减少一个服务器记录的restful api请求
type DeleteAddrRequest struct {
	Token string
	Data  ServerAddr
}

// ServerInfoRestAPIResponse 从status.tw获取的服务器列表响应
type ServerInfoRestAPIResponse struct {
	Servers []struct {
		ServerIP         string        `json:"server_ip"`
		ServerPort       string        `json:"server_port"`
		FirstSeen        string        `json:"first_seen"`
		LastSeen         string        `json:"last_seen"`
		Version          string        `json:"version"`
		Name             string        `json:"name"`
		Password         bool          `json:"password"`
		Ping             int           `json:"ping"`
		ServerLevel      int           `json:"server_level"`
		NumClients       int           `json:"num_clients"`
		MaxClients       int           `json:"max_clients"`
		NumPlayers       int           `json:"num_players"`
		MaxPlayers       int           `json:"max_players"`
		NumBotPlayers    int           `json:"num_bot_players"`
		NumBotSpectators int           `json:"num_bot_spectators"`
		Gamemode         string        `json:"gamemode"`
		Map              string        `json:"map"`
		Master           string        `json:"master"`
		Country          string        `json:"country"`
		IsVerified       bool          `json:"is_verified"`
		Players          []interface{} `json:"players"`
	} `json:"servers"`
}

// {
// 	"address": "49.232.3.102:8303",
// 	"version": "12.6.1",
// 	"servername": "Eki's DDNet Test Server",
// 	"mapname": "Just2Easy",
// 	"gametype": "ddnet",
// 	"flags": 0,
// 	"numplayers": 0,
// 	"maxplayers": 64,
// 	"numclients": 0,
// 	"maxclients": 64,
// 	"ping": 20,
// 	"clients": [
// 		{
// 			"name": "Eki",
// 			"clan": "",
// 			"country": 0,
// 			"score": 9999,
// 			"player": true
// 		}
// 	]
// }

// ServerInfoForWeb 给web客户端提供数据返回
type ServerInfoForWeb struct {
	Address    string             `json:"address"`
	Version    string             `json:"version"`
	Servername string             `json:"servername"`
	Mapname    string             `json:"mapname"`
	Gametype   string             `json:"gametype"`
	Flags      int                `json:"flags"`
	Numplayers int                `json:"numplayers"`
	Maxplayers int                `json:"maxplayers"`
	Numclients int                `json:"numclients"`
	Maxclients int                `json:"maxclients"`
	Ping       int                `json:"ping"`
	Clients    []ClientInfoForWeb `json:"clients"`
}

// ClientInfoForWeb 给web客户端提供数据返回
type ClientInfoForWeb struct {
	Name    string `json:"name"`
	Clan    string `json:"clan"`
	Country int    `json:"country"`
	Score   int    `json:"score"`
	Player  bool   `json:"player"`
}
