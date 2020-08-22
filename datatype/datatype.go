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
