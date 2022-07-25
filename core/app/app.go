package app

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	randMath "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/DanielRenne/GoCore/core/atomicTypes"
	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/fileCache"
	"github.com/DanielRenne/GoCore/core/ginServer"
	"github.com/DanielRenne/GoCore/core/gitWebHooks"
	"github.com/DanielRenne/GoCore/core/logger"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/store"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gopkg.in/go-playground/webhooks.v5/github"
)

type WebSocketRemoval func(info WebSocketConnectionMeta)
type customLog func(desc string, message string)

var BroadcastSockets bool
var CustomLog customLog
var totalWriteFailures sync.Map
var mutex sync.RWMutex

type WebSocketConnection struct {
	sync.RWMutex
	Id                   string
	Connection           *websocket.Conn
	Req                  *http.Request
	Context              interface{}
	ContextString        string
	ContextType          string
	ContextLock          sync.RWMutex
	WriteLock            sync.RWMutex
	LastResponseTime     time.Time
	LastResponseTimeLock sync.RWMutex

	GinContextSync GinContextSync
}

type GinContextSync struct {
	sync.RWMutex
	Initialized atomicTypes.AtomicBool
	Context     *gin.Context
}

type WebSocketConnectionMeta struct {
	Conn             *WebSocketConnection
	Context          interface{}
	ContextString    string
	ContextType      string
	LastResponseTime atomicTypes.AtomicTime
	TimeoutOverride  atomicTypes.AtomicInt
}

func init() {
	RegisterWebSocketDataCallback(handleWebSocketData)
	RegisterSecondaryWebSocketDataCallback(handleWebSocketSecondaryData)
	BroadcastSockets = true

	go func() {
		for {
			go func() {
				defer func() {
					if recover := recover(); recover != nil {
						log.Println("Panic Recovered at ReplyToWebSocketSynchronous():  ", recover)
						return
					}
				}()
				start := time.Now()
				totalPrimary := 0
				totalPrimaryDeleted := 0
				totalSecondary := 0
				totalSecondaryDeleted := 0
				GetAllWebSocketMeta().Range(func(key interface{}, value interface{}) bool {
					totalPrimary++
					meta, ok := value.(*WebSocketConnectionMeta)
					if ok {
						err := ReplyToWebSocketSynchronous(meta.Conn, []byte("."))
						if err != nil {
							val, ok := totalWriteFailures.Load(meta.Conn.Id)
							total := 0
							if ok {
								intVal, ok := val.(int)
								if ok {
									total = intVal + 1
									totalWriteFailures.Store(meta.Conn.Id, total)
								}
							} else {
								total = 1
								totalWriteFailures.Store(meta.Conn.Id, total)
							}
							if total == 5 {
								totalWriteFailures.Store(meta.Conn.Id, 0)
								totalPrimaryDeleted++
								deleteWebSocket(meta.Conn, false, "Write single byte of '.' failed 5 times: "+err.Error())
							}
						}
					}
					return true
				})
				GetAllSecondaryWebSocketMeta().Range(func(key interface{}, value interface{}) bool {
					totalSecondary++
					meta, ok := value.(*WebSocketConnectionMeta)
					if ok {
						err := ReplyToWebSocketSynchronous(meta.Conn, []byte("."))
						if err != nil {
							val, ok := totalWriteFailures.Load(meta.Conn.Id)
							total := 0
							if ok {
								intVal, ok := val.(int)
								if ok {
									total = intVal + 1
									totalWriteFailures.Store(meta.Conn.Id, total)
								}
							} else {
								total = 1
								totalWriteFailures.Store(meta.Conn.Id, total)
							}
							if total == 5 {
								totalSecondaryDeleted++
								deleteWebSocket(meta.Conn, true, "Write single byte of '.' failed 5 times: "+err.Error())
							}
						}
					}
					return true
				})
				t := time.Since(start)
				report := " took : " + extensions.Int64ToString(t.Milliseconds()) + "ms (Primary Total " + extensions.IntToString(totalPrimary) + " - deleted " + extensions.IntToString(totalPrimaryDeleted) + ") (Secondary Total " + extensions.IntToString(totalSecondary) + " - deleted " + extensions.IntToString(totalSecondaryDeleted) + ") "
				if CustomLog != nil {
					CustomLog("app->webSocketSendByte", report)
				}
				log.Println("app->webSocketSendByte: " + report)
			}()
			time.Sleep(time.Second * 10)
		}
	}()
}

