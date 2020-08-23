package handler

import (
	"net/http"

	"github.com/QingGo/teeworlds-master-cache/cache"
	"github.com/QingGo/teeworlds-master-cache/datatype"
	"github.com/QingGo/teeworlds-master-cache/myconst"
	"github.com/QingGo/teeworlds-master-cache/parser"
	"github.com/gin-gonic/gin"
)

func checkToken(token string) bool {
	return token == *myconst.PostToken
}

// Ping 测试接口是否正常运行
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// GetAddrList 获取服务器列表
func GetAddrList(c *gin.Context) {
	// 因为发送估计耗时较长，因此先加锁把数据取出，再发送
	cache.RWLock.RLock()
	serverAddrList := make([]datatype.ServerAddr, len(cache.ServerAddrList))
	copy(serverAddrList, cache.ServerAddrList)
	cache.RWLock.RUnlock()

	response := datatype.GetAddrListRespone{
		Code:    success,
		Message: codeMessageMap[success],
		Data:    serverAddrList,
	}
	c.JSON(200, response)

}

// PostAddrList 全量更新服务器列表
func PostAddrList(c *gin.Context) {
	var json datatype.PostAddrListRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		MessageResponse(c, http.StatusBadRequest, requestFormatError)
		return
	}
	if !checkToken(json.Token) {
		MessageResponse(c, http.StatusUnauthorized, tokenError)
		return
	}
	cache.RWLock.Lock()
	cache.ServerAddrList = json.Data
	cache.UDPResponseList = parser.ParseIPListToBytes(cache.ServerAddrList)
	cache.RWLock.Unlock()
	MessageResponse(c, http.StatusOK, success)
}

// PutAddr 增加一个服务器
func PutAddr(c *gin.Context) {
	var json datatype.PutAddrRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		MessageResponse(c, http.StatusBadRequest, requestFormatError)
		return
	}
	if !checkToken(json.Token) {
		MessageResponse(c, http.StatusUnauthorized, tokenError)
		return
	}
	isAppend := true
	for _, addr := range cache.ServerAddrList {
		if addr == json.Data {
			isAppend = false
			break
		}
	}
	if isAppend {
		cache.RWLock.Lock()
		cache.ServerAddrList = append(cache.ServerAddrList, json.Data)
		cache.UDPResponseList = parser.ParseIPListToBytes(cache.ServerAddrList)
		cache.RWLock.Unlock()
	}

	MessageResponse(c, http.StatusOK, success)
}

// DeleteAddr 减少一个服务器
func DeleteAddr(c *gin.Context) {
	var json datatype.PutAddrRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		MessageResponse(c, http.StatusBadRequest, requestFormatError)
		return
	}
	if !checkToken(json.Token) {
		MessageResponse(c, http.StatusUnauthorized, tokenError)
		return
	}

	deleteIndex := -1
	for i, addr := range cache.ServerAddrList {
		if addr == json.Data {
			deleteIndex = i
			break
		}
	}
	if deleteIndex != -1 {
		cache.RWLock.Lock()
		cache.ServerAddrList[deleteIndex] = cache.ServerAddrList[len(cache.ServerAddrList)-1]
		cache.ServerAddrList = cache.ServerAddrList[:len(cache.ServerAddrList)-1]
		cache.UDPResponseList = parser.ParseIPListToBytes(cache.ServerAddrList)
		cache.RWLock.Unlock()
	}

	MessageResponse(c, http.StatusOK, success)
}
