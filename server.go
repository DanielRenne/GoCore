package main

import (
	_ "core/app"
	"core/ginServer"
	"core/serverSettings"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketAPIObj struct {
	Data struct {
		//ServerPropertyID int    `json:"ServerPropertyId"`
		Controller string `json:"controller"`
		Method     string `json:"method"`
		CallBackID int    `json:"callBackId"`
	} `json:"data"`
}

func webSocketHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()

		var wsAPIObj WebSocketAPIObj
		json.Unmarshal(p, &wsAPIObj)

		// Loop over structs and display them.
		// for l := range languages {
		fmt.Printf("controller = %v, method = %v", wsAPIObj.Data.Controller, wsAPIObj.Data.Method)
		fmt.Println()
		// }

		log.Println(string(p[:]))

		if err != nil {
			return
		}
		if err = conn.WriteMessage(messageType, p); err != nil {
			return
		}
	}

}

func main() {

	if serverSettings.WebConfig.Application.ReleaseMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	loadHTMLTemplates(serverSettings.WebConfig.Application.Name)

	ginServer.Router.Static("/web", "web")

	ginServer.Router.GET("/ws", func(c *gin.Context) {
		webSocketHandler(c.Writer, c.Request)
	})

	go ginServer.Router.RunTLS(":"+strconv.Itoa(serverSettings.WebConfig.Application.HttpsPort), "keys/cert.pem", "keys/key.pem")
	ginServer.Router.Run(":" + strconv.Itoa(serverSettings.WebConfig.Application.HttpPort))

	ginServer.Router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://"+serverSettings.WebConfig.Application.Domain+":"+strconv.Itoa(serverSettings.WebConfig.Application.HttpsPort))
	})

}

func loadHTMLTemplates(appName string) {

	if serverSettings.WebConfig.Application.HtmlTemplates.Enabled {

		levels := "/*"
		dirLevel := ""

		switch serverSettings.WebConfig.Application.HtmlTemplates.DirectoryLevels {
		case 0:
			levels = "/*"
			dirLevel = ""
		case 1:
			levels = "/**/*"
			dirLevel = "root/"
		case 2:
			levels = "/**/**/*"
			dirLevel = "root/root/"
		}

		ginServer.Router.LoadHTMLGlob("web/" + appName + "/" + serverSettings.WebConfig.Application.HtmlTemplates.Directory + levels)

		ginServer.Router.GET("", func(c *gin.Context) {
			c.HTML(http.StatusOK, dirLevel+"index.tmpl", gin.H{})
		})
	} else {

		ginServer.Router.GET("", func(c *gin.Context) {
			ginServer.ReadHTMLFile("web/"+appName+"/index.html", c)
		})
	}
}
