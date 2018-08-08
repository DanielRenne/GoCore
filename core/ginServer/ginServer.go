package ginServer

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"

	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/fatih/color"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mmcloughlin/professor"
	"github.com/utrack/gin-csrf"
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
var ginCookieDomain string

func Initialize(mode string, cookieDomain string) {
	// Run a safe pprof localhost server.
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Panic Stack: " + string(debug.Stack()))
				log.Println("Recover Error:  " + fmt.Sprintf("%+v", r))
				return
			}
		}()
		professor.Launch("localhost:6897")
	}()
	gin.SetMode(mode)

	ginCookieDomain = cookieDomain

	if serverSettings.WebConfig.Application.CustomGinLogger {
		Router = gin.New()
		Router.Use(gin.Recovery())
	} else {
		Router = gin.Default()
	}

	store := sessions.NewCookieStore([]byte(serverSettings.WebConfig.Application.SessionKey))
	store.Options(sessions.Options{MaxAge: 86400 * serverSettings.WebConfig.Application.SessionExpirationDays,
		Secure: serverSettings.WebConfig.Application.SessionSecureCookie})

	if serverSettings.WebConfig.Application.SessionName != "" {
		Router.Use(sessions.Sessions(serverSettings.WebConfig.Application.SessionName, store))
	} else {
		Router.Use(sessions.Sessions("defaultSession", store))
	}

	//Protect from CSRF Hacking
	Router.Use(csrf.Middleware(csrf.Options{
		Secret: serverSettings.WebConfig.Application.CSRFSecret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	hasInitialized = true

	for _, group := range initializedRouterGroups {
		AddRouterGroup(group.group, group.route, group.method, group.fp)
	}
	initializedRouterGroups = nil
}

func InitializeLite(mode string) {
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
