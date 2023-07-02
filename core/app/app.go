// Package app is a package that contains the core functionality of the GoCore framework and websocket functionality
package app

import (
	"crypto/rand"
	"crypto/tls"
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
	"github.com/DanielRenne/GoCore/core/path"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/store"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gopkg.in/go-playground/webhooks.v5/github"
)

// WebSocketRemoval a type which removes a websocket callback
type WebSocketRemoval func(info WebSocketConnectionMeta)

// StaticWebLocation is the location of the static web files (defaults to "/web")
var StaticWebLocation string

type customLog func(desc string, message string)

// BroadcastSockets is a flag that determines if web sockets should broadcast to all clients (defaults to true with init())
var BroadcastSockets bool

// CustomLog allows you to set a custom log function for the web server logs that GoCore outputs
var CustomLog customLog

// PrimaryGoCoreHTTPServer is the primary http server that GoCore uses where if you needed to tear it down you could
var PrimaryGoCoreHTTPServer *http.Server

// WebSocketConnection is a websocket connection
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
	GinContextSync       GinContextSync
}

// GinContextSync is a sync wrapper for a gin context
type GinContextSync struct {
	sync.RWMutex
	Initialized atomicTypes.AtomicBool
	Context     *gin.Context
}

// WebSocketConnectionMeta is the meta data for a websocket connection
type WebSocketConnectionMeta struct {
	Conn             *WebSocketConnection
	Context          interface{}
	ContextString    string
	ContextType      string
	LastResponseTime atomicTypes.AtomicTime
	TimeoutOverride  atomicTypes.AtomicInt
}

func (obj *WebSocketConnectionMeta) SetTimeoutOverride(timeout int) {
	obj.TimeoutOverride.Set(timeout)
}

func (obj *WebSocketConnectionMeta) GetConnection() (conn *WebSocketConnection) {
	conn = obj.Conn
	return
}

// WebSocketConnectionCollection is a collection of websocket connections
type WebSocketConnectionCollection struct {
	sync.RWMutex
	Connections []*WebSocketConnection
}

// ConcurrentWebSocketConnectionItem is a concurrent websocket connection item
type ConcurrentWebSocketConnectionItem struct {
	Index int
	Conn  *WebSocketConnection
}

func (wscc *WebSocketConnectionCollection) Append(item *WebSocketConnection) {
	wscc.Lock()
	defer wscc.Unlock()
	wscc.Connections = append(wscc.Connections, item)

	if store.OnChange != nil {
		go func() {
			defer func() {
				if recover := recover(); recover != nil {
					log.Println("Panic Recovered at store.OnChange():  ", recover)
					return
				}
			}()

			store.OnChange(store.WebSocketStoreKey, "", store.PathAdd, nil, nil)
		}()
	}
}

func (wscc *WebSocketConnectionCollection) Iter() <-chan ConcurrentWebSocketConnectionItem {
	c := make(chan ConcurrentWebSocketConnectionItem)

	f := func() {
		wscc.RLock()
		defer wscc.RUnlock()
		for index := range wscc.Connections {
			value := wscc.Connections[index]
			c <- ConcurrentWebSocketConnectionItem{index, value}
		}
		close(c)
	}
	go f()

	return c
}

// WebSocketCallbackSync is a sync wrapper for a websocket callback
type WebSocketCallbackSync struct {
	sync.RWMutex
	callbacks []WebSocketCallback
}

// ConcurrentWebSocketCallbackItem is a concurrent websocket callback item
type ConcurrentWebSocketCallbackItem struct {
	Index    int
	Callback WebSocketCallback
}

func (obj *WebSocketCallbackSync) Append(item WebSocketCallback) {
	obj.RLock()
	defer obj.RUnlock()
	obj.callbacks = append(obj.callbacks, item)
}

func (obj *WebSocketCallbackSync) Iter() <-chan ConcurrentWebSocketCallbackItem {
	c := make(chan ConcurrentWebSocketCallbackItem)

	f := func() {
		obj.Lock()
		defer obj.Unlock()
		for index := range obj.callbacks {
			value := obj.callbacks[index]
			c <- ConcurrentWebSocketCallbackItem{index, value}
		}
		close(c)
	}
	go f()

	return c
}

// WebSocketPubSubPayload is a websocket pub sub payload
type WebSocketPubSubPayload struct {
	Key     string      `json:"Key"`
	Content interface{} `json:"Content"`
}