type WebSocketUpdateIdPayload struct {
	GoCoreEvent     string `json:"GoCoreEvent"`
	NewID           string `json:"NewID"`
	ResponseMessage string `json:"ResponseMessage"`
	ResponseErrors  bool   `json:"ResponseErrors"`
}

type WebSocketEventPayload struct {
	GoCoreEvent    string `json:"GoCoreEvent"`
	ResponseErrors bool   `json:"ResponseErrors"`
}

func handleWebSocketData(conn *WebSocketConnection, c *gin.Context, messageType int, id string, data []byte) {
	handleWebSocketBase(conn, c, messageType, id, data, false)
}

func handleWebSocketSecondaryData(conn *WebSocketConnection, c *gin.Context, messageType int, id string, data []byte) {
	handleWebSocketBase(conn, c, messageType, id, data, true)
}

func handleWebSocketBase(conn *WebSocketConnection, c *gin.Context, messageType int, id string, data []byte, secondary bool) {
	if strings.Index(string(data), "GoCoreEvent") == -1 {
		// bail on any payload without a GoCoreEvent
		return
	}
	var request WebSocketEventPayload
	errMarshal := json.Unmarshal(data, &request)
	if errMarshal != nil {
		request.GoCoreEvent = "Could not unmarshal JSON passed to websocket"
		request.ResponseErrors = true
		ReplyToWebSocketJSON(conn, request)
		return
	}
	if request.GoCoreEvent == "UpdateId" {
		// Websockets can update to IDs that the user cares about for easy identification to write to the particular socket you want.
		// Just Write to the socket a WebSocketUpdateIdPayload json instance with an Event == "UpdateId"
		var request WebSocketUpdateIdPayload
		errMarshal := json.Unmarshal(data, &request)
		if errMarshal != nil {
			request.GoCoreEvent = "UpdateId"
			request.ResponseErrors = true
			ReplyToWebSocketJSON(conn, request)
			return
		}

		var socket *WebSocketConnectionMeta
		var ok bool
		if !secondary {
			socket, ok = GetWebSocketMeta(id)
		} else {
			socket, ok = GetSecondaryWebSocketMeta(id)
		}
		if ok {
			socket.Conn.Id = request.NewID
			if !secondary {
				RemoveWebSocketMeta(id)
				SetWebSocketMeta(socket.Conn.Id, socket)
			} else {
				RemoveSecondaryWebSocketMeta(id)
				SetSecondaryWebSocketMeta(socket.Conn.Id, socket)
			}
			request.ResponseErrors = false
			ReplyToWebSocketJSON(conn, request)
		} else {
			request.ResponseMessage = "Could not find original ID (" + id + ") in memory of websocketMeta"
			request.ResponseErrors = true
			ReplyToWebSocketJSON(conn, request)
		}
	}
}

func (obj *WebSocketConnectionMeta) SetTimeoutOverride(timeout int) {
	obj.TimeoutOverride.Set(timeout)
}

func (obj *WebSocketConnectionMeta) GetConnection() (conn *WebSocketConnection) {
	conn = obj.Conn
	return
}

type WebSocketPubSubPayload struct {
	Key     string      `json:"Key"`
	Content interface{} `json:"Content"`
}

type WebSocketCallback func(conn *WebSocketConnection, c *gin.Context, messageType int, id string, data []byte)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var WebSocketConnections sync.Map
var WebSocketSecondaryConnections sync.Map
var webSocketConnectionsMeta sync.Map
var webSocketSecondaryConnectionsMeta sync.Map
var webSocketCallbacks sync.Map
var webSocketSecondaryCallbacks sync.Map
var WebSocketRemovalCallback WebSocketRemoval
var WebSocketSecondaryRemovalCallback WebSocketRemoval

