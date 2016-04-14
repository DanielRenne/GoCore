package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"

	//Application --------Change Below for the application you want to run-------
	_ "helloWorld"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
	// Redirect the incoming HTTP request. Note that "127.0.0.1:8081" will only work if you are accessing the server from your local machine.
	http.Redirect(w, r, "https://127.0.0.1:443"+r.RequestURI, http.StatusMovedPermanently)
}

type WebSocketAPIObj struct {
	CallBackID int `json:"callBackId"`
	Data       struct {
		//ServerPropertyID int    `json:"ServerPropertyId"`
		Controller string `json:"controller"`
		Method     string `json:"method"`
	} `json:"data"`
	Token string `json:"token"`
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
	go http.ListenAndServeTLS(":443", "keys/cert.pem", "keys/key.pem", nil)
	// Start the HTTP server and redirect all incoming connections to HTTPS
	http.ListenAndServe(":80", http.HandlerFunc(redirectToHttps))
}
