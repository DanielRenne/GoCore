package main

import (
	_ "core/app"
	"core/serverSettings"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://" + serverSettings.WebConfig.Application.Domain + ":" + strconv.Itoa(serverSettings.WebConfig.Application.HttpsPort) + r.RequestURI, http.StatusMovedPermanently)
}

type WebSocketAPIObj struct {
	
	Data       struct {
		//ServerPropertyID int    `json:"ServerPropertyId"`
		Controller string `json:"controller"`
		Method     string `json:"method"`
		CallBackID int 	`json:"callBackId"`
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

	http.HandleFunc("/websocket", webSocketHandler)

	http.Handle("/", http.FileServer(http.Dir(".")))

	// http.Handle("/websocket", websocket.Handler(WebSocketServer))
	go http.ListenAndServeTLS(":" + strconv.Itoa(serverSettings.WebConfig.Application.HttpsPort), "keys/cert.pem", "keys/key.pem", nil)
	// Start the HTTP server and redirect all incoming connections to HTTPS
	http.ListenAndServe(":" + strconv.Itoa(serverSettings.WebConfig.Application.HttpPort) , http.HandlerFunc(redirectToHttps))
}
