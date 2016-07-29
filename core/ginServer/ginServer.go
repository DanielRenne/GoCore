package ginServer

import (
	"github.com/fatih/color"
	// "github.com/gin-gonic/contrib/gzip"
	// "bufio"
	// "compress/gzip"
	"github.com/gin-gonic/gin"
	// "io"
	// "net"
	// "net/http"
	// "strings"
	"sync"
)

const (
	noWritten     = -1
	defaultStatus = 200
)

type routeGroup struct {
	enabled     bool
	routerGroup *gin.RouterGroup
}

var Router *gin.Engine

var groupRoutesSynced = struct {
	sync.RWMutex
	m map[string]routeGroup
}{m: make(map[string]routeGroup)}

type routerGroup struct {
	group  string
	route  string
	method string
	fp     func(*gin.Context)
}

var initializedRouterGroups []routerGroup
var hasInitialized bool

func Initialize(mode string) {
	gin.SetMode(mode)
	Router = gin.Default()
	hasInitialized = true

	for _, group := range initializedRouterGroups {
		AddRouterGroup(group.group, group.route, group.method, group.fp)
	}
	initializedRouterGroups = nil
}

func AddRouterGroup(group string, route string, method string, fp func(*gin.Context)) {

	if !hasInitialized {
		rgroup := routerGroup{group: group, route: route, method: method, fp: fp}
		initializedRouterGroups = append(initializedRouterGroups, rgroup)
		return
	}

	rg := getRouterSyncGroup(group)

	switch method {
	case "GET":
		rg.GET(route, fp)
		color.Green("Added GET to " + group + " with route " + route)
	case "POST":
		rg.POST(route, fp)
		color.Green("Added POST to " + group + " with route " + route)
	case "PUT":
		rg.PUT(route, fp)
		color.Green("Added PUT to " + group + " with route " + route)
	case "DELETE":
		rg.DELETE(route, fp)
		color.Green("Added DELETE to " + group + " with route " + route)
	case "PATCH":
		rg.PATCH(route, fp)
		color.Green("Added PATCH to " + group + " with route " + route)
	case "HEAD":
		rg.HEAD(route, fp)
		color.Green("Added HEAD to " + group + " with route " + route)
	case "OPTIONS":
		rg.OPTIONS(route, fp)
		color.Green("Added OPTIONS to " + group + " with route " + route)
	}
}

func getRouterSyncGroup(group string) *gin.RouterGroup {
	groupRoutesSynced.RLock()
	n := groupRoutesSynced.m[group]
	groupRoutesSynced.RUnlock()

	if n.enabled == false {
		r := Router.Group(group)
		rg := routeGroup{enabled: true, routerGroup: r}
		addRouterSyncGroup(group, rg)
		color.Green("Added Routing Group " + group)
		return rg.routerGroup
	}

	return n.routerGroup
}

func addRouterSyncGroup(group string, rg routeGroup) {
	groupRoutesSynced.Lock()
	groupRoutesSynced.m[group] = rg
	groupRoutesSynced.Unlock()
}

// type gzipResponseWriter struct {
// 	io.Writer
// 	http.ResponseWriter
// 	size   int
// 	status int
// }

// func (w *gzipResponseWriter) reset(writer http.ResponseWriter) {
// 	w.ResponseWriter = writer
// 	w.size = noWritten
// 	w.status = defaultStatus
// }

// func (w *gzipResponseWriter) WriteHeader(code int) {
// 	if code > 0 && w.status != code {
// 		w.status = code
// 	}
// }

// func (w *gzipResponseWriter) Write(b []byte) (int, error) {
// 	if "" == w.Header().Get("Content-Type") {
// 		// If no content type, apply sniffing algorithm to un-gzipped body.
// 		w.Header().Set("Content-Type", http.DetectContentType(b))
// 	}
// 	return w.Writer.Write(b)
// }

// func (w *gzipResponseWriter) WriteHeaderNow() {
// 	if !w.Written() {
// 		w.size = 0
// 		w.ResponseWriter.WriteHeader(w.status)
// 	}
// }

// func (w *gzipResponseWriter) WriteString(s string) (n int, err error) {
// 	w.WriteHeaderNow()
// 	n, err = io.WriteString(w.ResponseWriter, s)
// 	w.size += n
// 	return
// }

// func (w *gzipResponseWriter) Written() bool {
// 	return w.size != noWritten
// }

// func (w *gzipResponseWriter) Status() int {
// 	return w.status
// }

// func (w *gzipResponseWriter) Size() int {
// 	return w.size
// }

// // Implements the http.Hijacker interface
// func (w *gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
// 	if w.size < 0 {
// 		w.size = 0
// 	}
// 	return w.ResponseWriter.(http.Hijacker).Hijack()
// }

// // Implements the http.CloseNotify interface
// func (w *gzipResponseWriter) CloseNotify() <-chan bool {
// 	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
// }

// // Implements the http.Flush interface
// func (w *gzipResponseWriter) Flush() {
// 	w.ResponseWriter.(http.Flusher).Flush()
// }

// func makeGzipHandler() gin.HandlerFunc {

// 	return func(c *gin.Context) {
// 		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
// 			c.Next()
// 			return
// 		}
// 		c.Header("Content-Encoding", "gzip")

// 		if strings.Contains(c.Request.URL.String(), ".js") {
// 			c.Header("Content-Type", "application/javascript")
// 		}

// 		c.Next()
// 	}

// }
