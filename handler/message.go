package handler

import (
	"github.com/QingGo/teeworlds-master-cache/datatype"
	"github.com/gin-gonic/gin"
)

const (
	success            = iota
	unknownError       = 1
	requestFormatError = 10000
	tokenError         = 20000
)

var codeMessageMap = map[int]string{
	success:            "success",
	unknownError:       "unknown error",
	requestFormatError: "request format error",
	tokenError:         "token error",
}

// MessageResponse 返回只有错误码和对应信息的响应
func MessageResponse(c *gin.Context, httpcode int, errorCode int) {
	response := datatype.ChangeRespone{
		Code:    errorCode,
		Message: codeMessageMap[errorCode],
	}
	c.JSON(httpcode, response)
}
