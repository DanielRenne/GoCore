package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/appGen"
	"github.com/DanielRenne/GoCore/core/ginServer"
	"github.com/DanielRenne/GoCore/core/serverSettings"
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

	if serverSettings.WebConfig.Application.WebServiceOnly == false {

		loadHTMLTemplates()

		ginServer.Router.Static("/web", appGen.APP_LOCATION+"/web")

		ginServer.Router.GET("/ws", func(c *gin.Context) {
			webSocketHandler(c.Writer, c.Request)
		})
	}

	initializeStaticRoutes()

	go ginServer.Router.RunTLS(":"+strconv.Itoa(serverSettings.WebConfig.Application.HttpsPort), appGen.APP_LOCATION+"/keys/cert.pem", appGen.APP_LOCATION+"/keys/key.pem")

	ginServer.Router.Run(":" + strconv.Itoa(serverSettings.WebConfig.Application.HttpPort))

	// go ginServer.Router.GET("/", func(c *gin.Context) {
	// 	c.Redirect(http.StatusMovedPermanently, "https://"+serverSettings.WebConfig.Application.Domain+":"+strconv.Itoa(serverSettings.WebConfig.Application.HttpsPort))
	// })

}

func loadHTMLTemplates() {

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

		ginServer.Router.LoadHTMLGlob(appGen.APP_LOCATION + "web/" + serverSettings.WebConfig.Application.HtmlTemplates.Directory + levels)

		ginServer.Router.GET("", func(c *gin.Context) {
			c.HTML(http.StatusOK, dirLevel+"index.tmpl", gin.H{})
		})
	} else {

		ginServer.Router.GET("", func(c *gin.Context) {
			ginServer.ReadHTMLFile(appGen.APP_LOCATION+"/web/index.html", c)
		})
	}
}

func initializeStaticRoutes() {

	ginServer.Router.GET("/swagger", func(c *gin.Context) {
		// c.Redirect(http.StatusMovedPermanently, "https://"+serverSettings.WebConfig.Application.Domain+":"+strconv.Itoa(serverSettings.WebConfig.Application.HttpsPort)+"/web/swagger/dist/index.html")

		ginServer.ReadHTMLFile(appGen.APP_LOCATION+"/web/swagger/dist/index.html", c)
	})
}
