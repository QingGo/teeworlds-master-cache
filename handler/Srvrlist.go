package handler

import (
	"github.com/QingGo/teeworlds-master-cache/datatype"
	"github.com/gin-gonic/gin"
)

// Srvrlist 给web端提供服务器列表
func Srvrlist(c *gin.Context) {
	response := []datatype.ServerInfoForWeb{{
		Address:    "49.232.3.102:8303",
		Version:    "12.6.1",
		Servername: "Eki's DDNet Test Server",
		Mapname:    "Just2Easy",
		Gametype:   "ddnet",
		Flags:      0,
		Numplayers: 1,
		Maxplayers: 64,
		Numclients: 1,
		Maxclients: 64,
		Ping:       20,
		Clients: []datatype.ClientInfoForWeb{
			{
				Name:    "Eki",
				Clan:    "",
				Country: 0,
				Score:   9999,
				Player:  true,
			},
		},
	}}
	c.JSON(200, response)
}