func Initialize(path string, config string) (err error) {
	err = serverSettings.Initialize(path, config)
	if err != nil {
		return
	}

	serverSettings.WebConfigMutex.RLock()
	inRelease := serverSettings.WebConfig.Application.ReleaseMode == "release"
	serverSettings.WebConfigMutex.RUnlock()

	if inRelease {
		ginServer.Initialize(gin.ReleaseMode, serverSettings.WebConfig.Application.CookieDomain)
	} else {
		ginServer.Initialize(gin.DebugMode, serverSettings.WebConfig.Application.CookieDomain)
	}
	fileCache.Initialize()

	err = dbServices.Initialize()
	if err != nil {
		return
	}
	return
}

func InitializeLite(secureHeaders bool, allowedHosts []string) (err error) {
	ginServer.InitializeLite(gin.ReleaseMode, secureHeaders, allowedHosts)
	fileCache.Initialize()
	return
}

func RunLite(port int) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at RunLite():  ", r)
			time.Sleep(time.Millisecond * 3000)
			RunLite(port)
			return
		}
	}()

	ginServer.Router.GET("/ws", func(c *gin.Context) {
		webSocketHandler(c.Writer, c.Request, c, false)
	})

	secondaryPath := serverSettings.WebConfig.Application.SecondaryWebsocketPath
	if secondaryPath != "" {
		ginServer.Router.GET(secondaryPath, func(c *gin.Context) {
			webSocketHandler(c.Writer, c.Request, c, true)
		})
	}

	log.Println("GoCore Application Started")

	s := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      ginServer.Router,
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
	}
	s.ListenAndServe()

}

func Run() {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at Run():  ", r)
			time.Sleep(time.Millisecond * 3000)
			Run()
			return
		}
	}()

	if serverSettings.WebConfig.Application.MountGitWebHooks == true {
		hook, _ := github.New(github.Options.Secret(serverSettings.WebConfig.Application.GitWebHookSecretKey))
		http.HandleFunc(serverSettings.WebConfig.Application.GitWebHookPath, func(w http.ResponseWriter, r *http.Request) {

			// only these git hooks are supported right now to pass parsed github info to you
			payload, err := hook.Parse(r, github.PushEvent, github.IssuesEvent, github.IssueCommentEvent, github.CreateEvent, github.DeleteEvent, github.ProjectCardEvent, github.ProjectColumnEvent, github.ProjectEvent)
			if err != nil {
				if err == github.ErrEventNotFound {
					// ok event wasn;t one of the ones asked to be parsed
				}
			}
			switch payload.(type) {
			case github.ProjectCardPayload:
				info := payload.(github.ProjectCardPayload)
				gitWebHooks.RunEvent(gitWebHooks.PROJECT_CARD, info)
			case github.ProjectColumnPayload:
				info := payload.(github.ProjectColumnPayload)
				gitWebHooks.RunEvent(gitWebHooks.PROJECT_COLUMN, info)
			case github.ProjectPayload:
				info := payload.(github.ProjectPayload)
				gitWebHooks.RunEvent(gitWebHooks.PROJECT, info)
			case github.IssuesPayload:
				info := payload.(github.IssuesPayload)
				gitWebHooks.RunEvent(gitWebHooks.ISSUES, info)
			case github.IssueCommentPayload:
				info := payload.(github.IssueCommentPayload)
				gitWebHooks.RunEvent(gitWebHooks.ISSUE_COMMENT, info)
			case github.PushPayload:
				info := payload.(github.PushPayload)
				gitWebHooks.RunEvent(gitWebHooks.PUSH_TYPE, info)
			}
		})
		port := "12345"
		if serverSettings.WebConfig.Application.GitWebHookPort != "" {
			port = serverSettings.WebConfig.Application.GitWebHookPort
		}
		go http.ListenAndServe(":"+port, nil)
	}

	if serverSettings.WebConfig.Application.WebServiceOnly == false {

		loadHTMLTemplates()

		ginServer.Router.Static("/web", serverSettings.APP_LOCATION+"/web")

		ginServer.Router.GET("/ws", func(c *gin.Context) {
			webSocketHandler(c.Writer, c.Request, c, false)
		})

		secondaryPath := serverSettings.WebConfig.Application.SecondaryWebsocketPath
		if secondaryPath != "" {
			ginServer.Router.GET(secondaryPath, func(c *gin.Context) {
				webSocketHandler(c.Writer, c.Request, c, true)
			})
		}
	}

	initializeStaticRoutes()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Panic Recovered at Run TLS():  ", r)
				return
			}
		}()

		s := &http.Server{
			Addr: ":" + strconv.Itoa(serverSettings.WebConfig.Application.HttpsPort),
			TLSConfig: &tls.Config{
				PreferServerCipherSuites: true,
				CurvePreferences: []tls.CurveID{
					tls.CurveP256,
					tls.X25519,
				},
				MinVersion: tls.VersionTLS12,
			},
			Handler:      ginServer.Router,
			ReadTimeout:  900 * time.Second,
			WriteTimeout: 300 * time.Second,
		}

		err := s.ListenAndServeTLS(serverSettings.APP_LOCATION+"/keys/cert.pem", serverSettings.APP_LOCATION+"/keys/key.pem")
		if err != nil {
			log.Println("Application failed to ListenAndServeTLS:  " + err.Error())
		} else {
			log.Println("Application Listening on TLS port " + strconv.Itoa(serverSettings.WebConfig.Application.HttpsPort))
		}

	}()

	log.Println("GoCore Application Started")

	port := strconv.Itoa(serverSettings.WebConfig.Application.HttpPort)
	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
	}

	log.Println("Application Listening on port " + port)

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      ginServer.Router,
		ReadTimeout:  900 * time.Second,
		WriteTimeout: 300 * time.Second,
	}
	s.ListenAndServe()

	// ginServer.Router.Run(":" + strconv.Itoa(serverSettings.WebConfig.Application.HttpPort))

	// go ginServer.Router.GET("/", func(c *gin.Context) {
	// 	c.Redirect(http.StatusMovedPermanently, "https://"+serverSettings.WebConfig.Application.Domain+":"+strconv.Itoa(serverSettings.WebConfig.Application.HttpsPort))
	// })

}

