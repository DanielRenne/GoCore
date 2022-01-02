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

func (obj *WebSocketConnectionMeta) SetTimeoutOverride(timeout int) {
	obj.TimeoutOverride.Set(timeout)
}

func (obj *WebSocketConnectionMeta) GetConnection() (conn *WebSocketConnection) {
	conn = obj.Conn
	return
}

type WebSocketConnectionCollection struct {
	sync.RWMutex
	Connections []*WebSocketConnection
}

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

type WebSocketCallbackSync struct {
	sync.RWMutex
	callbacks []WebSocketCallback
}

type ConcurrentWebSocketCallbackItem struct {
	Index    int
	Callback WebSocketCallback
}

func (self *WebSocketCallbackSync) Append(item WebSocketCallback) {
	self.RLock()
	defer self.RUnlock()
	self.callbacks = append(self.callbacks, item)
}

func (self *WebSocketCallbackSync) Iter() <-chan ConcurrentWebSocketCallbackItem {
	c := make(chan ConcurrentWebSocketCallbackItem)

	f := func() {
		self.Lock()
		defer self.Unlock()
		for index := range self.callbacks {
			value := self.callbacks[index]
			c <- ConcurrentWebSocketCallbackItem{index, value}
		}
		close(c)
	}
	go f()

	return c
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
var webSocketConnectionsMeta sync.Map
var WebSocketCallbacks sync.Map
var WebSocketRemovalCallback WebSocketRemoval

func init() {
	BroadcastSockets = true
}

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

func InitializeLite() (err error) {
	ginServer.InitializeLite(gin.ReleaseMode)
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
		webSocketHandler(c.Writer, c.Request, c)
	})

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
			webSocketHandler(c.Writer, c.Request, c)
		})
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

func initializeStaticRoutes() {

	ginServer.Router.GET("/swagger", func(c *gin.Context) {
		// c.Redirect(http.StatusMovedPermanently, "https://"+serverSettings.WebConfig.Application.Domain+":"+strconv.Itoa(serverSettings.WebConfig.Application.HttpsPort)+"/web/swagger/dist/index.html")

		ginServer.ReadHTMLFile(serverSettings.APP_LOCATION+"/web/swagger/dist/index.html", c)
	})
}

func RegisterWebSocketDataCallback(callback WebSocketCallback) {
	uuid, _ := extensions.NewUUID()
	WebSocketCallbacks.Store(uuid, callback)
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

		conn.Connection.SetWriteDeadline(time.Now().Add(time.Duration(10000) * time.Millisecond))
		conn.WriteLock.Lock()
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

		conn.Connection.SetWriteDeadline(time.Now().Add(time.Duration(10000) * time.Millisecond))
		conn.WriteLock.Lock()
		conn.Connection.WriteJSON(payload)
		conn.WriteLock.Unlock()
		unlocked = true

	}()

}

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
	return
}

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