// WebSocketCallback is a websocket callback
type WebSocketCallback func(conn *WebSocketConnection, c *gin.Context, messageType int, id string, data []byte)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WebSocketConnections is a collection of websocket connections
var WebSocketConnections sync.Map
var webSocketConnectionsMeta sync.Map

// WebSocketCallbacks is a collection of websocket callbacks
var WebSocketCallbacks sync.Map

// WebSocketRemovalCallback is a collection of websocket removal callbacks
var WebSocketRemovalCallback WebSocketRemoval

func init() {
	StaticWebLocation = "web"
	BroadcastSockets = true
}

// Init initializes the web server with the default webConfig.json file
func Init() {
	Initialize(path.GetBinaryPath(), "webConfig.json")
}

// InitCustomWebCofig initializes the web server with a custom webConfig.json file
func InitCustomWebConfig(webConfig string) {
	Initialize(path.GetBinaryPath(), webConfig)
}

// Initialize initializes the web server with a full path to your proect and the name of your webConfig.json file.
func Initialize(path string, config string) {
	err := serverSettings.Initialize(path, config)
	if err != nil {
		log.Println("Failed to parse or read webConfig.json : " + err.Error())
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
	dbServices.Initialize()
}

// InitializeLite initilizes a basic gin server with no database coupling
func InitializeLite(secureHeaders bool, allowedHosts []string) (err error) {
	ginServer.InitializeLite(gin.ReleaseMode, secureHeaders, allowedHosts)
	// Why do we do this here.  Its not needed
	fileCache.Initialize()
	return
}

// RunLite is a lite version of Run that does not require a webConfig.json file or any serverSettings information
func RunLite(port int) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at RunLite():  ", r)
			time.Sleep(time.Millisecond * 3000)
			RunLite(port)
			return
		}
	}()
	if !serverSettings.WebConfig.Application.DisableWebSockets {
		ginServer.Router.GET("/ws", func(c *gin.Context) {
			webSocketHandler(c.Writer, c.Request, c)
		})
	}

	log.Println("GoCore Application Started")

	PrimaryGoCoreHTTPServer = &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      ginServer.Router,
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
	}
	err := PrimaryGoCoreHTTPServer.ListenAndServe()
	if err != nil {
		log.Println("GoCore Cannot open port " + strconv.Itoa(port) + " Reason: " + err.Error())
	}

}

// Run starts the server (Typically for a full GoCore server)
func Run() {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at Run():  ", r)
			time.Sleep(time.Millisecond * 3000)
			Run()
			return
		}
	}()
	if serverSettings.WebConfig.Application.MountGitWebHooks && serverSettings.WebConfig.Application.GitWebHookPath != "" {
		hook, _ := github.New(github.Options.Secret(serverSettings.WebConfig.Application.GitWebHookSecretKey))
		http.HandleFunc(serverSettings.WebConfig.Application.GitWebHookPath, func(w http.ResponseWriter, r *http.Request) {

			// only these git hooks are supported right now to pass parsed github info to you
			payload, _ := hook.Parse(r, github.PushEvent, github.IssuesEvent, github.IssueCommentEvent, github.CreateEvent, github.DeleteEvent, github.ProjectCardEvent, github.ProjectColumnEvent, github.ProjectEvent)
			switch payload := payload.(type) {
			case github.ProjectCardPayload:
				info := payload
				gitWebHooks.RunEvent(gitWebHooks.PROJECT_CARD, info)
			case github.ProjectColumnPayload:
				info := payload
				gitWebHooks.RunEvent(gitWebHooks.PROJECT_COLUMN, info)
			case github.ProjectPayload:
				info := payload
				gitWebHooks.RunEvent(gitWebHooks.PROJECT, info)
			case github.IssuesPayload:
				info := payload
				gitWebHooks.RunEvent(gitWebHooks.ISSUES, info)
			case github.IssueCommentPayload:
				info := payload
				gitWebHooks.RunEvent(gitWebHooks.ISSUE_COMMENT, info)
			case github.PushPayload:
				info := payload
				gitWebHooks.RunEvent(gitWebHooks.PUSH_TYPE, info)
			}
		})
		port := "12345"
		if serverSettings.WebConfig.Application.GitWebHookPort != "" {
			port = serverSettings.WebConfig.Application.GitWebHookPort
		}
		go func() {
			err := http.ListenAndServe(":"+port, nil)
			if err != nil {
				log.Println("GoCore Cannot open port " + port + " Reason: " + err.Error())
			}
		}()
	}

	if !serverSettings.WebConfig.Application.WebServiceOnly {

		loadHTMLTemplates()

		// Override to blank string if you dont want static files to be mounted
		if StaticWebLocation != "" {
			ginServer.Router.Static("/"+StaticWebLocation, serverSettings.APP_LOCATION+path.PathSeparator+StaticWebLocation)
		}
	}

	if !serverSettings.WebConfig.Application.DisableWebSockets {
		ginServer.Router.GET("/ws", func(c *gin.Context) {
			webSocketHandler(c.Writer, c.Request, c)
		})
	}

	if extensions.DoesFileExist(serverSettings.APP_LOCATION+"/keys/cert.pem") && extensions.DoesFileExist(serverSettings.APP_LOCATION+"/keys/key.pem") {
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
				log.Println("GoCore Application failed to ListenAndServeTLS:  " + err.Error())
			} else {
				log.Println("Application Listening on TLS port " + strconv.Itoa(serverSettings.WebConfig.Application.HttpsPort))
			}

		}()
	}
	log.Println("GoCore Application Started")
	// RunServer Blocking forever!!
	RunServer()
}