func webSocketHandler(w http.ResponseWriter, r *http.Request, c *gin.Context, secondary bool) {

	// return
	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at webSocketHandler():  ", recover)
			time.Sleep(time.Millisecond * 3000)
			webSocketHandler(w, r, c, secondary)
			return
		}
	}()

	if serverSettings.WebConfig.Application.AllowCrossOriginRequests {
		r.Header.Add("Access-Control-Allow-Origin", "*")
	}

	//log.Println("Web Socket Connection")
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		if CustomLog != nil {
			CustomLog("app->webSocketHandler", "Failed to upgrade http connection to websocket:  "+err.Error())
		}
		log.Println("Failed to upgrade http connection to websocket:  " + err.Error())
		return
	}

	//Start the Reader, listen for Close Message, and Add to the Connection Array.

	wsConn := new(WebSocketConnection)
	wsConn.Connection = conn
	wsConn.Req = r
	wsConn.GinContextSync.Context = c
	uuid, err := newUUID()
	if err == nil {
		wsConn.Id = uuid
	} else {
		uuid = randomString(20)
		wsConn.Id = uuid
	}

	socketMeta := new(WebSocketConnectionMeta)
	socketMeta.Conn = wsConn
	socketMeta.LastResponseTime.Set(time.Now())
	if !secondary {
		SetWebSocketMeta(uuid, socketMeta)
	} else {
		SetSecondaryWebSocketMeta(uuid, socketMeta)
	}

	logMsg := "Added Web Socket Connection from " + wsConn.Connection.RemoteAddr().String()
	if CustomLog != nil {
		CustomLog("app->webSocketHandler", logMsg)
	}
	log.Println(logMsg)

	//Reader
	go logger.GoRoutineLogger(func() {

		defer func() {
			if recover := recover(); recover != nil {
				log.Println("Panic Recovered at webSocketHandler-> Reader():  ", recover)
			}
		}()

		for {
			messageType, p, err := conn.ReadMessage()
			if err == nil {

				go func() {
					defer func() {
						if recover := recover(); recover != nil {
							log.Println("Panic Recovered at webSocketHandler-> Reader-> item.Callback():  ", recover)
						}
					}()

					precheckString := string(p)

					//Temp change for crashing web service.
					if strings.Contains(precheckString, "ProxyGateway") || strings.Contains(precheckString, "ProxyWebSocket") {
						return
					}
					var meta *WebSocketConnectionMeta
					var ok bool
					var callbacks sync.Map
					if !secondary {
						meta, ok = GetWebSocketMeta(uuid)
						callbacks = webSocketCallbacks
					} else {
						meta, ok = GetSecondaryWebSocketMeta(uuid)
						callbacks = webSocketSecondaryCallbacks
					}

					if ok {
						meta.LastResponseTime.Set(time.Now())
						if !secondary {
							SetWebSocketMeta(uuid, socketMeta)
						} else {
							SetSecondaryWebSocketMeta(uuid, socketMeta)
						}
					}

					callbacks.Range(func(key interface{}, value interface{}) bool {
						callback, parsed := value.(WebSocketCallback)
						if parsed {
							// if strings.Contains(meta.ContextString, "{\"Page\"") {
							// 	CustomLog("Websocket Request", string(p))
							// }
							callback(wsConn, c, messageType, wsConn.Id, p)
						}
						return true
					})

				}()

			} else {
				if CustomLog != nil {
					CustomLog("app->deleteWebSocket", "Deleting Web Socket "+wsConn.Id+" from read Timeout:  "+err.Error()+":  "+wsConn.Connection.RemoteAddr().String())
				}
				deleteWebSocket(wsConn, secondary, "Reader returned error: "+err.Error())
				return
			}
		}
	}, "GoCore/app.go->webSocketHandler[Reader]")

	if !secondary {
		WebSocketConnections.Store(wsConn.Id, wsConn)
	} else {
		WebSocketSecondaryConnections.Store(wsConn.Id, wsConn)
	}
}

