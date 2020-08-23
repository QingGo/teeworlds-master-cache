package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/QingGo/teeworlds-master-cache/cache"
	"github.com/QingGo/teeworlds-master-cache/handler"
	"github.com/QingGo/teeworlds-master-cache/udpserver"
)

func main() {
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{})
	gin.SetMode(gin.ReleaseMode)
	flag.Parse()
	cache.Init()

	udpServer, err := udpserver.NewUDPServer("0.0.0.0", 8300)
	if err != nil {
		log.Fatalf("初始化udp服务端失败：%s", err)
	}
	go udpServer.Run()

	r := gin.Default()
	r.GET("/ping", handler.Ping)
	group := r.Group("/api/v1")
	group.GET("server_list", handler.GetAddrList)
	group.POST("server_list", handler.PostAddrList)
	group.PUT("server_list", handler.PutAddr)
	group.DELETE("server_list", handler.DeleteAddr)
	r.Run("0.0.0.0:8080")
}