// RunServer Blocking forever with ListenAndServe!!
func RunServer() {
	port := strconv.Itoa(serverSettings.WebConfig.Application.HttpPort)
	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
	}

	PrimaryGoCoreHTTPServer = &http.Server{
		Addr:         ":" + port,
		Handler:      ginServer.Router,
		ReadTimeout:  900 * time.Second,
		WriteTimeout: 300 * time.Second,
	}
	err := PrimaryGoCoreHTTPServer.ListenAndServe()
	if err != nil {
		log.Println("GoCore Application failed to listen on port (" + port + "):  " + err.Error())
	} else {
		log.Println("Application Listening on port " + port)
	}
}

func webSocketHandler(w http.ResponseWriter, r *http.Request, c *gin.Context) {

	// return
	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at webSocketHandler():  ", recover)
			time.Sleep(time.Millisecond * 3000)
			webSocketHandler(w, r, c)
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

	SetWebSocketMeta(uuid, socketMeta)

	if CustomLog != nil {
		CustomLog("app->webSocketHandler", "Added Web Socket Connection from "+wsConn.Connection.RemoteAddr().String())
	}

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

					meta, ok := GetWebSocketMeta(uuid)
					if ok {
						meta.LastResponseTime.Set(time.Now())
						SetWebSocketMeta(uuid, meta)
					}

					WebSocketCallbacks.Range(func(key interface{}, value interface{}) bool {
						callback, parsed := value.(WebSocketCallback)
						if parsed {
							// if strings.Contains(meta.ContextString, "{\"Page\"") {
							// 	CustomLog("Websocket Request", string(p))
							// }
							callback(wsConn, c, messageType, uuid, p)
						}
						return true
					})

				}()

			} else {
				if CustomLog != nil {
					CustomLog("app->deleteWebSocket", "Deleting Web Socket from read Timeout:  "+err.Error()+":  "+wsConn.Connection.RemoteAddr().String())
				}
				deleteWebSocket(wsConn)
				return
			}
		}
	}, "GoCore/app.go->webSocketHandler[Reader]")

	WebSocketConnections.Store(wsConn.Id, wsConn)
}

// CloseAllSockets closes all web sockets
func CloseAllSockets() {

	items := []*WebSocketConnection{}

	WebSocketConnections.Range(func(key interface{}, value interface{}) bool {
		conn, _ := value.(*WebSocketConnection)
		items = append(items, conn)
		return true
	})

	for i := range items {
		connection := items[i]
		connection.Connection.UnderlyingConn().Close()
		WebSocketConnections.Delete(connection.Id)
	}

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

// RegisterWebSocketDataCallback registers a callback for a websocket data event
func RegisterWebSocketDataCallback(callback WebSocketCallback) {
	uuid, _ := extensions.NewUUID()
	WebSocketCallbacks.Store(uuid, callback)
}

// ReplyToWebSocket sends a message to a websocket connection
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

// ReplyToWebSocketJSON sends a JSON message to a websocket connection
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

		conn.Connection.SetWriteDeadline(time.Now().Add(time.Duration(10000) * time.Millisecond))
		conn.WriteLock.Lock()
		conn.Connection.WriteJSON(v)
		conn.WriteLock.Unlock()
		unlocked = true
	}()

}