func CloseAllSockets() {
	closeAllSockets(false)
}

func CloseAllSecondarySockets() {
	closeAllSockets(true)
}

func closeAllSockets(secondary bool) {

	items := []*WebSocketConnection{}
	var websockets sync.Map
	mutex.Lock()
	if !secondary {
		websockets = WebSocketConnections
	} else {
		websockets = WebSocketSecondaryConnections
	}
	websockets.Range(func(key interface{}, value interface{}) bool {
		conn, _ := value.(*WebSocketConnection)
		items = append(items, conn)
		return true
	})

	for i := range items {
		connection := items[i]
		connection.Connection.Close()
		websockets.Delete(connection.Id)
	}
	mutex.Unlock()

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

		ginServer.Router.LoadHTMLGlob(serverSettings.APP_LOCATION + "/web/" + serverSettings.WebConfig.Application.HtmlTemplates.Directory + levels)

		ginServer.Router.GET("", func(c *gin.Context) {
			c.HTML(http.StatusOK, dirLevel+"index.tmpl", gin.H{})
		})
	} else {

		if serverSettings.WebConfig.Application.DisableRootIndex {
			return
		}

		ginServer.Router.GET("", func(c *gin.Context) {
			if serverSettings.WebConfig.Application.RootIndexPath == "" {
				ginServer.ReadHTMLFile(serverSettings.APP_LOCATION+"/web/index.htm", c)
			} else {
				ginServer.ReadHTMLFile(serverSettings.APP_LOCATION+"/web/"+serverSettings.WebConfig.Application.RootIndexPath, c)
			}
		})
	}
}

func initializeStaticRoutes() {

	ginServer.Router.GET("/swagger", func(c *gin.Context) {
		// c.Redirect(http.StatusMovedPermanently, "https://"+serverSettings.WebConfig.Application.Domain+":"+strconv.Itoa(serverSettings.WebConfig.Application.HttpsPort)+"/web/swagger/dist/index.html")

		ginServer.ReadHTMLFile(serverSettings.APP_LOCATION+"/web/swagger/dist/index.html", c)
	})
}

func RegisterWebSocketDataCallback(callback WebSocketCallback) {
	uuid, _ := extensions.NewUUID()
	webSocketCallbacks.Store(uuid, callback)
}

