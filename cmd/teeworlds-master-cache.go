package main

import (
	"flag"
	"math/rand"
	"os"
	"time"

	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/QingGo/teeworlds-master-cache/cache"
	"github.com/QingGo/teeworlds-master-cache/handler"
	"github.com/QingGo/teeworlds-master-cache/myconst"
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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func proxy(c *gin.Context) {
	// https://49.232.3.102:10443/api/v1/server_list
	remote, err := url.Parse(*myconst.ProxyURL)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	//Define the director func
	//This is a good place to log, for example
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = c.Param("proxyPath")
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func main() {
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{})
	gin.SetMode(gin.ReleaseMode)
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = "18080"
	}

	r := gin.Default()

	if *myconst.ProxyURL == "" {
		r.Use(CORSMiddleware())
		if *myconst.PostToken == "" {
			rand.Seed(time.Now().UnixNano())
			_PostToken := randStringRunes(16)
			myconst.PostToken = &_PostToken
			log.Infof("没有指定token，自动生成：%s", _PostToken)
		}
		cache.Init()

		udpServer, err := udpserver.NewUDPServer("0.0.0.0", 8300)
		if err != nil {
			log.Fatalf("初始化udp服务端失败：%s", err)
		}
		go udpServer.Run()

		r.GET("/ping", handler.Ping)
		r.GET("/srvrlist.json", handler.Srvrlist)
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
	} else {
		// 反向代理模式
		//Create a catchall route
		log.Info("使用反向代理模式")
		r.Any("/*proxyPath", proxy)
		r.Run(":" + port)
	}
}