// ReplyToWebSocketPubSub sends a message to all websocket connections
func ReplyToWebSocketPubSub(conn *WebSocketConnection, key string, v interface{}) {
	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at ReplyToWebSocketPubSub():  ", recover)
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

		conn.Connection.SetWriteDeadline(time.Now().Add(time.Duration(10000) * time.Millisecond))
		conn.WriteLock.Lock()
		conn.Connection.WriteJSON(payload)
		conn.WriteLock.Unlock()
		unlocked = true

	}()

}

// BroadcastWebSocketData sends a message to all websocket connections
func BroadcastWebSocketData(data []byte) {

	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at WebSocketConnections():  ", recover)
			return
		}
	}()

	if !BroadcastSockets {
		return
	}
	WebSocketConnections.Range(func(key interface{}, value interface{}) bool {
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
}

// BroadcastWebSocketJSON sends a JSON message to all websocket connections
func BroadcastWebSocketJSON(v interface{}) {
	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at BroadcastWebSocketJSON():  ", recover)
			return
		}
	}()

	if !BroadcastSockets {
		return
	}
	WebSocketConnections.Range(func(key interface{}, value interface{}) bool {
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

			conn.Connection.SetWriteDeadline(time.Now().Add(time.Duration(10000) * time.Millisecond))
			conn.WriteLock.Lock()
			conn.Connection.WriteJSON(v)
			conn.WriteLock.Unlock()
			unlocked = true
		}()
		return true
	})

}

// PublishWebSocketJSON sends a JSON message to all websocket connections
func PublishWebSocketJSON(key string, v interface{}) {
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

	WebSocketConnections.Range(func(key interface{}, value interface{}) bool {
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

			conn.Connection.SetWriteDeadline(time.Now().Add(time.Duration(10000) * time.Millisecond))
			conn.WriteLock.Lock()
			conn.Connection.WriteJSON(payload)
			conn.WriteLock.Unlock()
			unlocked = true
			// deadLockChan <- 0
		}()
		return true
	})
}

// SetWebSocketTimeout sets the timeout for the websocket connections and will remove ones who havent sent a message within that time frame
func SetWebSocketTimeout(timeout int) {
	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at SetWebSocketTimeout():  ", recover)
			time.Sleep(time.Millisecond * 3000)
			SetWebSocketTimeout(timeout)
			return
		}
	}()

	// if CustomLog != nil {
	// 	CustomLog("app->SetWebSocketTimeout", "Checking for Web Socket Timeouts.")
	// }

	for {
		webSocketConnectionsMeta.Range(func(key interface{}, value interface{}) bool {
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
					deleteWebSocket(meta.GetConnection())
				}
			}
			return true
		})
		time.Sleep(time.Millisecond * 2500)
	}

}

func deleteWebSocket(c *WebSocketConnection) {

	go func() {
		defer func() {
			if recover := recover(); recover != nil {
				CustomLog("app->deleteWebSocket", "Panic Recovered at deleteWebSocket():  "+fmt.Sprintf("%+v", recover))
				return
			}
		}()

		if CustomLog != nil {
			CustomLog("app->deleteWebSocket", "Deleting Web Socket from client:  "+c.Connection.RemoteAddr().String())
		}
		c.Connection.Close()

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

	}()

}

// GetWebSocketMeta returns the meta data for a websocket connection
func GetWebSocketMeta(id string) (info *WebSocketConnectionMeta, ok bool) {
	result, ok := webSocketConnectionsMeta.Load(id)
	if ok {
		info = result.(*WebSocketConnectionMeta)
		return
	}
	return
}

// SetWebSocketMeta sets the meta data for a websocket connection
func SetWebSocketMeta(id string, info *WebSocketConnectionMeta) {
	webSocketConnectionsMeta.Store(id, info)
}

// RemoveWebSocketMeta removes the meta data for a websocket connection
func RemoveWebSocketMeta(id string) {
	webSocketConnectionsMeta.Delete(id)
}

// GetAllWebSocketMeta returns all the meta data for all websocket connections
func GetAllWebSocketMeta() (items *sync.Map) {
	return &webSocketConnectionsMeta
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