func RegisterSecondaryWebSocketDataCallback(callback WebSocketCallback) {
	uuid, _ := extensions.NewUUID()
	webSocketSecondaryCallbacks.Store(uuid, callback)
}

func ReplyToWebSocketSynchronous(conn *WebSocketConnection, data []byte) (err error) {
	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at ReplyToWebSocketSynchronous():  ", recover)
			return
		}

	}()

	conn.WriteLock.Lock()
	conn.Connection.SetWriteDeadline(time.Now().Add(time.Duration(10000) * time.Millisecond))
	err = conn.Connection.WriteMessage(websocket.TextMessage, data)
	conn.WriteLock.Unlock()
	return
}

func ReplyToWebSocket(conn *WebSocketConnection, data []byte) {
	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at ReplyToWebSocket():  ", recover)
			return
		}

	}()

	go func() {

		unlocked := false
		defer func() {
			if recover := recover(); recover != nil {
				if !unlocked {
					conn.WriteLock.Unlock()
				}
				CustomLog("app->ReplyToWebSocket", "Panic Recovered at ReplyToWebSocket():  "+fmt.Sprintf("%+v", recover))
			}
		}()
		conn.WriteLock.Lock()
		conn.Connection.WriteMessage(websocket.TextMessage, data)
		conn.WriteLock.Unlock()
		unlocked = true

	}()
}

func ReplyToWebSocketJSON(conn *WebSocketConnection, v interface{}) {

	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at ReplyToWebSocketJSON():  ", recover)
			return
		}
	}()

	if !BroadcastSockets {
		return
	}
	go func() {

		unlocked := false

		defer func() {
			if recover := recover(); recover != nil {
				if !unlocked {
					conn.WriteLock.Unlock()
				}
				CustomLog("app->ReplyToWebSocketJSON", "Panic Recovered at ReplyToWebSocketJSON():  "+fmt.Sprintf("%+v", recover))
			}
		}()

		conn.WriteLock.Lock()
		conn.Connection.SetWriteDeadline(time.Now().Add(time.Duration(10000) * time.Millisecond))
		conn.Connection.WriteJSON(v)
		conn.WriteLock.Unlock()
		unlocked = true
	}()

}

func ReplyToWebSocketPubSub(conn *WebSocketConnection, key string, v interface{}) {
	defer func() {
		if recover := recover(); recover != nil {
		}
	}()

	if !BroadcastSockets {
		return
	}
	var payload WebSocketPubSubPayload
	payload.Key = key
	payload.Content = v

	go func() {

		unlocked := false
		defer func() {
			if recover := recover(); recover != nil {
				if !unlocked {
					conn.WriteLock.Unlock()
				}
				CustomLog("app->ReplyToWebSocketPubSub", "Panic Recovered at ReplyToWebSocketPubSub():  "+fmt.Sprintf("%+v", recover))
			}
		}()

		conn.WriteLock.Lock()
		conn.Connection.SetWriteDeadline(time.Now().Add(time.Duration(10000) * time.Millisecond))
		conn.Connection.WriteJSON(payload)
		conn.WriteLock.Unlock()
		unlocked = true

	}()

}

func BroadcastWebSocketData(data []byte) {
	broadcastWebSocketData(data, false)
}

func BroadcastSecondaryWebSocketData(data []byte) {
	broadcastWebSocketData(data, true)
}

func broadcastWebSocketData(data []byte, secondary bool) {

	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at WebSocketConnections():  ", recover)
			return
		}
	}()

	if !BroadcastSockets {
		return
	}
	mutex.Lock()
	var websockets sync.Map
	if !secondary {
		websockets = WebSocketConnections
	} else {
		websockets = WebSocketSecondaryConnections
	}
	websockets.Range(func(key interface{}, value interface{}) bool {
		conn, _ := value.(*WebSocketConnection)

		go func() {

			unlocked := false
			defer func() {
				if recover := recover(); recover != nil {
					if !unlocked {
						conn.WriteLock.Unlock()
					}
					CustomLog("app->BroadcastWebSocketData", "Panic Recovered at BroadcastWebSocketData():  "+fmt.Sprintf("%+v", recover))
				}
			}()
			conn.WriteLock.Lock()
			conn.Connection.WriteMessage(websocket.BinaryMessage, data)
			conn.WriteLock.Unlock()
			unlocked = true
			// deadLockChan <- 0
		}()
		return true
	})
	mutex.Unlock()
	return
}

