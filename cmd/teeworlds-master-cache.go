package main

import (
	"flag"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/QingGo/teeworlds-master-cache/cache"
	"github.com/QingGo/teeworlds-master-cache/handler"
	"github.com/QingGo/teeworlds-master-cache/udpserver"
)

// CORSMiddleware 解决跨域问题
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	gin.SetMode(gin.ReleaseMode)
	flag.Parse()
	cache.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "18080"
	}

	udpServer, err := udpserver.NewUDPServer("0.0.0.0", 8300)
	if err != nil {
		log.Fatalf("初始化udp服务端失败：%s", err)
	}
	go udpServer.Run()

	r := gin.Default()
	r.Use(CORSMiddleware())
	r.GET("/ping", handler.Ping)
	group := r.Group("/api/v1")
	group.GET("server_list", handler.GetAddrList)
	group.POST("server_list", handler.PostAddrList)
	group.PUT("server_list", handler.PutAddr)
	group.DELETE("server_list", handler.DeleteAddr)

	log.Infof("启动http服务器：%s", ":"+port)
	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("启动服务器失败：%s", err)
	}
}