func BroadcastWebSocketJSON(v interface{}) {
	broadcastWebSocketJSON(v, false)
}

func BroadcastSecondaryWebSocketJSON(v interface{}) {
	broadcastWebSocketJSON(v, true)
}

func broadcastWebSocketJSON(v interface{}, secondary bool) {
	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at BroadcastWebSocketJSON():  ", recover)
			return
		}
	}()

	if !BroadcastSockets {
		return
	}

	mutex.Lock()
	var websockets sync.Map
	if !secondary {
		websockets = WebSocketConnections
	} else {
		websockets = WebSocketSecondaryConnections
	}
	websockets.Range(func(key interface{}, value interface{}) bool {
		conn, _ := value.(*WebSocketConnection)

		go func() {

			unlocked := false
			defer func() {
				if recover := recover(); recover != nil {
					if !unlocked {
						conn.WriteLock.Unlock()
					}
					CustomLog("app->BroadcastWebSocketJSON", "Panic Recovered at BroadcastWebSocketJSON():  "+fmt.Sprintf("%+v", recover))
				}
			}()

			conn.WriteLock.Lock()
			conn.Connection.SetWriteDeadline(time.Now().Add(time.Duration(10000) * time.Millisecond))
			conn.Connection.WriteJSON(v)
			conn.WriteLock.Unlock()
			unlocked = true
		}()
		return true
	})
	mutex.Unlock()
}

func PublishWebSocketJSON(key string, v interface{}) {
	publishWebSocketJSON(key, v, false)
}

func PublishSecondaryWebSocketJSON(key string, v interface{}) {
	publishWebSocketJSON(key, v, true)
}

func publishWebSocketJSON(key string, v interface{}, secondary bool) {
	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at PublishWebSocketJSON():  ", recover)
			return
		}
	}()
	if !BroadcastSockets {
		return
	}
	var payload WebSocketPubSubPayload
	payload.Key = key
	payload.Content = v
	var websockets sync.Map
	mutex.Lock()
	if !secondary {
		websockets = WebSocketConnections
	} else {
		websockets = WebSocketSecondaryConnections
	}
	websockets.Range(func(key interface{}, value interface{}) bool {
		conn, _ := value.(*WebSocketConnection)
		go func() {

			unlocked := false
			defer func() {
				if recover := recover(); recover != nil {
					if !unlocked {
						conn.WriteLock.Unlock()
					}
					CustomLog("app->PublishWebSocketJSON", "Panic Recovered at PublishWebSocketJSON():  "+fmt.Sprintf("%+v", recover)+" payload: "+fmt.Sprintf("%+v", payload))
				}
			}()

			conn.WriteLock.Lock()
			conn.Connection.SetWriteDeadline(time.Now().Add(time.Duration(10000) * time.Millisecond))
			conn.Connection.WriteJSON(payload)
			conn.WriteLock.Unlock()
			unlocked = true
			// deadLockChan <- 0
		}()
		return true
	})
	mutex.Unlock()
}

func SetWebSocketTimeout(timeout int) {
	setWebSocketTimeout(timeout, false)
}

func SetSecondaryWebSocketTimeout(timeout int) {
	setWebSocketTimeout(timeout, true)
}

func setWebSocketTimeout(timeout int, secondary bool) {
	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at SetWebSocketTimeout():  ", recover)
			time.Sleep(time.Millisecond * 3000)
			setWebSocketTimeout(timeout, secondary)
			return
		}
	}()

	// if CustomLog != nil {
	// 	CustomLog("app->SetWebSocketTimeout", "Checking for Web Socket Timeouts.")
	// }

	for {
		var websockets sync.Map
		mutex.Lock()
		// even with this mutex, there still is a fatal error: concurrent map iteration and map write in this loop
		if !secondary {
			websockets = webSocketConnectionsMeta
		} else {
			websockets = webSocketSecondaryConnectionsMeta
		}
		websockets.Range(func(key interface{}, value interface{}) bool {
			meta, ok := value.(*WebSocketConnectionMeta)
			if ok {
				duration := time.Millisecond * time.Duration(timeout)
				timeoutOverride := meta.TimeoutOverride.Get()
				if timeoutOverride != 0 {
					duration = time.Millisecond * time.Duration(timeoutOverride)
				}
				if meta.LastResponseTime.Get().Add(duration).Before(time.Now()) {
					if CustomLog != nil {
						CustomLog("app->SetWebSocketTimeout", "Removed Websocket due to timeout from :  "+meta.Conn.Connection.RemoteAddr().String())
					}
					log.Println("Removed Websocket due to timeout from :  " + meta.GetConnection().Connection.RemoteAddr().String())
					deleteWebSocket(meta.GetConnection(), secondary, "Last Response Timeout")
				}
			}
			return true
		})
		mutex.Unlock()
		time.Sleep(time.Second * 10)
	}

}

func deleteWebSocket(c *WebSocketConnection, secondary bool, reason string) {

	go func() {
		secondaryLog := ""
		if secondary {
			secondaryLog = "(secondary) "
		}
		defer func() {
			if recover := recover(); recover != nil {
				CustomLog("app->deleteWebSocket", "Panic Recovered at deleteWebSocket():  "+fmt.Sprintf("%+v", recover))
				return
			}
		}()
		c.Connection.Close()

		logDeletion := "Deleting " + secondaryLog + "Web Socket from client:  " + c.Connection.RemoteAddr().String() + " Reason: " + reason
		if CustomLog != nil {
			CustomLog("app->deleteWebSocket", logDeletion)
		}
		log.Println(logDeletion)

		if !secondary {
			WebSocketConnections.Delete(c.Id)

			if store.OnChange != nil {
				go func() {
					defer func() {
						if recover := recover(); recover != nil {
							log.Println("Panic Recovered at store.OnChange():  ", recover)
							return
						}
					}()

					store.OnChange(store.WebSocketStoreKey, "", store.PathRemove, nil, nil)
				}()
			}

			if WebSocketRemovalCallback != nil {
				info, ok := GetWebSocketMeta(c.Id)
				if ok {
					go func(c *WebSocketConnection) {
						defer func() {
							if recover := recover(); recover != nil {
								log.Println("Panic Recovered at deleteWebSocket():  ", recover)
								return
							}
						}()
						WebSocketRemovalCallback(*info)
					}(c)
				}
			}

			RemoveWebSocketMeta(c.Id)
		} else {
			WebSocketSecondaryConnections.Delete(c.Id)
			RemoveSecondaryWebSocketMeta(c.Id)
		}

	}()

}

func GetWebSocketMeta(id string) (info *WebSocketConnectionMeta, ok bool) {
	result, ok := webSocketConnectionsMeta.Load(id)
	if ok {
		info = result.(*WebSocketConnectionMeta)
		return
	}
	return
}

func SetWebSocketMeta(id string, info *WebSocketConnectionMeta) {
	webSocketConnectionsMeta.Store(id, info)
}

func RemoveWebSocketMeta(id string) {
	webSocketConnectionsMeta.Delete(id)
}

func GetAllWebSocketMeta() (items *sync.Map) {
	return &webSocketConnectionsMeta
}

func GetSecondaryWebSocketMeta(id string) (info *WebSocketConnectionMeta, ok bool) {
	result, ok := webSocketSecondaryConnectionsMeta.Load(id)
	if ok {
		info = result.(*WebSocketConnectionMeta)
		return
	}
	return
}

func SetSecondaryWebSocketMeta(id string, info *WebSocketConnectionMeta) {
	webSocketSecondaryConnectionsMeta.Store(id, info)
}

func RemoveSecondaryWebSocketMeta(id string) {
	webSocketSecondaryConnectionsMeta.Delete(id)
}

func GetAllSecondaryWebSocketMeta() (items *sync.Map) {
	return &webSocketSecondaryConnectionsMeta
}

func removeWebSocket(s []*WebSocketConnection, i int) []*WebSocketConnection {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func randomString(strlen int) string {
	randMath.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[randMath.Intn(len(chars))]
	}
	return string(result)
}
